package url

import (
	"fmt"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/path"
)

func generateUrl(baseUrl string, p *path.Path, fm *FrontMatter, existingUrls map[string]bool) (string, error) {
	// index files are handled first
	if p.IsIndex() {
		return generateIndexUrl(baseUrl, p, existingUrls), nil
	}

	category := string(p.Category())

	// Handle root files (e.g., content/ja/search.md)
	if category == "overall" {
		url := generateRootUrl(baseUrl, p, existingUrls)
		if url != "" {
			return url, nil
		}
		return "", fmt.Errorf("no valid URL found for top-level path: %s", p.Original())
	}

	switch category {
	case "blog":
		url := generateBlogUrl(baseUrl, p, fm, existingUrls)
		if url != "" {
			return url, nil
		}
	case "docs":
		url := generateDocsUrl(baseUrl, p, fm, existingUrls)
		if url != "" {
			return url, nil
		}
	case "case-studies":
		url := generateCaseStudyUrl(baseUrl, p, existingUrls)
		if url != "" {
			return url, nil
		}
	case "includes":
		// includes files do not have URLs
		return "", nil
	default:
		// careers, community, examples, partners, releases, training, _common-resources
		url := generateFallbackUrl(baseUrl, p, existingUrls)
		if url != "" {
			return url, nil
		}
	}

	return "", fmt.Errorf("no valid URL found for path: %s", p.Original())
}

func generateIndexUrl(baseUrl string, p *path.Path, existingUrls map[string]bool) string {
	lang := string(p.Language())
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	// e.g. content/en/blog/_index.md -> https://kubernetes.io/blog/
	// e.g. content/ja/blog/_index.md -> https://kubernetes.io/ja/blog/
	prefix := fmt.Sprintf("content/%s/", lang)
	remainder := strings.TrimPrefix(p.Original(), prefix)
	remainder = strings.TrimSuffix(remainder, "/_index.md")
	remainder = strings.TrimSuffix(remainder, "/_index.html")
	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, remainder)

	return matchUrl(candidate, existingUrls)
}

func generateCaseStudyUrl(baseUrl string, p *path.Path, existingUrls map[string]bool) string {
	lang := string(p.Language())
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	segments := p.Segments()

	// If only "case-studies" (no subdirectories), return root URL
	// e.g., content/en/case-studies/_index.md -> https://kubernetes.io/case-studies/
	if len(segments) == 3 || (len(segments) == 4 && p.IsIndex()) {
		candidate := fmt.Sprintf("%s/%scase-studies/", baseUrl, langPrefix)
		return matchUrl(candidate, existingUrls)
	}

	// Remove last segment (filename) from path
	// e.g., content/en/case-studies/example/index.html -> case-studies/example/
	// Python: parts[1:-1] means skip first (after content/en) and last (filename)
	pathSegments := segments[2 : len(segments)-1] // Skip "content", "lang", and last filename
	casePath := strings.Join(pathSegments, "/")

	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, casePath)
	return matchUrl(candidate, existingUrls)
}

func generateFallbackUrl(baseUrl string, p *path.Path, existingUrls map[string]bool) string {
	lang := string(p.Language())
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	prefix := fmt.Sprintf("content/%s/", lang)
	remainder := strings.TrimPrefix(p.Original(), prefix)
	remainder = strings.TrimSuffix(remainder, ".md")
	remainder = strings.TrimSuffix(remainder, ".html")
	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, remainder)

	return matchUrl(candidate, existingUrls)
}

func matchUrl(candidate string, existingUrls map[string]bool) string {
	// Exact match
	if _, ok := existingUrls[candidate]; ok {
		return candidate
	}

	// lower-case match
	lowerCandidate := strings.ToLower(candidate)
	if _, ok := existingUrls[lowerCandidate]; ok {
		return lowerCandidate
	}

	return ""
}

// generateRootUrl generates URLs for root content files
// e.g., content/ja/search.md -> https://kubernetes.io/ja/search/
func generateRootUrl(baseUrl string, p *path.Path, existingUrls map[string]bool) string {
	lang := string(p.Language())
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	prefix := fmt.Sprintf("content/%s/", lang)
	remainder := strings.TrimPrefix(p.Original(), prefix)
	remainder = strings.TrimSuffix(remainder, ".md")
	remainder = strings.TrimSuffix(remainder, ".html")

	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, remainder)

	return matchUrl(candidate, existingUrls)
}
