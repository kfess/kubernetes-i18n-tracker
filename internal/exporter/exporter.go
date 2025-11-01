package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"github.com/kfess/kubernetes-i18n-tracker/internal/translation"
)

// DiffEntry represents a single diff entry for a file.
type DiffEntry struct {
	EnglishPath             string `json:"englishPath"`
	Language                string `json:"language"`
	RefEnglishCommitHash    string `json:"refEnglishCommitHash"`
	EnglishLatestCommitHash string `json:"englishLatestCommitHash"`
	Diff                    string `json:"diff"`
}

// MatrixPR represents a pull request in the matrix format.
type MatrixPR struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	Files  []string `json:"files"`
}

// MatrixIssue represents an issue in the matrix format.
type MatrixIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

// MatrixTranslation represents translation status for a specific language in the matrix.
type MatrixTranslation struct {
	Status                  string        `json:"status"`
	Severity                string        `json:"severity"`
	DaysBehind              int           `json:"daysBehind"`
	CommitsBehind           int           `json:"commitsBehind"`
	TotalChangeLines        int           `json:"totalChangeLines"`
	TargetLatestDate        *string       `json:"targetLatestDate"`
	EnglishLatestDate       *string       `json:"englishLatestDate"`
	TranslationUrl          *string       `json:"translationUrl"`
	Views                   int           `json:"views"`
	NewUsers                int           `json:"newUsers"`
	AverageSessionDuration  float64       `json:"averageSessionDuration"`
	Issues                  []MatrixIssue `json:"issues"`
	PRs                     []MatrixPR    `json:"prs"`
	EnglishLatestCommitHash string        `json:"englishLatestCommitHash"`
	RefEnglishCommitHash    *string       `json:"refEnglishCommitHash"`
	RefEnglishCommitDate    *string       `json:"refEnglishCommitDate"`
}

// MatrixArticle represents a single article in the matrix format.
type MatrixArticle struct {
	EnglishPath  string                       `json:"englishPath"`
	EnglishUrl   string                       `json:"englishUrl"`
	Translations map[string]MatrixTranslation `json:"translations"`
}

// MatrixOutput represents the complete matrix output format.
type MatrixOutput struct {
	LastUpdated string          `json:"lastUpdated"`
	Articles    []MatrixArticle `json:"articles"`
}

// ExportOptions contains options for exporting translation status.
type ExportOptions struct {
	OutputDir string
}

// Exporter exports translation status to various formats.
type Exporter struct {
	options ExportOptions
}

// NewExporter creates a new Exporter with the given options.
func NewExporter(options ExportOptions) *Exporter {
	return &Exporter{
		options: options,
	}
}

// Export exports translation status to diff and matrix formats.
func (e *Exporter) Export(results map[string]*translation.TranslationStatus) error {
	// Create output directories
	diffDir := filepath.Join(e.options.OutputDir, "diff")
	matrixDir := filepath.Join(e.options.OutputDir, "matrix")

	if err := os.MkdirAll(diffDir, 0755); err != nil {
		return fmt.Errorf("failed to create diff directory: %w", err)
	}
	if err := os.MkdirAll(matrixDir, 0755); err != nil {
		return fmt.Errorf("failed to create matrix directory: %w", err)
	}

	// Group by category
	byCategory := e.groupByCategory(results)

	// Export diff files
	logger.Info("Exporting diff files...")
	if err := e.exportDiffs(byCategory, diffDir); err != nil {
		return fmt.Errorf("failed to export diffs: %w", err)
	}

	// Export matrix files
	logger.Info("Exporting matrix files...")
	if err := e.exportMatrices(byCategory, matrixDir); err != nil {
		return fmt.Errorf("failed to export matrices: %w", err)
	}

	return nil
}

// groupByCategory groups translation statuses by category.
func (e *Exporter) groupByCategory(results map[string]*translation.TranslationStatus) map[string]map[string]*translation.TranslationStatus {
	byCategory := make(map[string]map[string]*translation.TranslationStatus)

	for path, status := range results {
		category := e.buildCategoryName(status)

		if byCategory[category] == nil {
			byCategory[category] = make(map[string]*translation.TranslationStatus)
		}

		byCategory[category][path] = status
	}

	return byCategory
}

