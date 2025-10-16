package url

import (
	"fmt"
	"path/filepath"
	"strings"
)

func generateDocsUrl(baseUrl string, cp contentPath, fm *FrontMatter, existingUrls map[string]bool) string {
	langPrefix := ""
	if cp.language != "en" {
		langPrefix = cp.language + "/"
	}

	// Check if this is a special case
	if len(cp.segments) >= 5 {
		subSection := cp.segments[3] // e.g., "reference", "contribute"

		// Special case: docs/reference/glossary
		if subSection == "reference" && len(cp.segments) >= 5 && cp.segments[4] == "glossary" {
			return buildGlossaryUrl(baseUrl, langPrefix, fm)
		}

		// Special case: docs/contribute/blog
		if subSection == "contribute" && len(cp.segments) >= 5 && cp.segments[4] == "blog" {
			return buildContributeBlogUrl(baseUrl, langPrefix, cp, fm, existingUrls)
		}
	}

	// Normal docs URL handling
	// e.g., content/en/docs/concepts/overview.md -> docs/concepts/overview
	pathAfterLang := strings.Join(cp.segments[2:], "/") // Skip "content" and "{lang}"
	pathAfterLang = strings.TrimSuffix(pathAfterLang, ".md")
	pathAfterLang = strings.TrimSuffix(pathAfterLang, ".html")

	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, pathAfterLang)

	if matched := matchUrl(candidate, existingUrls); matched != "" {
		return matched
	}

	return ""
}

func buildGlossaryUrl(baseUrl string, langPrefix string, fm *FrontMatter) string {
	fullLink := fm.FullLink
	if fullLink == "" {
		return ""
	}

	// full_link must start with "/"
	if !strings.HasPrefix(fullLink, "/") {
		return ""
	}

	// Build URL: baseUrl + langPrefix + full_link (without leading /)
	url := fmt.Sprintf("%s/%s%s", baseUrl, langPrefix, strings.TrimPrefix(fullLink, "/"))

	// If full_link contains a fragment (#), don't add trailing slash
	if strings.Contains(fullLink, "#") {
		return strings.TrimSuffix(url, "/")
	}

	// Otherwise, ensure trailing slash
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	return url
}

func buildContributeBlogUrl(baseUrl string, langPrefix string, cp contentPath, fm *FrontMatter, existingUrls map[string]bool) string {
	// Priority 1: Use slug from front matter
	if fm.HasSlug() {
		candidate := fmt.Sprintf("%s/%sdocs/contribute/blog/%s/", baseUrl, langPrefix, fm.Slug)
		if matched := matchUrl(candidate, existingUrls); matched != "" {
			return matched
		}
	}

	// Priority 2: Use file path (remove content/{lang}/ and extension)
	// e.g., content/en/docs/contribute/blog/example.md -> docs/contribute/blog/example
	pathAfterLang := strings.Join(cp.segments[2:], "/")
	pathAfterLang = strings.TrimSuffix(pathAfterLang, filepath.Ext(pathAfterLang))

	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, pathAfterLang)
	if matched := matchUrl(candidate, existingUrls); matched != "" {
		return matched
	}

	return ""
}
