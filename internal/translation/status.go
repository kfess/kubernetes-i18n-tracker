package translation

import (
	"math"
	"strings"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
)

// Status represents the translation status of a file.
type Status string

const (
	StatusUpToDate         Status = "up_to_date"
	StatusOutdated         Status = "outdated"
	StatusPossiblyOutdated Status = "possibly_outdated"
	StatusNotTranslated    Status = "not_translated"
	StatusNoEnglishVersion Status = "no_english_version"
	StatusUnknown          Status = "unknown"
)

// calculateStatus determines the translation status of a file based on its history.
// It now includes header count comparison to detect structural differences.
func calculateStatus(
	language string,
	englishCommits []*git.Commit,
	translationCommits []*git.Commit,
	englishContent string,
	translationContent string,
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

	// Translation appears to be up-to-date, but re-evaluate by:
	// 1. Checking if header counts match (structural similarity)
	// 2. Looking at translation commits that happened after the English latest commit
	//    If all such commits are minor, treat as possibly_outdated
	if translationLatest.Date.Equal(englishLatest.Date) || translationLatest.Date.After(englishLatest.Date) {
		// Check header counts - if they differ, the translation is possibly outdated
		if englishContent != "" && translationContent != "" {
			englishHeaders := CountHeadersFromContent(englishContent)
			translationHeaders := CountHeadersFromContent(translationContent)

			if !HeadersMatch(englishHeaders, translationHeaders) {
				logger.Infof("Header mismatch detected (EN: %d, Trans: %d) - marking as possibly_outdated",
					englishHeaders.Total, translationHeaders.Total)
				return StatusPossiblyOutdated
			}
		}
		// Collect translation commits that occurred after the English latest commit
		var afterEnglish []*git.Commit
		for _, c := range translationCommits {
			if c.Date.After(englishLatest.Date) {
				afterEnglish = append(afterEnglish, c)
			}
		}

		// If there are no translation commits after English, it's truly up-to-date
		if len(afterEnglish) == 0 {
			return StatusUpToDate
		}

		// If any commit after English is a major change, the translation is truly up-to-date
		for _, c := range afterEnglish {
			if !isMinorCommit(c) {
				return StatusUpToDate
			}
		}

		// All commits that happened after the English latest are minor -> possibly outdated
		logger.Infof("Possibly outdated translation for language %s, English hash: %s, translation latest: %s", language, englishLatest.Hash, translationLatest.Hash)
		return StatusPossiblyOutdated
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

func maybeMinorChange(commit git.Commit) bool {
	msg := commit.Message
	minorKeywords := []string{
		"typo", "spelling", "grammar", "format", "whitespace",
		"punctuation", "minor", "chore", "fix",
	}

	for _, keyword := range minorKeywords {
		if strings.Contains(strings.ToLower(msg), keyword) {
			return true
		}
	}

	return false
}

func mustBeMajorChange(commit git.Commit) bool {
	msg := commit.Message
	majorKeywords := []string{"translate", "sync", "create", "add", "update"}

	for _, keyword := range majorKeywords {
		if strings.Contains(strings.ToLower(msg), keyword) {
			return true
		}
	}

	return false
}

// isMinorCommit determines if a commit should be considered minor and ignored for status calculation.
// This is heuristic-based to filter out trivial changes in translation files.
// So, there may be false positives/negatives.
// A commit is minor if:
// 1. Contains major keywords (translate, sync, etc.) → NOT minor (major takes precedence)
// 2. Contains minor keywords (typo, formatting, etc.) → minor
// 3. Has small number of changes (<=10 lines) → minor
// 4. Otherwise → NOT minor
func isMinorCommit(commit *git.Commit) bool {
	hasMinor := maybeMinorChange(*commit)
	hasMajor := mustBeMajorChange(*commit)

	// If major keywords present, it's definitely NOT minor (major takes precedence)
	if hasMajor {
		return false
	}

	totalChanges := commit.Insertions + commit.Deletions

	if hasMinor || totalChanges <= 10 {
		return true
	}

	// Default to NOT minor (stricter heuristic to avoid false positives)
	return false
}
