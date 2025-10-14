package url

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/logger"
	"golang.org/x/sync/errgroup"
)

type urlset struct {
	XMLName xml.Name  `xml:"urlset"`
	URLs    []urlInfo `xml:"url"`
}

type urlInfo struct {
	Loc string `xml:"loc"`
}

// Client is a client for fetching URLs from the Kubernetes website.
type Client struct {
	client  *http.Client
	baseUrl string
}

// NewClient creates a new Client with the given base URL.
func NewClient(baseUrl string) *Client {
	return &Client{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseUrl: baseUrl,
	}
}

func (c *Client) fetchSitemap(ctx context.Context, lang string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+"/"+lang+"/sitemap.xml", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch sitemap: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var urlset urlset
	if err := xml.Unmarshal(body, &urlset); err != nil {
		return nil, fmt.Errorf("failed to decode sitemap XML: %w", err)
	}

	urls := make([]string, len(urlset.URLs))
	for i, u := range urlset.URLs {
		urls[i] = u.Loc
	}

	return urls, nil
}

// FetchAllSitemaps fetches sitemaps for all supported languages concurrently
// and returns a combined list of URLs.
func (c *Client) FetchAllSitemaps(ctx context.Context) ([]string, error) {

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	var mu sync.Mutex
	var allUrls = []string{}

	for _, lang := range language.SupportedLanguages {
		lang := lang

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			urls, err := c.fetchSitemap(ctx, lang)
			if err != nil {
				logger.Warnf("Failed to fetch sitemap for language %s: %v", lang, err)
				return err
			}

			mu.Lock()
			allUrls = append(allUrls, urls...)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return allUrls, nil
}
