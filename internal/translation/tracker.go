package translation

import (
	"context"
	"strings"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/diff"
	"github.com/kfess/kubernetes-i18n-tracker/internal/history"
	"github.com/kfess/kubernetes-i18n-tracker/internal/issue"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pr"
	"github.com/kfess/kubernetes-i18n-tracker/internal/url"
)

// Tracker tracks and provides translation status by combining multiple data sources.
type Tracker struct {
	history       *history.History
	urlConverter  *url.Converter
	prIndex       *pr.Index
	issueIndex    *issue.Index
	repoPath      string
	existingPaths map[string]bool
}

// Config contains configuration for creating a Tracker.
type Config struct {
	RepoPath      string
	ExistingPaths []string
}

// NewTracker creates a new Tracker with the given dependencies.
func NewTracker(
	historyTracker *history.History,
	urlConverter *url.Converter,
	prIndex *pr.Index,
	issueIndex *issue.Index,
	config Config,
) *Tracker {
	pathMap := make(map[string]bool, len(config.ExistingPaths))
	for _, path := range config.ExistingPaths {
		pathMap[path] = true
	}

	return &Tracker{
		history:       historyTracker,
		urlConverter:  urlConverter,
		prIndex:       prIndex,
		issueIndex:    issueIndex,
		repoPath:      config.RepoPath,
		existingPaths: pathMap,
	}
}

// GetStatus returns comprehensive translation status for a single file.
func (t *Tracker) GetStatus(ctx context.Context, translationPath string) (*TranslationStatus, error) {
	pathInfo := parsePath(translationPath)
	englishPath := pathInfo.ToEnglishPath()

	status := &TranslationStatus{
		Path:        translationPath,
		EnglishPath: englishPath,
		Language:    pathInfo.Language,
		Category:    pathInfo.Category,
		CreatedAt:   time.Now(),
	}

	// Build each component
	status.History = t.buildHistory(englishPath, translationPath)
	status.PullRequests = t.buildPullRequests(translationPath)
	status.Issues = t.buildIssues(translationPath)
	status.URL = t.buildURL(ctx, translationPath)

	// Build diff only if outdated (heavy operation)
	if status.History != nil && status.History.Status == StatusOutdated {
		status.Diff = t.buildDiff(ctx, status.History, englishPath)
	}

	return status, nil
}

// buildHistory builds history analysis by comparing English and translation files.
func (t *Tracker) buildHistory(englishPath string, translationPath string) *HistoryAnalysis {
	// Get English file history
	englishCommits := t.history.GetCommits(englishPath)
	if len(englishCommits) == 0 {
		return &HistoryAnalysis{
			Status:   StatusNoEnglishVersion,
			Severity: SeverityCurrent,
		}
	}

	englishLatest := englishCommits[len(englishCommits)-1]

	// Get translation file history
	translationCommits := t.history.GetCommits(translationPath)
	if len(translationCommits) == 0 || !t.existingPaths[translationPath] {
		return t.buildNotTranslatedHistory(englishCommits, englishLatest)
	}

	translationLatest := translationCommits[len(translationCommits)-1]

	// Compare and build analysis
	return t.buildComparisonHistory(
		englishPath,
		englishLatest,
		translationCommits,
		translationLatest,
	)
}

// buildNotTranslatedHistory creates history analysis for files that haven't been translated.
func (t *Tracker) buildNotTranslatedHistory(
	englishCommits []*history.Commit,
	englishLatest *history.Commit,
) *HistoryAnalysis {
	stats := calculateChangeStats(englishCommits)
	daysBehind := int(time.Since(englishLatest.Date).Hours() / 24)

	return &HistoryAnalysis{
		Status:              StatusNotTranslated,
		Severity:            calculateSeverity(stats.Total),
		EnglishLastModified: &englishLatest.Date,
		EnglishLatestCommit: englishLatest,
		DaysBehind:          daysBehind,
		CommitsBehind:       len(englishCommits),
		LinesBehind:         stats.Total,
		InsertionsBehind:    stats.Insertions,
		DeletionsBehind:     stats.Deletions,
		MissingCommits:      englishCommits,
	}
}

