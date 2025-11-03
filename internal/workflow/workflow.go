package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kfess/kubernetes-i18n-tracker/internal/exporter"
	"github.com/kfess/kubernetes-i18n-tracker/internal/history"
	"github.com/kfess/kubernetes-i18n-tracker/internal/issue"
	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pageview"
	"github.com/kfess/kubernetes-i18n-tracker/internal/path"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pr"
	"github.com/kfess/kubernetes-i18n-tracker/internal/translation"
	"github.com/kfess/kubernetes-i18n-tracker/internal/url"
)

// Workflow orchestrates the entire translation tracking process.
type Workflow struct {
	config Config
}

// New creates a new Workflow instance.
func New(config Config) *Workflow {
	return &Workflow{
		config: config,
	}
}

// Run executes the complete workflow.
func (w *Workflow) Run(ctx context.Context) error {
	// Load git history
	logger.Info("Loading git history...")
	events, err := loadEvents(w.config.GitHistoryFile)
	if err != nil {
		return fmt.Errorf("failed to load git history: %w", err)
	}
	logger.Infof("Loaded %d events", len(events))

	// Build history tracker
	logger.Info("Building history tracker...")
	historyTracker := history.Build(events)
	logger.Infof("Built history for %d files", len(historyTracker.AllPaths()))

	// Load existing file paths
	logger.Info("Loading existing file paths...")
	existingPathsMap, err := loadExistingPathsMap(w.config.AllFilesPath)
	if err != nil {
		return fmt.Errorf("failed to load file paths: %w", err)
	}
	logger.Infof("Loaded %d existing file paths", len(existingPathsMap))

	// Build PR index
	var prIndex *pr.Index
	if w.config.GitHubToken != "" {
		logger.Info("Fetching pull requests...")
		prIndex, err = w.buildPRIndex(ctx)
		if err != nil {
			logger.Errorf("Failed to fetch PRs: %v", err)
		} else {
			logger.Infof("Indexed %d PRs covering %d files", prIndex.TotalPRs(), prIndex.TotalFiles())
		}
	}

	// Build issue index
	var issueIndex *issue.Index
	if w.config.GitHubToken != "" {
		logger.Info("Fetching issues...")
		issueIndex, err = w.buildIssueIndex(ctx, existingPathsMap)
		if err != nil {
			logger.Errorf("Failed to fetch issues: %v", err)
		} else {
			logger.Infof("Indexed issues covering %d files", issueIndex.TotalFiles())
		}
	}

	// Build URL converter
	logger.Info("Building URL converter...")
	urlConverter, existingURLsMap, err := w.buildURLConverter(ctx)
	if err != nil {
		return fmt.Errorf("failed to build URL converter: %w", err)
	}

	// Build pageview index
	logger.Info("Loading page view data...")
	pageviewIndex, err := w.buildPageviewIndex(existingURLsMap)
	if err != nil {
		logger.Errorf("Failed to load page view data: %v", err)
	} else {
		logger.Infof("Loaded page view data for %d URLs", pageviewIndex.TotalURLs())
	}

	// Analyze translations
	logger.Info("Analyzing translation status...")
	results, err := w.analyzeTranslations(ctx, historyTracker, urlConverter, prIndex, issueIndex, pageviewIndex, existingPathsMap)
	if err != nil {
		return fmt.Errorf("failed to analyze translations: %w", err)
	}
	logger.Infof("Analyzed %d translation files", len(results))

	// Export results
	logger.Info("Exporting results...")
	if err := w.exportResults(results); err != nil {
		return fmt.Errorf("failed to export results: %w", err)
	}

	// Save complete results
	logger.Info("Saving complete results...")
	if err := w.saveResults(results); err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}

	logger.Info("Translation analysis complete!")
	return nil
}

// buildPRIndex fetches and indexes pull requests.
func (w *Workflow) buildPRIndex(ctx context.Context) (*pr.Index, error) {
	prClient := pr.NewClient(w.config.GitHubToken, w.config.RepoOwner, w.config.RepoName)
	prFetcher := pr.NewFetcher(prClient)
	prs, err := prFetcher.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	return pr.BuildPRIndex(prs), nil
}

