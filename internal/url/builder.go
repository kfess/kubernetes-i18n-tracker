package url

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

// UrlBuilder defines the interface for building URLs
type UrlBuilder interface {
	Build(ctx context.Context, path string, lang string)
}

func BuildUrl(ctx context.Context, baseUrl string, path string, existingUrls []string, parser FrontMatterParser) (string, error) {
	if err := validateFilePath(path); err != nil {
		return "", err
	}

	// Not-public files should not have URLs generated
	fm, _ := parser.Parse(path)
	if fm != nil && !fm.IsPublic() {
		return "", fmt.Errorf("file %s is not public, skipping URL generation", path)
	}

	section, lang, err := extractLangAndSection(path)
	if err != nil {
		return "", err
	}

	if isIndexFile(path) {
		langPrefix := ""
		if lang != "en" {
			langPrefix = lang + "/"
		}
		// e.g. content/en/blog/_index.md -> https://kubernetes.io/blog/
		// e.g. content/ja/blog/_index.md -> https://kubernetes.io/ja/blog/
		prefix := fmt.Sprintf("content/%s/", lang)
		remainder := strings.TrimPrefix(path, prefix)
		remainder = strings.TrimSuffix(remainder, "/_index.md")
		remainder = strings.TrimSuffix(remainder, "/_index.html")

		return fmt.Sprintf("%s/%s%s/", baseUrl, langPrefix, remainder), nil
	}

	// directly instantiate the appropriate builder based on section
	switch section {
	case "blog":
		builder := NewBlogUrlBuilder(baseUrl, parser, existingUrls)
		return builder.Build(ctx, path, lang)
	case "docs":
		// builder := NewDocsUrlBuilder(baseUrl, parser, existingUrls)
		// return builder.Build(ctx, path, lang)
	default:
		return "", fmt.Errorf("unsupported section: %s", section)
	}

	return "", fmt.Errorf("failed to build URL for path: %s", path)

}

func validateFilePath(path string) error {
	// Ensure the file path starts with a valid prefix (content/<lang>/)
	isValidPrefix := isValidContentPath(path)
	if !isValidPrefix {
		return fmt.Errorf("invalid file path: %s (must start with content/<lang>/<section>/<file>)", path)
	}

	// Ensure the file has an allowed extension (.md, .html)
	if !hasAllowedExtension(path) {
		return fmt.Errorf("invalid file extension: %s (allowed: .md, .html)", filepath.Ext(path))
	}

	return nil
}

func isValidContentPath(path string) bool {
	validCategories := []string{
		"blog",
		"careers",
		"case-studies",
		"community",
		"docs",
		"examples",
		"includes",
		"partners",
		"releases",
		"training",
		"_common-resources",
	}

	for _, lang := range language.SupportedLanguages {
		prefix := "content/" + lang + "/"
		if !strings.HasPrefix(path, prefix) {
			continue
		}

		remainder := strings.TrimPrefix(path, prefix)
		parts := strings.Split(remainder, "/")

		if len(parts) == 0 {
			return false
		}

		category := parts[0]
		for _, validCat := range validCategories {
			if category == validCat {
				return true
			}
		}
	}

	return false
}

func hasAllowedExtension(path string) bool {
	allowedExtensions := []string{".md", ".html"}
	for _, ext := range allowedExtensions {
		if filepath.Ext(path) == ext {
			return true
		}
	}

	return false
}

func extractLangAndSection(path string) (string, string, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return "", "", fmt.Errorf("invalid file path: %s", path)
	}

	lang := parts[1]
	section := parts[2]
	return lang, section, nil
}

func isIndexFile(path string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	return name == "_index" && (ext == ".md" || ext == ".html")
}
