package translation

import (
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/issue"
	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pageview"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pr"
)

// TranslationStatus holds comprehensive information about a translation file.
type TranslationStatus struct {
	// Target file path
	Path string `json:"path"`

	// English file path representing the source of the translation
	EnglishPath string `json:"english_path"`

	// Language of the translation
	Language language.Language `json:"language"`

	// Content category (e.g., docs, blog, tutorials, etc.)
	Category string `json:"category"`

	// Git history analysis results
	History *HistoryAnalysis `json:"history,omitempty"`

	// Pull request information
	PullRequests []*pr.PullRequest `json:"pull_requests,omitempty"`

	// Issue information
	Issues []*issue.Issue `json:"issues,omitempty"`

	// URL information
	URL *URL `json:"url,omitempty"`

	// Page view statistics
	PageViewStats *pageview.PageViewStats `json:"page_view_stats,omitempty"`

	// Timestamp when this status was created
	CreatedAt time.Time `json:"created_at"`
}

// HistoryAnalysis contains git history analysis results.
type HistoryAnalysis struct {
	// Translation status
	Status Status `json:"status"`

	// Severity of being out-of-date
	Severity Severity `json:"severity"`

	// Last modified time of the translation file
	LastModified *time.Time `json:"last_modified,omitempty"`
	// Latest commit of the translation file
	LatestCommit *git.Commit `json:"latest_commit,omitempty"`

	// Full commit history of the translation file
	CommitHistory []*git.Commit `json:"commit_history,omitempty"`

	// Last modified time of the English file
	EnglishLastModified *time.Time `json:"english_last_modified,omitempty"` // Last modified time of the English file

	// Latest commit of the English file
	EnglishLatestCommit *git.Commit `json:"english_latest_commit,omitempty"` // Latest commit of the English file

	// Reference English commit at time of translation
	ReferenceCommit *git.Commit `json:"reference_commit,omitempty"`

	// Days since last translation update
	DaysBehind int `json:"days_behind"`

	// Number of English commits since last translation update
	CommitsBehind int `json:"commits_behind"`

	// List of English commits not yet reflected in translation
	MissingCommits []*git.Commit `json:"missing_commits,omitempty"`

	// Diff between reference English version and current English version
	Diff *Diff `json:"diff,omitempty"`
}

// URL contains URL information for the file.
type URL struct {
	// Public Kubernetes Official website URL (e.g., https://kubernetes.io/ja/docs/...)
	Website string `json:"website"`

	// GitHub repository URL (e.g., https://github.com/kubernetes/website/blob/main/content/ja/docs/...)
	GitHub string `json:"github"`
}

// Diff represents the difference between the current English version and the reference English version
type Diff struct {
	// Diff content
	Content string `json:"content"`

	// Total number of lines changed (insertions + deletions)
	LinesChanged int `json:"lines_changed"`

	// Number of lines added
	Insertions int `json:"insertions"`

	// Number of lines deleted
	Deletions int `json:"deletions"`

	// Commits involved in the diff
	OldCommitHash string `json:"old_commit_hash"`

	// Current English commit hash
	NewCommitHash string `json:"new_commit_hash"`
}
