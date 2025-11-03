package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"github.com/kfess/kubernetes-i18n-tracker/internal/workflow"
)

func main() {
	logger.Init()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warnf("No .env file found: %v", err)
	}

	// Create workflow config
	config := workflow.DefaultConfig()
	config.GitHubToken = os.Getenv("KUBERNETES_WEBSITE_READ_GITHUB_TOKEN")

	if config.GitHubToken == "" {
		logger.Warn("KUBERNETES_WEBSITE_READ_GITHUB_TOKEN not set, GitHub data will be skipped")
	}

	// Create and run workflow
	ctx := context.Background()
	wf := workflow.New(config)

	if err := wf.Run(ctx); err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}
}
