package main

import (
	// "context"
	// "fmt"
	// "os"

	"context"

	"github.com/joho/godotenv"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"github.com/kfess/kubernetes-i18n-tracker/internal/url"
)

func main() {
	logger.Init()
	err := godotenv.Load()
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
		return
	}

	// token := os.Getenv("KUBERNETES_WEBSITE_READ_GITHUB_TOKEN")
	// client := pr.NewClient(token, "kubernetes", "website")
	// fetcher := pr.NewFetcher(client)
	// prs, err := fetcher.FetchAll(context.Background())
	// if err != nil {
	// 	logger.Errorf("Error fetching PRs: %v", err)
	// 	return
	// }
	// for _, pr := range prs {
	// 	logger.Infof("PR #%d: %s", pr.Number, pr.Title)

	// 	// fmt.Printf("PR #%d: %s\n", pr.Number, pr.Title)
	// 	// fmt.Println("Files changed in this PR:")
	// 	for _, file := range pr.Files {
	// 		fmt.Printf("- %s\n", file)
	// 	}
	// }

	client := url.NewClient("https://kubernetes.io")
	urls, err := client.FetchAllSitemaps(context.Background())
	if err != nil {
		logger.Errorf("Error fetching sitemaps: %v", err)
		return
	}
	logger.Infof("Fetched URLs for English: %v", urls["en"])
}
