package translation

import (
	"testing"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/structure"
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
			got := calculateStatus(tt.language, tt.englishCommits, tt.translationCommits)
			if got != tt.want {
				t.Errorf("CalculateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCombineStatus(t *testing.T) {
	signal := func(s structure.Signal) *structure.Result {
		return &structure.Result{Signal: s}
	}

	tests := []struct {
		name       string
		gitStatus  Status
		structural *structure.Result
		want       Status
	}{
		{"up to date, no comparison", StatusUpToDate, nil, StatusUpToDate},
		{"up to date, no signal", StatusUpToDate, signal(structure.SignalNone), StatusUpToDate},
		{"up to date, moderate signal", StatusUpToDate, signal(structure.SignalModerate), StatusPossiblyOutdated},
		{"up to date, strong signal", StatusUpToDate, signal(structure.SignalStrong), StatusPossiblyOutdated},
		{"outdated stays outdated without a signal", StatusOutdated, signal(structure.SignalNone), StatusOutdated},
		{"outdated stays outdated with a signal", StatusOutdated, signal(structure.SignalStrong), StatusOutdated},
		{"not translated is untouched", StatusNotTranslated, nil, StatusNotTranslated},
		{"no English version is untouched", StatusNoEnglishVersion, signal(structure.SignalStrong), StatusNoEnglishVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineStatus(tt.gitStatus, tt.structural); got != tt.want {
				t.Errorf("combineStatus() = %v, want %v", got, tt.want)
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
