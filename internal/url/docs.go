package url

import (
	"fmt"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/path"
)

func generateDocsUrl(baseUrl string, p *path.Path, fm *FrontMatter, existingUrls map[string]bool) string {
	lang := string(p.Language())
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	segments := p.Segments()

	// Check if this is a special case
	if len(segments) >= 5 {
		subSection := segments[3] // e.g., "reference", "contribute"

		// Special case: docs/reference/glossary
		if subSection == "reference" && len(segments) >= 5 && segments[4] == "glossary" {
			return buildGlossaryUrl(baseUrl, langPrefix, fm)
		}

		// Special case: docs/contribute/blog
		if subSection == "contribute" && len(segments) >= 5 && segments[4] == "blog" {
			return buildContributeBlogUrl(baseUrl, langPrefix, p, fm, existingUrls)
		}
	}

	// Normal docs URL handling
	// e.g., content/en/docs/concepts/overview.md -> docs/concepts/overview
	pathAfterLang := strings.Join(segments[2:], "/") // Skip "content" and "{lang}"
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

func buildContributeBlogUrl(baseUrl string, langPrefix string, p *path.Path, fm *FrontMatter, existingUrls map[string]bool) string {
	// Priority 1: Use slug from front matter
	if fm.HasSlug() {
		candidate := fmt.Sprintf("%s/%sdocs/contribute/blog/%s/", baseUrl, langPrefix, fm.Slug)
		if matched := matchUrl(candidate, existingUrls); matched != "" {
			return matched
		}
	}

	// Priority 2: Use filename
	filename := p.Filename()
	filename = strings.TrimSuffix(filename, ".md")
	filename = strings.TrimSuffix(filename, ".html")

	candidate := fmt.Sprintf("%s/%sdocs/contribute/blog/%s/", baseUrl, langPrefix, filename)
	if matched := matchUrl(candidate, existingUrls); matched != "" {
		return matched
	}

	return ""
}
