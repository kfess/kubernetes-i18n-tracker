package issue

import (
	"context"
	"fmt"

	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
)

// Fetcher fetches issues from a GitHub repository.
type Fetcher struct {
	client *Client
}

// NewFetcher creates a new Fetcher with the given GitHub client.
func NewFetcher(client *Client) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

// FetchAll fetches all open issues from the repository.
func (f *Fetcher) FetchAll(ctx context.Context) ([]Issue, error) {
	rawIssues, err := f.client.FetchIssueList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}

	logger.Debugf("Total issues fetched: %d", len(rawIssues))

	issues := make([]Issue, 0, len(rawIssues))
	for _, rawIssue := range rawIssues {
		if rawIssue.Number == nil || rawIssue.Title == nil || rawIssue.HTMLURL == nil {
			logger.Warnf("Skipping issue with missing data: %+v", rawIssue)
			continue
		}

		issue := Issue{
			Number: *rawIssue.Number,
			Title:  *rawIssue.Title,
			URL:    *rawIssue.HTMLURL,
			Labels: extractLabels(rawIssue),
		}

		issues = append(issues, issue)
	}

	logger.Debugf("Processed %d issues", len(issues))

	return issues, nil
}

// extractLabels extracts label names from a GitHub issue.
func extractLabels(issue *GitHubIssue) []string {
	labels := make([]string, 0, len(issue.Labels))
	for _, label := range issue.Labels {
		if label.Name != nil {
			labels = append(labels, *label.Name)
		}
	}
	return labels
}
