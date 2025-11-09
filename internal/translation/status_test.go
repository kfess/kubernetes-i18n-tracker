package translation

import (
	"testing"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
)

func TestCalculateStatus(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)

	tests := []struct {
		name               string
		language           string
		englishCommits     []*git.Commit
		translationCommits []*git.Commit
		want               Status
	}{
		{
			name:               "English file is always up-to-date",
			language:           "en",
			englishCommits:     []*git.Commit{{Date: now}},
			translationCommits: nil,
			want:               StatusUpToDate,
		},
		{
			name:               "No English version exists",
			language:           "ja",
			englishCommits:     []*git.Commit{},
			translationCommits: []*git.Commit{{Date: now}},
			want:               StatusNoEnglishVersion,
		},
		{
			name:               "Translation doesn't exist (no commits)",
			language:           "ja",
			englishCommits:     []*git.Commit{{Date: now}},
			translationCommits: []*git.Commit{},
			want:               StatusNotTranslated,
		},
		{
			name:     "Translation is outdated",
			language: "ja",
			englishCommits: []*git.Commit{
				{Date: twoDaysAgo},
				{Date: now},
			},
			translationCommits: []*git.Commit{
				{Date: yesterday},
			},
			want: StatusOutdated,
		},
		{
			name:     "Translation is up-to-date",
			language: "ja",
			englishCommits: []*git.Commit{
				{Date: twoDaysAgo, Message: "Add new content"},
			},
			translationCommits: []*git.Commit{
				{Date: yesterday, Message: "Translate page"},
				{Date: now, Message: "Update translation"}, // Major change after English
			},
			want: StatusUpToDate,
		},
		{
			name:     "Translation is up-to-date (same date)",
			language: "fr",
			englishCommits: []*git.Commit{
				{Date: now},
			},
			translationCommits: []*git.Commit{
				{Date: now},
			},
			want: StatusUpToDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For unit tests, we pass empty strings for content to skip header checking
			got := calculateStatus(tt.language, tt.englishCommits, tt.translationCommits, "", "")
			if got != tt.want {
				t.Errorf("CalculateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateStatusWithHeaderCheck(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	tests := []struct {
		name               string
		language           string
		englishCommits     []*git.Commit
		translationCommits []*git.Commit
		englishContent     string
		translationContent string
		want               Status
	}{
		{
			name:     "Header mismatch marks as possibly_outdated",
			language: "ja",
			englishCommits: []*git.Commit{
				{Date: yesterday, Message: "Add content"},
			},
			translationCommits: []*git.Commit{
				{Date: now, Message: "Translate content"},
			},
			englishContent: `---
title: Test
---
# Header 1
## Header 2
### Header 3
Content here
`,
			translationContent: `---
title: Test
---
# Header 1
## Header 2
Content here
`,
			want: StatusPossiblyOutdated,
		},
		{
			name:     "Header match confirms up_to_date",
			language: "ja",
			englishCommits: []*git.Commit{
				{Date: yesterday, Message: "Add content"},
			},
			translationCommits: []*git.Commit{
				{Date: now, Message: "Update translation"},
			},
			englishContent: `---
title: Test
---
# Header 1
## Header 2
### Header 3
Content here
`,
			translationContent: `---
title: テスト
---
# ヘッダー 1
## ヘッダー 2
### ヘッダー 3
翻訳内容
`,
			want: StatusUpToDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateStatus(tt.language, tt.englishCommits, tt.translationCommits, tt.englishContent, tt.translationContent)
			if got != tt.want {
				t.Errorf("CalculateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDaysBehind(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)
	threeDaysAgo := now.Add(-72 * time.Hour)

	tests := []struct {
		name               string
		englishCommits     []*git.Commit
		translationCommits []*git.Commit
		want               int
	}{
		{
			name: "Translation up-to-date",
			englishCommits: []*git.Commit{
				{Date: twoDaysAgo},
			},
			translationCommits: []*git.Commit{
				{Date: now},
			},
			want: 0,
		},
		{
			name: "Translation one day behind",
			englishCommits: []*git.Commit{
				{Date: yesterday},
			},
			translationCommits: []*git.Commit{
				{Date: now},
			},
			want: 0,
		},
		{
			name: "Translation two days behind",
			englishCommits: []*git.Commit{
				{Date: now},
			},
			translationCommits: []*git.Commit{
				{Date: twoDaysAgo},
			},
			want: 2,
		},
		{
			name:               "No English commits",
			englishCommits:     []*git.Commit{},
			translationCommits: []*git.Commit{{Date: now}},
			want:               0,
		},
		{
			name: "No translation commits",
			englishCommits: []*git.Commit{
				{Date: now},
			},
			translationCommits: []*git.Commit{},
			want:               0,
		},
		{
			name: "Multiple English commits, translation behind",
			englishCommits: []*git.Commit{
				{Date: threeDaysAgo},
				{Date: yesterday},
				{Date: now},
			},
			translationCommits: []*git.Commit{
				{Date: twoDaysAgo},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateDaysBehind(tt.englishCommits, tt.translationCommits)
			if got != tt.want {
				t.Errorf("GetDaysBehind() = %v, want %v", got, tt.want)
			}
		})
	}
}
