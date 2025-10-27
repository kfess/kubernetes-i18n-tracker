package translation

import (
	"context"
	"fmt"
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
	fmParser      url.FrontMatterParser
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
	fmParser url.FrontMatterParser,
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
		fmParser:      fmParser,
		prIndex:       prIndex,
		issueIndex:    issueIndex,
		repoPath:      config.RepoPath,
		existingPaths: pathMap,
	}
}

// GetTranslationStatus returns comprehensive translation status for a single file.
func (t *Tracker) GetTranslationStatus(
	ctx context.Context,
	translationPath string,
	translationContent string,
) (*TranslationStatus, error) {
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
	history, err := t.buildHistory(ctx, englishPath, translationPath)
	if err != nil {
		return nil, fmt.Errorf("build history for %s: %w", translationPath, err)
	}
	translationStatus.History = history
	translationStatus.PullRequests = t.buildPullRequests(translationPath)
	translationStatus.Issues = t.buildIssues(translationPath)
	translationStatus.URL = t.buildURL(ctx, translationPath, translationContent)

	return translationStatus, nil
}

// buildHistory builds history by comparing English and translation file.
func (t *Tracker) buildHistory(ctx context.Context, englishPath string, translationPath string) (*HistoryAnalysis, error) {
	pathInfo := parsePath(translationPath)

	englishCommits := t.history.GetCommits(englishPath)
	translationCommits := t.history.GetCommits(translationPath)

	status := calculateStatus(
		string(pathInfo.Language),
		englishCommits,
		translationCommits,
	)

	logger.Debugf("Translation status for %s: %s (EN commits: %d, Translation commits: %d)",
		translationPath, status, len(englishCommits), len(translationCommits))

	// Handle no English version case
	if status == StatusNoEnglishVersion {
		logger.Warnf("No English version found for %s", translationPath)
		return t.buildNoEnglishVersionHistory(), nil
	}

	var englishLatest, translationLatest *git.Commit
	if len(englishCommits) > 0 {
		englishLatest = englishCommits[len(englishCommits)-1]
	}
	if len(translationCommits) > 0 {
		translationLatest = translationCommits[len(translationCommits)-1]
	}

	if status == StatusNotTranslated {
		return t.buildNotTranslatedHistory(englishCommits, englishLatest), nil
	}

	if status == StatusUpToDate {
		return t.buildUpToDateHistory(
			englishCommits,
			englishLatest,
			translationCommits,
			translationLatest,
		), nil
	}

	outdatedHistory, err := t.buildOutdatedHistory(
		ctx,
		englishPath,
		englishCommits,
		englishLatest,
		translationCommits,
		translationLatest,
	)
	if err != nil {
		return nil, err
	}

	return outdatedHistory, nil
}

func (t *Tracker) buildNoEnglishVersionHistory() *HistoryAnalysis {
	return &HistoryAnalysis{
		Status:   StatusNoEnglishVersion,
		Severity: SeverityCurrent,
	}
}

// buildNotTranslatedHistory creates history analysis for files that haven't been translated.
func (t *Tracker) buildNotTranslatedHistory(
	englishCommits []*git.Commit,
	englishLatest *git.Commit,
) *HistoryAnalysis {
	daysBehind := calculateDaysBehind(englishCommits, []*git.Commit{})

	return &HistoryAnalysis{
		Status:              StatusNotTranslated,
		Severity:            SeverityCritical, // Not translated is always critical
		EnglishLastModified: &englishLatest.Date,
		EnglishLatestCommit: englishLatest,
		DaysBehind:          daysBehind,
		CommitsBehind:       len(englishCommits),
		MissingCommits:      englishCommits,
	}
}

// buildUpToDateHistory creates history analysis for up-to-date translations.
func (t *Tracker) buildUpToDateHistory(
	englishCommits []*git.Commit,
	englishLatest *git.Commit,
	translationCommits []*git.Commit,
	translationLatest *git.Commit,
) *HistoryAnalysis {
	daysBehind := calculateDaysBehind(englishCommits, translationCommits)

	return &HistoryAnalysis{
		Status:              StatusUpToDate,
		Severity:            SeverityCurrent,
		LastModified:        &translationLatest.Date,
		LatestCommit:        translationLatest,
		CommitHistory:       translationCommits,
		EnglishLastModified: &englishLatest.Date,
		EnglishLatestCommit: englishLatest,
		ReferenceCommit:     translationLatest,
		DaysBehind:          daysBehind,
		CommitsBehind:       0,
		MissingCommits:      []*git.Commit{},
	}
}

