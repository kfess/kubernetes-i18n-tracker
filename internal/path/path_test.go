package path

import (
	"testing"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected *Path
		wantErr  bool
	}{
		{
			name: "Japanese docs path",
			path: "content/ja/docs/concepts/overview.md",
			expected: &Path{
				original:  "content/ja/docs/concepts/overview.md",
				language:  language.Japanese,
				category:  Docs,
				segments:  []string{"content", "ja", "docs", "concepts", "overview.md"},
				filename:  "overview",
				extension: ".md",
			},
			wantErr: false,
		},
		{
			name: "English blog path",
			path: "content/en/blog/2023/kubernetes-release.md",
			expected: &Path{
				original:  "content/en/blog/2023/kubernetes-release.md",
				language:  language.English,
				category:  Blog,
				segments:  []string{"content", "en", "blog", "2023", "kubernetes-release.md"},
				filename:  "kubernetes-release",
				extension: ".md",
			},
			wantErr: false,
		},
		{
			name: "HTML file",
			path: "content/en/docs/index.html",
			expected: &Path{
				original:  "content/en/docs/index.html",
				language:  language.English,
				category:  Docs,
				segments:  []string{"content", "en", "docs", "index.html"},
				filename:  "index",
				extension: ".html",
			},
			wantErr: false,
		},
		{
			name: "Index file path",
			path: "content/ja/docs/_index.md",
			expected: &Path{
				original:  "content/ja/docs/_index.md",
				language:  language.Japanese,
				category:  Docs,
				segments:  []string{"content", "ja", "docs", "_index.md"},
				filename:  "_index",
				extension: ".md",
			},
			wantErr: false,
		},
		{
			name: "Root level file (overall category)",
			path: "content/ja/_index.html",
			expected: &Path{
				original:  "content/ja/_index.html",
				language:  language.Japanese,
				category:  "overall",
				segments:  []string{"content", "ja", "_index.html"},
				filename:  "_index",
				extension: ".html",
			},
			wantErr: false,
		},
		{
			name: "Chinese path",
			path: "content/zh-cn/docs/setup/production-environment.md",
			expected: &Path{
				original:  "content/zh-cn/docs/setup/production-environment.md",
				language:  language.Chinese,
				category:  Docs,
				segments:  []string{"content", "zh-cn", "docs", "setup", "production-environment.md"},
				filename:  "production-environment",
				extension: ".md",
			},
			wantErr: false,
		},
		{
			name: "Portuguese path",
			path: "content/pt-br/blog/2024/announcement.md",
			expected: &Path{
				original:  "content/pt-br/blog/2024/announcement.md",
				language:  language.PortugueseBR,
				category:  Blog,
				segments:  []string{"content", "pt-br", "blog", "2024", "announcement.md"},
				filename:  "announcement",
				extension: ".md",
			},
			wantErr: false,
		},
		{
			name: "Non-content path",
			path: "static/images/logo.png",
			expected: &Path{
				original:  "static/images/logo.png",
				language:  "",
				category:  "overall",
				segments:  []string{"static", "images", "logo.png"},
				filename:  "logo",
				extension: ".png",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got.original != tt.expected.original {
				t.Errorf("Parse().original = %v, want %v", got.original, tt.expected.original)
			}
			if got.language != tt.expected.language {
				t.Errorf("Parse().language = %v, want %v", got.language, tt.expected.language)
			}
			if got.category != tt.expected.category {
				t.Errorf("Parse().category = %v, want %v", got.category, tt.expected.category)
			}
			if len(got.segments) != len(tt.expected.segments) {
				t.Errorf("Parse().segments length = %v, want %v", len(got.segments), len(tt.expected.segments))
			} else {
				for i, seg := range got.segments {
					if seg != tt.expected.segments[i] {
						t.Errorf("Parse().segments[%d] = %v, want %v", i, seg, tt.expected.segments[i])
					}
				}
			}
			if got.filename != tt.expected.filename {
				t.Errorf("Parse().filename = %v, want %v", got.filename, tt.expected.filename)
			}
			if got.extension != tt.expected.extension {
				t.Errorf("Parse().extension = %v, want %v", got.extension, tt.expected.extension)
			}
		})
	}
}

