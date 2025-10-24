package translation

import (
	"context"
	"strings"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/diff"
	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/history"
	"github.com/kfess/kubernetes-i18n-tracker/internal/issue"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
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

// GetTranslationStatus returns comprehensive translation status for a single file.
func (t *Tracker) GetTranslationStatus(ctx context.Context, translationPath string) (*TranslationStatus, error) {
	pathInfo := parsePath(translationPath)
	englishPath := pathInfo.ToEnglishPath()

	translationStatus := &TranslationStatus{
		Path:        translationPath,
		EnglishPath: englishPath,
		Language:    pathInfo.Language,
		Category:    pathInfo.Category,
		CreatedAt:   time.Now(),
	}

	// Build each component
	translationStatus.History = t.buildHistory(englishPath, translationPath)
	translationStatus.PullRequests = t.buildPullRequests(translationPath)
	translationStatus.Issues = t.buildIssues(translationPath)
	translationStatus.URL = t.buildURL(ctx, translationPath)

	// Build diff only if outdated
	if translationStatus.History != nil && translationStatus.History.Status == StatusOutdated {
		translationStatus.Diff = t.buildDiff(ctx, translationStatus.History, englishPath)
	}

	return translationStatus, nil
}

// buildHistory builds history by comparing English and translation file.
func (t *Tracker) buildHistory(englishPath string, translationPath string) *HistoryAnalysis {
	pathInfo := parsePath(translationPath)

	englishCommits := t.history.GetCommits(englishPath)
	translationCommits := t.history.GetCommits(translationPath)

	status := calculateStatus(
		string(pathInfo.Language),
		englishCommits,
		translationCommits,
	)

	// Handle no English version case
	if status == StatusNoEnglishVersion {
		return &HistoryAnalysis{
			Status:   StatusNoEnglishVersion,
			Severity: SeverityCurrent,
		}
	}

	var englishLatest, translationLatest *git.Commit
	if len(englishCommits) > 0 {
		englishLatest = englishCommits[len(englishCommits)-1]
	}
	if len(translationCommits) > 0 {
		translationLatest = translationCommits[len(translationCommits)-1]
	}

	// Handle not translated case
	if status == StatusNotTranslated {
		return t.buildNotTranslatedHistory(englishCommits, englishLatest)
	}

	// Handle up-to-date or outdated cases
	return t.buildTranslatedHistory(
		englishPath,
		englishCommits,
		englishLatest,
		translationCommits,
		translationLatest,
	)
}

// buildNotTranslatedHistory creates history analysis for files that haven't been translated.
func (t *Tracker) buildNotTranslatedHistory(
	englishCommits []*git.Commit,
	englishLatest *git.Commit,
) *HistoryAnalysis {
	stats := calculateChangeStats(englishCommits)
	daysBehind := calculateDaysBehind(englishCommits, []*git.Commit{})

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

// buildTranslatedHistory creates history analysis by comparing English and translation files.
func (t *Tracker) buildTranslatedHistory(
	englishPath string,
	englishCommits []*git.Commit,
	englishLatest *git.Commit,
	translationCommits []*git.Commit,
	translationLatest *git.Commit,
) *HistoryAnalysis {
	missingCommits := t.history.GetCommitsAfter(englishPath, translationLatest.Date)
	referenceCommit := t.history.GetCommitBeforeOrAt(englishPath, translationLatest.Date)
	
	stats := calculateChangeStats(missingCommits)
	daysBehind := calculateDaysBehind(englishCommits, translationCommits)

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
		return []*pr.PullRequest{} // Return empty slice instead of nil
	}

	prs := t.prIndex.GetPRsForFile(path)
	if prs == nil {
		return []*pr.PullRequest{} // Return empty slice instead of nil
	}
	if len(prs) > 0 {
		logger.Debugf("Found %d PRs for %s", len(prs), path)
	}
	return prs
}

func (t *Tracker) buildIssues(path string) []*issue.Issue {
	if t.issueIndex == nil {
		return []*issue.Issue{} // Return empty slice instead of nil
	}

	issues := t.issueIndex.GetIssuesForFile(path)
	if issues == nil {
		return []*issue.Issue{} // Return empty slice instead of nil
	}
	return issues
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
func calculateChangeStats(commits []*git.Commit) struct {
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
