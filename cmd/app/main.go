package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kfess/kubernetes-i18n-tracker/internal/exporter"
	"github.com/kfess/kubernetes-i18n-tracker/internal/history"
	"github.com/kfess/kubernetes-i18n-tracker/internal/issue"
	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pr"
	"github.com/kfess/kubernetes-i18n-tracker/internal/translation"
	"github.com/kfess/kubernetes-i18n-tracker/internal/url"
)

const (
	repoOwner = "kubernetes"
	repoName  = "website"
	repoPath  = "./k8s-repo/website"

	gitHistoryFile = "./data/master/git_history.jsonl"
	allFilesPath   = "./data/master/all_files.csv"
	outputDir      = "./data/output"
)

// loadEvents reads events from a JSONL file
func loadEvents(path string) ([]*history.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var events []*history.Event
	scanner := bufio.NewScanner(file)

	// Increase buffer size for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var event history.Event
		if err := json.Unmarshal(line, &event); err != nil {
			log.Printf("Warning: failed to parse line %d: %v", lineNum, err)
			continue
		}

		events = append(events, &event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return events, nil
}

func loadAllPaths(path string) ([]string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var allPaths []string
	lines := string(file)
	for _, line := range strings.Split(lines, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			allPaths = append(allPaths, line)
		}
	}
	return allPaths, nil
}

// saveResults saves translation status results to JSON files
func saveResults(results map[string]*translation.TranslationStatus) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save complete results
	outputPath := fmt.Sprintf("%s/translation_status.json", outputDir)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to encode results: %w", err)
	}

	logger.Infof("Results saved to %s", outputPath)
	return nil
}

func main() {
	logger.Init()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warnf("No .env file found: %v", err)
	}

	ctx := context.Background()

	// Step 1: Load git history
	logger.Info("Loading git history...")
	events, err := loadEvents(gitHistoryFile)
	if err != nil {
		log.Fatalf("Failed to load git history: %v", err)
	}
	logger.Infof("Loaded %d events", len(events))

	// Step 2: Build history tracker
	logger.Info("Building history tracker...")
	historyTracker := history.Build(events)
	logger.Infof("Built history for %d files", len(historyTracker.AllPaths()))

	// Step 3: Load existing file paths
	logger.Info("Loading existing file paths...")
	allPaths, err := loadAllPaths(allFilesPath)
	if err != nil {
		log.Fatalf("Failed to load file paths: %v", err)
	}
	existingPathsMap := make(map[string]bool, len(allPaths))
	for _, path := range allPaths {
		existingPathsMap[path] = true
	}
	logger.Infof("Loaded %d existing file paths", len(allPaths))

	// Step 4: Fetch PRs
	token := os.Getenv("KUBERNETES_WEBSITE_READ_GITHUB_TOKEN")
	if token == "" {
		logger.Warn("KUBERNETES_WEBSITE_READ_GITHUB_TOKEN not set, skipping PR fetch")
	}

	var prIndex *pr.Index
	if token != "" {
		logger.Info("Fetching pull requests...")
		prClient := pr.NewClient(token, repoOwner, repoName)
		prFetcher := pr.NewFetcher(prClient)
		prs, err := prFetcher.FetchAll(ctx)
		if err != nil {
			logger.Errorf("Failed to fetch PRs: %v", err)
		} else {
			prIndex = pr.BuildPRIndex(prs)
			logger.Infof("Fetched and indexed %d PRs covering %d files", len(prs), prIndex.TotalFiles())
		}
	}

	// Step 5: Fetch Issues
	var issueIndex *issue.Index
	if token != "" {
		logger.Info("Fetching issues...")
		issueClient := issue.NewClient(token, repoOwner, repoName)
		issueFetcher := issue.NewFetcher(issueClient)
		issues, err := issueFetcher.FetchAll(ctx)
		if err != nil {
			logger.Errorf("Failed to fetch issues: %v", err)
		} else {
			issueIndex = issue.NewIndex(issues, existingPathsMap)
			logger.Infof("Fetched and indexed %d issues covering %d files", len(issues), issueIndex.TotalFiles())
		}
	}

	// Step 6: Fetch sitemaps for URL conversion
	logger.Info("Fetching sitemaps...")
	urlClient := url.NewClient("https://kubernetes.io")
	sitemapURLs, err := urlClient.FetchAllSitemaps(ctx)
	if err != nil {
		logger.Errorf("Failed to fetch sitemaps: %v", err)
		sitemapURLs = []string{}
	}
	existingURLsMap := make(map[string]bool, len(sitemapURLs))
	for _, u := range sitemapURLs {
		existingURLsMap[u] = true
	}
	logger.Infof("Fetched %d URLs from sitemaps", len(sitemapURLs))

	// Step 7: Create URL converter
	parser := url.NewYAMLFrontMatterParser(repoPath)
	urlConfig := url.Config{
		BaseUrl:        "https://kubernetes.io",
		ExistingUrls:   existingURLsMap,
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections: []string{
			"docs", "blog", "case-studies", "careers", "community",
			"examples", "partners", "releases", "training",
			"_common-resources", "includes",
		},
	}
	urlConverter := url.NewConverter(urlConfig, parser)

	// Step 8: Create translation tracker
	logger.Info("Creating translation tracker...")
	trackerConfig := translation.Config{
		RepoPath:      repoPath,
		ExistingPaths: allPaths,
	}
	tracker := translation.NewTracker(historyTracker, urlConverter, prIndex, issueIndex, trackerConfig)

	// Step 9: Analyze translation status for all files
	logger.Info("Analyzing translation status...")
	results := make(map[string]*translation.TranslationStatus)

	// Collect English files
	englishFiles := []string{}
	for _, path := range allPaths {
		if strings.HasPrefix(path, "content/en/") {
			ext := filepath.Ext(path)
			if ext == ".md" || ext == ".html" {
				englishFiles = append(englishFiles, path)
			}
		}
	}

	// For each English file, check all supported languages
	supportedLangs := []string{"bn", "de", "es", "fr", "hi", "id", "it", "ja", "ko", "pl", "pt-br", "ru", "uk", "vi", "zh-cn"}
	for _, englishPath := range englishFiles {
		for _, lang := range supportedLangs {
			// Convert English path to translation path
			translationPath := strings.Replace(englishPath, "content/en/", "content/"+lang+"/", 1)

			status, err := tracker.GetStatus(ctx, translationPath)
			if err != nil {
				logger.Warnf("Failed to get status for %s: %v", translationPath, err)
				continue
			}

			results[translationPath] = status
		}
	}

	logger.Infof("Analyzed %d translation files", len(results))

	// Step 10: Export results to diff_go and matrix_go
	logger.Info("Exporting results...")
	exp := exporter.NewExporter(exporter.ExportOptions{
		OutputDir: outputDir,
	})
	if err := exp.Export(results); err != nil {
		log.Fatalf("Failed to export results: %w", err)
	}

	// Step 11: Save complete results to JSON
	logger.Info("Saving complete results...")
	if err := saveResults(results); err != nil {
		log.Fatalf("Failed to save results: %v", err)
	}

	logger.Info("Translation analysis complete!")
}
