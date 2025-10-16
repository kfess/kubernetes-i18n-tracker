package url

import (
	"fmt"
	"strings"
)

func generateUrl(baseUrl string, cp contentPath, fm *FrontMatter, existingUrls map[string]bool) (string, error) {
	// index files are handled first
	if cp.isIndex {
		return generateIndexUrl(baseUrl, cp, existingUrls), nil
	}

	switch cp.section {
	case "blog":
		url := generateBlogUrl(baseUrl, cp, fm, existingUrls)
		if url != "" {
			return url, nil
		}
	case "docs":
		url := generateDocsUrl(baseUrl, cp, fm, existingUrls)
		if url != "" {
			return url, nil
		}
	case "case-studies":
		url := generateCaseStudyUrl(baseUrl, cp, existingUrls)
		if url != "" {
			return url, nil
		}
	case "includes":
		// includes files do not have URLs
		return "", nil
	default:
		// careers, community, examples, partners, releases, training, _common-resources
		url := generateFallbackUrl(baseUrl, cp, existingUrls)
		if url != "" {
			return url, nil
		}
	}

	return "", fmt.Errorf("no valid URL found for path: %s", cp.raw)
}

func generateIndexUrl(baseUrl string, cp contentPath, existingUrls map[string]bool) string {
	lang := cp.language
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	// e.g. content/en/blog/_index.md -> https://kubernetes.io/blog/
	// e.g. content/ja/blog/_index.md -> https://kubernetes.io/ja/blog/
	prefix := fmt.Sprintf("content/%s/", lang)
	remainder := strings.TrimPrefix(cp.raw, prefix)
	remainder = strings.TrimSuffix(remainder, "/_index.md")
	remainder = strings.TrimSuffix(remainder, "/_index.html")
	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, remainder)

	return matchUrl(candidate, existingUrls)
}

func generateCaseStudyUrl(baseUrl string, cp contentPath, existingUrls map[string]bool) string {
	lang := cp.language
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	// If only "case-studies" (no subdirectories), return root URL
	// e.g., content/en/case-studies/_index.md -> https://kubernetes.io/case-studies/
	if len(cp.segments) == 3 || (len(cp.segments) == 4 && cp.isIndex) {
		candidate := fmt.Sprintf("%s/%scase-studies/", baseUrl, langPrefix)
		return matchUrl(candidate, existingUrls)
	}

	// Remove last segment (filename) from path
	// e.g., content/en/case-studies/example/index.html -> case-studies/example/
	// Python: parts[1:-1] means skip first (after content/en) and last (filename)
	pathSegments := cp.segments[2 : len(cp.segments)-1] // Skip "content", "lang", and last filename
	casePath := strings.Join(pathSegments, "/")

	candidate := fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, casePath)
	return matchUrl(candidate, existingUrls)
}

func generateFallbackUrl(baseUrl string, cp contentPath, existingUrls map[string]bool) string {
	lang := cp.language
	langPrefix := ""
	if lang != "en" {
		langPrefix = lang + "/"
	}

	prefix := fmt.Sprintf("content/%s/", lang)
	remainder := strings.TrimPrefix(cp.raw, prefix)
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
