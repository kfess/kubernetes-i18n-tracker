package structure

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	versionExpr             = regexp.MustCompile(`(?:^|[^A-Za-z0-9])v([0-9]+)\.([0-9]{1,3})`)
	featureStateVersionExpr = regexp.MustCompile(`for_k8s_version="(v\d+\.\d+)"`)
	featureStateGateExpr    = regexp.MustCompile(`feature_gate_name="([A-Za-z_][A-Za-z0-9_]*)"`)
	anchorExpr              = regexp.MustCompile(`\{#([^}]+)\}`)
	frontMatterExpr         = regexp.MustCompile(`(?s)^---\s*\n.*?\n---\s*\n`)
	codeExpr                = regexp.MustCompile("(?s)```.*?```")
	commentExpr             = regexp.MustCompile(`(?s)<!--.*?-->`)
	inlineExpr              = regexp.MustCompile("`[^`\n]+`")
	bodyWordExpr            = regexp.MustCompile(`[\p{L}\p{N}]{2,}`)
	paragraphSplitExpr      = regexp.MustCompile(`\n{2,}`)
	apiVersionLineExpr      = regexp.MustCompile(`(?m)^[ \t]*-?[ \t]*apiVersion:\s*\"?([A-Za-z0-9./_-]+)\"?\s*$`)
	kindLineExpr            = regexp.MustCompile(`(?m)^[ \t]*-?[ \t]*kind:\s*\"?([A-Z][A-Za-z0-9]+)\"?\s*$`)
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
	// Normalize line endings so CRLF files are parsed like LF files.
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

// countVisibleLines counts the non-blank lines outside the front matter and
// HTML comments. Code-block lines count: code volume is legitimate content.
func countVisibleLines(content string) int {
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

// countBodyWords counts the words of prose, leaving out front matter, HTML
// comments, code blocks and inline code. It tells a fully translated page
// that merely wraps its lines differently from one that lost content.
func countBodyWords(content string) int {
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

// extractStructure returns the number of H2 headings, H3 headings and code
// blocks, and the set of heading anchors ({#id}, lowercased). Headings and
// anchors inside code blocks are skipped.
func extractStructure(content string) (int, int, int, map[string]struct{}) {
	// Strip comments first: some localization teams add comments, e.g., zh-cn uses
	// `<!-- -->` blocks to record the English original, which would desync inCode.
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

// extractFeatureStateTokens collects the version and feature gate named by
// feature-state shortcodes, as "version:v1.31" and "gate:Name" tokens.
func extractFeatureStateTokens(content string) map[string]struct{} {
	// Strip comments first: some localization teams add comments, e.g.,
	// zh-cn uses `<!-- -->` blocks to record the English original, which
	// would mask token drift.
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

// extractApiKindTokens collects the apiVersion and kind values of YAML
// examples, as "api:apps/v1" and "kind:Deployment" tokens.
func extractApiKindTokens(content string) map[string]struct{} {
	// Comments are stripped for the same reason as in extractFeatureStateTokens.
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
