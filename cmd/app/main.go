package main

import (
	"context"
	"fmt"

	"encoding/json"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pageview"
	"github.com/kfess/kubernetes-i18n-tracker/internal/url"
)

func main() {
	logger.Init()
	err := godotenv.Load()
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
		return
	}

	// Pull Request
	// token := os.Getenv("KUBERNETES_WEBSITE_READ_GITHUB_TOKEN")
	// client := pr.NewClient(token, "kubernetes", "website")
	// fetcher := pr.NewFetcher(client)
	// prs, err := fetcher.FetchAll(context.Background())
	// if err != nil {
	// 	logger.Errorf("Error fetching PRs: %v", err)
	// 	return
	// }
	// for _, pr := range prs {
	// 	logger.Infof("PR #%d: %s", pr.Number, pr.Title)

	// 	// fmt.Printf("PR #%d: %s\n", pr.Number, pr.Title)
	// 	// fmt.Println("Files changed in this PR:")
	// 	for _, file := range pr.Files {
	// 		fmt.Printf("- %s\n", file)
	// 	}
	// }

	// URL
	// 全パスの読み込み
	f, err := os.ReadFile("./data/master/all_files.csv")
	if err != nil {
		logger.Errorf("Error reading all_files.csv: %v", err)
		return
	}
	lines := string(f)
	var allPaths []string
	for _, line := range strings.Split(lines, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			allPaths = append(allPaths, line)
		}
	}

	client := url.NewClient("https://kubernetes.io")
	urls, err := client.FetchAllSitemaps(context.Background())
	if err != nil {
		logger.Errorf("Error fetching sitemaps: %v", err)
		return
	}

	urlMap := make(map[string]bool, len(urls))
	for _, url := range urls {
		urlMap[url] = true
	}

	parser := url.NewYAMLFrontMatterParser("./k8s-repo/website")
	config := url.Config{
		BaseUrl:        "https://kubernetes.io",
		ExistingUrls:   urlMap,
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections: []string{
			"docs",
			"blog",
			"case-studies",
			"careers",
			"community",
			"examples",
			"partners",
			"releases",
			"training",
			"_common-resources",
			"includes",
		},
	}
	converter := url.NewConverter(config, parser)
	generatedurl, err := converter.Convert(context.Background(), "content/ko/blog/_posts/2018-11-07-grpc-load-balancing-with-linkerd.md")
	fmt.Println(generatedurl)

	// URL変換結果を保存する構造体
	type URLResult struct {
		Path  string `json:"path"`
		URL   string `json:"url,omitempty"`
		Error string `json:"error,omitempty"`
	}

	var results []URLResult
	successCount := 0
	errorCount := 0

	logger.Info("Starting URL conversion...")
	for _, path := range allPaths {
		generatedURL, err := converter.Convert(context.Background(), path)
		if err != nil {
			results = append(results, URLResult{
				Path:  path,
				Error: err.Error(),
			})
			errorCount++
			logger.Warnf("Error converting URL for path %s: %v", path, err)
			continue
		}
		results = append(results, URLResult{
			Path: path,
			URL:  generatedURL,
		})
		successCount++
	}

	logger.Infof("Conversion completed: %d successful, %d errors", successCount, errorCount)

	// 結果をJSONファイルに出力
	outputPath := "./data/output/url_conversion_results.json"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		logger.Errorf("Error creating output file: %v", err)
		return
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		logger.Errorf("Error encoding results to JSON: %v", err)
		return
	}

	logger.Infof("Results saved to %s", outputPath)

	// Page View データの集計
	logger.Info("Starting page view aggregation...")
	pageViewPath := "./data/master/page_view.csv"
	pageViewStats, err := pageview.AggregatePageViews(pageViewPath, urlMap, "https://kubernetes.io")
	if err != nil {
		logger.Errorf("Error aggregating page views: %v", err)
		return
	}

	logger.Infof("Page view data aggregated: %d unique URLs", len(pageViewStats))

	// ページビュー結果をJSONファイルに出力
	pageViewOutputPath := "./data/output/page_view_stats.json"
	pageViewOutputFile, err := os.Create(pageViewOutputPath)
	if err != nil {
		logger.Errorf("Error creating page view output file: %v", err)
		return
	}
	defer pageViewOutputFile.Close()

	pageViewEncoder := json.NewEncoder(pageViewOutputFile)
	pageViewEncoder.SetIndent("", "  ")
	if err := pageViewEncoder.Encode(pageViewStats); err != nil {
		logger.Errorf("Error encoding page view stats to JSON: %v", err)
		return
	}

	logger.Infof("Page view stats saved to %s", pageViewOutputPath)

	generatedUrl, err := converter.Convert(context.Background(), "content/en/blog/_posts/2020-06-30-SIG-Windows-Spotlight/index.md")
	if err != nil {
		logger.Errorf("Error converting URL: %v", err)
		return
	}
	logger.Infof("Generated URL: %s", generatedUrl)

	// Diff
	// oldCommitHash := "17f080049278a691d417c44decddbdb297f744b5"
	// newCommitHash := "43b5c46f2a625ee95fd65e770d7a7ac66ea1d440"
	// diff, err := diff.CalculateDiff(context.Background(), "./k8s-repo/website", oldCommitHash, newCommitHash, "content/ja/docs/concepts/architecture/_index.md")
	// if err != nil {
	// 	logger.Errorf("Error calculating diff: %v", err)
	// 	return
	// }
	// logger.Infof("Diff between %s..%s:\n%s\nEnglish File Path: %s", oldCommitHash, newCommitHash, diff.Language, diff.EnglishFilePath)
}
