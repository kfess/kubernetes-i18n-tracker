package language

type Language string

const (
	LanguageEnglish      Language = "en"
	LanguageChinese      Language = "zh-cn"
	LanguageKorean       Language = "ko"
	LanguageJapanese     Language = "ja"
	LanguageBengali      Language = "bn"
	LanguageGerman       Language = "de"
	LanguageSpanish      Language = "es"
	LanguageFrench       Language = "fr"
	LanguageHindi        Language = "hi"
	LanguageIndonesian   Language = "id"
	LanguageItalian      Language = "it"
	LanguagePolish       Language = "pl"
	LanguagePortugueseBR Language = "pt-br"
	LanguageRussian      Language = "ru"
	LanguageUkrainian    Language = "uk"
	LanguageVietnamese   Language = "vi"

	// Deprecated languages
	LanguageChineseDeprecated Language = "cn"
	LanguageNorwegian         Language = "no"
	LanguagePortuguese        Language = "pt"

	// Not supported yet languages
	LanguageArabic Language = "ar"
)

// SupportedLanguages lists the languages currently supported by the Kubernetes documentation site.
// Languages are sorted by their website traffic (highest to lowest).
var SupportedLanguages = []string{
	string(LanguageEnglish),
	string(LanguageChinese),
	string(LanguageKorean),
	string(LanguageJapanese),
	string(LanguageBengali),
	string(LanguageGerman),
	string(LanguageSpanish),
	string(LanguageFrench),
	string(LanguageHindi),
	string(LanguageIndonesian),
	string(LanguageItalian),
	string(LanguagePolish),
	string(LanguagePortugueseBR),
	string(LanguageRussian),
	string(LanguageUkrainian),
	string(LanguageVietnamese),
}

var DeprecatedLanguages = []string{
	string(LanguageChineseDeprecated),
	string(LanguageNorwegian),
	string(LanguagePortuguese),
}

var NotSupportedYetLanguages = []string{
	string(LanguageArabic),
}

func IsSupportedLanguage(lang string) bool {
	for _, l := range SupportedLanguages {
		if l == lang {
			return true
		}
	}

	return false
}
