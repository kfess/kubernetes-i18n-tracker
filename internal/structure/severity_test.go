package structure

import (
	"testing"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

func versions(minors ...int) map[[2]int]struct{} {
	m := make(map[[2]int]struct{}, len(minors))
	for _, minor := range minors {
		m[[2]int{1, minor}] = struct{}{}
	}
	return m
}

// The implementation was checked against upstream's analyze_file_pair() on
// 200,000 generated page pairs and on real pages; these cases pin one example
// per rule.
func TestCompare(t *testing.T) {
	// A 120-line English page with every kind of structure.
	en := Features{
		VisibleLines:       120,
		H2:                 4,
		H3:                 6,
		CodeBlocks:         4,
		Anchors:            set("a", "b", "c", "d", "e", "f"),
		BodyWords:          1000,
		Versions:           versions(28, 29, 30, 31),
		FeatureStateTokens: set("version:v1.31"),
		ApiKindTokens:      set("api:apps/v1", "kind:Deployment"),
	}
	// with returns a copy of en changed by fn, standing in for a translation.
	with := func(fn func(f *Features)) Features {
		l10n := en
		fn(&l10n)
		return l10n
	}

	tests := []struct {
		name string
		lang language.Language
		en   Features
		l10n Features
		want Signal
	}{
		{
			name: "identical structure",
			lang: language.Japanese,
			en:   en,
			l10n: en,
			want: SignalNone,
		},
		{
			name: "empty localized page",
			lang: language.Japanese,
			en:   en,
			l10n: Features{},
			want: SignalStrong,
		},
		{
			name: "feature state and API both differ",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) {
				f.FeatureStateTokens = set("version:v1.28")
				f.ApiKindTokens = set("api:apps/v1beta1", "kind:Deployment")
			}),
			want: SignalStrong,
		},
		{
			name: "feature state alone differs",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) { f.FeatureStateTokens = set("version:v1.28") }),
			want: SignalNone,
		},
		{
			name: "two severe losses",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) { f.H2, f.CodeBlocks = 2, 1 }),
			want: SignalStrong,
		},
		{
			name: "one severe and one moderate loss",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) { f.H2, f.CodeBlocks = 2, 3 }),
			want: SignalStrong,
		},
		{
			name: "one severe loss alone",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) { f.Versions = versions(28) }),
			want: SignalModerate,
		},
		{
			name: "one missing H3 is below the moderate cutoff",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) { f.H3 = 5 }),
			want: SignalNone,
		},
		{
			name: "Japanese ignores a length gap with no other loss",
			lang: language.Japanese,
			en:   en,
			l10n: with(func(f *Features) { f.VisibleLines = 60 }),
			want: SignalNone,
		},
		{
			name: "Russian keeps a length gap with no other loss",
			lang: language.Russian,
			en:   en,
			l10n: with(func(f *Features) { f.VisibleLines = 60 }),
			want: SignalModerate,
		},
		{
			name: "Japanese length gap backed by a missing code block",
			lang: language.Japanese,
			en:   en,
			l10n: with(func(f *Features) { f.VisibleLines, f.CodeBlocks = 70, 3 }),
			want: SignalModerate,
		},
		{
			name: "Spanish full word volume suppresses a length-only gap",
			lang: language.Spanish,
			en:   en,
			l10n: with(func(f *Features) { f.VisibleLines, f.BodyWords = 50, 950 }),
			want: SignalNone,
		},
		{
			name: "Spanish thin word volume keeps a length-only gap",
			lang: language.Spanish,
			en:   en,
			l10n: with(func(f *Features) { f.VisibleLines, f.BodyWords = 50, 500 }),
			want: SignalModerate,
		},
		{
			name: "large gap with a missing heading and anchors",
			lang: language.Spanish,
			en:   en,
			l10n: with(func(f *Features) {
				f.VisibleLines, f.H2 = 50, 3
				f.Anchors = set("a", "b", "c", "d")
			}),
			want: SignalStrong,
		},
		{
			name: "Spanish translated anchors do not promote a large gap",
			lang: language.Spanish,
			en:   en,
			l10n: with(func(f *Features) {
				f.VisibleLines, f.BodyWords = 50, 950
				f.Anchors = set("a", "b", "c", "d")
			}),
			want: SignalModerate,
		},
		{
			name: "Russian large gap with missing anchors is promoted",
			lang: language.Russian,
			en:   en,
			l10n: with(func(f *Features) {
				f.VisibleLines, f.BodyWords = 50, 950
				f.Anchors = set("a", "b", "c", "d")
			}),
			want: SignalStrong,
		},
		{
			name: "small gap alone",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) { f.VisibleLines = 90 }),
			want: SignalNone,
		},
		{
			name: "small gap backed by a feature state difference",
			lang: language.Korean,
			en:   en,
			l10n: with(func(f *Features) {
				f.VisibleLines = 90
				f.FeatureStateTokens = set("version:v1.28")
			}),
			want: SignalModerate,
		},
		{
			name: "short English page protects a length-only gap",
			lang: language.Polish,
			en:   Features{VisibleLines: 30, BodyWords: 200},
			l10n: Features{VisibleLines: 10, BodyWords: 60},
			want: SignalNone,
		},
		{
			name: "short English page under the CJK cutoff",
			lang: language.Chinese,
			en:   Features{VisibleLines: 50, BodyWords: 300},
			l10n: Features{VisibleLines: 20, BodyWords: 100},
			want: SignalNone,
		},
		{
			name: "same page over the default cutoff",
			lang: language.Polish,
			en:   Features{VisibleLines: 50, BodyWords: 300},
			l10n: Features{VisibleLines: 20, BodyWords: 100},
			want: SignalModerate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compare(tt.en, tt.l10n, tt.lang)
			if got.Signal != tt.want {
				t.Errorf("Compare() signal = %q, want %q (gap %+v)", got.Signal, tt.want, got.Gap)
			}
		})
	}
}
