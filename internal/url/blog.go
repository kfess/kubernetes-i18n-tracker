package url

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",                // Date only
		"2006-01-02 15:04:05 -0700", // Date with time and timezone
		"2006-01-02 15:04:05",       // Date with time
		time.RFC3339,                // ISO 8601 format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// buildBlogUrl builds a blog URL with Hugo priority: slug+date > url > filename > title+date
func generateBlogUrl(baseURL string, cp contentPath, fm *FrontMatter, existingUrls map[string]bool) string {
	langPrefix := ""
	if cp.language != "en" {
		langPrefix = cp.language + "/"
	}

	// Priority 1: slug + date
	if fm.HasSlug() && fm.HasDate() {
		if url := trySlugWithDate(fm, langPrefix, baseURL, existingUrls); url != "" {
			return url
		}
	}

	// Priority 2: Explicit URL
	if fm.HasExplicitUrl() {
		if url := tryExplicitURL(fm, langPrefix, baseURL, existingUrls); url != "" {
			return url
		}
	}

	// Priority 3: filename
	if url := tryFilename(cp.raw, langPrefix, baseURL, existingUrls); url != "" {
		return url
	}

	// Priority 4: title + date
	if fm.HasTitle() && fm.HasDate() {
		if url := tryTitleWithDate(fm, langPrefix, baseURL, existingUrls); url != "" {
			return url
		}
	}

	// Priority 5: blog category (fallback)
	if url := tryFilenameAsIs(cp.raw, langPrefix, baseURL, existingUrls); url != "" {
		return url
	}

	return ""
}

func trySlugWithDate(fm *FrontMatter, langPrefix string, baseURL string, existingUrls map[string]bool) string {
	t, err := parseDate(fm.Date)
	if err != nil {
		return ""
	}

	year := fmt.Sprintf("%04d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	day := fmt.Sprintf("%02d", t.Day())

	var candidate string
	if day == "00" {
		candidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/", baseURL, langPrefix, year, month, fm.Slug)
	} else {
		candidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/%s/", baseURL, langPrefix, year, month, day, fm.Slug)
	}

	return matchUrl(candidate, existingUrls)
}

func tryExplicitURL(fm *FrontMatter, langPrefix string, baseURL string, existingUrls map[string]bool) string {
	urlPath := strings.Trim(fm.Url, "/")
	candidate := fmt.Sprintf("%s/%s%s/", baseURL, langPrefix, urlPath)

	return matchUrl(candidate, existingUrls)
}

func tryFilename(path string, langPrefix string, baseURL string, existingUrls map[string]bool) string {
	filename := filepath.Base(path)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	// Match pattern: YYYY-MM-DD-title
	datePattern := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})-(.+)`)
	matches := datePattern.FindStringSubmatch(filename)

	if matches == nil || len(matches) != 5 {
		return ""
	}

	year, month, day, slug := matches[1], matches[2], matches[3], matches[4]

	var candidate string
	if day == "00" {
		candidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/", baseURL, langPrefix, year, month, slug)
	} else {
		candidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/%s/", baseURL, langPrefix, year, month, day, slug)
	}

	return matchUrl(candidate, existingUrls)
}

func tryTitleWithDate(fm *FrontMatter, langPrefix string, baseURL string, existingUrls map[string]bool) string {
	title := textToSlug(fm.Title)

	t, err := parseDate(fm.Date)
	if err != nil {
		return ""
	}

	year := fmt.Sprintf("%04d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	day := fmt.Sprintf("%02d", t.Day())

	var candidate string
	if day == "00" {
		candidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/", baseURL, langPrefix, year, month, title)
	} else {
		candidate = fmt.Sprintf("%s/%sblog/%s/%s/%s/%s/", baseURL, langPrefix, year, month, day, title)
	}

	return matchUrl(candidate, existingUrls)
}

func tryFilenameAsIs(path string, langPrefix string, baseURL string, existingUrls map[string]bool) string {
	filename := filepath.Base(path)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	candidate := fmt.Sprintf("%s/%sblog/%s/", baseURL, langPrefix, filename)

	return matchUrl(candidate, existingUrls)
}

// textToSlug converts text to URL-friendly slug by removing symbols & replacing spaces with -
func textToSlug(text string) string {
	// Remove apostrophes before processing
	text = strings.ReplaceAll(text, "’", "")

	// Also remove straight single quotes
	text = strings.ReplaceAll(text, "'", "")

	reSymbolsAndSpaces := regexp.MustCompile(`[^\w./\-]+`)
	hyphenated := reSymbolsAndSpaces.ReplaceAllString(text, "-")

	hyphenated = strings.ToLower(hyphenated)

	reHyphens := regexp.MustCompile(`-+`)
	normalized := reHyphens.ReplaceAllString(hyphenated, "-")

	normalized = strings.Trim(normalized, "-")

	return normalized
}