// buildOutdatedHistory creates history analysis for outdated translations.
func (t *Tracker) buildOutdatedHistory(
	ctx context.Context,
	englishPath string,
	englishCommits []*git.Commit,
	englishLatest *git.Commit,
	translationCommits []*git.Commit,
	translationLatest *git.Commit,
) (*HistoryAnalysis, error) {
	daysBehind := calculateDaysBehind(englishCommits, translationCommits)

	missingCommits := t.history.GetCommitsAfter(englishPath, translationLatest.Date)
	referenceCommit := t.history.GetCommitBeforeOrAt(englishPath, translationLatest.Date)

	logger.Debugf("Outdated translation %s: %d days behind, %d missing commits",
		englishPath, daysBehind, len(missingCommits))

	analysis := &HistoryAnalysis{
		ReferenceCommit:     referenceCommit,
		EnglishLatestCommit: englishLatest,
	}

	diff, err := t.buildDiff(ctx, analysis, englishPath)
	if err != nil {
		return nil, fmt.Errorf("calculate diff for %s: %w", englishPath, err)
	}

	severity := SeverityCurrent
	if diff != nil {
		severity = calculateSeverity(diff.LinesChanged)
		logger.Debugf("Diff calculated for %s: %d lines changed, severity: %s",
			englishPath, diff.LinesChanged, severity)
	}

	return &HistoryAnalysis{
		Status:              StatusOutdated,
		Severity:            severity,
		LastModified:        &translationLatest.Date,
		LatestCommit:        translationLatest,
		CommitHistory:       translationCommits,
		EnglishLastModified: &englishLatest.Date,
		EnglishLatestCommit: englishLatest,
		ReferenceCommit:     referenceCommit,
		DaysBehind:          daysBehind,
		CommitsBehind:       len(missingCommits),
		MissingCommits:      missingCommits,
		Diff:                diff,
	}, nil
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

	if len(issues) > 0 {
		logger.Debugf("Found %d issues for %s", len(issues), path)
	}

	return issues
}

// buildURL builds URL information for the file.
func (t *Tracker) buildURL(ctx context.Context, path string, content string) *URL {
	githubURL := toGitHubURL(path)

	if content == "" || t.urlConverter == nil || t.fmParser == nil {
		return &URL{
			GitHub: githubURL,
		}
	}

	fm, err := t.fmParser.Parse(content)
	if err != nil {
		logger.Debugf("Front matter parsing failed for %s, using GitHub URL only: %v", path, err)
		return &URL{
			GitHub: githubURL,
		}
	}

	websiteURL, err := t.urlConverter.Convert(ctx, path, fm)
	if err != nil {
		logger.Debugf("Website URL conversion failed for %s, using GitHub URL only: %v", path, err)
		return &URL{
			GitHub: githubURL,
		}
	}

	return &URL{
		Website: websiteURL,
		GitHub:  githubURL,
	}
}

// buildDiff builds diff between reference English commit and latest English commit.
func (t *Tracker) buildDiff(
	ctx context.Context,
	historyAnalysis *HistoryAnalysis,
	englishPath string,
) (*Diff, error) {
	if historyAnalysis.ReferenceCommit == nil || historyAnalysis.EnglishLatestCommit == nil {
		logger.Debugf("Skipping diff for %s: missing reference or latest commit", englishPath)
		return nil, nil
	}

	logger.Debugf("Calculating diff for %s (%s..%s)",
		englishPath,
		historyAnalysis.ReferenceCommit.Hash[:7],
		historyAnalysis.EnglishLatestCommit.Hash[:7])

	result, err := diff.CalculateDiff(
		ctx,
		t.repoPath,
		historyAnalysis.ReferenceCommit.Hash,
		historyAnalysis.EnglishLatestCommit.Hash,
		englishPath,
	)
	if err != nil {
		return nil, err
	}

	return &Diff{
		Content:       result.Content,
		LinesChanged:  result.LinesChanged,
		Insertions:    result.Insertions,
		Deletions:     result.Deletions,
		OldCommitHash: result.OldCommitHash,
		NewCommitHash: result.NewCommitHash,
	}, nil
}
