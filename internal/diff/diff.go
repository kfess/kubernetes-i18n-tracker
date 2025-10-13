package diff

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

type Diff struct {
	// Commit
	OldCommitHash string `json:"old_commit_hash"`
	NewCommitHash string `json:"new_commit_hash"`

	// File
	FilePath        string `json:"file_path"`
	EnglishFilePath string `json:"english_file_path"`
	Language        string `json:"language"`

	// Diff Content
	Content string `json:"content"`
}

func CalculateDiff(ctx context.Context, repoPath, oldCommitHash, newCommitHash, filePath string) (Diff, error) {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return Diff{}, fmt.Errorf("failed to resolve repository path: %w", err)
	}

	language, err := extractLanguageFromFilePath(filePath)
	if err != nil {
		return Diff{}, fmt.Errorf("failed to extract language from file path: %w", err)
	}

	// Validate commit hashes (basic check)
	if err := validateCommitHash(oldCommitHash); err != nil {
		return Diff{}, fmt.Errorf("invalid old commit hash: %w", err)
	}
	if err := validateCommitHash(newCommitHash); err != nil {
		return Diff{}, fmt.Errorf("invalid new commit hash: %w", err)
	}

	cmd := exec.CommandContext(ctx, "git", "diff", oldCommitHash, newCommitHash, "--", filePath)
	cmd.Dir = absRepoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return Diff{}, fmt.Errorf("git diff failed: %w, output: %s", err, string(output))
	}

	return Diff{
		OldCommitHash:   oldCommitHash,
		NewCommitHash:   newCommitHash,
		FilePath:        filePath,
		EnglishFilePath: generateEnglishFilePath(filePath, language),
		Content:         string(output),
		Language:        language,
	}, nil
}

func validateCommitHash(hash string) error {
	// Git commit hash is 40 characters (SHA-1) or 7+ for short hash
	hash = strings.TrimSpace(hash)
	if len(hash) < 7 || len(hash) > 40 {
		return fmt.Errorf("invalid hash length: %d", len(hash))
	}
	// Check if it contains only hexadecimal characters
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("hash contains invalid character: %c", c)
		}
	}
	return nil
}

func extractLanguageFromFilePath(filePath string) (string, error) {
	parts := strings.Split(filePath, "/")
	for i, part := range parts {
		if part == "content" && i+1 < len(parts) {
			if slices.Contains(language.SupportedLanguages, parts[i+1]) {
				return parts[i+1], nil
			}
		}
	}

	return "", fmt.Errorf("language not found in file path: %s", filePath)
}

func generateEnglishFilePath(filePath, language string) string {
	parts := strings.Split(filePath, "/")
	for i, part := range parts {
		if part == "content" && i+1 < len(parts) && parts[i+1] == language {
			parts[i+1] = "en"
			break
		}
	}
	return strings.Join(parts, "/")
}
