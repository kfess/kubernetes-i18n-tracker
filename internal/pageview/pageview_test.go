package pageview

import (
	"strings"
	"testing"
)

func TestAggregateFromReader(t *testing.T) {
	csvData := `Page path,Views,New users,Engagement rate,Average session duration
/search/,571107,437,0.9521388203107544,244.54877420890625
/,194334,80468,0.6589300700886883,100.457034106046
/search/,100000,200,0.95,150.5`

	existingUrls := map[string]bool{
		"https://kubernetes.io/search/": true,
		"https://kubernetes.io/":        true,
	}

	reader := strings.NewReader(csvData)
	result, err := AggregateFromReader(reader, existingUrls, "https://kubernetes.io")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(result))
	}

	// Check aggregated /search/ data (two rows combined)
	searchStats := result["https://kubernetes.io/search/"]
	if searchStats == nil {
		t.Fatal("Expected stats for /search/, got nil")
	}

	expectedViews := 571107 + 100000
	if searchStats.Views != expectedViews {
		t.Errorf("Expected views %d, got %d", expectedViews, searchStats.Views)
	}

	expectedNewUsers := 437 + 200
	if searchStats.NewUsers != expectedNewUsers {
		t.Errorf("Expected new users %d, got %d", expectedNewUsers, searchStats.NewUsers)
	}

	// Check root path data
	rootStats := result["https://kubernetes.io/"]
	if rootStats == nil {
		t.Fatal("Expected stats for /, got nil")
	}

	if rootStats.Views != 194334 {
		t.Errorf("Expected views 194334, got %d", rootStats.Views)
	}
}

func TestAggregateFromReader_InvalidURL(t *testing.T) {
	csvData := `Page path,Views,New users,Engagement rate,Average session duration
/invalid/,100,50,0.95,150.5`

	existingUrls := map[string]bool{
		"https://kubernetes.io/": true,
	}

	reader := strings.NewReader(csvData)
	result, err := AggregateFromReader(reader, existingUrls, "https://kubernetes.io")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Invalid URLs should be skipped
	if len(result) != 0 {
		t.Errorf("Expected 0 URLs (invalid URL should be skipped), got %d", len(result))
	}
}

func TestAggregateFromReader_EmptyCSV(t *testing.T) {
	csvData := `Page path,Views,New users,Engagement rate,Average session duration`

	existingUrls := map[string]bool{
		"https://kubernetes.io/": true,
	}

	reader := strings.NewReader(csvData)
	_, err := AggregateFromReader(reader, existingUrls, "https://kubernetes.io")

	if err == nil {
		t.Error("Expected error for empty CSV, got nil")
	}
}

func TestAggregateFromReader_InvalidFormat(t *testing.T) {
	csvData := `Page path,Views
/search/,100`

	existingUrls := map[string]bool{
		"https://kubernetes.io/search/": true,
	}

	reader := strings.NewReader(csvData)
	_, err := AggregateFromReader(reader, existingUrls, "https://kubernetes.io")

	// Should return error for invalid CSV format (less than 5 columns)
	if err == nil {
		t.Error("Expected error for invalid CSV format, got nil")
	}
}

func TestBuildURLFromPath(t *testing.T) {
	existingUrls := map[string]bool{
		"https://kubernetes.io/":        true,
		"https://kubernetes.io/search/": true,
	}

	tests := []struct {
		path     string
		expected string
	}{
		{"/", "https://kubernetes.io/"},
		{"/search/", "https://kubernetes.io/search/"},
		{"/invalid/", ""},
	}

	for _, tt := range tests {
		result := buildURLFromPath(tt.path, existingUrls, "https://kubernetes.io")
		if result != tt.expected {
			t.Errorf("buildURLFromPath(%q) = %q, expected %q", tt.path, result, tt.expected)
		}
	}
}
