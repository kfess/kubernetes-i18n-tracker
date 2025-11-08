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
	RepoPath        string   // Path to the git repository
	Workers         int      // Number of parallel workers for fetching history
	ValidExtensions []string // Valid file extensions to consider
}

// Fetcher fetches git history for files in a repository.
type Fetcher struct {
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

// isAncestor checks if commit is an ancestor of base.
func (f *Fetcher) isAncestor(ctx context.Context, commit string, base string) (bool, error) {
	cmd := exec.CommandContext(ctx,
		"git", "merge-base", "--is-ancestor", commit, base)
	cmd.Dir = f.options.RepoPath

	err := cmd.Run()
	return err == nil, nil
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
	// Get the first commit message from the feature branch
	firstCommitMessage, err := f.getFirstFeatureBranchCommitMessage(ctx, commitHash)
	if err != nil {
		// Fallback to merge commit subject if we can't get feature branch message
		logger.Warn(fmt.Sprintf("Failed to get first feature branch commit message for %s: %v, using merge commit message", commitHash, err))
		firstCommitMessage, err = f.getCommitSubject(ctx, commitHash)
		if err != nil {
			return nil, err
		}
	}

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

		// Find the last commit that actually modified this file
		lastCommit, err := f.findLastCommit(ctx, commitHash, file, "")
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to find last commit for %s in %s: %v", file, commitHash, err))
			lastCommit = commitHash
		}

		// Get commit info - use lastCommit hash but firstCommitMessage from feature branch
		event, err := f.createEventFromCommit(ctx, lastCommit, file, firstCommitMessage, parts)
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

// getFirstFeatureBranchCommitMessage gets the message of the first (oldest) commit from the feature branch.
func (f *Fetcher) getFirstFeatureBranchCommitMessage(ctx context.Context, mergeCommit string) (string, error) {
	// Get commits from feature branch (^2) excluding main branch (^1)
	cmd := exec.CommandContext(ctx,
		"git", "log", "--pretty=format:%s", "--reverse",
		fmt.Sprintf("%s^2", mergeCommit),
		fmt.Sprintf("^%s^1", mergeCommit))
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w", err)
	}

	lines := bytes.Split(output, []byte("\n"))
	if len(lines) > 0 && len(lines[0]) > 0 {
		return string(lines[0]), nil
	}

	return "", fmt.Errorf("no commits found in feature branch")
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

// findLastCommit finds the last commit that modified a file in a merge.
func (f *Fetcher) findLastCommit(ctx context.Context, mergeCommit string, file string, originalMainParent string) (string, error) {
	// Set original main parent on first call
	if originalMainParent == "" {
		originalMainParent = fmt.Sprintf("%s^1", mergeCommit)
	}

	// Search from ^1 excluding ^2
	fromFirst, err := f.getCommitFromBranch(ctx, mergeCommit, file, "^1", "^2")
	if err != nil {
		return mergeCommit, err
	}

	// Search from ^2 excluding ^1
	fromSecond, err := f.getCommitFromBranch(ctx, mergeCommit, file, "^2", "^1")
	if err != nil {
		return mergeCommit, err
	}

	// Both empty
	if fromFirst == "" && fromSecond == "" {
		return mergeCommit, nil
	}

	// Select the commit from feature branch (not in main)
	last := ""

	if fromFirst != "" && fromSecond != "" {
		// Both found - choose the one not in main
		firstInMain, _ := f.isAncestor(ctx, fromFirst, originalMainParent)
		secondInMain, _ := f.isAncestor(ctx, fromSecond, originalMainParent)

		if firstInMain && !secondInMain {
			last = fromSecond
		} else if !firstInMain && secondInMain {
			last = fromFirst
		} else if !firstInMain && !secondInMain {
			// Both not in main - choose newer
			firstDate, _ := f.getCommitDate(ctx, fromFirst)
			secondDate, _ := f.getCommitDate(ctx, fromSecond)
			if secondDate > firstDate {
				last = fromSecond
			} else {
				last = fromFirst
			}
		}
	} else if fromFirst != "" {
		inMain, _ := f.isAncestor(ctx, fromFirst, originalMainParent)
		if !inMain {
			last = fromFirst
		}
	} else if fromSecond != "" {
		inMain, _ := f.isAncestor(ctx, fromSecond, originalMainParent)
		if !inMain {
			last = fromSecond
		}
	}

	// No feature branch commit found
	if last == "" {
		return mergeCommit, nil
	}

	// If it's also a merge commit, recurse
	isMerge, _ := f.isMergeCommit(ctx, last)
	if isMerge {
		return f.findLastCommit(ctx, last, file, originalMainParent)
	}

	return last, nil
}

// getCommitFromBranch gets the first commit for a file from a specific parent.
func (f *Fetcher) getCommitFromBranch(ctx context.Context, mergeCommit, file, parent, exclude string) (string, error) {
	cmd := exec.CommandContext(ctx,
		"git", "log", "--pretty=format:%H",
		fmt.Sprintf("%s%s", mergeCommit, parent),
		fmt.Sprintf("^%s%s", mergeCommit, exclude),
		"--", file)
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return "", nil // Not an error, just no commits found
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", nil
}

// getCommitDate returns the commit date as unix timestamp.
func (f *Fetcher) getCommitDate(ctx context.Context, commitHash string) (int64, error) {
	cmd := exec.CommandContext(ctx,
		"git", "log", "--pretty=format:%at", "-n", "1", commitHash)
	cmd.Dir = f.options.RepoPath

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(bytes.TrimSpace(output)), 10, 64)
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
