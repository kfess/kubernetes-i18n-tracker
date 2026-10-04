package translation

import (
	"context"
	"testing"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/history"
	"github.com/kfess/kubernetes-i18n-tracker/internal/structure"
)

func TestGetTranslationStatusCombinesStructure(t *testing.T) {
	const (
		englishPath     = "content/en/docs/concepts/foo.md"
		syncedPath      = "content/ja/docs/concepts/foo.md"
		untranslated    = "content/ko/docs/concepts/foo.md"
		englishContent  = "---\ntitle: Foo\n---\n## One\n\ntext\n\n## Two\n\ntext\n\n## Three\n\ntext\n"
		completeContent = "---\ntitle: Foo\n---\n## いち\n\n本文\n\n## に\n\n本文\n\n## さん\n\n本文\n"
		// Two of the three H2 sections are missing.
		thinContent = "---\ntitle: Foo\n---\n## いち\n\n本文\n"
	)

	event := func(path, date string) *git.Event {
		n := 1
		return &git.Event{
			Hash:    "hash-" + path + date,
			Date:    date,
			Message: "update",
			File:    git.FileInfo{Path: path, Insertions: &n, Deletions: &n},
		}
	}
	// The Japanese page was committed after the latest English commit.
	h := history.Build([]*git.Event{
		event(englishPath, "2026-01-01 00:00:00 +0000"),
		event(syncedPath, "2026-02-01 00:00:00 +0000"),
	})
	tracker := NewTracker(h, nil, nil, nil, nil, nil, Config{})
	english := structure.Parse(englishContent)

	tests := []struct {
		name            string
		path            string
		englishFeatures *structure.Features
		content         string
		wantStatus      Status
		wantGitStatus   Status
		wantSignal      structure.Signal // empty when no comparison is expected
	}{
		{
			name:            "git up to date and structure matches",
			path:            syncedPath,
			englishFeatures: &english,
			content:         completeContent,
			wantStatus:      StatusUpToDate,
			wantGitStatus:   StatusUpToDate,
			wantSignal:      structure.SignalNone,
		},
		{
			name:            "git up to date but sections are missing",
			path:            syncedPath,
			englishFeatures: &english,
			content:         thinContent,
			wantStatus:      StatusPossiblyOutdated,
			wantGitStatus:   StatusUpToDate,
			wantSignal:      structure.SignalModerate,
		},
		{
			name:            "unreadable English page skips the comparison",
			path:            syncedPath,
			englishFeatures: nil,
			content:         thinContent,
			wantStatus:      StatusUpToDate,
			wantGitStatus:   StatusUpToDate,
		},
		{
			name:            "unreadable translation skips the comparison",
			path:            syncedPath,
			englishFeatures: &english,
			content:         "",
			wantStatus:      StatusUpToDate,
			wantGitStatus:   StatusUpToDate,
		},
		{
			name:            "not translated has no comparison",
			path:            untranslated,
			englishFeatures: &english,
			content:         "",
			wantStatus:      StatusNotTranslated,
			wantGitStatus:   StatusNotTranslated,
		},
		{
			name:            "English itself has no comparison",
			path:            englishPath,
			englishFeatures: &english,
			content:         englishContent,
			wantStatus:      StatusUpToDate,
			wantGitStatus:   StatusUpToDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tracker.GetTranslationStatus(context.Background(), tt.path, tt.englishFeatures, tt.content)
			if err != nil {
				t.Fatalf("GetTranslationStatus() error = %v", err)
			}

			if got.History.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", got.History.Status, tt.wantStatus)
			}
			if got.History.GitStatus != tt.wantGitStatus {
				t.Errorf("GitStatus = %v, want %v", got.History.GitStatus, tt.wantGitStatus)
			}

			if tt.wantSignal == "" {
				if got.History.Structure != nil {
					t.Errorf("Structure = %+v, want nil", got.History.Structure)
				}
				return
			}
			if got.History.Structure == nil {
				t.Fatalf("Structure = nil, want signal %v", tt.wantSignal)
			}
			if got.History.Structure.Signal != tt.wantSignal {
				t.Errorf("Signal = %v, want %v", got.History.Structure.Signal, tt.wantSignal)
			}
		})
	}
}
