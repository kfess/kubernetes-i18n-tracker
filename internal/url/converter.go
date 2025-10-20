package url

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
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

	fm, err := c.parser.Parse(path)
	if err != nil {
		return "", err
	}
	if !fm.IsPublic() {
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
	ext      string
}

// parseContentPath parses and validates a content file path
// Expected format: content/{lang}/{section}/{...path} or content/{lang}/{file} (top-level)
func parseContentPath(path string, validLangs []string, validExts []string, validSections []string) (*contentPath, error) {
	if !hasValidExtension(path, validExts) {
		return nil, fmt.Errorf("invalid file extension: %s (allowed: %v)", filepath.Ext(path), validExts)
	}

	if !strings.HasPrefix(path, "content/") {
		return nil, fmt.Errorf("path must start with 'content/', got: %s", path)
	}

	segments := strings.Split(path, "/")

	// Need at least: content/{lang}/{file} (3 segments minimum)
	if len(segments) < 3 {
		return nil, fmt.Errorf("invalid path structure: %s (expected content/{lang}/{file} or content/{lang}/{section}/{file})", path)
	}

	lang := segments[1]
	var section string

	// root files: content/{lang}/{file}
	if len(segments) == 3 {
		section = ""
	} else {
		// regular files: content/{lang}/{section}/{...path}
		section = segments[2]
		if !isValidSection(section, validSections) {
			return nil, fmt.Errorf("unsupported section '%s' in path: %s", section, path)
		}
	}

	if !isValidLanguage(lang, validLangs) {
		return nil, fmt.Errorf("unsupported language code '%s' in path: %s", lang, path)
	}

	isIndex := isIndexFile(path)
	ext := filepath.Ext(path)

	return &contentPath{
		raw:      path,
		language: lang,
		section:  section,
		segments: segments,
		isIndex:  isIndex,
		ext:      ext,
	}, nil
}

func hasValidExtension(path string, validExts []string) bool {
	ext := filepath.Ext(path)
	return slices.Contains(validExts, ext)
}

func isValidLanguage(lang string, validLangs []string) bool {
	return slices.Contains(validLangs, lang)
}

func isValidSection(section string, validSections []string) bool {
	return slices.Contains(validSections, section)
}

func isIndexFile(path string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	return name == "_index" && (ext == ".md" || ext == ".html")
}
