package structure

import "github.com/kfess/kubernetes-i18n-tracker/internal/language"

// These thresholds are empirically set per locale;
// a statistical approach (median+MAD from structurally-clean files)
// could replace them with auto-calibrated per-locale values.
const (
	// Below this, length gap requires a companion indicator
	shortEnThreshold = 40

	// EN files shorter than this skip all length-gap checks
	lengthGapMinEnLines = 15

	// same guard for CJK (denser text, files run longer)
	cjkShortEnThreshold = 56

	// Latin compactness guard activates only above this
	latinMinEnLines = 55

	// Empirical boundary: Latin false alarms sit at >=0.94, real ru mismatch at <=0.85.
	// body-word ratio floor for "full translation" in Latin guard
	latinBodyRatioMin = 0.90
)

type LocaleProfile struct {
	Script    string
	Direction string

	// Length Gap behabior
	ShortEnThreshold       int
	MinEnLinesForLengthGap int
	IgnoreOnlyLengthGap    bool

	// Optional compactness guard (Latin-script: word counts behave like EN)
	CompactnessGuard      bool
	CompactnessMinEnLines int
	BodyWordRatioMin      float64

	// Free-form note recording empirical rationale or future calibration plans.
	Note string
}

// defaultLocaleProfile returns the profile used for every field a locale does
// not override.
func defaultLocaleProfile() LocaleProfile {
	return LocaleProfile{
		Script:                 "unknown",
		Direction:              "ltr",
		ShortEnThreshold:       shortEnThreshold,
		MinEnLinesForLengthGap: lengthGapMinEnLines,
		CompactnessMinEnLines:  latinMinEnLines,
		BodyWordRatioMin:       latinBodyRatioMin,
	}
}

// getLocaleProfile returns the calibration for a locale. Locales without a
// specific calibration get the default profile.
func getLocaleProfile(lang language.Language) LocaleProfile {
	p := defaultLocaleProfile()

	switch lang {
	case language.Korean, language.Chinese:
		p.Script = "cjk"
		p.ShortEnThreshold = cjkShortEnThreshold
	case language.Japanese:
		p.Script = "cjk"
		p.ShortEnThreshold = cjkShortEnThreshold
		p.IgnoreOnlyLengthGap = true
		p.Note = "Japanese can produce line-count-only false alarms because " +
			"physical wrapping differs from English; non-length structural " +
			"signals are still checked."
	case language.PortugueseBR, language.Spanish, language.German, language.French, language.Italian:
		p.Script = "latin"
		p.CompactnessGuard = true
		p.Note = "Latin-script compactness guard can suppress line-count-only false " +
			"alarms when a localized page has high body-word ratio, because " +
			"full translations may wrap differently from English. Non-length " +
			"structural signals are still checked."
	case language.Russian, language.Ukrainian:
		// Deliberately no compactness guard: real mismatches sit at low body-word ratios.
		p.Script = "cyrillic"
	case language.Persian:
		p.Script = "arabic"
		p.Direction = "rtl"
	}

	return p
}