// buildIssueIndex fetches and indexes issues.
func (w *Workflow) buildIssueIndex(ctx context.Context, existingPaths map[string]bool) (*issue.Index, error) {
	issueClient := issue.NewClient(w.config.GitHubToken, w.config.RepoOwner, w.config.RepoName)
	issueFetcher := issue.NewFetcher(issueClient)
	issues, err := issueFetcher.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	return issue.NewIndex(issues, existingPaths), nil
}

// buildURLConverter creates a URL converter with sitemap data.
func (w *Workflow) buildURLConverter(ctx context.Context) (*url.Converter, map[string]bool, error) {
	urlClient := url.NewClient("https://kubernetes.io")
	sitemapURLs, err := urlClient.FetchAllSitemaps(ctx, language.SupportedLanguages)
	if err != nil {
		logger.Errorf("Failed to fetch sitemaps: %v", err)
		sitemapURLs = []string{}
	}

	existingURLsMap := make(map[string]bool, len(sitemapURLs))
	for _, u := range sitemapURLs {
		existingURLsMap[u] = true
	}
	logger.Infof("Fetched %d URLs from sitemaps", len(sitemapURLs))

	urlConfig := url.Config{
		BaseUrl:        "https://kubernetes.io",
		ExistingUrls:   existingURLsMap,
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  path.SupportedCategories,
	}
	return url.NewConverter(urlConfig), existingURLsMap, nil
}

// buildPageviewIndex loads and indexes pageview data.
func (w *Workflow) buildPageviewIndex(existingURLs map[string]bool) (*pageview.Index, error) {
	stats, err := pageview.AggregatePageViews(w.config.PageViewFile, existingURLs, "https://kubernetes.io")
	if err != nil {
		return pageview.NewIndex(make(map[string]*pageview.PageViewStats)), err
	}
	return pageview.NewIndex(stats), nil
}

// analyzeTranslations performs the translation analysis.
func (w *Workflow) analyzeTranslations(
	ctx context.Context,
	historyTracker *history.History,
	urlConverter *url.Converter,
	prIndex *pr.Index,
	issueIndex *issue.Index,
	pageviewIndex *pageview.Index,
	existingPaths map[string]bool,
) (map[string]*translation.TranslationStatus, error) {
	parser := url.NewYAMLFrontMatterParser()
	trackerConfig := translation.Config{
		RepoPath:      w.config.RepoPath,
		ExistingPaths: mapToSlice(existingPaths),
	}
	tracker := translation.NewTracker(historyTracker, urlConverter, parser, prIndex, issueIndex, pageviewIndex, trackerConfig)

	results := make(map[string]*translation.TranslationStatus)

	for englishPath := range existingPaths {
		pathInfo, err := path.Parse(englishPath)
		if err != nil {
			logger.Errorf("Failed to parse path %s: %v", englishPath, err)
			continue
		}
		if pathInfo.Language() != language.English || !pathInfo.IsContentFile() {
			continue
		}

		for _, lang := range language.SupportedLanguages {
			translationPath := pathInfo.ToLanguagePath(language.Language(lang))

			contentBytes, err := os.ReadFile(filepath.Join(w.config.RepoPath, translationPath))
			if err != nil {
				contentBytes = []byte{}
			}

			status, err := tracker.GetTranslationStatus(ctx, translationPath, string(contentBytes))
			if err != nil {
				logger.Errorf("Failed to get status for %s: %v", translationPath, err)
				continue
			}

			results[translationPath] = status
		}
	}

	return results, nil
}

// exportResults exports the results using the exporter.
func (w *Workflow) exportResults(results map[string]*translation.TranslationStatus) error {
	exp := exporter.NewExporter(exporter.ExportOptions{
		OutputDir: w.config.OutputDir,
	})
	return exp.Export(results)
}

// saveResults saves the complete results to JSON.
func (w *Workflow) saveResults(results map[string]*translation.TranslationStatus) error {
	if err := os.MkdirAll(w.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(w.config.OutputDir, "translation_status.json")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to encode results: %w", err)
	}

	logger.Infof("Results saved to %s", outputPath)
	return nil
}

// mapToSlice converts a map[string]bool to []string.
func mapToSlice(m map[string]bool) []string {
	slice := make([]string, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	return slice
}
