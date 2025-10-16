package url

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

type Converter struct {
	config Config
	parser FrontMatterParser
}

type Config struct {
	BaseUrl        string
	ExistingUrls   map[string]bool
	SupportedLangs []string
	SupportedExts  []string
	ValidSections  []string
}

func NewConverter(config Config, parser FrontMatterParser) *Converter {
	return &Converter{
		config: config,
		parser: parser,
	}
}

func (c *Converter) Convert(ctx context.Context, path string) (string, error) {
	cp, err := parseContentPath(path, c.config.SupportedLangs, c.config.SupportedExts, c.config.ValidSections)
	if err != nil {
		return "", err
	}

	fm, _ := c.parser.Parse(path)
	if fm != nil && !fm.IsPublic() {
		return "", fmt.Errorf("file %s is not public, skipping URL generation", path)
	}

	url, err := generateUrl(c.config.BaseUrl, *cp, fm, c.config.ExistingUrls)
	if err != nil {
		return "", err
	}
	return url, nil
}

type contentPath struct {
	raw      string
	language string
	section  string
	segments []string
	isIndex  bool
}

// parseContentPath parses and validates a content file path
// Expected format: content/{lang}/{section}/{...path}
func parseContentPath(path string, validLangs []string, validExts []string, validSections []string) (*contentPath, error) {
	if !hasValidExtension(path, validExts) {
		return nil, fmt.Errorf("invalid file extension: %s (allowed: %v)", filepath.Ext(path), validExts)
	}

	if !strings.HasPrefix(path, "content/") {
		return nil, fmt.Errorf("path must start with 'content/', got: %s", path)
	}

	segments := strings.Split(path, "/")

	// Need at least: content/{lang}/{category}/{file}
	if len(segments) < 4 {
		return nil, fmt.Errorf("invalid path structure: %s (expected content/{lang}/{category}/{file})", path)
	}

	lang := segments[1]
	section := segments[2]

	if !isValidLanguage(lang, validLangs) {
		return nil, fmt.Errorf("unsupported language code '%s' in path: %s", lang, path)
	}

	if !isValidSection(section, validSections) {
		return nil, fmt.Errorf("unsupported section '%s' in path: %s", section, path)
	}

	isIndex := isIndexFile(path)

	return &contentPath{
		raw:      path,
		language: lang,
		section:  section,
		segments: segments,
		isIndex:  isIndex,
	}, nil
}

func hasValidExtension(path string, validExts []string) bool {
	ext := filepath.Ext(path)
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func isValidLanguage(lang string, validLangs []string) bool {
	for _, validLang := range validLangs {
		if lang == validLang {
			return true
		}
	}
	return false
}

func isValidSection(section string, validSections []string) bool {
	for _, validSec := range validSections {
		if section == validSec {
			return true
		}
	}
	return false
}

func isIndexFile(path string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	return name == "_index" && (ext == ".md" || ext == ".html")
}
