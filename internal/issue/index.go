package issue

// Index provides efficient lookup of issues by file path.
type Index struct {
	issuesByFile map[string][]*Issue
}

// NewIndex creates a new Index from a list of issues.
func NewIndex(issues []Issue, existingPaths map[string]bool) *Index {
	issuesByFile := make(map[string][]*Issue)

	for i := range issues {
		issue := &issues[i]
		lang := GuessLanguage(*issue)
		if lang == "" {
			continue
		}

		path := GuessPath(*issue, lang, existingPaths)
		if path != "" {
			issuesByFile[path] = append(issuesByFile[path], issue)
		}
	}

	return &Index{
		issuesByFile: issuesByFile,
	}
}

// GetIssuesForFile returns all issues associated with a specific file path.
func (idx *Index) GetIssuesForFile(path string) []*Issue {
	return idx.issuesByFile[path]
}

// AllPaths returns all file paths that have associated issues.
func (idx *Index) AllPaths() []string {
	paths := make([]string, 0, len(idx.issuesByFile))
	for path := range idx.issuesByFile {
		paths = append(paths, path)
	}
	return paths
}

// TotalFiles returns the number of unique files that have associated issues.
func (idx *Index) TotalFiles() int {
	return len(idx.issuesByFile)
}
