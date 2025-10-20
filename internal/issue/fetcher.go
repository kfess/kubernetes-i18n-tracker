package issue

import (
	"context"
	"fmt"
	"sync"

	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"golang.org/x/sync/errgroup"
)

const activeFetchers = 10

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

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(activeFetchers)

	var mu sync.Mutex
	issues := make([]Issue, 0, len(rawIssues))
	var processed, skipped int

	for _, rawIssue := range rawIssues {
		rawIssue := rawIssue

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if rawIssue.Number == nil || rawIssue.Title == nil || rawIssue.HTMLURL == nil {
				mu.Lock()
				skipped++
				mu.Unlock()
				logger.Warnf("Skipping issue with missing data: %+v", rawIssue)
				return nil
			}

			issue := Issue{
				Number: *rawIssue.Number,
				Title:  *rawIssue.Title,
				URL:    *rawIssue.HTMLURL,
				Labels: extractLabels(rawIssue),
			}

			mu.Lock()
			issues = append(issues, issue)
			processed++
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logger.Errorf("Error processing issues: %v", err)
		return nil, err
	}

	logger.Infof("Total issues processed: %d, skipped: %d", processed, skipped)

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
