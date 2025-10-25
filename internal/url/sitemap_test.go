package url

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

func TestFetchAllSitemaps(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := strings.Split(r.URL.Path, "/")[1]

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url><loc>https://kubernetes.io/` + lang + `/docs/</loc></url>
</urlset>`))
	}))

	defer server.Close()

	client := NewClient(server.URL)

	urls, err := client.FetchAllSitemaps(context.Background(), language.SupportedLanguages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(urls) == 0 {
		t.Fatalf("expected non-empty URLs, got empty")
	}
}

func TestFetchAllSitemaps_PartialFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := strings.Split(r.URL.Path, "/")[1]

		// Simulate failure for one language
		if lang == "ja" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url><loc>https://kubernetes.io/` + lang + `/docs/</loc></url>
</urlset>`))
	}))

	defer server.Close()

	client := NewClient(server.URL)

	urls, err := client.FetchAllSitemaps(context.Background(), language.SupportedLanguages)

	if err == nil {
		t.Error("expected error when one sitemap fails, got nil")
	}

	if len(urls) != 0 {
		t.Errorf("expected empty URLs on error, got %d URLs", len(urls))
	}
}

func TestFetchAllSitemaps_InvalidXML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid xml content`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.FetchAllSitemaps(context.Background(), language.SupportedLanguages)

	if err == nil {
		t.Error("expected error for invalid XML, got nil")
	}
}

func TestFetchAllSitemaps_ContextCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url><loc>https://kubernetes.io/test/</loc></url>
</urlset>`))
	}))
	defer server.Close()

	client := NewClient(server.URL)

	ctx, cancel := context.WithCancel(context.Background())

	// Immediately cancel the context
	cancel()

	_, err := client.FetchAllSitemaps(ctx, language.SupportedLanguages)

	if err == nil {
		t.Error("expected context canceled error, got nil")
	}
}

func TestFetchSitemap_WithRetry(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)

		// return 503 for the first two attempts
		if count <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// return valid sitemap on the third attempt
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url><loc>https://kubernetes.io/ja/docs/</loc></url>
</urlset>`))
	}))

	defer server.Close()

	client := NewClient(server.URL).WithRetry(3, 100*time.Millisecond)

	urls, err := client.FetchAllSitemaps(context.Background(), []string{"ja"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(urls) != 1 {
		t.Errorf("expected 1 URL, got %d", len(urls))
	}

	finalCount := atomic.LoadInt32(&attemptCount)
	if finalCount != 3 {
		t.Errorf("expected 3 attempts, got %d", finalCount)
	}
}

func TestFetchSitemap_ExceedsMaxRetries(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attemptCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithRetry(2, 50*time.Millisecond)

	_, err := client.FetchAllSitemaps(context.Background(), []string{"ja"})
	if err == nil {
		t.Error("expected error after max retries, got nil")
	}

	finalCount := atomic.LoadInt32(&attemptCount)
	if finalCount != 3 {
		t.Errorf("expected 3 attempts, got %d", finalCount)
	}
}

func TestFetchSitemap_NonRetryableError(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attemptCount, 1)
		w.WriteHeader(http.StatusNotFound) // 404 is not retriable
	}))
	defer server.Close()

	client := NewClient(server.URL).WithRetry(3, 50*time.Millisecond)

	_, err := client.FetchAllSitemaps(context.Background(), []string{"ja"})
	if err == nil {
		t.Error("expected error for 404, got nil")
	}

	finalCount := atomic.LoadInt32(&attemptCount)
	// do not retry on 404
	if finalCount != 1 {
		t.Errorf("expected 1 attempt for non-retryable error, got %d", finalCount)
	}
}
