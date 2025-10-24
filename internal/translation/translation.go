package translation

import (
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/issue"
	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pr"
)

// TranslationStatus holds comprehensive information about a translation file.
type TranslationStatus struct {
	// Basic
	Path        string            `json:"path"`
	EnglishPath string            `json:"english_path"`
	Language    language.Language `json:"language"`
	Category    string            `json:"category"` // docs, blog, tutorials, etc.

	// Analysis results (nil if not available/applicable)
	History      *HistoryAnalysis  `json:"history,omitempty"`
	PullRequests []*pr.PullRequest `json:"pull_requests,omitempty"`
	Issues       []*issue.Issue    `json:"issues,omitempty"`
	URL          *URL              `json:"url,omitempty"`
	Diff         *Diff             `json:"diff,omitempty"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
}

// HistoryAnalysis contains git history analysis results.
type HistoryAnalysis struct {
	Status   Status   `json:"status"`   // Translation status
	Severity Severity `json:"severity"` // Severity of being out-of-date

	// Translation file metrics
	LastModified  *time.Time    `json:"last_modified,omitempty"`  // Last modified time of the translation file
	LatestCommit  *git.Commit   `json:"latest_commit,omitempty"`  // Latest commit of the translation file
	CommitHistory []*git.Commit `json:"commit_history,omitempty"` // Full commit history of the translation file

	// English file metrics
	EnglishLastModified *time.Time  `json:"english_last_modified,omitempty"` // Last modified time of the English file
	EnglishLatestCommit *git.Commit `json:"english_latest_commit,omitempty"` // Latest commit of the English file
	ReferenceCommit     *git.Commit `json:"reference_commit,omitempty"`      // English commit at time of translation

	// Comparison metrics
	DaysBehind       int           `json:"days_behind"`               // Days since last translation update
	CommitsBehind    int           `json:"commits_behind"`            // Number of English commits since last translation update
	LinesBehind      int           `json:"lines_behind"`              // Total changed lines
	InsertionsBehind int           `json:"insertions_behind"`         // Insertions in English since last translation update
	DeletionsBehind  int           `json:"deletions_behind"`          // Deletions in English since last translation update
	MissingCommits   []*git.Commit `json:"missing_commits,omitempty"` // List of English commits not yet reflected in translation
}

// URL contains URL information for the file.
type URL struct {
	Website string `json:"website"` // Public website URL (e.g., https://kubernetes.io/ja/docs/...)
	GitHub  string `json:"github"`  // GitHub repository URL (e.g., https://github.com/kubernetes/website/blob/main/content/ja/docs/...)
}

// Diff contains the diff between the current English version and the reference English version
// that was used when the translation was last updated.
type Diff struct {
	Content      string `json:"content"`       // Diff content
	LinesChanged int    `json:"lines_changed"` // Number of lines changed
	OldCommit    string `json:"old_commit"`    // Reference English commit hash (at translation time)
	NewCommit    string `json:"new_commit"`    // Current English commit hash
}
