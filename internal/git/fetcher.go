package git

import (
	"container/heap"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
)

const (
	// Separators used in git pretty formats
	recordSeparator = "\x1E"
	fieldSeparator  = "\x1F"
)

// FetchOptions contains options for fetching git history.
type FetchOptions struct {
	// Path to the git repository
	RepoPath string

	// Valid file extensions to consider
	ValidExtensions []string
}

// Fetcher fetches git history for files in a repository.
type Fetcher struct {
	// Fetch options
	options FetchOptions
}

// commitMeta contains the commit information used to build events.
type commitMeta struct {
	hash    string
	author  string
	date    string
	message string
}

// firstParentCommit is a commit on the first-parent chain of main with its
// numstat lines (diff against the first parent) for content/ files.
type firstParentCommit struct {
	commitMeta
	numstat []string
}

// contentCommit is a non-merge commit with the content/ files it touched.
type contentCommit struct {
	commitMeta
	files []string
}

// commitGraph holds the parents and commit timestamps of all commits reachable from main.
type commitGraph struct {
	// Tip of main
	head       string
	parents    map[string][]string
	timestamps map[string]int64
}

// NewFetcher creates a new Fetcher with the given options.
func NewFetcher(options FetchOptions) *Fetcher {
	return &Fetcher{
		options: options,
	}
}

// UpdateRepo updates the repository to the latest version.
func (f *Fetcher) UpdateRepo() error {
	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = f.options.RepoPath

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update repository: %w, output: %s", err, string(output))
	}

	return nil
}

// FetchHistory fetches git history using git log --first-parent approach.
//
// The whole history is loaded with a few git invocations and merge commits are
// attributed to feature branch commits in memory, instead of spawning git
// processes per commit and per file.
func (f *Fetcher) FetchHistory(ctx context.Context) ([]*Event, error) {
	logger.Info("Fetching commits from git log --first-parent...")

	// Get all first-parent commits that touched content/
	commits, err := f.getFirstParentCommits(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	logger.Info(fmt.Sprintf("Found %d commits to process", len(commits)))

	logger.Info("Loading commit graph...")
	graph, err := f.getCommitGraph(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit graph: %w", err)
	}

	logger.Info("Loading files changed by each commit...")
	contentCommits, err := f.getContentCommits(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	mergeCommits := make(map[string]*firstParentCommit)
	for _, commit := range commits {
		if len(graph.parents[commit.hash]) > 1 {
			mergeCommits[commit.hash] = commit
		}
	}

	// Resolve merge commits while walking the first-parent chain from the oldest commit
	mergeEvents := resolveMergeCommits(graph, mergeCommits, contentCommits)

	var allEvents []*Event
	for _, commit := range commits {
		if _, isMerge := mergeCommits[commit.hash]; isMerge {
			allEvents = append(allEvents, mergeEvents[commit.hash]...)
			continue
		}

		events, err := processRegularCommit(commit)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to process commit %s: %v", commit.hash, err))
			continue
		}
		allEvents = append(allEvents, events...)
	}

	logger.Info(fmt.Sprintf("Completed: %d commits processed, %d events generated", len(commits), len(allEvents)))

	return allEvents, nil
}

// runGit runs a git command in the repository and returns its stdout.
func (f *Fetcher) runGit(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w", args[0], err)
	}

	return string(output), nil
}

