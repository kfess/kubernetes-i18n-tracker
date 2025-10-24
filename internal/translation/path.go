package translation

import (
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

// PathInfo represents parsed path information.
type PathInfo struct {
	// Original path
	Original string

	// Language code extracted from the path
	Language language.Language

	// Category of the content (docs, blog, etc.)
	Category string
}

// parsePath extracts language and category information from a file path.
func parsePath(path string) *PathInfo {
	info := &PathInfo{
		Original: path,
		Category: "overall", // Default to overall instead of unknown
	}

	// Extract language code from path like "content/ja/docs/..."
	parts := strings.Split(path, "/")
	if len(parts) >= 2 && parts[0] == "content" {
		info.Language = language.Language(parts[1])

		// Extract category (docs, blog, etc.)
		if len(parts) >= 3 {
			category := parts[2]
			// Skip if it's a file at the root (like _index.html)
			if !strings.HasSuffix(category, ".html") && !strings.HasSuffix(category, ".md") {
				info.Category = category
			}
			// Otherwise keep "overall" as default
		}
	}

	return info
}

// ToEnglishPath converts a translated path to its English equivalent.
func (pi *PathInfo) ToEnglishPath() string {
	return strings.Replace(pi.Original, "content/"+string(pi.Language)+"/", "content/en/", 1)
}
