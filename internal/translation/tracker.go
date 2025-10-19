package translation

import (
	// "context"

	"github.com/kfess/kubernetes-i18n-tracker/internal/history"
	"github.com/kfess/kubernetes-i18n-tracker/internal/pr"
	"github.com/kfess/kubernetes-i18n-tracker/internal/url"
)

type Tracker struct {
	history      *history.History
	urlConv      *url.Converter
	prIndex      *pr.Index
	repoPath     string
	existingUrls map[string]bool

	// add issue index later
}

type Config struct {
	RepoPath     string
	existingUrls []string
}

func NewTracker(history *history.History, urlConverter *url.Converter, prIndex *pr.Index, config Config) *Tracker {
	pathMap := make(map[string]bool, len(config.existingUrls))
	for _, url := range config.existingUrls {
		pathMap[url] = true
	}

	return &Tracker{
		history:      history,
		urlConv:      urlConverter,
		prIndex:      prIndex,
		repoPath:     config.RepoPath,
		existingUrls: pathMap,
	}
}

// func (t *Tracker) GetStatus(ctx context.Context, path string) (*Status, error) {}