// splitRecord splits a git log record into its header fields and the non-empty lines following the header.
func splitRecord(record string) ([]string, []string) {
	header, body, _ := strings.Cut(record, "\n")

	var lines []string
	for _, line := range strings.Split(body, "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return strings.Split(header, fieldSeparator), lines
}

// getFirstParentCommits returns all commits from git log --first-parent main that touched content/,
// with the numstat against their first parent.
func (f *Fetcher) getFirstParentCommits(ctx context.Context) ([]*firstParentCommit, error) {
	// With -m --first-parent, merge commits are diffed against their first parent only
	output, err := f.runGit(ctx,
		"log", "--first-parent", "-m", "main",
		"--pretty=format:"+recordSeparator+"%H"+fieldSeparator+"%an"+fieldSeparator+"%ad"+fieldSeparator+"%s",
		"--numstat", "--date=iso", "--", "content/")
	if err != nil {
		return nil, err
	}

	var commits []*firstParentCommit
	for _, record := range strings.Split(output, recordSeparator) {
		fields, lines := splitRecord(record)
		if len(fields) < 4 {
			continue
		}

		commits = append(commits, &firstParentCommit{
			commitMeta: commitMeta{hash: fields[0], author: fields[1], date: fields[2], message: fields[3]},
			numstat:    lines,
		})
	}

	return commits, nil
}

// getCommitGraph returns the parents and commit timestamps of all commits reachable from main.
func (f *Fetcher) getCommitGraph(ctx context.Context) (*commitGraph, error) {
	// Each line: "timestamp hash parent1 parent2 ..."
	output, err := f.runGit(ctx, "rev-list", "--timestamp", "--parents", "main")
	if err != nil {
		return nil, err
	}

	graph := &commitGraph{
		parents:    make(map[string][]string),
		timestamps: make(map[string]int64),
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		timestamp, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unexpected git rev-list output: %q", line)
		}

		hash := fields[1]
		if graph.head == "" {
			graph.head = hash
		}
		graph.parents[hash] = fields[2:]
		graph.timestamps[hash] = timestamp
	}

	return graph, nil
}

// getContentCommits returns all non-merge commits reachable from main that touched content/,
// with the files they touched.
func (f *Fetcher) getContentCommits(ctx context.Context) (map[string]*contentCommit, error) {
	// Renames are listed as both the old and the new path, and paths are not quoted
	output, err := f.runGit(ctx,
		"-c", "core.quotePath=false",
		"log", "main", "--no-merges", "--full-history", "--no-renames",
		"--pretty=format:"+recordSeparator+"%H"+fieldSeparator+"%an"+fieldSeparator+"%ad"+fieldSeparator+"%s",
		"--name-only", "--date=iso", "--", "content/")
	if err != nil {
		return nil, err
	}

	commits := make(map[string]*contentCommit)
	for _, record := range strings.Split(output, recordSeparator) {
		fields, lines := splitRecord(record)
		if len(fields) < 4 {
			continue
		}

		commits[fields[0]] = &contentCommit{
			commitMeta: commitMeta{hash: fields[0], author: fields[1], date: fields[2], message: fields[3]},
			files:      lines,
		}
	}

	return commits, nil
}

// processRegularCommit processes a regular (non-merge) commit.
func processRegularCommit(commit *firstParentCommit) ([]*Event, error) {
	header := strings.Join([]string{commit.hash, commit.author, commit.date, commit.message}, fieldSeparator)
	lines := append([]string{header}, commit.numstat...)

	return ParseGitLog([]byte(strings.Join(lines, "\n")))
}

// resolveMergeCommits returns the events of the given merge commits, keyed by merge commit hash.
//
// The first-parent chain of main is walked from the oldest commit while tracking the commits
// reachable from the current position. For a merge commit M, the commits reachable from M^2 that
// have not been visited yet are exactly the feature branch commits (M^2 ^M^1).
func resolveMergeCommits(graph *commitGraph, mergeCommits map[string]*firstParentCommit, contentCommits map[string]*contentCommit) map[string][]*Event {
	var chain []string
	for hash := graph.head; hash != ""; {
		chain = append(chain, hash)
		parents := graph.parents[hash]
		if len(parents) == 0 {
			break
		}
		hash = parents[0]
	}

	visited := make(map[string]bool, len(graph.parents))
	events := make(map[string][]*Event, len(mergeCommits))

	for i := len(chain) - 1; i >= 0; i-- {
		hash := chain[i]
		visited[hash] = true

		parents := graph.parents[hash]
		if len(parents) < 2 {
			continue
		}

		mergeCommit, ok := mergeCommits[hash]
		if !ok {
			// Not a merge commit to process, only mark its branches as visited
			for _, parent := range parents[1:] {
				graph.walk(parent, visited, nil)
			}
			continue
		}

		// Find the most recent feature branch (^2) commit touching each file
		lastTouched := make(map[string]*contentCommit)
		graph.walk(parents[1], visited, func(hash string) {
			commit, ok := contentCommits[hash]
			if !ok {
				return
			}
			for _, file := range commit.files {
				if _, found := lastTouched[file]; !found {
					lastTouched[file] = commit
				}
			}
		})
		for _, parent := range parents[2:] {
			graph.walk(parent, visited, nil)
		}

		events[hash] = processMergeCommit(mergeCommit, lastTouched)
	}

	return events
}

