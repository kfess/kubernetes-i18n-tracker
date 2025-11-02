package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/translation"
)

func TestBuildCategoryName(t *testing.T) {
	e := NewExporter(ExportOptions{})

	tests := []struct {
		name     string
		status   *translation.TranslationStatus
		expected string
	}{
		{
			name: "blog category",
			status: &translation.TranslationStatus{
				Category:    "blog",
				EnglishPath: "content/en/blog/_posts/2024-10-01-test.md",
			},
			expected: "blog",
		},
		{
			name: "docs with subcategory",
			status: &translation.TranslationStatus{
				Category:    "docs",
				EnglishPath: "content/en/docs/concepts/overview.md",
			},
			expected: "docs_concepts",
		},
		{
			name: "docs misc",
			status: &translation.TranslationStatus{
				Category:    "docs",
				EnglishPath: "content/en/docs/overview.md",
			},
			expected: "docs_misc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.buildCategoryName(tt.status)
			if result != tt.expected {
				t.Errorf("buildCategoryName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExportDiffs(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	diffDir := filepath.Join(tmpDir, "diff_go")
	if err := os.MkdirAll(diffDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	e := NewExporter(ExportOptions{OutputDir: tmpDir})

	// Create test data
	now := time.Now()
	byCategory := map[string]map[string]*translation.TranslationStatus{
		"blog": {
			"content/ja/blog/_posts/2024-10-01-test.md": {
				Path:        "content/ja/blog/_posts/2024-10-01-test.md",
				EnglishPath: "content/en/blog/_posts/2024-10-01-test.md",
				Language:    language.Japanese,
				Category:    "blog",
				History: &translation.HistoryAnalysis{
					Status:   translation.StatusOutdated,
					Severity: translation.SeverityMinor,
					ReferenceCommit: &git.Commit{
						Hash: "abc123",
						Date: now,
					},
					EnglishLatestCommit: &git.Commit{
						Hash: "def456",
						Date: now,
					},
					Diff: &translation.Diff{
						Content:       "diff content here",
						LinesChanged:  10,
						OldCommitHash: "abc123",
						NewCommitHash: "def456",
					},
				},
			},
		},
	}

	// Export diffs
	if err := e.exportDiffs(byCategory, diffDir); err != nil {
		t.Fatalf("exportDiffs() failed: %v", err)
	}

	// Verify file was created
	diffFile := filepath.Join(diffDir, "blog_diff.json")
	if _, err := os.Stat(diffFile); os.IsNotExist(err) {
		t.Errorf("Expected diff file was not created: %s", diffFile)
	}

	// Verify content
	data, err := os.ReadFile(diffFile)
	if err != nil {
		t.Fatalf("Failed to read diff file: %v", err)
	}

	var diffs map[string]DiffEntry
	if err := json.Unmarshal(data, &diffs); err != nil {
		t.Fatalf("Failed to parse diff file: %v", err)
	}

	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff entry, got %d", len(diffs))
	}
}

func TestExportMatrices(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	matrixDir := filepath.Join(tmpDir, "matrix_go")
	if err := os.MkdirAll(matrixDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	e := NewExporter(ExportOptions{OutputDir: tmpDir})

	// Create test data
	now := time.Now()
	byCategory := map[string]map[string]*translation.TranslationStatus{
		"blog": {
			"content/ja/blog/_posts/2024-10-01-test.md": {
				Path:        "content/ja/blog/_posts/2024-10-01-test.md",
				EnglishPath: "content/en/blog/_posts/2024-10-01-test.md",
				Language:    language.Japanese,
				Category:    "blog",
				History: &translation.HistoryAnalysis{
					Status:        translation.StatusUpToDate,
					Severity:      translation.SeverityCurrent,
					DaysBehind:    0,
					CommitsBehind: 0,
					LatestCommit: &git.Commit{
						Hash: "abc123",
						Date: now,
					},
					EnglishLatestCommit: &git.Commit{
						Hash: "abc123",
						Date: now,
					},
					ReferenceCommit: &git.Commit{
						Hash: "abc123",
						Date: now,
					},
				},
			},
		},
	}

	// Export matrices
	if err := e.exportMatrices(byCategory, matrixDir); err != nil {
		t.Fatalf("exportMatrices() failed: %v", err)
	}

	// Verify file was created
	matrixFile := filepath.Join(matrixDir, "blog.json")
	if _, err := os.Stat(matrixFile); os.IsNotExist(err) {
		t.Errorf("Expected matrix file was not created: %s", matrixFile)
	}

	// Verify content
	data, err := os.ReadFile(matrixFile)
	if err != nil {
		t.Fatalf("Failed to read matrix file: %v", err)
	}

	var matrix MatrixOutput
	if err := json.Unmarshal(data, &matrix); err != nil {
		t.Fatalf("Failed to parse matrix file: %v", err)
	}

	if len(matrix.Articles) != 1 {
		t.Errorf("Expected 1 article, got %d", len(matrix.Articles))
	}

	if matrix.LastUpdated == "" {
		t.Error("LastUpdated should not be empty")
	}
}
