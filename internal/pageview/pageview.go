package pageview

import (
	"encoding/csv"
	"fmt"
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

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

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

		path := row[0]

		views, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("invalid views value at row %d: %w", i+2, err)
		}

		newUsers, err := strconv.Atoi(row[2])
		if err != nil {
			return nil, fmt.Errorf("invalid new users value at row %d: %w", i+2, err)
		}

		avgSessionDuration, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid average session duration at row %d: %w", i+2, err)
		}

		url := buildURLFromPath(path, existingUrls, baseURL)
		if url == "" {
			continue
		}

		if stats, exists := data[url]; exists {
			stats.Views += views
			stats.NewUsers += newUsers
			stats.AverageSessionDuration += avgSessionDuration
		} else {
			data[url] = &PageViewStats{
				Views:                  views,
				NewUsers:               newUsers,
				AverageSessionDuration: avgSessionDuration,
			}
		}
	}

	return data, nil
}

func buildURLFromPath(path string, existingUrls map[string]bool, baseURL string) string {
	fullURL := baseURL + path

	if _, exists := existingUrls[fullURL]; exists {
		return fullURL
	}

	return ""
}
