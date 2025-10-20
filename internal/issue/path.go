package issue

import (
	"regexp"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
)

var languageLabels = map[string]language.Language{
	"language/en": language.LanguageEnglish,
	"language/ko": language.LanguageKorean,
	"language/ja": language.LanguageJapanese,
	"language/zh": language.LanguageChinese,
	"language/pt": language.LanguagePortugueseBR,
	"language/es": language.LanguageSpanish,
	"language/hi": language.LanguageHindi,
	"language/id": language.LanguageIndonesian,
	"language/de": language.LanguageGerman,
	"language/fr": language.LanguageFrench,
	"language/it": language.LanguageItalian,
	"language/vi": language.LanguageVietnamese,
	"language/ru": language.LanguageRussian,
	"language/uk": language.LanguageUkrainian,
	"language/pl": language.LanguagePolish,
	"language/bn": language.LanguageBengali,
}

// GuessLanguage guesses the language from issue labels or title.
func GuessLanguage(issue Issue) language.Language {
	// Check labels first
	for _, label := range issue.Labels {
		if lang, ok := languageLabels[label]; ok {
			return lang
		}
	}

	// Try to extract from title like "[ja]", "[ko]", etc.
	titlePattern := regexp.MustCompile(`^\[([^\]]+)\]`)
	if matches := titlePattern.FindStringSubmatch(strings.TrimSpace(issue.Title)); matches != nil {
		langStr := strings.ToLower(strings.TrimSpace(matches[1]))

		// Map common title patterns to language codes
		switch langStr {
		case "en":
			return language.LanguageEnglish
		case "ko":
			return language.LanguageKorean
		case "ja":
			return language.LanguageJapanese
		case "zh", "zh-cn":
			return language.LanguageChinese
		case "pt", "pt-br":
			return language.LanguagePortugueseBR
		case "es":
			return language.LanguageSpanish
		case "hi":
			return language.LanguageHindi
		case "id":
			return language.LanguageIndonesian
		case "de":
			return language.LanguageGerman
		case "fr":
			return language.LanguageFrench
		case "it":
			return language.LanguageItalian
		case "vi":
			return language.LanguageVietnamese
		case "ru":
			return language.LanguageRussian
		case "uk":
			return language.LanguageUkrainian
		case "pl":
			return language.LanguagePolish
		case "bn":
			return language.LanguageBengali
		}
	}

	return language.Language("")
}

// ExtractPathLikeString extracts a path-like string from the issue title.
func ExtractPathLikeString(title string) string {
	patterns := []string{
		`/[a-zA-Z0-9/_.-]+`,                 // starts with slash
		`[a-zA-Z0-9/_.-]+/[a-zA-Z0-9/_.-]+`, // contains slash
		`[a-zA-Z0-9/_.-]+\.(?:md|html)`,     // .md or .html file
	}

	var longestMatch string
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindString(title); match != "" {
			match = strings.TrimPrefix(match, "/")
			if len(match) > len(longestMatch) {
				longestMatch = match
			}
		}
	}

	return longestMatch
}

// GeneratePathCandidates generates a list of path candidates based on the given path and language.
func GeneratePathCandidates(path string, lang language.Language) []string {
	path = strings.ToLower(strings.TrimSpace(path))
	if path == "" {
		return nil
	}

	// Replace language codes in path
	langPattern := regexp.MustCompile(`\b(en|ko|ja|zh-cn|zh|pt-br|pt|es|hi|id|de|fr|it|vi|ru|uk|pl|bn)\b`)
	path = langPattern.ReplaceAllString(path, string(lang))

	// Remove k8s.io/ prefix
	path = strings.TrimPrefix(path, "k8s.io/")

	// Convert underscores to hyphens (except for _index files)
	hyphenatedPath := convertUnderscoresToHyphens(path)

	// Common path prefixes (matching Python implementation)
	prefixes := []string{
		"",
		"content/",
		"content/" + string(lang) + "/",
		"content/" + string(lang) + "/docs/",
		"content/" + string(lang) + "/docs/concepts/",
		"content/" + string(lang) + "/docs/contribute/",
		"content/" + string(lang) + "/docs/doc-contributor-tools/",
		"content/" + string(lang) + "/docs/home/",
		"content/" + string(lang) + "/docs/images/",
		"content/" + string(lang) + "/docs/reference/",
		"content/" + string(lang) + "/docs/setup/",
		"content/" + string(lang) + "/docs/tasks/",
		"content/" + string(lang) + "/docs/tutorials/",
		"content/" + string(lang) + "/blog/",
		"content/" + string(lang) + "/blog/_posts/",
		"content/" + string(lang) + "/careers/",
		"content/" + string(lang) + "/case-studies/",
		"content/" + string(lang) + "/community/",
		"content/" + string(lang) + "/examples/",
		"content/" + string(lang) + "/includes/",
		"content/" + string(lang) + "/partners/",
		"content/" + string(lang) + "/releases/",
		"content/" + string(lang) + "/training/",
	}

	candidates := make(map[string]bool)

	// Add candidates with file extensions (matching Python)
	addCandidatesWithExtensions := func(base string) {
		candidates[base] = true
		if !strings.HasSuffix(base, ".md") && !strings.HasSuffix(base, ".html") {
			candidates[base+".md"] = true
			candidates[base+"/index.md"] = true
			candidates[base+"/_index.md"] = true
			candidates[base+".html"] = true
			candidates[base+"/index.html"] = true
			candidates[base+"/_index.html"] = true
		}
	}

	// prefix + path
	for _, prefix := range prefixes {
		base := prefix + path
		addCandidatesWithExtensions(base)

		base2 := prefix + hyphenatedPath
		if base2 != base {
			addCandidatesWithExtensions(base2)
		}
	}

	// shorten path (progressively remove leading components)
	parts := strings.Split(path, "/")
	for i := 1; i < len(parts); i++ {
		shortened := strings.Join(parts[i:], "/")
		addCandidatesWithExtensions(shortened)
	}

	parts2 := strings.Split(hyphenatedPath, "/")
	for i := 1; i < len(parts2); i++ {
		shortened := strings.Join(parts2[i:], "/")
		addCandidatesWithExtensions(shortened)
	}

	// Convert map to slice
	result := make([]string, 0, len(candidates))
	for candidate := range candidates {
		result = append(result, candidate)
	}

	return result
}

// convertUnderscoresToHyphens converts underscores to hyphens, preserving _index files.
func convertUnderscoresToHyphens(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if !strings.HasPrefix(part, "_") {
			parts[i] = strings.ReplaceAll(part, "_", "-")
		}
	}
	return strings.Join(parts, "/")
}

// GuessPath attempts to guess the file path for an issue.
// Returns the first matching path from existing paths, or empty string if no match.
func GuessPath(issue Issue, lang language.Language, existingPaths map[string]bool) string {
	pathLike := ExtractPathLikeString(issue.Title)
	if pathLike == "" {
		return ""
	}

	candidates := GeneratePathCandidates(pathLike, lang)

	// Try to find the first existing file match
	for _, candidate := range candidates {
		if existingPaths[candidate] {
			return candidate
		}
	}

	// No match found
	return ""
}
