package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
)

func main() {
	repoPath := flag.String("repo", "./k8s-repo/website", "Path to the git repository")
	outputFile := flag.String("output", "./data/master/git_history.jsonl", "Output JSONL file path")
	workers := flag.Int("workers", 8, "Number of parallel workers")
	skipUpdate := flag.Bool("skip-update", false, "Skip git pull")
	flag.Parse()

	// Ensure output directory exists
	outputDir := filepath.Dir(*outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create output directory: %v", err))
		os.Exit(1)
	}

	// Validate repository path
	if _, err := os.Stat(*repoPath); os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Repository not found: %s", *repoPath))
		os.Exit(1)
	}

	logger.Info("Starting git history extraction...")

	// Create fetcher
	fetcher := git.NewFetcher(git.FetchOptions{
		RepoPath: *repoPath,
		Workers:  *workers, // Note: Workers not used in new approach
	})

	// Update repository
	if !*skipUpdate {
		logger.Info("Updating the repository to the latest version...")
		if err := fetcher.UpdateRepo(); err != nil {
			logger.Warn(fmt.Sprintf("Failed to update repository: %v", err))
		}
	}

	// Fetch history using git log --first-parent approach
	ctx := context.Background()
	events, err := fetcher.FetchHistory(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to fetch history: %v", err))
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Fetched %d events", len(events)))

	// Write to JSONL file
	logger.Info(fmt.Sprintf("Writing to %s...", *outputFile))
	writer, err := git.NewWriter(*outputFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create writer: %v", err))
		os.Exit(1)
	}
	defer writer.Close()

	if err := writer.WriteAll(events); err != nil {
		logger.Error(fmt.Sprintf("Failed to write events: %v", err))
		os.Exit(1)
	}

	if err := writer.Flush(); err != nil {
		logger.Error(fmt.Sprintf("Failed to flush writer: %v", err))
		os.Exit(1)
	}

	logger.Info("Git history extraction completed successfully.")
}
