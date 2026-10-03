package structure

import (
	"strconv"
	"strings"
)

type Features struct {
	VisibleLines       int
	H2                 int
	H3                 int
	CodeBlocks         int
	Anchors            map[string]struct{}
	BodyWords          int
	Versions           map[[2]int]struct{}
	FeatureStateTokens map[string]struct{}
	ApiKindTokens      map[string]struct{}
}

// Parse extracts the structural features of a Markdown page.
func Parse(content string) Features {
	// Python reads files in text mode, which turns CRLF and lone CR into LF.
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	h2, h3, codeBlocks, anchors := extractStructure(content)

	return Features{
		VisibleLines:       countVisibleLines(content),
		H2:                 h2,
		H3:                 h3,
		CodeBlocks:         codeBlocks,
		Anchors:            anchors,
		BodyWords:          countBodyWords(content),
		Versions:           extractVersions(content),
		FeatureStateTokens: extractFeatureStateTokens(content),
		ApiKindTokens:      extractApiKindTokens(content),
	}
}

// extractVersions collects Kubernetes version references such as v1.31 as (major, minor) pairs.
func extractVersions(content string) map[[2]int]struct{} {
	// Upstream matches on the raw text, comments included. Kept as is for parity.
	versions := make(map[[2]int]struct{})
	for _, loc := range versionExpr.FindAllStringSubmatchIndex(content, -1) {
		// Emulate the trailing (?!\d): a minor followed by another digit is
		// not a version reference (e.g. v1.2845).
		minorEnd := loc[5]
		if minorEnd < len(content) && content[minorEnd] >= '0' && content[minorEnd] <= '9' {
			continue
		}
		major, err := strconv.Atoi(content[loc[2]:loc[3]])
		if err != nil {
			continue
		}
		minor, err := strconv.Atoi(content[loc[4]:minorEnd])
		if err != nil {
			continue
		}
		versions[[2]int{major, minor}] = struct{}{}
	}
	return versions
}

// countVisibleLines counts the number of visible lines.
func countVisibleLines(content string) int {
	// Code-block contents kept: code volume is legitimate content.
	loc := frontMatterExpr.FindStringIndex(content)
	if loc != nil {
		content = content[loc[1]:]
	}
	content = commentExpr.ReplaceAllString(content, "")
	visibleLines := 0
	for _, line := range strings.Split(content, "\n") {
		if len(strings.TrimSpace(line)) > 0 {
			visibleLines += 1
		}
	}

	return visibleLines
}

func countBodyWords(content string) int {
	// For the Latin-compactness check: separates loose-wrapped full
	// translations from thinned content.
	loc := frontMatterExpr.FindStringIndex(content)
	if loc != nil {
		content = content[loc[1]:]
	}
	content = commentExpr.ReplaceAllString(content, "")
	content = codeExpr.ReplaceAllString(content, "")

	// Drop paragraphs that are entirely indented code blocks.
	kept := []string{}
	for _, raw := range paragraphSplitExpr.Split(content, -1) {
		if isIndentedCodeBlock(raw) {
			continue
		}
		kept = append(kept, raw)
	}
	content = inlineExpr.ReplaceAllString(strings.Join(kept, "\n\n"), "")

	return len(bodyWordExpr.FindAllStringIndex(content, -1))
}

// isIndentedCodeBlock reports whether every non-blank line of the paragraph
// starts with four spaces. A paragraph with no non-blank line is not one.
func isIndentedCodeBlock(paragraph string) bool {
	hasNonBlank := false
	for _, line := range strings.Split(paragraph, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		hasNonBlank = true
		if !strings.HasPrefix(line, "    ") {
			return false
		}
	}
	return hasNonBlank
}

func extractStructure(content string) (int, int, int, map[string]struct{}) {
	// Strip comments first: some localization teams add comments, e.g., zh-cn uses
	// `<!-- -->` blocks to record the English original, which would desync `in_code`.
	// Toggle on 0-3-space fences (CommonMark); only count column-0 fences.
	content = commentExpr.ReplaceAllString(content, "")
	lines := strings.Split(content, "\n")

	h2, h3, fences := 0, 0, 0
	anchors := make(map[string]struct{})
	inCode := false
	for _, line := range lines {
		stripped := strings.TrimLeft(line, " ")
		leading := len(line) - len(stripped)
		if leading < 4 && strings.HasPrefix(stripped, "```") {
			if leading == 0 {
				fences += 1
			}
			inCode = !inCode
			continue
		}
		if inCode {
			continue
		}

		if strings.HasPrefix(line, "### ") {
			h3 += 1
		} else if strings.HasPrefix(line, "## ") {
			h2 += 1
		}

		if strings.HasPrefix(line, "#") {
			matches := anchorExpr.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				anchors[strings.ToLower(strings.TrimSpace(match[1]))] = struct{}{}
			}
		}
	}
	return h2, h3, fences / 2, anchors
}

func extractFeatureStateTokens(content string) map[string]struct{} {
	// Strip comments first: some localization teams add comments, e.g.,
	// zh-cn uses `<!-- -->` blocks to record the English original, which
	// would mask token drift.
	// Prefixes keep version/gate name spaces from colliding.
	content = commentExpr.ReplaceAllString(content, "")
	tokens := make(map[string]struct{})
	matches := featureStateVersionExpr.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		tokens["version:"+match[1]] = struct{}{}
	}
	matches = featureStateGateExpr.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		tokens["gate:"+match[1]] = struct{}{}
	}

	return tokens
}

func extractApiKindTokens(content string) map[string]struct{} {
	// Same comment-stripping rationale as feature_state. Prefixes keep
	// apiVersion and kind value spaces separate.
	content = commentExpr.ReplaceAllString(content, "")
	tokens := make(map[string]struct{})
	matches := apiVersionLineExpr.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		tokens["api:"+match[1]] = struct{}{}
	}
	matches = kindLineExpr.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		tokens["kind:"+match[1]] = struct{}{}
	}

	return tokens
}
