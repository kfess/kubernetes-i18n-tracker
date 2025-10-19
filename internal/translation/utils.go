package translation

import "fmt"

// toGitHubURL generates GitHub repository URL for a file path.
func toGitHubURL(path string) string {
	return fmt.Sprintf("https://github.com/kubernetes/website/blob/main/%s", path)
}