func TestParseWithValidation(t *testing.T) {
	validLangs := []language.Language{language.English, language.Japanese, language.Korean, language.Chinese}
	validExts := []string{".md", ".html"}
	validCategories := []Category{Docs, Blog}

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid Japanese docs path",
			path:    "content/ja/docs/overview.md",
			wantErr: false,
		},
		{
			name:    "Valid English blog path",
			path:    "content/en/blog/post.md",
			wantErr: false,
		},
		{
			name:    "Invalid extension",
			path:    "content/ja/docs/file.txt",
			wantErr: true,
			errMsg:  "invalid file extension",
		},
		{
			name:    "Does not start with content/",
			path:    "static/images/logo.md",
			wantErr: true,
			errMsg:  "path must start with 'content/'",
		},
		{
			name:    "Invalid path structure (too short)",
			path:    "content/ja.md",
			wantErr: true,
			errMsg:  "invalid path structure",
		},
		{
			name:    "Unsupported language",
			path:    "content/fr/docs/overview.md",
			wantErr: true,
			errMsg:  "unsupported language code",
		},
		{
			name:    "Unsupported category",
			path:    "content/ja/community/overview.md",
			wantErr: true,
			errMsg:  "unsupported category",
		},
		{
			name:    "Valid overall category",
			path:    "content/ja/_index.md",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseWithValidation(tt.path, validLangs, validExts, validCategories)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWithValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseWithValidation() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestPath_ToEnglishPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Japanese to English",
			path:     "content/ja/docs/concepts/overview.md",
			expected: "content/en/docs/concepts/overview.md",
		},
		{
			name:     "Korean to English",
			path:     "content/ko/blog/2023/release.md",
			expected: "content/en/blog/2023/release.md",
		},
		{
			name:     "Chinese to English",
			path:     "content/zh-cn/docs/setup.md",
			expected: "content/en/docs/setup.md",
		},
		{
			name:     "Already English (no change)",
			path:     "content/en/docs/overview.md",
			expected: "content/en/docs/overview.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse(tt.path)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			got := p.ToEnglishPath()
			if got != tt.expected {
				t.Errorf("ToEnglishPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPath_ToLanguagePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		lang     language.Language
		expected string
	}{
		{
			name:     "English to Japanese",
			path:     "content/en/docs/concepts/overview.md",
			lang:     language.Japanese,
			expected: "content/ja/docs/concepts/overview.md",
		},
		{
			name:     "English to Korean",
			path:     "content/en/blog/2023/release.md",
			lang:     language.Korean,
			expected: "content/ko/blog/2023/release.md",
		},
		{
			name:     "English to Chinese",
			path:     "content/en/docs/setup.md",
			lang:     language.Chinese,
			expected: "content/zh-cn/docs/setup.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse(tt.path)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			got := p.ToLanguagePath(tt.lang)
			if got != tt.expected {
				t.Errorf("ToLanguagePath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPath_Getters(t *testing.T) {
	path := "content/ja/docs/concepts/overview.md"
	p, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if p.Language() != language.Japanese {
		t.Errorf("Language() = %v, want %v", p.Language(), language.Japanese)
	}

	if p.Category() != Docs {
		t.Errorf("Category() = %v, want %v", p.Category(), Docs)
	}

	if p.Original() != path {
		t.Errorf("Original() = %v, want %v", p.Original(), path)
	}

	if p.Filename() != "overview" {
		t.Errorf("Filename() = %v, want %v", p.Filename(), "overview")
	}

	if p.Extension() != ".md" {
		t.Errorf("Extension() = %v, want %v", p.Extension(), ".md")
	}

	segments := p.Segments()
	expectedSegments := []string{"content", "ja", "docs", "concepts", "overview.md"}
	if len(segments) != len(expectedSegments) {
		t.Errorf("Segments() length = %v, want %v", len(segments), len(expectedSegments))
	}
}

func TestPath_IsContentFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Markdown content file",
			path:     "content/ja/docs/overview.md",
			expected: true,
		},
		{
			name:     "HTML content file",
			path:     "content/en/blog/post.html",
			expected: true,
		},
		{
			name:     "Non-content file",
			path:     "static/images/logo.png",
			expected: false,
		},
		{
			name:     "Non-markdown/html file in content",
			path:     "content/ja/docs/data.json",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse(tt.path)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			got := p.IsContentFile()
			if got != tt.expected {
				t.Errorf("IsContentFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPath_IsIndex(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Index markdown file",
			path:     "content/ja/docs/_index.md",
			expected: true,
		},
		{
			name:     "Index HTML file",
			path:     "content/en/blog/_index.html",
			expected: true,
		},
		{
			name:     "Regular markdown file",
			path:     "content/ja/docs/overview.md",
			expected: false,
		},
		{
			name:     "Index-like name but wrong extension",
			path:     "content/ja/docs/_index.txt",
			expected: false,
		},
		{
			name:     "File containing index in name",
			path:     "content/ja/docs/index-page.md",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse(tt.path)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			got := p.IsIndex()
			if got != tt.expected {
				t.Errorf("IsIndex() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
