package issue

import (
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

// Index provides efficient lookup of issues by file path.
type Index struct {
	issuesByFile map[string][]*Issue
}

// NewIndex creates a new Index from a list of issues.
func NewIndex(issues []Issue, existingPaths map[string]bool) *Index {
	issuesByFile := make(map[string][]*Issue)

	// Expand existing paths to include all language variants
	// This matches Python's behavior: all_paths includes en -> ja, en -> ko, etc.
	expandedPaths := expandPathsForAllLanguages(existingPaths)

	for i := range issues {
		issue := &issues[i]
		lang := GuessLanguage(*issue)
		if lang == "" {
			continue
		}

		path := GuessPath(*issue, lang, expandedPaths)
		if path != "" {
			issuesByFile[path] = append(issuesByFile[path], issue)
		}
	}

	return &Index{
		issuesByFile: issuesByFile,
	}
}

// expandPathsForAllLanguages expands existing paths to include all language variants.
// For example, if "content/en/docs/concepts/overview.md" exists,
// it adds "content/ja/docs/concepts/overview.md", "content/ko/docs/concepts/overview.md", etc.
func expandPathsForAllLanguages(existingPaths map[string]bool) map[string]bool {
	expanded := make(map[string]bool)

	// Copy original paths
	for path := range existingPaths {
		expanded[path] = true
	}

	// Add language variants
	languages := []language.Language{
		language.English,
		language.Korean,
		language.Japanese,
		language.Chinese,
		language.PortugueseBR,
		language.Spanish,
		language.Hindi,
		language.Indonesian,
		language.German,
		language.French,
		language.Italian,
		language.Vietnamese,
		language.Russian,
		language.Ukrainian,
		language.Polish,
		language.Bengali,
	}

	for path := range existingPaths {
		// If path starts with content/en/, generate variants for all languages
		if strings.HasPrefix(path, "content/en/") {
			for _, lang := range languages {
				variant := strings.Replace(path, "content/en/", "content/"+string(lang)+"/", 1)
				expanded[variant] = true
			}
		}
	}

	return expanded
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
