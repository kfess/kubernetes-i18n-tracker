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

	// Translation appears to be up-to-date, but check if it's truly up-to-date
	// or possibly outdated (latest translation commits are only minor changes)
	if translationLatest.Date.Equal(englishLatest.Date) || translationLatest.Date.After(englishLatest.Date) {
		// Filter out minor commits from translation files
		majorTranslationCommits := filterMajorCommits(translationCommits)

		// If there are no major translation commits, it's outdated (no real translation work)
		if len(majorTranslationCommits) == 0 {
			// This is not true, but just to be safe
			return StatusUpToDate
		}

		// Compare English latest with the latest *major* translation commit
		majorTranslationLatest := majorTranslationCommits[len(majorTranslationCommits)-1]

		// If major translation commit is older than English, it's possibly outdated
		// (the translation appears up-to-date only because of minor commits like typo fixes)
		if englishLatest.Date.After(majorTranslationLatest.Date) {
			logger.Infof("Possibly outdated translation for language %s, English hash: %s, major translation hash: %s", language, englishLatest.Hash, majorTranslationLatest.Hash)
			return StatusPossiblyOutdated
		}

		// Major translation commit is same or newer than English, truly up-to-date
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
//  1. Contains minor keywords (typo, spelling, etc.) - takes precedence even if major keywords present
//  2. OR total changes <= 20 lines AND does NOT contain major keywords
func isMinorCommit(commit *git.Commit) bool {
	hasMinor := maybeMinorChange(*commit)
	hasMajor := mustBeMajorChange(*commit)

	// If minor keywords present, evaluate as minor regardless of major keywords
	if hasMinor {
		return true
	}

	// If major keywords without minor keywords, definitely not minor
	if hasMajor {
		return false
	}

	// No keywords - check modification size
	totalChanges := commit.Insertions + commit.Deletions

	// Very small changes (1-10 lines) are considered minor
	if totalChanges > 0 && totalChanges <= 10 {
		return true
	}

	return false
}

// filterMajorCommits returns only the commits that are considered major (important).
// Minor commits are filtered out for translation files to provide more accurate status.
func filterMajorCommits(commits []*git.Commit) []*git.Commit {
	major := make([]*git.Commit, 0, len(commits))
	for _, commit := range commits {
		if !isMinorCommit(commit) {
			major = append(major, commit)
		}
	}

	return major
}
