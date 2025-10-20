package pr

import (
	"context"
	"fmt"
	"sync"

	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"golang.org/x/sync/errgroup"
)

const activeFetchers = 10

// Fetcher fetches and filters pull requests from a GitHub repository.
type Fetcher struct {
	client *Client
}

// NewFetcher creates a new Fetcher with the given GitHub client.
func NewFetcher(client *Client) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

// FetchAll fetches all pull requests, filters them, and returns a list of relevant PRs.
func (f *Fetcher) FetchAll(ctx context.Context) ([]PullRequest, error) {
	rawPRs, err := f.client.FetchPRList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PRs: %w", err)
	}

	logger.Debugf("Total PRs fetched: %d", len(rawPRs))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(activeFetchers)

	var mu sync.Mutex
	prs := make([]PullRequest, 0, len(rawPRs))
	var processed, skipped int

	for _, rawPr := range rawPRs {
		rawPr := rawPr

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if rawPr.Number == nil || rawPr.Title == nil || rawPr.HTMLURL == nil {
				logger.Warnf("Skipping PR with missing data: %+v", rawPr)
				return nil
			}

			commits, files, err := f.client.FetchPRDetails(ctx, *rawPr.Number)
			if err != nil {
				logger.Warnf("Failed to fetch details for PR %d: %v", *rawPr.Number, err)
				return nil
			}

			if shouldSkipPR(rawPr, len(files), commits) {
				mu.Lock()
				skipped++
				mu.Unlock()
				logger.Warnf("Skipping PR #%d: %s", *rawPr.Number, *rawPr.Title)
				return nil
			}

			pr := PullRequest{
				Number:  *rawPr.Number,
				Title:   *rawPr.Title,
				Url:     *rawPr.HTMLURL,
				Commits: commits,
				Files:   files,
				Labels:  extractLabels(rawPr),
			}

			mu.Lock()
			prs = append(prs, pr)
			processed++
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logger.Errorf("Error processing PRs: %v", err)
		return nil, err
	}

	logger.Infof("Total PRs processed: %d, skipped: %d", processed, skipped)

	return prs, nil
}
