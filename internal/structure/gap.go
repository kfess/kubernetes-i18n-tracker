package structure

// Gap describes what a localized page lacks compared with the English one.
type Gap struct {
	// Localized-to-English ratio of visible lines, capped at 2.
	L10nToEnLineRatio float64 `json:"lineRatio"`
	// Localized-to-English ratio of prose words.
	L10nToEnBodyWordRatio float64 `json:"bodyWordRatio"`

	MissingH2         int `json:"missingH2"`
	MissingH3         int `json:"missingH3"`
	MissingCodeBlocks int `json:"missingCodeBlocks"`
	MissingAnchors    int `json:"missingAnchors"`
	// English version references newer than the newest one the localized page mentions.
	MissingNewVersions  int `json:"missingNewVersions"`
	MissingFeatureState int `json:"missingFeatureState"`
	MissingApiOrKind    int `json:"missingApiOrKind"`
}

// featureGap measures what the localized page lacks compared with the English
// one. Counts are one-directional: extra content on the localized side is
// ignored. An empty localized page reports no missing feature-state or API
// tokens, since it is graded as an empty stub instead.
func featureGap(en, l10n Features) Gap {
	var lineRatio = 1.0
	if en.VisibleLines > 0 {
		lineRatio = min(2.0, float64(l10n.VisibleLines)/float64(en.VisibleLines))
	}

	var bodyWordRatio = 1.0
	if en.BodyWords > 0 {
		bodyWordRatio = float64(l10n.BodyWords) / float64(en.BodyWords)
	}

	missingAnchors := countMissing(en.Anchors, l10n.Anchors)

	missingFeatureState := 0
	missingApiOrKind := 0
	if l10n.VisibleLines > 0 {
		missingFeatureState = countMissing(en.FeatureStateTokens, l10n.FeatureStateTokens)
		missingApiOrKind = countMissing(en.ApiKindTokens, l10n.ApiKindTokens)
	}

	return Gap{
		L10nToEnLineRatio:     lineRatio,
		L10nToEnBodyWordRatio: bodyWordRatio,
		MissingH2:             max(0, en.H2-l10n.H2),
		MissingH3:             max(0, en.H3-l10n.H3),
		MissingCodeBlocks:     max(0, en.CodeBlocks-l10n.CodeBlocks),
		MissingAnchors:        missingAnchors,
		MissingNewVersions:    countMissingNewVersions(en.Versions, l10n.Versions),
		MissingFeatureState:   missingFeatureState,
		MissingApiOrKind:      missingApiOrKind,
	}
}

// countMissing counts the entries of en that l10n does not have (set difference).
func countMissing(en, l10n map[string]struct{}) int {
	missing := 0
	for token := range en {
		if _, ok := l10n[token]; !ok {
			missing++
		}
	}
	return missing
}

// countMissingNewVersions counts the English versions that are newer than the
// newest version the localized page mentions.
func countMissingNewVersions(enVersion, l10nVersion map[[2]int]struct{}) int {
	if len(l10nVersion) == 0 {
		return len(enVersion)
	}

	var l10nMax [2]int
	first := true
	for v := range l10nVersion {
		if first || versionLess(l10nMax, v) {
			l10nMax = v
			first = false
		}
	}

	missing := 0
	for v := range enVersion {
		if versionLess(l10nMax, v) {
			missing++
		}
	}
	return missing
}

// versionLess reports whether a is older than b, comparing major then minor.
func versionLess(a, b [2]int) bool {
	if a[0] != b[0] {
		return a[0] < b[0]
	}
	return a[1] < b[1]
}