// processMergeCommit processes a merge commit.
// lastTouched maps a file to the most recent feature branch commit that modified it.
func processMergeCommit(mergeCommit *firstParentCommit, lastTouched map[string]*contentCommit) []*Event {
	var events []*Event

	// Parse each file in the diff between commit^1 and commit
	for _, line := range mergeCommit.numstat {
		// Parse: "insertions\tdeletions\tfilepath"
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}

		file := parts[2]

		// Extract the new path from rename notation: {old => new} or just use the path as-is
		searchPath := extractNewPath(file)

		// Use the actual commit from feature branch that modified this file
		actualCommit := mergeCommit.commitMeta
		// "path/{old => }/file.md" leaves a double slash, which git treats as a single one
		if commit, ok := lastTouched[strings.ReplaceAll(searchPath, "//", "/")]; ok {
			actualCommit = commit.commitMeta
		} else {
			logger.Warn(fmt.Sprintf("Failed to find feature branch commit for %s in %s: no commits found for file %s in feature branch, using merge commit", file, mergeCommit.hash, searchPath))
		}

		actualCommit.message = strings.TrimSpace(actualCommit.message)

		events = append(events, createEvent(actualCommit, file, parts[0], parts[1]))
	}

	return events
}

// commitQueue is a priority queue of commits ordered like git log: newest commit timestamp first.
type commitQueue []queuedCommit

type queuedCommit struct {
	hash      string
	timestamp int64
	// Insertion order, used to break ties
	order int
}

func (q commitQueue) Len() int { return len(q) }
func (q commitQueue) Less(i, j int) bool {
	if q[i].timestamp != q[j].timestamp {
		return q[i].timestamp > q[j].timestamp
	}
	return q[i].order < q[j].order
}
func (q commitQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *commitQueue) Push(x any)   { *q = append(*q, x.(queuedCommit)) }
func (q *commitQueue) Pop() any {
	old := *q
	item := old[len(old)-1]
	*q = old[:len(old)-1]
	return item
}

// walk visits the commits reachable from start that have not been visited yet,
// in the same order as git log, and marks them as visited.
func (g *commitGraph) walk(start string, visited map[string]bool, visit func(hash string)) {
	if visited[start] {
		return
	}

	queue := &commitQueue{}
	order := 0
	push := func(hash string) {
		visited[hash] = true
		heap.Push(queue, queuedCommit{hash: hash, timestamp: g.timestamps[hash], order: order})
		order++
	}

	push(start)
	for queue.Len() > 0 {
		commit := heap.Pop(queue).(queuedCommit)
		if visit != nil {
			visit(commit.hash)
		}

		for _, parent := range g.parents[commit.hash] {
			if !visited[parent] {
				push(parent)
			}
		}
	}
}

// extractNewPath extracts the new path from a rename notation.
// Examples:
//   - "path/{old => new}/file.md" -> "path/new/file.md"
//   - "regular/path.md" -> "regular/path.md"
func extractNewPath(path string) string {
	// Check if path contains rename notation: {old => new}
	start := strings.Index(path, "{")
	end := strings.Index(path, "}")
	if start == -1 || end == -1 || end <= start {
		return path
	}

	// Extract the part inside braces
	parts := strings.Split(path[start+1:end], " => ")
	if len(parts) != 2 {
		return path
	}

	// Replace {old => new} with just new
	return path[:start] + parts[1] + path[end+1:]
}

// createEvent creates an Event from commit info and a numstat line.
func createEvent(commit commitMeta, file string, insertionsStr string, deletionsStr string) *Event {
	var insertions, deletions *int
	if insertionsStr != "-" {
		if val, err := strconv.Atoi(insertionsStr); err == nil {
			insertions = &val
		}
	}
	if deletionsStr != "-" {
		if val, err := strconv.Atoi(deletionsStr); err == nil {
			deletions = &val
		}
	}

	// Parse rename path
	oldPath, newPath := parseRenamePath(file)

	return &Event{
		Hash:    commit.hash,
		Author:  commit.author,
		Date:    commit.date,
		Message: commit.message,
		File: FileInfo{
			Path:       newPath,
			Insertions: insertions,
			Deletions:  deletions,
			OldPath:    oldPath,
		},
	}
}
