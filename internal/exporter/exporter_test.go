package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/structure"
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
					Status:    translation.StatusOutdated,
					GitStatus: translation.StatusOutdated,
					Severity:  translation.SeverityMinor,
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

func TestBuildMatrixTranslationStructure(t *testing.T) {
	e := NewExporter(ExportOptions{})

	tests := []struct {
		name     string
		history  *translation.HistoryAnalysis
		wantJSON string
	}{
		{
			name: "structural comparison is exported with the git status",
			history: &translation.HistoryAnalysis{
				Status:    translation.StatusPossiblyOutdated,
				GitStatus: translation.StatusUpToDate,
				Severity:  translation.SeverityCurrent,
				Structure: &structure.Result{
					Signal: structure.SignalModerate,
					Gap: structure.Gap{
						L10nToEnLineRatio:     0.6,
						L10nToEnBodyWordRatio: 0.7,
						MissingH2:             2,
						MissingCodeBlocks:     1,
					},
				},
			},
			wantJSON: `{"signal":"moderate","gap":{"lineRatio":0.6,"bodyWordRatio":0.7,"missingH2":2,"missingH3":0,` +
				`"missingCodeBlocks":1,"missingAnchors":0,"missingNewVersions":0,"missingFeatureState":0,"missingApiOrKind":0}}`,
		},
		{
			name: "no comparison is exported as null",
			history: &translation.HistoryAnalysis{
				Status:    translation.StatusNotTranslated,
				GitStatus: translation.StatusNotTranslated,
				Severity:  translation.SeverityCritical,
			},
			wantJSON: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := e.buildMatrixTranslation(&translation.TranslationStatus{History: tt.history})

			if mt.Status != string(tt.history.Status) || mt.GitStatus != string(tt.history.GitStatus) {
				t.Errorf("status = %q, gitStatus = %q, want %q and %q",
					mt.Status, mt.GitStatus, tt.history.Status, tt.history.GitStatus)
			}

			data, err := json.Marshal(mt)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got := string(fields["structure"]); got != tt.wantJSON {
				t.Errorf("structure JSON = %s, want %s", got, tt.wantJSON)
			}
			if _, ok := fields["gitStatus"]; !ok {
				t.Error("gitStatus is missing from the JSON")
			}
		})
	}
}
