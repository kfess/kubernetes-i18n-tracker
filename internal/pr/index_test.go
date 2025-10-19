package pr

import (
	"testing"
)

func TestBuildPRIndex(t *testing.T) {
	prs := []PullRequest{
		{
			Number: 100,
			Title:  "Update file1",
			Files:  []string{"content/ja/docs/file1.md", "content/ja/docs/file2.md"},
		},
		{
			Number: 101,
			Title:  "Update file1 again",
			Files:  []string{"content/ja/docs/file1.md"},
		},
		{
			Number: 102,
			Title:  "Update file3",
			Files:  []string{"content/ja/docs/file3.md"},
		},
	}

	index := buildPRIndex(prs)

	// Test file1 has 2 PRs
	file1PRs := index.GetPRsForFile("content/ja/docs/file1.md")
	if len(file1PRs) != 2 {
		t.Errorf("expected 2 PRs for file1, got %d", len(file1PRs))
	}

	// Test PRs are sorted by number
	if file1PRs[0].Number != 100 || file1PRs[1].Number != 101 {
		t.Errorf("PRs not sorted correctly: got %d, %d", file1PRs[0].Number, file1PRs[1].Number)
	}

	// Test file2 has 1 PR
	file2PRs := index.GetPRsForFile("content/ja/docs/file2.md")
	if len(file2PRs) != 1 {
		t.Errorf("expected 1 PR for file2, got %d", len(file2PRs))
	}

	// Test file3 has 1 PR
	file3PRs := index.GetPRsForFile("content/ja/docs/file3.md")
	if len(file3PRs) != 1 {
		t.Errorf("expected 1 PR for file3, got %d", len(file3PRs))
	}

	// Test nonexistent file
	nonexistentPRs := index.GetPRsForFile("content/ja/docs/nonexistent.md")
	if nonexistentPRs != nil {
		t.Errorf("expected nil for nonexistent file, got %d PRs", len(nonexistentPRs))
	}
}

func TestPRIndex_GetLatestPR(t *testing.T) {
	prs := []PullRequest{
		{Number: 100, Files: []string{"file1.md"}},
		{Number: 101, Files: []string{"file1.md"}},
		{Number: 102, Files: []string{"file1.md"}},
	}

	index := buildPRIndex(prs)

	latest := index.GetLatestPR("file1.md")
	if latest == nil {
		t.Fatal("expected latest PR, got nil")
	}
	if latest.Number != 102 {
		t.Errorf("expected latest PR to be #102, got #%d", latest.Number)
	}

	// Test nonexistent file
	latest = index.GetLatestPR("nonexistent.md")
	if latest != nil {
		t.Errorf("expected nil for nonexistent file, got PR #%d", latest.Number)
	}
}

func TestPRIndex_GetRecentPRs(t *testing.T) {
	prs := []PullRequest{
		{Number: 100, Files: []string{"file1.md"}},
		{Number: 101, Files: []string{"file1.md"}},
		{Number: 102, Files: []string{"file1.md"}},
		{Number: 103, Files: []string{"file1.md"}},
		{Number: 104, Files: []string{"file1.md"}},
	}

	index := buildPRIndex(prs)

	// Get recent 3 PRs
	recent := index.GetRecentPRs("file1.md", 3)
	if len(recent) != 3 {
		t.Fatalf("expected 3 recent PRs, got %d", len(recent))
	}
	if recent[0].Number != 102 || recent[1].Number != 103 || recent[2].Number != 104 {
		t.Errorf("unexpected recent PRs: %d, %d, %d", recent[0].Number, recent[1].Number, recent[2].Number)
	}

	// Get more PRs than available
	recent = index.GetRecentPRs("file1.md", 10)
	if len(recent) != 5 {
		t.Errorf("expected 5 PRs (all available), got %d", len(recent))
	}

	// Test nonexistent file
	recent = index.GetRecentPRs("nonexistent.md", 3)
	if recent != nil {
		t.Errorf("expected nil for nonexistent file, got %d PRs", len(recent))
	}
}

func TestPRIndex_HasPRs(t *testing.T) {
	prs := []PullRequest{
		{Number: 100, Files: []string{"file1.md"}},
	}

	index := buildPRIndex(prs)

	if !index.HasPRs("file1.md") {
		t.Error("expected HasPRs to return true for file1.md")
	}

	if index.HasPRs("nonexistent.md") {
		t.Error("expected HasPRs to return false for nonexistent file")
	}
}

func TestPRIndex_Stats(t *testing.T) {
	prs := []PullRequest{
		{
			Number: 100,
			Files:  []string{"file1.md", "file2.md"},
		},
		{
			Number: 101,
			Files:  []string{"file1.md", "file3.md"},
		},
		{
			Number: 102,
			Files:  []string{"file4.md"},
		},
	}

	index := buildPRIndex(prs)

	// Test total files
	if index.TotalFiles() != 4 {
		t.Errorf("expected 4 total files, got %d", index.TotalFiles())
	}

	// Test total PRs
	if index.TotalPRs() != 3 {
		t.Errorf("expected 3 total PRs, got %d", index.TotalPRs())
	}
}
