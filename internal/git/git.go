package git

import "time"

// Event represents a single file change event from git history (one line in JSONL).
type Event struct {
	Hash    string   `json:"hash"`
	Author  string   `json:"author"`
	Date    string   `json:"date"`
	Message string   `json:"message"`
	File    FileInfo `json:"file"`
}

// FileInfo contains the file change details within an event.
type FileInfo struct {
	Path       string `json:"path"`
	Insertions *int   `json:"insertions"` // nil for rename-only events
	Deletions  *int   `json:"deletions"`  // nil for rename-only events
	OldPath    string `json:"old_path,omitempty"`
}

// Commit represents a processed commit for a specific file.
type Commit struct {
	Hash        string
	Author      string
	Date        time.Time
	Message     string
	Path        string
	Insertions  int
	Deletions   int
	RenamedFrom string // Empty if not a rename
}

// Stats contains statistics about a file's history.
type Stats struct {
	TotalCommits    int
	TotalInsertions int
	TotalDeletions  int
	RenameCount     int
	FirstCommit     time.Time
	LatestCommit    time.Time
}
