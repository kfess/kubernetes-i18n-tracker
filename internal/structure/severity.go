package structure

import "github.com/kfess/kubernetes-i18n-tracker/internal/language"

// Signal is the overall verdict for one localized page.
type Signal string

const (
	SignalNone     Signal = "none"
	SignalModerate Signal = "moderate"
	SignalStrong   Signal = "strong"
)

// severity is the verdict for a single factor (headings, code blocks, ...).
type severity int

const (
	severityNone severity = iota
	severityModerate
	severitySevere
)

// lengthGapSeverity grades how much shorter the localized page is.
type lengthGapSeverity int

const (
	lengthGapNone lengthGapSeverity = iota
	lengthGapSmall
	lengthGapModerate
	lengthGapLarge
)

// Cutoffs for when missing items count as a "strong" outdatedness signal;
// see classify() for how strong signals affect the final status.
const (
	// missing H2 headings
	h2Threshold = 2

	// missing H3s, plus at least one missing H2
	h3WithH2Threshold = 5

	// missing code blocks
	codeThreshold = 3

	// missing section anchors
	anchorThreshold = 5

	// missing newer Kubernetes version refs
	versionThreshold = 3
)

// indicators holds the per-factor verdicts the overall signal is derived from.
type indicators struct {
	emptyStub     bool
	apiAndFeature bool
	heading       severity
	code          severity
	anchor        severity
	version       severity
	lengthGap     lengthGapSeverity
}

// Result is the outcome of comparing a localized page with its English source.
type Result struct {
	Signal Signal `json:"signal"`
	Gap    Gap    `json:"gap"`
}

// Compare grades how far the localized page has structurally fallen behind the English one.
func Compare(en, l10n Features, lang language.Language) Result {
	profile := getLocaleProfile(lang)
	gap := featureGap(en, l10n)
	ind := buildIndicators(gap, en, l10n, profile)

	return Result{
		Signal: classify(ind, gap, profile),
		Gap:    gap,
	}
}

// classifyLengthGap grades the localized-to-English line ratio: below 0.50 is
// large, below 0.65 moderate, below 0.80 small.
func classifyLengthGap(lineRatio float64) lengthGapSeverity {
	switch {
	case lineRatio < 0.50:
		return lengthGapLarge
	case lineRatio < 0.65:
		return lengthGapModerate
	case lineRatio < 0.80:
		return lengthGapSmall
	default:
		return lengthGapNone
	}
}

// hasLengthGapSupport reports whether a non-length loss backs up a length gap.
// Missing anchors do not count: they are too often a naming difference.
func hasLengthGapSupport(gap Gap) bool {
	return gap.MissingH2 > 0 || gap.MissingH3 > 0 ||
		gap.MissingCodeBlocks > 0 || gap.MissingNewVersions > 0
}

// isOnlyLengthGap reports whether no structural loss is present at all.
func isOnlyLengthGap(gap Gap) bool {
	return gap.MissingH2 == 0 && gap.MissingH3 == 0 &&
		gap.MissingCodeBlocks == 0 && gap.MissingAnchors == 0 &&
		gap.MissingNewVersions == 0
}

// shouldIgnoreLengthGap drops a length gap that is unreliable (short English
// page with nothing else firing) or expected for the locale (Japanese line
// wrapping, loosely wrapped Latin-script translations with full word volume).
func shouldIgnoreLengthGap(gap Gap, en, l10n Features, profile LocaleProfile, lg lengthGapSeverity) bool {
	if lg == lengthGapNone || l10n.VisibleLines == 0 {
		return false
	}
	if !hasLengthGapSupport(gap) && en.VisibleLines < profile.ShortEnThreshold {
		return true
	}
	if (lg == lengthGapModerate || lg == lengthGapLarge) && isOnlyLengthGap(gap) {
		if profile.IgnoreOnlyLengthGap && en.VisibleLines >= profile.ShortEnThreshold {
			return true
		}
		if profile.CompactnessGuard &&
			en.VisibleLines >= profile.CompactnessMinEnLines &&
			gap.L10nToEnBodyWordRatio >= profile.BodyWordRatioMin {
			return true
		}
	}
	return false
}

