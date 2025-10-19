package history

import (
	"sort"
	"time"
)

// renameTransition represents a single rename event.
type renameTransition struct {
	date    time.Time
	oldPath string
	newPath string
}

// renameChain tracks all rename events and provides methods to resolve paths.
type renameChain struct {
	// All transitions sorted by date (ascending)
	transitions []renameTransition
	// Index: oldPath -> list of transitions from that path (sorted by date)
	byOldPath map[string][]renameTransition
}

// buildRenameChain extracts and indexes all rename events from the input events.
func buildRenameChain(events []*Event) *renameChain {
	var transitions []renameTransition

	for _, event := range events {
		if event.File.OldPath != "" {
			date, err := time.Parse("2006-01-02 15:04:05 -0700", event.Date)
			if err != nil {
				// Skip invalid dates
				continue
			}

			transitions = append(transitions, renameTransition{
				date:    date,
				oldPath: event.File.OldPath,
				newPath: event.File.Path,
			})
		}
	}

	// Sort all transitions by date (ascending)
	sort.Slice(transitions, func(i, j int) bool {
		return transitions[i].date.Before(transitions[j].date)
	})

	// Build index by oldPath
	byOldPath := make(map[string][]renameTransition)
	for _, t := range transitions {
		byOldPath[t.oldPath] = append(byOldPath[t.oldPath], t)
	}

	return &renameChain{
		transitions: transitions,
		byOldPath:   byOldPath,
	}
}

// resolveFinalPath determines the final (most recent) path for a file
// given its path at a specific commit time.
// It follows renames forward in time to find where the file ultimately ends up.
func (rc *renameChain) resolveFinalPath(pathAtCommit string, commitDate time.Time) string {
	currentPath := pathAtCommit

	for {
		// Find transitions from currentPath that occur after commitDate
		transitions := rc.byOldPath[currentPath]
		if len(transitions) == 0 {
			return currentPath
		}

		// Find the first transition that happens after commitDate
		found := false
		for _, t := range transitions {
			if t.date.After(commitDate) {
				currentPath = t.newPath
				commitDate = t.date
				found = true
				break
			}
		}

		if !found {
			return currentPath
		}
	}
}

// getHistoricalPaths returns all historical paths for a given current path,
// from newest to oldest.
func (rc *renameChain) getHistoricalPaths(currentPath string) []string {
	paths := []string{currentPath}
	path := currentPath

	// Walk backwards through transitions (newest to oldest)
	for i := len(rc.transitions) - 1; i >= 0; i-- {
		t := rc.transitions[i]
		if t.newPath == path {
			path = t.oldPath
			paths = append(paths, path)
		}
	}

	return paths
}
