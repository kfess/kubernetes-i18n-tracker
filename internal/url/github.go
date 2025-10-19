package url

import (
	"fmt"
	"strings"
)

type GitHubRepo struct {
	Owner  string // e.g., "kubernetes"
	Name   string // e.g., "website"
	Branch string // e.g., "main"
}

// ConvertToGitHubUrl converts a file path to a GitHub URL for the given repository.
func ConvertToGitHubUrl(repo GitHubRepo, filePath string) string {
	cleanPath := strings.TrimPrefix(filePath, "/")

	return fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s",
		repo.Owner, repo.Name, repo.Branch, cleanPath)
}