// buildComparisonHistory creates history analysis by comparing English and translation files.
func (t *Tracker) buildComparisonHistory(
	englishPath string,
	englishLatest *history.Commit,
	translationCommits []*history.Commit,
	translationLatest *history.Commit,
) *HistoryAnalysis {
	missingCommits := t.history.GetCommitsAfter(englishPath, translationLatest.Date)
	stats := calculateChangeStats(missingCommits)

	referenceCommit := t.history.GetCommitBeforeOrAt(englishPath, translationLatest.Date)

	daysBehind := max(0, int(englishLatest.Date.Sub(translationLatest.Date).Hours()/24))

	status := StatusUpToDate
	if len(missingCommits) > 0 {
		status = StatusOutdated
	}

	return &HistoryAnalysis{
		Status:              status,
		Severity:            calculateSeverity(stats.Total),
		LastModified:        &translationLatest.Date,
		LatestCommit:        translationLatest,
		CommitHistory:       translationCommits,
		EnglishLastModified: &englishLatest.Date,
		EnglishLatestCommit: englishLatest,
		ReferenceCommit:     referenceCommit,
		DaysBehind:          daysBehind,
		CommitsBehind:       len(missingCommits),
		LinesBehind:         stats.Total,
		InsertionsBehind:    stats.Insertions,
		DeletionsBehind:     stats.Deletions,
		MissingCommits:      missingCommits,
	}
}

// buildPullRequests builds PR list for the file.
func (t *Tracker) buildPullRequests(path string) []*pr.PullRequest {
	if t.prIndex == nil {
		return nil
	}

	return t.prIndex.GetPRsForFile(path)
}

func (t *Tracker) buildIssues(path string) []*issue.Issue {
	if t.issueIndex == nil {
		return nil
	}

	return t.issueIndex.GetIssuesForFile(path)
}

// buildURL builds URL information for the file.
func (t *Tracker) buildURL(ctx context.Context, path string) *URL {
	if t.urlConverter == nil {
		return nil
	}

	githubURL := toGitHubURL(path)
	websiteURL, err := t.urlConverter.Convert(ctx, path)
	if err != nil {
		// If URL generation fails, return GitHub URL only
		return &URL{
			GitHub: githubURL,
		}
	}

	return &URL{
		Website: websiteURL,
		GitHub:  githubURL,
	}
}

// buildDiff builds diff information (heavy operation, called only when needed).
func (t *Tracker) buildDiff(
	ctx context.Context,
	historyAnalysis *HistoryAnalysis,
	englishPath string,
) *Diff {
	if historyAnalysis.ReferenceCommit == nil || historyAnalysis.EnglishLatestCommit == nil {
		return nil
	}

	result, err := diff.CalculateDiff(
		ctx,
		t.repoPath,
		historyAnalysis.ReferenceCommit.Hash,
		historyAnalysis.EnglishLatestCommit.Hash,
		englishPath,
	)
	if err != nil {
		return nil
	}

	return &Diff{
		Content:      result.Content,
		LinesChanged: countDiffLines(result.Content),
		OldCommit:    result.OldCommitHash,
		NewCommit:    result.NewCommitHash,
	}
}

// Helper functions

// calculateChangeStats calculates insertion/deletion statistics from commits.
func calculateChangeStats(commits []*history.Commit) struct {
	Insertions int
	Deletions  int
	Total      int
} {
	var stats struct {
		Insertions int
		Deletions  int
		Total      int
	}

	for _, commit := range commits {
		stats.Insertions += commit.Insertions
		stats.Deletions += commit.Deletions
	}
	stats.Total = stats.Insertions + stats.Deletions

	return stats
}

// countDiffLines counts the number of changed lines in a diff.
func countDiffLines(diffContent string) int {
	lines := 0
	for _, line := range strings.Split(diffContent, "\n") {
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			// Exclude file headers (+++, ---)
			if !strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "---") {
				lines++
			}
		}
	}
	return lines
}