// buildCategoryName builds the category name for a translation status.
func (e *Exporter) buildCategoryName(status *translation.TranslationStatus) string {
	if status.Category == "docs" {
		subcategory := e.extractDocsSubcategory(status.EnglishPath)
		if subcategory != "" {
			return "docs_" + subcategory
		}
	}
	return status.Category
}

// extractDocsSubcategory extracts the subcategory from a docs path.
func (e *Exporter) extractDocsSubcategory(path string) string {
	re := regexp.MustCompile(`content/[^/]+/docs/([^/]+)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		subcategory := matches[1]
		// Check if it's a filename (has .md extension)
		if strings.HasSuffix(subcategory, ".md") || strings.HasSuffix(subcategory, ".html") {
			return "misc"
		}
		return subcategory
	}
	return "misc"
}

// exportDiffs exports diff files for each category.
func (e *Exporter) exportDiffs(byCategory map[string]map[string]*translation.TranslationStatus, diffDir string) error {
	for category, statuses := range byCategory {
		diffs := make(map[string]DiffEntry)

		for path, status := range statuses {
			if status.History != nil &&
				(status.History.Status == translation.StatusOutdated || status.History.Status == translation.StatusNotTranslated) &&
				status.History.Diff != nil {
				refCommit := ""
				if status.History.ReferenceCommit != nil {
					refCommit = status.History.ReferenceCommit.Hash
				}

				latestCommit := ""
				if status.History.EnglishLatestCommit != nil {
					latestCommit = status.History.EnglishLatestCommit.Hash
				}

				diffs[path] = DiffEntry{
					EnglishPath:             status.EnglishPath,
					Language:                string(status.Language),
					RefEnglishCommitHash:    refCommit,
					EnglishLatestCommitHash: latestCommit,
					Diff:                    status.History.Diff.Content,
				}
			}
		}

		// Skip if no diffs
		if len(diffs) == 0 {
			continue
		}

		// Write to file
		filename := filepath.Join(diffDir, fmt.Sprintf("%s_diff.json", category))
		if err := e.writeJSON(filename, diffs); err != nil {
			return fmt.Errorf("failed to write diff file for %s: %w", category, err)
		}

		logger.Infof("Exported diff for %s: %d files", category, len(diffs))
	}

	return nil
}

// exportMatrices exports matrix files for each category.
func (e *Exporter) exportMatrices(byCategory map[string]map[string]*translation.TranslationStatus, matrixDir string) error {
	for category, statuses := range byCategory {
		// Group by English path
		byEnglishPath := make(map[string]map[language.Language]*translation.TranslationStatus)

		for _, status := range statuses {
			if byEnglishPath[status.EnglishPath] == nil {
				byEnglishPath[status.EnglishPath] = make(map[language.Language]*translation.TranslationStatus)
			}
			byEnglishPath[status.EnglishPath][status.Language] = status
		}

		// Build articles
		articles := []MatrixArticle{}
		for englishPath, translations := range byEnglishPath {
			// Use the existing generated website URL from the English status.
			englishURL := ""
			if engStatus, ok := translations[language.LanguageEnglish]; ok && engStatus.URL != nil {
				englishURL = engStatus.URL.Website
			}

			article := MatrixArticle{
				EnglishPath:  englishPath,
				EnglishUrl:   englishURL,
				Translations: make(map[string]MatrixTranslation),
			}

			for lang, status := range translations {
				article.Translations[string(lang)] = e.buildMatrixTranslation(status)
			}

			articles = append(articles, article)
		}

		// Sort articles
		e.sortArticles(articles, category)

		// Create matrix output
		matrix := MatrixOutput{
			LastUpdated: time.Now().UTC().Format(time.RFC3339),
			Articles:    articles,
		}

		// Write to file
		filename := filepath.Join(matrixDir, fmt.Sprintf("%s.json", category))
		if err := e.writeJSON(filename, matrix); err != nil {
			return fmt.Errorf("failed to write matrix file for %s: %w", category, err)
		}

		logger.Infof("Exported matrix for %s: %d articles", category, len(articles))
	}

	return nil
}

// buildMatrixTranslation builds a MatrixTranslation from TranslationStatus.
func (e *Exporter) buildMatrixTranslation(status *translation.TranslationStatus) MatrixTranslation {
	// Calculate total change lines from diff if available
	totalChangeLines := 0
	if status.History.Diff != nil {
		totalChangeLines = status.History.Diff.LinesChanged
	}

	mt := MatrixTranslation{
		Status:                 string(status.History.Status),
		Severity:               string(status.History.Severity),
		DaysBehind:             status.History.DaysBehind,
		CommitsBehind:          status.History.CommitsBehind,
		TotalChangeLines:       totalChangeLines,
		Views:                  0,
		NewUsers:               0,
		AverageSessionDuration: 0.0,
		Issues:                 []MatrixIssue{},
		PRs:                    []MatrixPR{},
	}

	// Add page view data from TranslationStatus if available
	if status.PageViewStats != nil {
		mt.Views = status.PageViewStats.Views
		mt.NewUsers = status.PageViewStats.NewUsers
		mt.AverageSessionDuration = status.PageViewStats.AverageSessionDuration
	}

	// Add dates
	if status.History.LatestCommit != nil {
		date := status.History.LatestCommit.Date.Format(time.RFC3339)
		mt.TargetLatestDate = &date
	}

	if status.History.EnglishLatestCommit != nil {
		date := status.History.EnglishLatestCommit.Date.Format(time.RFC3339)
		mt.EnglishLatestDate = &date
		mt.EnglishLatestCommitHash = status.History.EnglishLatestCommit.Hash
	}

	if status.History.ReferenceCommit != nil {
		hash := status.History.ReferenceCommit.Hash
		date := status.History.ReferenceCommit.Date.Format(time.RFC3339)
		mt.RefEnglishCommitHash = &hash
		mt.RefEnglishCommitDate = &date
	}

	// Add URL
	if status.URL != nil && status.URL.Website != "" {
		mt.TranslationUrl = &status.URL.Website
	}

	// Add PRs
	if status.PullRequests != nil {
		for _, pr := range status.PullRequests {
			mt.PRs = append(mt.PRs, MatrixPR{
				Number: pr.Number,
				Title:  pr.Title,
				URL:    pr.Url,
				Files:  pr.Files,
			})
		}
	}

	// Add Issues
	if status.Issues != nil {
		for _, issue := range status.Issues {
			mt.Issues = append(mt.Issues, MatrixIssue{
				Number: issue.Number,
				Title:  issue.Title,
				URL:    issue.URL,
			})
		}
	}

	return mt
}

// writeJSON writes data to a JSON file with indentation.
func (e *Exporter) writeJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// Ignore error, refactor exporter.go
	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// sortArticles sorts articles based on the category.
// For blog category, sort _index.md first, then by date (newest first) extracted from filename.
// For other categories, sort by English path alphabetically.
func (e *Exporter) sortArticles(articles []MatrixArticle, category string) {
	if category == "blog" {
		// Sort blog articles: _index.md first, then by date (newest first)
		sort.Slice(articles, func(i, j int) bool {
			isIndexI := strings.HasSuffix(articles[i].EnglishPath, "/_index.md")
			isIndexJ := strings.HasSuffix(articles[j].EnglishPath, "/_index.md")

			// _index.md always comes first
			if isIndexI && !isIndexJ {
				return true
			}
			if !isIndexI && isIndexJ {
				return false
			}

			// If both are _index.md or both are not, sort by date
			dateI := extractDateFromBlogPath(articles[i].EnglishPath)
			dateJ := extractDateFromBlogPath(articles[j].EnglishPath)
			// Descending order (newest first)
			return dateI > dateJ
		})
	} else {
		// Sort other articles by English path alphabetically
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].EnglishPath < articles[j].EnglishPath
		})
	}
}

// extractDateFromBlogPath extracts the date (yyyy-mm-dd) from a blog file path.
// Blog paths typically look like: content/en/blog/_posts/yyyy-mm-dd-title.md
// Returns the date string, or empty string if not found.
func extractDateFromBlogPath(path string) string {
	// Match yyyy-mm-dd pattern in the filename
	re := regexp.MustCompile(`/(\d{4}-\d{2}-\d{2})`)
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1]
	}
	// If no date found, return empty string (will be sorted to the end)
	return ""
}
