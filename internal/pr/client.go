package pr

import (
	"context"

	"github.com/google/go-github/v75/github"
)

// Client is a GitHub client for fetching pull requests and their files
type Client struct {
	gh    *github.Client
	owner string
	repo  string
}

// NewClient creates a new GitHub client for the given repository
func NewClient(token, owner, repo string) *Client {
	return &Client{
		gh:    github.NewClient(nil).WithAuthToken(token),
		owner: owner,
		repo:  repo,
	}
}

// Fetches the list of pull requests from the repository
func (c *Client) FetchPRList(ctx context.Context) ([]*github.PullRequest, error) {
	var allPrs []*github.PullRequest

	opts := &github.PullRequestListOptions{
		State:     "open",
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		prs, resp, err := c.gh.PullRequests.List(ctx, c.owner, c.repo, opts)
		if err != nil {
			return nil, err
		}

		allPrs = append(allPrs, prs...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allPrs, nil
}

// FetchPRDetails fetches both the commit count and files for a PR
func (c *Client) FetchPRDetails(ctx context.Context, prNumber int) (commits int, files []string, err error) {
	pr, _, err := c.gh.PullRequests.Get(ctx, c.owner, c.repo, prNumber)
	if err != nil {
		return 0, nil, err
	}

	commits = 0
	if pr.Commits != nil {
		commits = *pr.Commits
	}

	var fileNames []string
	opts := &github.ListOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		fileList, resp, err := c.gh.PullRequests.ListFiles(ctx, c.owner, c.repo, prNumber, opts)
		if err != nil {
			return 0, nil, err
		}

		for _, file := range fileList {
			fileNames = append(fileNames, file.GetFilename())
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return commits, fileNames, nil
}
