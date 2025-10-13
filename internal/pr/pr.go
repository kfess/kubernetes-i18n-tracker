package pr

// PullRequest represetnts a GitHub Pull Request.
type PullRequest struct {
	Number  int      `json:"number"`
	Title   string   `json:"title"`
	Url     string   `json:"url"`
	Commits int      `json:"commits"`
	Files   []string `json:"files"`
	Labels  []string `json:"labels"`
}
