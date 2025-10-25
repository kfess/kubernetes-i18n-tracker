package translation

import (
	"math"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
)

// Status represents the translation status of a file.
type Status string

const (
	StatusUpToDate         Status = "up_to_date"
	StatusOutdated         Status = "outdated"
	StatusNotTranslated    Status = "not_translated"
	StatusNoEnglishVersion Status = "no_english_version"
	StatusUnknown          Status = "unknown"
)

// calculateStatus determines the translation status of a file based on its history.
func calculateStatus(
	language string,
	englishCommits []*git.Commit,
	translationCommits []*git.Commit,
) Status {
	if len(englishCommits) == 0 {
		return StatusNoEnglishVersion
	}

	// For English files, they are always up-to-date with themselves
	if language == "en" {
		return StatusUpToDate
	}

	if len(translationCommits) == 0 {
		return StatusNotTranslated
	}

	englishLatest := englishCommits[len(englishCommits)-1]
	translationLatest := translationCommits[len(translationCommits)-1]

	if englishLatest.Date.After(translationLatest.Date) {
		return StatusOutdated
	}

	if translationLatest.Date.Equal(englishLatest.Date) || translationLatest.Date.After(englishLatest.Date) {
		return StatusUpToDate
	}

	// This should not be reachable, but just in case
	return StatusUnknown
}

func calculateDaysBehind(englishCommits []*git.Commit, translationCommits []*git.Commit) int {
	// If there are no English commits, return 0 days behind
	if len(englishCommits) == 0 {
		return 0
	}

	var diffHours float64
	if len(translationCommits) == 0 {
		// If there are no translation commits, calculate from the first English commit
		englishFirst := englishCommits[0]
		diffHours = time.Since(englishFirst.Date).Hours()
	} else {
		// Compare the English latest and translation latest commits
		englishLatest := englishCommits[len(englishCommits)-1]
		translationLatest := translationCommits[len(translationCommits)-1]
		diffHours = englishLatest.Date.Sub(translationLatest.Date).Hours()
	}

	// Round to the nearest whole day
	// Without this, translations that are 23 hours behind would show as 0 days behind
	daysBehind := int(math.Round(diffHours / 24))
	if daysBehind < 0 {
		return 0
	}

	return daysBehind
}
