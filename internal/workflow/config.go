package workflow

// Config holds all configuration needed for the workflow.
type Config struct {
	// Repository owner
	RepoOwner string

	// Repository name
	RepoName string

	// Local path to the repository
	RepoPath string

	// Path to git history file
	GitHistoryFile string

	// Path to all files CSV
	AllFilesPath string

	// Path to page view CSV
	PageViewFile string

	// Output directory
	OutputDir string

	// GitHub token
	GitHubToken string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		RepoOwner:      "kubernetes",
		RepoName:       "website",
		RepoPath:       "./k8s-repo/website",
		GitHistoryFile: "./data/master/git_history.jsonl",
		AllFilesPath:   "./data/master/all_files.csv",
		PageViewFile:   "./data/master/page_view.csv",
		OutputDir:      "./data/output",
	}
}
