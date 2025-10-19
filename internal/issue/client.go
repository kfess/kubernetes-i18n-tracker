package issue

import (
	"context"

	"github.com/google/go-github/v75/github"
)

// GitHubIssue represents a raw GitHub issue from the API.
type GitHubIssue struct {
	Number  *int            `json:"number"`
	Title   *string         `json:"title"`
	HTMLURL *string         `json:"html_url"`
	Labels  []*github.Label `json:"labels"`
}

// Client is a GitHub client for fetching issues.
type Client struct {
	gh    *github.Client
	owner string
	repo  string
}

// NewClient creates a new GitHub client for the given repository.
func NewClient(token string, owner string, repo string) *Client {
	return &Client{
		gh:    github.NewClient(nil).WithAuthToken(token),
		owner: owner,
		repo:  repo,
	}
}

// FetchIssueList fetches all open issues from the repository.
func (c *Client) FetchIssueList(ctx context.Context) ([]*GitHubIssue, error) {
	var allIssues []*GitHubIssue

	opts := &github.IssueListByRepoOptions{
		State:     "open",
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		issues, resp, err := c.gh.Issues.ListByRepo(ctx, c.owner, c.repo, opts)
		if err != nil {
			return nil, err
		}

		for _, issue := range issues {
			// Skip pull requests (GitHub API returns PRs as issues)
			if issue.PullRequestLinks != nil {
				continue
			}

			allIssues = append(allIssues, &GitHubIssue{
				Number:  issue.Number,
				Title:   issue.Title,
				HTMLURL: issue.HTMLURL,
				Labels:  issue.Labels,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	return allIssues, nil
}
