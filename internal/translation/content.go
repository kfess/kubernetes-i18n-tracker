package translation

import (
	"regexp"
	"strings"
)

// HeaderCount represents the count of markdown headers at each level.
type HeaderCount struct {
	H1    int
	H2    int
	H3    int
	H4    int
	H5    int
	H6    int
	Total int
}

// CountHeadersFromContent counts markdown headers from content string, excluding those in front matter and comments.
func CountHeadersFromContent(content string) *HeaderCount {
	count := &HeaderCount{}

	if content == "" {
		return count
	}

	lines := strings.Split(content, "\n")
	inFrontMatter := false
	inComment := false
	frontMatterCount := 0

	// Regex pattern for markdown headers: ^(#{1,6})\s+.+
	headerPattern := regexp.MustCompile(`^(#{1,6})\s+.+`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle YAML front matter (--- ... ---)
		if trimmed == "---" {
			frontMatterCount++
			if frontMatterCount == 1 {
				inFrontMatter = true
				continue
			} else if frontMatterCount == 2 {
				inFrontMatter = false
				continue
			}
		}

		// Skip lines inside front matter
		if inFrontMatter {
			continue
		}

		// Handle HTML comments (<!-- -->)
		if strings.Contains(trimmed, "<!--") {
			inComment = true
			// Check if comment ends on the same line
			if strings.Contains(trimmed, "-->") {
				inComment = false
			}
			continue
		}

		// Skip lines inside comment
		if inComment {
			if strings.Contains(trimmed, "-->") {
				inComment = false
			}
			continue
		}

		// Check for markdown header
		matches := headerPattern.FindStringSubmatch(trimmed)
		if len(matches) >= 2 {
			headerLevel := len(matches[1])
			switch headerLevel {
			case 1:
				count.H1++
			case 2:
				count.H2++
			case 3:
				count.H3++
			case 4:
				count.H4++
			case 5:
				count.H5++
			case 6:
				count.H6++
			}
			count.Total++
		}
	}

	return count
}

// HeadersMatch returns true if two header counts have the same total.
func HeadersMatch(a, b *HeaderCount) bool {
	if a == nil || b == nil {
		return false
	}

	return a.H1 == b.H1 &&
		a.H2 == b.H2 &&
		a.H3 == b.H3 &&
		a.H5 == b.H5 &&
		a.H6 == b.H6 &&
		a.Total == b.Total // to validate the existence of H7, H8, ...
}

// HeaderDiff returns the difference in total header count (translation - english).
func HeaderDiff(translation, english *HeaderCount) int {
	if translation == nil || english == nil {
		return 0
	}
	return translation.Total - english.Total
}
