package language

// SupportedLanguages lists the languages currently supported by the Kubernetes documentation site.
// Languagess are sorted by their website traffic (highest to lowest).
var SupportedLanguages = []string{
	"en",
	"zh-cn",
	"ko",
	"ja",
	"bn",
	"de",
	"es",
	"fr",
	"hi",
	"id",
	"it",
	"pl",
	"pt-br",
	"ru",
	"uk",
	"vi",
}

var DeprecatedLanguages = []string{
	"cn", // Chinese
	"zh", // Chinese
	"pt", // Portuguese
	"no", // Norwegian
}

var NotSupportedYetLanguages = []string{
	"ar", // Arabic
}
