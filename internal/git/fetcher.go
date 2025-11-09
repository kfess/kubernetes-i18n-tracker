package git

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
)

// FetchOptions contains options for fetching git history.
type FetchOptions struct {
	// Path to the git repository
	RepoPath string

	// Number of parallel workers for fetching history
	Workers int

	// Valid file extensions to consider
	ValidExtensions []string
}

// Fetcher fetches git history for files in a repository.
type Fetcher struct {
	// Fetch options
	options FetchOptions
}

// NewFetcher creates a new Fetcher with the given options.
func NewFetcher(options FetchOptions) *Fetcher {
	if options.Workers <= 0 {
		options.Workers = 4
	}
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
func (f *Fetcher) FetchHistory(ctx context.Context) ([]*Event, error) {
	logger.Info("Fetching commit hashes from git log --first-parent...")

	// Get all commits that touched content/
	commits, err := f.getFirstParentCommits(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	logger.Info(fmt.Sprintf("Found %d commits to process", len(commits)))

	var (
		allEvents        []*Event
		processedCommits int64
		mu               sync.Mutex
		wg               sync.WaitGroup
	)

	// Create a semaphore to limit concurrent git processes
	semaphore := make(chan struct{}, f.options.Workers)

	for _, commitHash := range commits {
		wg.Add(1)
		go func(hash string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			events, err := f.processCommit(ctx, hash)
			if err != nil {
				logger.Warn(fmt.Sprintf("Failed to process commit %s: %v", hash, err))
				return
			}

			// Thread-safe append
			mu.Lock()
			allEvents = append(allEvents, events...)
			mu.Unlock()

			// Progress logging
			processed := atomic.AddInt64(&processedCommits, 1)
			if processed%100 == 0 || processed == int64(len(commits)) {
				percent := (processed * 100) / int64(len(commits))
				logger.Info(fmt.Sprintf("Progress: %d/%d commits (%d%%)", processed, len(commits), percent))
			}
		}(commitHash)
	}

	wg.Wait()

	logger.Info(fmt.Sprintf("Completed: %d commits processed, %d events generated", len(commits), len(allEvents)))

	return allEvents, nil
}

// getFirstParentCommits returns all commit hashes from git log --first-parent main.
func (f *Fetcher) getFirstParentCommits(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx,
		"git", "log", "--first-parent", "main",
		"--pretty=format:%H", "--", "content/")
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var commits []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			commits = append(commits, line)
		}
	}

	return commits, scanner.Err()
}

// isMergeCommit checks if a commit has more than 2 parents (merge commit).
func (f *Fetcher) isMergeCommit(ctx context.Context, commitHash string) (bool, error) {
	cmd := exec.CommandContext(ctx,
		"git", "rev-list", "--parents", "-n", "1", commitHash)
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("git rev-list failed: %w", err)
	}

	// Count parents: "hash parent1 parent2 ..." -> split and count
	parts := bytes.Fields(output)
	return len(parts) > 2, nil
}

// processCommit processes a single commit and returns events.
func (f *Fetcher) processCommit(ctx context.Context, commitHash string) ([]*Event, error) {
	// Check if it's a merge commit
	isMerge, err := f.isMergeCommit(ctx, commitHash)
	if err != nil {
		return nil, err
	}

	if isMerge {
		return f.processMergeCommit(ctx, commitHash)
	}

	return f.processRegularCommit(ctx, commitHash)
}

// processRegularCommit processes a regular (non-merge) commit.
func (f *Fetcher) processRegularCommit(ctx context.Context, commitHash string) ([]*Event, error) {
	cmd := exec.CommandContext(ctx,
		"git", "show", "--pretty=format:%H\x1F%an\x1F%ad\x1F%s",
		"--numstat", "--date=iso", commitHash, "--", "content/")
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git show failed: %w", err)
	}

	return ParseGitLog(output)
}

