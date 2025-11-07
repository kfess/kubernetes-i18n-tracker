package pr

import (
	"sort"
)

// Index provides efficient lookup of pull requests by file path.
type Index struct {
	// Map from file path to list of PRs that modified that file
	byFile map[string][]*PullRequest

	// Map from PR number to PR for quick lookup
	byNumber map[int]*PullRequest
}

// BuildPRIndex creates an Index from a list of pull requests.
// It indexes PRs by their modified files for quick lookup.
func BuildPRIndex(prs []PullRequest) *Index {
	index := &Index{
		byFile:   make(map[string][]*PullRequest),
		byNumber: make(map[int]*PullRequest),
	}

	for i := range prs {
		prPtr := &prs[i]

		// Index by number
		index.byNumber[prPtr.Number] = prPtr

		for _, file := range prPtr.Files {
			index.byFile[file] = append(index.byFile[file], prPtr)
		}
	}

	// Sort PRs by number (ascending) for each file
	// This ensures latest PR is at the end
	for _, prs := range index.byFile {
		sort.Slice(prs, func(i, j int) bool {
			return prs[i].Number < prs[j].Number
		})
	}

	return index
}

// GetPRsForFile returns all PRs that modified the given file path.
// Returns nil if no PRs are found.
func (idx *Index) GetPRsForFile(path string) []*PullRequest {
	return idx.byFile[path]
}

// GetLatestPR returns the most recent PR that modified the given file.
// Returns nil if no PRs are found.
func (idx *Index) GetLatestPR(path string) *PullRequest {
	prs := idx.byFile[path]
	if len(prs) == 0 {
		return nil
	}
	return prs[len(prs)-1]
}

// GetRecentPRs returns the N most recent PRs that modified the given file.
// If there are fewer PRs than n, returns all available PRs.
func (idx *Index) GetRecentPRs(path string, n int) []*PullRequest {
	prs := idx.byFile[path]
	if len(prs) == 0 {
		return nil
	}

	start := len(prs) - n
	if start < 0 {
		start = 0
	}

	return prs[start:]
}

// HasPRs returns true if there are any PRs for the given file.
func (idx *Index) HasPRs(path string) bool {
	return len(idx.byFile[path]) > 0
}

// TotalFiles returns the number of unique files across all indexed PRs.
func (idx *Index) TotalFiles() int {
	return len(idx.byFile)
}

// TotalPRs returns the total number of unique PRs in the index.
func (idx *Index) TotalPRs() int {
	return len(idx.byNumber)
}

// GetPRByNumber returns a PR by its number.
// Returns nil if the PR is not found.
func (idx *Index) GetPRByNumber(number int) *PullRequest {
	return idx.byNumber[number]
}
