package url

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// BlogUrlBuilder implements UrlBuilder for blog posts
type BlogUrlBuilder struct {
	baseUrl      string
	parser       FrontMatterParser
	existingUrls map[string]struct{} // map for quick lookup
}

// NewBlogUrlBuilder creates a new instance of BlogUrlBuilder
func NewBlogUrlBuilder(baseUrl string, parser FrontMatterParser, existingUrls []string) *BlogUrlBuilder {
	urlMap := make(map[string]struct{}, len(existingUrls))
	for _, url := range existingUrls {
		urlMap[url] = struct{}{}
	}

	return &BlogUrlBuilder{
		baseUrl:      baseUrl,
		parser:       parser,
		existingUrls: urlMap,
	}
}

// Build constructs the URL based on front matter and existing URLs
func (b *BlogUrlBuilder) Build(ctx context.Context, path string, lang string) (string, error) {
	fm, err := b.parser.Parse(path)
	if err != nil {
		return "", err
	}

	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	// Priority 1: slug + date
	if fm.HasSlug() && fm.HasDate() {
		if url := b.trySlugWithDate(fm, langPrefix); url != "" {
			return url, nil
		}
	}

	// Priority 2: Explicit URL
	if fm.HasExplicitUrl() {
		if url := b.tryExplicitUrl(fm, langPrefix); url != "" {
			return url, nil
		}
	}

	// Priority 3: filename
	if url := b.tryFilename(path, langPrefix); url != "" {
		return url, nil
	}

	// Priority 4: title + date
	if fm.HasTitle() && fm.HasDate() {
		if url := b.tryTitleWithDate(fm, langPrefix); url != "" {
			return url, nil
		}
	}

	// Priority 5: filename as-is
	if url := b.tryFilenameAsIs(path, langPrefix); url != "" {
		return url, nil
	}

	// If no URL could be constructed
	return "", fmt.Errorf("no valid URL found for path: %s", path)
}

// resolveUrl checks if the candidate URL exists in the registry (as-is or lowercase)
func (b *BlogUrlBuilder) resolveUrl(candidateURL string) (string, bool) {
	if _, exists := b.existingUrls[candidateURL]; exists {
		return candidateURL, true
	}

	loweredCandidateUrl := strings.ToLower(candidateURL)
	if _, exists := b.existingUrls[loweredCandidateUrl]; exists {
		return loweredCandidateUrl, true
	}

	return "", false
}

func (b *BlogUrlBuilder) trySlugWithDate(fm *FrontMatter, langPrefix string) string {
	date := fm.Date // Expecting format "YYYY-MM-DD"
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}
	year := fmt.Sprintf("%04d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	day := fmt.Sprintf("%02d", t.Day())

	var urlCandidate string

	if day == "00" {
		urlCandidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/", b.baseUrl, langPrefix, year, month, fm.Slug)
	} else {
		urlCandidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/%s/", b.baseUrl, langPrefix, year, month, day, fm.Slug)
	}

	// Check if this URL exists in the registry
	if validUrl, ok := b.resolveUrl(urlCandidate); ok {
		return validUrl
	}

	return ""
}

func (b *BlogUrlBuilder) tryExplicitUrl(fm *FrontMatter, langPrefix string) string {
	urlPath := strings.Trim(fm.Url, "/")
	urlCandidate := fmt.Sprintf("%s/%s%s/", b.baseUrl, langPrefix, urlPath)

	// Check if this URL exists in the registry
	if validUrl, ok := b.resolveUrl(urlCandidate); ok {
		return validUrl
	}

	return ""
}

func (b *BlogUrlBuilder) tryFilename(path string, langPrefix string) string {
	filename := filepath.Base(path)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	datePattern := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})-(.+)`)
	matches := datePattern.FindStringSubmatch(filename)

	if matches == nil || len(matches) != 5 {
		return ""
	}

	year, month, day, slug := matches[1], matches[2], matches[3], matches[4]

	var urlCandidate string

	if day == "00" {
		urlCandidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/", b.baseUrl, langPrefix, year, month, slug)
	} else {
		urlCandidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/%s/", b.baseUrl, langPrefix, year, month, day, slug)
	}

	// Check if this URL exists in the registry
	if validUrl, ok := b.resolveUrl(urlCandidate); ok {
		return validUrl
	}

	return ""
}

// textToSlug converts text to URL-friendly slug by removing symbols & replacing spaces with -.
func (b *BlogUrlBuilder) textToSlug(text string) string {
	// Remove all characters except alphanumeric, spaces, dots, slashes, and hyphens
	reClean := regexp.MustCompile(`[^a-zA-Z0-9\s\./\-]`)
	cleaned := reClean.ReplaceAllString(text, "")

	// Replace spaces with hyphens and convert to lowercase
	reSpaces := regexp.MustCompile(`\s+`)
	hyphenated := reSpaces.ReplaceAllString(cleaned, "-")
	hyphenated = strings.ToLower(hyphenated)

	// Normalize multiple consecutive hyphens to single hyphen
	reHyphens := regexp.MustCompile(`-+`)
	normalized := reHyphens.ReplaceAllString(hyphenated, "-")

	// Trim leading and trailing hyphens
	normalized = strings.Trim(normalized, "-")

	return normalized
}

func (b *BlogUrlBuilder) tryTitleWithDate(fm *FrontMatter, langPrefix string) string {
	title := b.textToSlug(fm.Title) // Convert title to URL-friendly slug
	date := fm.Date                 // Expecting format "YYYY-MM-DD"
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}
	year := fmt.Sprintf("%04d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	day := fmt.Sprintf("%02d", t.Day())

	var urlCandidate string

	if day == "00" {
		urlCandidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/", b.baseUrl, langPrefix, year, month, title)
	} else {
		urlCandidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/%s/", b.baseUrl, langPrefix, year, month, day, title)
	}

	// Check if this URL exists in the registry
	if validUrl, ok := b.resolveUrl(urlCandidate); ok {
		return validUrl
	}

	return ""
}

func (b *BlogUrlBuilder) tryFilenameAsIs(path string, langPrefix string) string {
	filename := filepath.Base(path)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	urlCandidate := fmt.Sprintf("%s/%sblog/%s/", b.baseUrl, langPrefix, filename)

	// Check if this URL exists in the registry
	if validUrl, ok := b.resolveUrl(urlCandidate); ok {
		return validUrl
	}

	return ""
}
