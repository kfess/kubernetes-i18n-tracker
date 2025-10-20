package issue

import (
	"testing"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

func TestGuessLanguage(t *testing.T) {
	tests := []struct {
		name     string
		issue    Issue
		expected language.Language
	}{
		{
			name: "from language label",
			issue: Issue{
				Title:  "Update docs",
				Labels: []string{"language/ja", "kind/localization"},
			},
			expected: language.LanguageJapanese,
		},
		{
			name: "from title prefix",
			issue: Issue{
				Title:  "[ko] Update documentation",
				Labels: []string{},
			},
			expected: language.LanguageKorean,
		},
		{
			name: "from title prefix with pt-br",
			issue: Issue{
				Title:  "[pt-br] Update docs",
				Labels: []string{},
			},
			expected: language.LanguagePortugueseBR,
		},
		{
			name: "no language found",
			issue: Issue{
				Title:  "Update docs",
				Labels: []string{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GuessLanguage(tt.issue)
			if result != tt.expected {
				t.Errorf("GuessLanguage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractPathLikeString(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "path with slash prefix",
			title:    "[ja] Update /content/ja/docs/concepts/overview.md",
			expected: "content/ja/docs/concepts/overview.md",
		},
		{
			name:     "path without slash prefix",
			title:    "[ko] Fix content/ko/docs/tasks/configure-pod.md",
			expected: "content/ko/docs/tasks/configure-pod.md",
		},
		{
			name:     "path with underscores",
			title:    "Update docs/concepts/workloads/pods/pod_lifecycle.md",
			expected: "docs/concepts/workloads/pods/pod_lifecycle.md",
		},
		{
			name:     "no path found",
			title:    "[ja] General documentation update",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractPathLikeString(tt.title)
			if result != tt.expected {
				t.Errorf("ExtractPathLikeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeneratePathCandidates(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		lang          language.Language
		shouldContain []string
	}{
		{
			name: "simple path",
			path: "docs/concepts/overview.md",
			lang: language.LanguageJapanese,
			shouldContain: []string{
				"docs/concepts/overview.md",
				"content/ja/docs/concepts/overview.md",
			},
		},
		{
			name: "path with underscores",
			path: "docs/pod_lifecycle.md",
			lang: language.LanguageKorean,
			shouldContain: []string{
				"content/ko/docs/pod-lifecycle.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates := GeneratePathCandidates(tt.path, tt.lang)
			candidateMap := make(map[string]bool)
			for _, c := range candidates {
				candidateMap[c] = true
			}

			for _, expected := range tt.shouldContain {
				if !candidateMap[expected] {
					t.Errorf("GeneratePathCandidates() missing expected candidate: %v", expected)
				}
			}
		})
	}
}

func TestGuessPath(t *testing.T) {
	existingPaths := map[string]bool{
		"content/ja/docs/concepts/overview.md":   true,
		"content/ko/docs/tasks/configure-pod.md": true,
		"content/en/docs/reference/kubectl.md":   true,
	}

	tests := []struct {
		name     string
		issue    Issue
		lang     language.Language
		expected string
	}{
		{
			name: "exact match",
			issue: Issue{
				Title: "[ja] Update /content/ja/docs/concepts/overview.md",
			},
			lang:     language.LanguageJapanese,
			expected: "content/ja/docs/concepts/overview.md",
		},
		{
			name: "partial match with prefix",
			issue: Issue{
				Title: "[ko] Fix configure-pod.md",
			},
			lang:     language.LanguageKorean,
			expected: "content/ko/docs/tasks/configure-pod.md",
		},
		{
			name: "new file creation issue (file doesn't exist yet)",
			issue: Issue{
				Title: "[ja] Translated /docs/concepts/cluster-administration/compatibility-version.md into Japanese",
			},
			lang:     language.LanguageJapanese,
			expected: "", // File doesn't exist, so no match expected
		},
		{
			name: "no path extractable from title",
			issue: Issue{
				Title: "[ja] Update nonexistent file",
			},
			lang:     language.LanguageJapanese,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GuessPath(tt.issue, tt.lang, existingPaths)
			if result != tt.expected {
				t.Errorf("GuessPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}
