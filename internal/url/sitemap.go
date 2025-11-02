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
	client     *http.Client
	baseUrl    string
	maxRetries int
	retryDelay time.Duration
}

// NewClient creates a new Client with the given base URL.
func NewClient(baseUrl string) *Client {
	return &Client{
		client:     &http.Client{Timeout: 10 * time.Second},
		baseUrl:    baseUrl,
		maxRetries: 3,               // default 3 retries
		retryDelay: 1 * time.Second, // default 1 second delay
	}
}

// WithRetry sets the retry configuration
func (c *Client) WithRetry(maxRetries int, retryDelay time.Duration) *Client {
	c.maxRetries = maxRetries
	c.retryDelay = retryDelay
	return c
}

func (c *Client) fetchSitemap(ctx context.Context, lang string) ([]string, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := c.retryDelay * time.Duration(1<<(attempt-1))
			logger.Infof("Retrying sitemap fetch for %s (attempt %d/%d) after %v",
				lang, attempt, c.maxRetries, waitTime)

			select {
			case <-time.After(waitTime):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		urls, statusCode, err := c.fetchSitemapOnce(ctx, lang)

		if err == nil {
			if attempt > 0 {
				logger.Infof("Successfully fetched sitemap for %s after %d retries", lang, attempt)
			}
			return urls, nil
		}

		lastErr = err

		if !isRetryableError(statusCode) {
			logger.Warnf("Non-retryable error for %s: %v", lang, err)
			return nil, err
		}

		logger.Warnf("Retryable error for %s (attempt %d/%d): %v",
			lang, attempt+1, c.maxRetries+1, err)
	}

	return nil, fmt.Errorf("failed after %d retries: %w", c.maxRetries, lastErr)
}

func (c *Client) fetchSitemapOnce(ctx context.Context, lang string) ([]string, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+"/"+lang+"/sitemap.xml", nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("failed to fetch sitemap: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	var urlset urlset
	if err := xml.Unmarshal(body, &urlset); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to decode sitemap XML: %w", err)
	}

	urls := make([]string, len(urlset.URLs))
	for i, u := range urlset.URLs {
		urls[i] = u.Loc
	}

	return urls, resp.StatusCode, nil
}

// FetchAllSitemaps fetches sitemaps for all supported languages concurrently
// and returns a combined list of URLs.
func (c *Client) FetchAllSitemaps(ctx context.Context, langs []language.Language) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	var mu sync.Mutex
	var allUrls = []string{}

	for _, lang := range langs {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			urls, err := c.fetchSitemap(ctx, string(lang))

			if err != nil {
				logger.Warnf("Failed to fetch sitemap for language %s after retries: %v", lang, err)
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

// isRetryableError determines if an error should trigger a retry
func isRetryableError(statusCode int) bool {
	if statusCode == 0 {
		return true
	}

	// Retryable HTTP status codes
	switch statusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	}

	// 4xx errors other than 429 are not retryable
	return false
}