// processMergeCommit processes a merge commit.
func (f *Fetcher) processMergeCommit(ctx context.Context, commitHash string) ([]*Event, error) {
	// Get diff between commit^1 and commit
	diffOutput, err := f.getDiffNumstat(ctx, fmt.Sprintf("%s^1", commitHash), commitHash)
	if err != nil {
		return nil, err
	}

	var events []*Event

	// Parse each file in the diff
	scanner := bufio.NewScanner(bytes.NewReader(diffOutput))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse: "insertions\tdeletions\tfilepath"
		parts := bytes.Split([]byte(line), []byte("\t"))
		if len(parts) < 3 {
			continue
		}

		file := string(parts[2])

		// Extract the new path from rename notation: {old => new} or just use the path as-is
		searchPath := extractNewPath(file)

		// Find the actual commit from feature branch that modified this file
		// Use git log on the feature branch (^2) to find the last commit touching this file
		actualCommit, err := f.getFeatureBranchCommitForFile(ctx, commitHash, searchPath)
		if err != nil || actualCommit == "" {
			logger.Warn(fmt.Sprintf("Failed to find feature branch commit for %s in %s: %v, using merge commit", file, commitHash, err))
			actualCommit = commitHash
		}

		// Get the correct commit message for this specific file's commit
		commitMessage, err := f.getCommitSubject(ctx, actualCommit)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to get commit message for %s: %v, using fallback", actualCommit, err))
			commitMessage = ""
		}

		// Get commit info - use the actual commit that modified this file
		event, err := f.createEventFromCommit(ctx, actualCommit, file, commitMessage, parts)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to create event for %s: %v", file, err))
			continue
		}

		if event != nil {
			events = append(events, event)
		}
	}

	return events, scanner.Err()
}

// getCommitSubject returns the subject (first line of message) of a commit.
func (f *Fetcher) getCommitSubject(ctx context.Context, commitHash string) (string, error) {
	cmd := exec.CommandContext(ctx,
		"git", "log", "--pretty=format:%s", "-n", "1", commitHash)
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w", err)
	}

	return string(bytes.TrimSpace(output)), nil
}

// extractNewPath extracts the new path from a rename notation.
// Examples:
//   - "path/{old => new}/file.md" -> "path/new/file.md"
//   - "regular/path.md" -> "regular/path.md"
func extractNewPath(path string) string {
	// Check if path contains rename notation: {old => new}
	if !bytes.Contains([]byte(path), []byte("{")) {
		return path
	}

	// Use regex or simple string manipulation
	// Pattern: {old => new}
	start := bytes.Index([]byte(path), []byte("{"))
	end := bytes.Index([]byte(path), []byte("}"))
	if start == -1 || end == -1 || end <= start {
		return path
	}

	// Extract the part inside braces
	inside := path[start+1 : end]
	parts := bytes.Split([]byte(inside), []byte(" => "))
	if len(parts) != 2 {
		return path
	}

	// Replace {old => new} with just new
	newPath := path[:start] + string(parts[1]) + path[end+1:]
	return newPath
}

// getFeatureBranchCommitForFile finds the actual commit from the feature branch that modified a file.
// For a merge commit M, this looks in M^2 (feature branch) for commits that touched the file.
func (f *Fetcher) getFeatureBranchCommitForFile(ctx context.Context, mergeCommit string, file string) (string, error) {
	// Get commits from feature branch (^2) that are not in main (^1)
	// This gives us all commits in the PR/feature branch
	cmd := exec.CommandContext(ctx,
		"git", "log", "--pretty=format:%H", "--follow",
		fmt.Sprintf("%s^2", mergeCommit),
		fmt.Sprintf("^%s^1", mergeCommit),
		"--", file)
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w", err)
	}

	// Return the first (most recent) commit that touched this file
	scanner := bufio.NewScanner(bytes.NewReader(output))
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", fmt.Errorf("no commits found for file %s in feature branch", file)
}

// getDiffNumstat returns the numstat diff between two commits for content/ files.
func (f *Fetcher) getDiffNumstat(ctx context.Context, from string, to string) ([]byte, error) {
	cmd := exec.CommandContext(ctx,
		"git", "diff", "--numstat", from, to, "--", "content/")
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w", err)
	}

	return output, nil
}

// createEventFromCommit creates an Event from commit info and numstat line.
func (f *Fetcher) createEventFromCommit(ctx context.Context, commitHash string, file string, message string, numstatParts [][]byte) (*Event, error) {
	// Get commit author and date
	cmd := exec.CommandContext(ctx,
		"git", "show", "--pretty=format:%an\x1F%ad",
		"--date=iso", "-s", commitHash)
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git show failed: %w", err)
	}

	parts := bytes.Split(bytes.TrimSpace(output), []byte("\x1F"))
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected git show output")
	}

	author := string(parts[0])
	date := string(parts[1])

	// Parse numstat
	insertionsStr := string(numstatParts[0])
	deletionsStr := string(numstatParts[1])

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
		Hash:    commitHash,
		Author:  author,
		Date:    date,
		Message: message,
		File: FileInfo{
			Path:       newPath,
			Insertions: insertions,
			Deletions:  deletions,
			OldPath:    oldPath,
		},
	}, nil
}
