package path

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

// Path represents parsed path information.
type Path struct {
	// Original path
	original string

	// Language code extracted from the path
	language language.Language

	// Category of the content (docs, blog, etc.)
	category Category

	// Path segments
	segments []string

	// File name without extension
	filename string

	// Extension of the file
	extension string
}

// Parse extracts language and category information from a file path.
func Parse(path string) (*Path, error) {
	info := &Path{
		original: path,
		segments: strings.Split(path, "/"),
		category: "overall", // Default to overall instead of unknown
	}

	// Extract language code from path like "content/ja/docs/..."
	parts := strings.Split(path, "/")
	if len(parts) >= 2 && parts[0] == "content" {
		info.language = language.Language(parts[1])

		// Extract category (docs, blog, etc.)
		if len(parts) >= 3 {
			category := parts[2]
			// Skip if it's a file at the root (like _index.html)
			if !strings.HasSuffix(category, ".html") && !strings.HasSuffix(category, ".md") {
				info.category = Category(category)
			}
			// Otherwise keep "overall" as default
		}
	}

	// Extract file name without extension
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	info.filename = strings.TrimSuffix(base, ext)

	// Extract file extension
	info.extension = strings.ToLower(filepath.Ext(path))

	return info, nil
}

// ParseWithValidation parses and validates a content file path with additional constraints.
func ParseWithValidation(path string, validLangs []language.Language, validExts []string, validCategories []Category) (*Path, error) {
	ext := filepath.Ext(path)
	if !slices.Contains(validExts, ext) {
		return nil, fmt.Errorf("invalid file extension: %s (allowed: %v)", ext, validExts)
	}

	if !strings.HasPrefix(path, "content/") {
		return nil, fmt.Errorf("path must start with 'content/', got: %s", path)
	}

	p, err := Parse(path)
	if err != nil {
		return nil, err
	}

	if len(p.segments) < 3 {
		return nil, fmt.Errorf("invalid path structure: %s (expected content/{lang}/{file} or content/{lang}/{section}/{file})", path)
	}

	if !slices.Contains(validLangs, p.language) {
		return nil, fmt.Errorf("unsupported language code '%s' in path: %s", p.language, path)
	}

	if p.category != "overall" && !slices.Contains(validCategories, p.category) {
		return nil, fmt.Errorf("unsupported category '%s' in path: %s", p.category, path)
	}

	return p, nil
}

// ToEnglishPath converts a translated path to its English equivalent.
func (p *Path) ToEnglishPath() string {
	return strings.Replace(p.original, "content/"+string(p.language)+"/", "content/en/", 1)
}

// ToLanguagePath converts an English path to its equivalent in the specified language.
func (p *Path) ToLanguagePath(lang language.Language) string {
	return strings.Replace(p.original, "content/en/", "content/"+string(lang)+"/", 1)
}

// Language returns the language code of the path.
func (p *Path) Language() language.Language {
	return p.language
}

// Category returns the category of the content.
func (p *Path) Category() Category {
	return p.category
}

// Segments returns the path segments.
func (p *Path) Segments() []string {
	return p.segments
}

// Original returns the original file path.
func (p *Path) Original() string {
	return p.original
}

// Filename returns the file name without extension.
func (p *Path) Filename() string {
	if p.filename == "" {
		base := filepath.Base(p.original)
		ext := filepath.Ext(base)
		p.filename = strings.TrimSuffix(base, ext)
	}
	return p.filename
}

// Extension returns the file extension.
func (p *Path) Extension() string {
	return p.extension
}

// IsContentFile checks whether the file is a content file (Markdown or HTML) under the content directory.
func (p *Path) IsContentFile() bool {
	return strings.HasPrefix(p.original, "content/") && (p.isMarkdown() || p.isHTML())
}

// IsMarkdown checks whether the file is a Markdown file.
func (p *Path) isMarkdown() bool {
	return p.extension == ".md"
}

// IsHTML checks whether the file is an HTML file.
func (p *Path) isHTML() bool {
	return p.extension == ".html"
}

// IsIndex checks whether the file is an index file (_index.md or _index.html).
func (p *Path) IsIndex() bool {
	name := p.Filename()
	ext := filepath.Ext(p.original)

	return name == "_index" && (ext == ".md" || ext == ".html")
}
