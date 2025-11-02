package language

import "slices"

type Language string

const (
	English      Language = "en"
	Chinese      Language = "zh-cn"
	Korean       Language = "ko"
	Japanese     Language = "ja"
	Bengali      Language = "bn"
	German       Language = "de"
	Spanish      Language = "es"
	French       Language = "fr"
	Hindi        Language = "hi"
	Indonesian   Language = "id"
	Italian      Language = "it"
	Polish       Language = "pl"
	PortugueseBR Language = "pt-br"
	Russian      Language = "ru"
	Ukrainian    Language = "uk"
	Vietnamese   Language = "vi"

	// Deprecated languages
	ChineseDeprecated Language = "cn"
	Norwegian         Language = "no"
	Portuguese        Language = "pt"

	// Not supported yet languages
	Arabic Language = "ar"
)

// SupportedLanguages lists the languages currently supported by the Kubernetes documentation site.
// Languages are sorted by their website traffic (highest to lowest).
var SupportedLanguages = []Language{
	English,
	Japanese,
	Korean,
	Chinese,
	PortugueseBR,
	Spanish,
	Hindi,
	Indonesian,
	German,
	French,
	Italian,
	Vietnamese,
	Russian,
	Ukrainian,
	Polish,
	Bengali,
}

var DeprecatedLanguages = []Language{
	ChineseDeprecated,
	Norwegian,
	Portuguese,
}

var NotSupportedYetLanguages = []Language{
	Arabic,
}

// IsSupported checks if the given language is in the SupportedLanguages list.
func IsSupported(lang string) bool {
	return slices.Contains(SupportedLanguages, Language(lang))
}

// IsDeprecated checks if the given language is in the DeprecatedLanguages list.
func IsDeprecated(lang string) bool {
	return slices.Contains(DeprecatedLanguages, Language(lang))
}

// IsNotSupportedYet checks if the given language is in the NotSupportedYetLanguages list.
func IsNotSupportedYet(lang string) bool {
	return slices.Contains(NotSupportedYetLanguages, Language(lang))
}
