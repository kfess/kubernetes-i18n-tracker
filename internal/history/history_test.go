package history

import (
	"testing"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
)

func TestBuild_SimpleHistory(t *testing.T) {
	events := []*git.Event{
		{
			Hash:    "abc123",
			Author:  "Alice",
			Date:    "2025-01-01 10:00:00 +0000",
			Message: "Add file",
			File: git.FileInfo{
				Path:       "content/en/docs/guide.md",
				Insertions: intPtr(10),
				Deletions:  intPtr(0),
			},
		},
		{
			Hash:    "def456",
			Author:  "Bob",
			Date:    "2025-02-01 11:00:00 +0000",
			Message: "Update file",
			File: git.FileInfo{
				Path:       "content/en/docs/guide.md",
				Insertions: intPtr(5),
				Deletions:  intPtr(2),
			},
		},
	}

	history := Build(events)

	commits := history.GetCommits("content/en/docs/guide.md")
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}

	// Check chronological order (oldest first)
	if commits[0].Hash != "abc123" {
		t.Errorf("expected first commit to be abc123, got %s", commits[0].Hash)
	}
	if commits[1].Hash != "def456" {
		t.Errorf("expected second commit to be def456, got %s", commits[1].Hash)
	}
}

func TestBuild_WithRename(t *testing.T) {
	events := []*git.Event{
		{
			Hash:    "commit1",
			Author:  "Alice",
			Date:    "2025-01-01 10:00:00 +0000",
			Message: "Add file",
			File: git.FileInfo{
				Path:       "old_name.md",
				Insertions: intPtr(10),
				Deletions:  intPtr(0),
			},
		},
		{
			Hash:    "commit2",
			Author:  "Bob",
			Date:    "2025-02-01 11:00:00 +0000",
			Message: "Rename file",
			File: git.FileInfo{
				Path:       "new_name.md",
				Insertions: nil, // Rename events often have nil insertions/deletions
				Deletions:  nil,
				OldPath:    "old_name.md",
			},
		},
		{
			Hash:    "commit3",
			Author:  "Charlie",
			Date:    "2025-03-01 12:00:00 +0000",
			Message: "Update file",
			File: git.FileInfo{
				Path:       "new_name.md",
				Insertions: intPtr(3),
				Deletions:  intPtr(1),
			},
		},
	}

	history := Build(events)

	// All commits should be under the final path
	commits := history.GetCommits("new_name.md")
	if len(commits) != 3 {
		t.Fatalf("expected 3 commits under new_name.md, got %d", len(commits))
	}

	// Old path should have no commits
	oldCommits := history.GetCommits("old_name.md")
	if len(oldCommits) != 0 {
		t.Errorf("expected 0 commits under old_name.md, got %d", len(oldCommits))
	}

	// Check rename info
	if commits[1].RenamedFrom != "old_name.md" {
		t.Errorf("expected RenamedFrom to be old_name.md, got %s", commits[1].RenamedFrom)
	}
}

func TestBuild_ChainedRenames(t *testing.T) {
	events := []*git.Event{
		{
			Hash:    "commit1",
			Author:  "Alice",
			Date:    "2025-01-01 10:00:00 +0000",
			Message: "Add file A",
			File: git.FileInfo{
				Path:       "file_a.md",
				Insertions: intPtr(10),
				Deletions:  intPtr(0),
			},
		},
		{
			Hash:    "commit2",
			Author:  "Bob",
			Date:    "2025-02-01 11:00:00 +0000",
			Message: "Rename A to B",
			File: git.FileInfo{
				Path:    "file_b.md",
				OldPath: "file_a.md",
			},
		},
		{
			Hash:    "commit3",
			Author:  "Charlie",
			Date:    "2025-03-01 12:00:00 +0000",
			Message: "Rename B to C",
			File: git.FileInfo{
				Path:    "file_c.md",
				OldPath: "file_b.md",
			},
		},
		{
			Hash:    "commit4",
			Author:  "Dave",
			Date:    "2025-04-01 13:00:00 +0000",
			Message: "Update C",
			File: git.FileInfo{
				Path:       "file_c.md",
				Insertions: intPtr(5),
				Deletions:  intPtr(2),
			},
		},
	}

	history := Build(events)

	// All commits should be under final path file_c.md
	commits := history.GetCommits("file_c.md")
	if len(commits) != 4 {
		t.Fatalf("expected 4 commits under file_c.md, got %d", len(commits))
	}

	// Verify chronological order
	expectedHashes := []string{"commit1", "commit2", "commit3", "commit4"}
	for i, expected := range expectedHashes {
		if commits[i].Hash != expected {
			t.Errorf("commit %d: expected %s, got %s", i, expected, commits[i].Hash)
		}
	}

	// Historical paths should be: file_c.md -> file_b.md -> file_a.md
	paths := history.GetHistoricalPaths("file_c.md")
	expectedPaths := []string{"file_c.md", "file_b.md", "file_a.md"}
	if len(paths) != len(expectedPaths) {
		t.Fatalf("expected %d paths, got %d", len(expectedPaths), len(paths))
	}
	for i, expected := range expectedPaths {
		if paths[i] != expected {
			t.Errorf("path %d: expected %s, got %s", i, expected, paths[i])
		}
	}
}