// buildIndicators grades each factor of the gap. The length gap is graded
// first and then discarded when the page is too short, when the locale profile
// says it is unreliable, or, for a small gap, when no other factor backs it.
func buildIndicators(gap Gap, en, l10n Features, profile LocaleProfile) indicators {
	var ind indicators

	ind.emptyStub = l10n.VisibleLines == 0 && en.VisibleLines >= 1

	// Very short English pages skip the length check; an empty localized page
	// bypasses that floor.
	lg := lengthGapNone
	if en.VisibleLines >= profile.MinEnLinesForLengthGap || l10n.VisibleLines == 0 {
		lg = classifyLengthGap(gap.L10nToEnLineRatio)
	}
	if shouldIgnoreLengthGap(gap, en, l10n, profile, lg) {
		lg = lengthGapNone
	}
	// An empty localized page is reported as emptyStub, not as a length gap.
	if l10n.VisibleLines == 0 {
		lg = lengthGapNone
	}

	switch {
	case gap.MissingH2 >= h2Threshold || (gap.MissingH2 >= 1 && gap.MissingH3 >= h3WithH2Threshold):
		ind.heading = severitySevere
	case gap.MissingH2 >= 1 || gap.MissingH3 >= 2:
		ind.heading = severityModerate
	}

	ind.code = gradeMissing(gap.MissingCodeBlocks, codeThreshold)
	ind.anchor = gradeMissing(gap.MissingAnchors, anchorThreshold)
	ind.version = gradeMissing(gap.MissingNewVersions, versionThreshold)

	ind.apiAndFeature = gap.MissingFeatureState > 0 && gap.MissingApiOrKind > 0

	// A small gap only counts when something else points the same way.
	if lg == lengthGapSmall {
		backed := ind.heading != severityNone || ind.code != severityNone ||
			ind.anchor != severityNone || ind.version != severityNone ||
			gap.MissingFeatureState > 0 || gap.MissingApiOrKind > 0
		if !backed {
			lg = lengthGapNone
		}
	}
	ind.lengthGap = lg

	return ind
}

// gradeMissing grades a missing-item count: severe at the threshold, moderate
// from one missing item.
func gradeMissing(missing, severeThreshold int) severity {
	switch {
	case missing >= severeThreshold:
		return severitySevere
	case missing >= 1:
		return severityModerate
	default:
		return severityNone
	}
}

// isExpectedTranslatedAnchorLoss covers Latin-script locales that translate
// anchor IDs: the page looks compact with a small anchor mismatch but carries
// the full word volume, so it must not be promoted to strong.
func isExpectedTranslatedAnchorLoss(ind indicators, moderateFactors int, gap Gap, profile LocaleProfile) bool {
	onlyAnchor := ind.anchor == severityModerate && moderateFactors == 1
	return profile.CompactnessGuard &&
		gap.L10nToEnBodyWordRatio >= profile.BodyWordRatioMin &&
		gap.MissingAnchors <= 2 &&
		onlyAnchor
}

// classify derives the overall signal from the indicators. Rules apply in order:
//
//  1. empty stub, or API and feature-state mismatch together -> strong
//  2. two or more strong indicators -> strong
//  3. a strong and a supporting indicator -> strong
//  4. large length gap plus a non-length supporting indicator -> strong
//     (unless it is the expected translated-anchor case)
//  5. any strong or supporting indicator -> moderate
//  6. a backed small length gap -> moderate
//  7. otherwise -> none
func classify(ind indicators, gap Gap, profile LocaleProfile) Signal {
	factors := [...]severity{ind.heading, ind.code, ind.anchor, ind.version}

	// Strong: empty stub and every severe factor. Supporting: every moderate
	// factor and a moderate or large length gap.
	strong, moderateFactors := 0, 0
	if ind.emptyStub {
		strong++
	}
	for _, f := range factors {
		switch f {
		case severitySevere:
			strong++
		case severityModerate:
			moderateFactors++
		}
	}
	supporting := moderateFactors
	if ind.lengthGap == lengthGapModerate || ind.lengthGap == lengthGapLarge {
		supporting++
	}

	switch {
	case ind.emptyStub, ind.apiAndFeature:
		return SignalStrong
	case strong >= 2:
		return SignalStrong
	case strong >= 1 && supporting >= 1:
		return SignalStrong
	case ind.lengthGap == lengthGapLarge && moderateFactors >= 1 &&
		!isExpectedTranslatedAnchorLoss(ind, moderateFactors, gap, profile):
		return SignalStrong
	case strong >= 1 || supporting >= 1:
		return SignalModerate
	case ind.lengthGap == lengthGapSmall:
		return SignalModerate
	default:
		return SignalNone
	}
}
