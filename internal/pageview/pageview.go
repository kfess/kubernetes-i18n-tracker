package pageview

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// PageViewStats represents aggregated page view statistics for a URL
type PageViewStats struct {
	Views                  int     `json:"views"`
	NewUsers               int     `json:"newUsers"`
	AverageSessionDuration float64 `json:"averageSessionDuration"`
}

// AggregatePageViews reads a CSV file and aggregates page view data by URL
// CSV format: Page path,Views,New users,Engagement rate,Average session duration
func AggregatePageViews(csvPath string, existingUrls map[string]bool, baseURL string) (map[string]*PageViewStats, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	return AggregateFromReader(file, existingUrls, baseURL)
}

// AggregateFromReader reads from an io.Reader and aggregates page view data
// This function is more testable as it doesn't require file system access
func AggregateFromReader(r io.Reader, existingUrls map[string]bool, baseURL string) (map[string]*PageViewStats, error) {
	reader := csv.NewReader(r)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	return aggregateRows(rows, existingUrls, baseURL)
}

// aggregateRows processes CSV rows and aggregates page view data
func aggregateRows(rows [][]string, existingUrls map[string]bool, baseURL string) (map[string]*PageViewStats, error) {
	if len(rows) < 2 {
		return nil, fmt.Errorf("CSV file is empty or has no data rows")
	}

	header := rows[0]
	if len(header) < 5 {
		return nil, fmt.Errorf("invalid CSV format: expected at least 5 columns, got %d", len(header))
	}

	data := make(map[string]*PageViewStats)

	for i, row := range rows[1:] {
		if len(row) < 5 {
			continue // Skip malformed rows
		}

		stats, url, err := parseRow(row, i+2, existingUrls, baseURL)
		if err != nil {
			return nil, err
		}

		if url == "" {
			continue
		}

		// Aggregate data for the same URL
		if existing, exists := data[url]; exists {
			existing.Views += stats.Views
			existing.NewUsers += stats.NewUsers
			existing.AverageSessionDuration += stats.AverageSessionDuration
		} else {
			data[url] = stats
		}
	}

	return data, nil
}

// parseRow parses a single CSV row and returns stats and URL
func parseRow(row []string, lineNum int, existingUrls map[string]bool, baseURL string) (*PageViewStats, string, error) {
	path := row[0]

	views, err := strconv.Atoi(row[1])
	if err != nil {
		return nil, "", fmt.Errorf("invalid views value at row %d: %w", lineNum, err)
	}

	newUsers, err := strconv.Atoi(row[2])
	if err != nil {
		return nil, "", fmt.Errorf("invalid new users value at row %d: %w", lineNum, err)
	}

	avgSessionDuration, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, "", fmt.Errorf("invalid average session duration at row %d: %w", lineNum, err)
	}

	url := buildURLFromPath(path, existingUrls, baseURL)

	stats := &PageViewStats{
		Views:                  views,
		NewUsers:               newUsers,
		AverageSessionDuration: avgSessionDuration,
	}

	return stats, url, nil
}

// buildURLFromPath constructs a full URL from a path
func buildURLFromPath(path string, existingUrls map[string]bool, baseURL string) string {
	fullURL := baseURL + path

	if _, exists := existingUrls[fullURL]; exists {
		return fullURL
	}

	return ""
}