func TestBuild_MultipleFiles(t *testing.T) {
	events := []*git.Event{
		{
			Hash:    "commit1",
			Author:  "Alice",
			Date:    "2025-01-01 10:00:00 +0000",
			Message: "Add file1",
			File: git.FileInfo{
				Path:       "file1.md",
				Insertions: intPtr(10),
				Deletions:  intPtr(0),
			},
		},
		{
			Hash:    "commit2",
			Author:  "Bob",
			Date:    "2025-01-02 10:00:00 +0000",
			Message: "Add file2",
			File: git.FileInfo{
				Path:       "file2.md",
				Insertions: intPtr(20),
				Deletions:  intPtr(0),
			},
		},
		{
			Hash:    "commit3",
			Author:  "Charlie",
			Date:    "2025-01-03 10:00:00 +0000",
			Message: "Update file1",
			File: git.FileInfo{
				Path:       "file1.md",
				Insertions: intPtr(5),
				Deletions:  intPtr(1),
			},
		},
	}

	history := Build(events)

	// Check that we have 2 distinct files
	paths := history.AllPaths()
	if len(paths) != 2 {
		t.Fatalf("expected 2 files, got %d", len(paths))
	}

	// Check file1 has 2 commits
	file1Commits := history.GetCommits("file1.md")
	if len(file1Commits) != 2 {
		t.Errorf("expected 2 commits for file1.md, got %d", len(file1Commits))
	}

	// Check file2 has 1 commit
	file2Commits := history.GetCommits("file2.md")
	if len(file2Commits) != 1 {
		t.Errorf("expected 1 commit for file2.md, got %d", len(file2Commits))
	}
}

func TestStats(t *testing.T) {
	events := []*git.Event{
		{
			Hash:    "commit1",
			Author:  "Alice",
			Date:    "2025-01-01 10:00:00 +0000",
			Message: "Add file",
			File: git.FileInfo{
				Path:       "test.md",
				Insertions: intPtr(10),
				Deletions:  intPtr(0),
			},
		},
		{
			Hash:    "commit2",
			Author:  "Bob",
			Date:    "2025-02-01 11:00:00 +0000",
			Message: "Update file",
			File: git.FileInfo{
				Path:       "test.md",
				Insertions: intPtr(5),
				Deletions:  intPtr(2),
			},
		},
		{
			Hash:    "commit3",
			Author:  "Charlie",
			Date:    "2025-03-01 12:00:00 +0000",
			Message: "Rename file",
			File: git.FileInfo{
				Path:       "test_renamed.md",
				Insertions: intPtr(0),
				Deletions:  intPtr(0),
				OldPath:    "test.md",
			},
		},
	}

	history := Build(events)
	stats := history.Stats("test_renamed.md")

	if stats == nil {
		t.Fatal("expected stats, got nil")
	}

	if stats.TotalCommits != 3 {
		t.Errorf("expected 3 commits, got %d", stats.TotalCommits)
	}

	if stats.TotalInsertions != 15 {
		t.Errorf("expected 15 insertions, got %d", stats.TotalInsertions)
	}

	if stats.TotalDeletions != 2 {
		t.Errorf("expected 2 deletions, got %d", stats.TotalDeletions)
	}

	if stats.RenameCount != 1 {
		t.Errorf("expected 1 rename, got %d", stats.RenameCount)
	}

	expectedFirst, _ := time.Parse("2006-01-02 15:04:05 -0700", "2025-01-01 10:00:00 +0000")
	expectedLatest, _ := time.Parse("2006-01-02 15:04:05 -0700", "2025-03-01 12:00:00 +0000")

	if !stats.FirstCommit.Equal(expectedFirst) {
		t.Errorf("expected first commit at %v, got %v", expectedFirst, stats.FirstCommit)
	}

	if !stats.LatestCommit.Equal(expectedLatest) {
		t.Errorf("expected latest commit at %v, got %v", expectedLatest, stats.LatestCommit)
	}
}

func TestStats_NonexistentFile(t *testing.T) {
	events := []*git.Event{
		{
			Hash:    "commit1",
			Author:  "Alice",
			Date:    "2025-01-01 10:00:00 +0000",
			Message: "Add file",
			File: git.FileInfo{
				Path:       "exists.md",
				Insertions: intPtr(10),
				Deletions:  intPtr(0),
			},
		},
	}

	history := Build(events)
	stats := history.Stats("nonexistent.md")

	if stats != nil {
		t.Error("expected nil stats for nonexistent file")
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}
