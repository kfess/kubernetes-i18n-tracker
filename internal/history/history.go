package history

import (
	"sort"
	"time"
)

// Tracker builds and maintains file histories from git events.
type History struct {
	// Map from final file path to its commit history (sorted by date, oldest to newest)
	histories map[string][]*Commit
	// Rename chain for path resolution
	renames *renameChain
}

// Build creates a new history from a list of events.
// It processes all events, resolves renames, groups commits by final file path,
// and sorts them chronologically.
func Build(events []*Event) *History {
	history := &History{
		histories: make(map[string][]*Commit),
		renames:   buildRenameChain(events),
	}

	for _, event := range events {
		commit := history.eventToCommit(event)
		if commit == nil {
			continue
		}

		finalPath := history.renames.resolveFinalPath(commit.Path, commit.Date)
		history.histories[finalPath] = append(history.histories[finalPath], commit)
	}

	for _, commits := range history.histories {
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Date.Before(commits[j].Date)
		})
	}

	return history
}

// eventToCommit converts an Event to a Commit.
func (h *History) eventToCommit(event *Event) *Commit {
	date, err := time.Parse("2006-01-02 15:04:05 -0700", event.Date)
	if err != nil {
		return nil
	}

	insertions := 0
	if event.File.Insertions != nil {
		insertions = *event.File.Insertions
	}

	deletions := 0
	if event.File.Deletions != nil {
		deletions = *event.File.Deletions
	}

	commit := &Commit{
		Hash:        event.Hash,
		Author:      event.Author,
		Date:        date,
		Message:     event.Message,
		Path:        event.File.Path,
		Insertions:  insertions,
		Deletions:   deletions,
		RenamedFrom: event.File.OldPath,
	}

	return commit
}

// GetCommits returns the commit history for a specific file path.
// Returns nil if the path has no history.
func (h *History) GetCommits(path string) []*Commit {
	return h.histories[path]
}

// AllPaths returns all file paths that have history.
func (h *History) AllPaths() []string {
	paths := make([]string, 0, len(h.histories))
	for path := range h.histories {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

// Stats calculates statistics for a specific file path.
func (h *History) Stats(path string) *Stats {
	commits := h.histories[path]
	if len(commits) == 0 {
		return nil
	}

	stats := &Stats{
		TotalCommits: len(commits),
		FirstCommit:  commits[0].Date,
		LatestCommit: commits[len(commits)-1].Date,
	}

	renameCount := 0
	for _, commit := range commits {
		stats.TotalInsertions += commit.Insertions
		stats.TotalDeletions += commit.Deletions
		if commit.RenamedFrom != "" {
			renameCount++
		}
	}
	stats.RenameCount = renameCount

	return stats
}

// GetHistoricalPaths returns all historical paths for a given current path,
// from newest to oldest.
func (h *History) GetHistoricalPaths(currentPath string) []string {
	return h.renames.getHistoricalPaths(currentPath)
}

// GetOldestCommit returns the oldest commit for a specific file path.
func (h *History) GetOldestCommit(path string) *Commit {
	commits := h.histories[path]
	if len(commits) == 0 {
		return nil
	}

	return commits[0]
}

// GetLatestCommit returns the latest commit for a specific file path.
func (h *History) GetLatestCommit(path string) *Commit {
	commits := h.histories[path]
	if len(commits) == 0 {
		return nil
	}

	return commits[len(commits)-1]
}

// GetCommitsAfter returns all commits for a file path that occurred after the specified time.
// Returns nil if the path has no history or no commits match the time filter.
func (h *History) GetCommitsAfter(path string, after time.Time) []*Commit {
	commits := h.histories[path]
	if len(commits) == 0 {
		return nil
	}

	var result []*Commit
	for _, c := range commits {
		if c.Date.After(after) {
			result = append(result, c)
		}
	}

	return result
}

// GetCommitBeforeOrAt finds the latest commit for a file path that occurred before or at the specified time.
// Returns nil if no such commit exists or the path has no history.
func (h *History) GetCommitBeforeOrAt(path string, date time.Time) *Commit {
	commits := h.histories[path]
	if len(commits) == 0 {
		return nil
	}

	var result *Commit
	for _, commit := range commits {
		if commit.Date.Before(date) || commit.Date.Equal(date) {
			result = commit
		} else {
			// Since commits are sorted chronologically, we can break early
			break
		}
	}
	return result
}
