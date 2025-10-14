package pr

import (
	"slices"

	"github.com/google/go-github/v75/github"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
)

const (
	tooManyFilesChanged = 1000
	tooManyCommits      = 100
)

// shouldSkipPR determines if a PR should be filtered out based on labels and change metrics.
func shouldSkipPR(pr *github.PullRequest, fileCount int, commitCount int) bool {
	labels := extractLabels(pr)

	// Skip draft PRs (still being worked on)
	if pr.GetDraft() {
		logger.Warnf("Skipping PR #%d: draft PR", pr.GetNumber())
		return true
	}

	// Skip PRs without CLA signature
	if !hasLabel(labels, "cncf-cla: yes") {
		logger.Warnf("Skipping PR #%d: missing CLA signature", pr.GetNumber())
		return true
	}

	// Skip PRs not related to localization
	if !hasLabel(labels, "area/localization") {
		logger.Warnf("Skipping PR #%d: not related to localization", pr.GetNumber())
		return true
	}

	// Skip PRs related to Arabic localization
	// Arabic localization is still under discussion and has not started yet
	if hasLabel(labels, "lang/ar") {
		logger.Warnf("Skipping PR #%d: Arabic localization is still under discussion", pr.GetNumber())
		return true
	}

	// Skip PRs with excessive file changes (likely bulk updates or automated changes)
	if fileCount > tooManyFilesChanged {
		logger.Warnf("Skipping PR #%d: too many files changed (%d files)", pr.GetNumber(), fileCount)
		return true
	}

	// Skip PRs with excessive commits
	if commitCount > tooManyCommits {
		logger.Warnf("Skipping PR #%d: too many commits (%d commits)", pr.GetNumber(), commitCount)
		return true
	}

	return false
}

// hasLabel checks if a target label exists in the label list.
func hasLabel(labels []string, target string) bool {
	return slices.Contains(labels, target)
}

// extractLabels extracts label names from GitHub PR labels.
func extractLabels(pr *github.PullRequest) []string {
	labels := make([]string, 0, len(pr.Labels))
	for _, label := range pr.Labels {
		labels = append(labels, *label.Name)
	}
	return labels
}
