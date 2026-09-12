package path

import "slices"

type Category string

const (
	// Docs represents the "docs" category.
	Docs Category = "docs"

	// Blog represents the "blog" category.
	Blog Category = "blog"

	// Career represents the "careers" category.
	Career Category = "careers"

	// Community represents the "community" category.
	Community Category = "community"

	// Example represents the "examples" category.
	Example Category = "examples"

	// Partner represents the "partners" category.
	Partner Category = "partners"

	// Release represents the "releases" category.
	Release Category = "releases"

	// Training represents the "training" category.
	Training Category = "training"

	// CommonResources represents the "_common-resources" category.
	CommonResources Category = "_common-resources"

	// Includes represents the "includes" category.
	Includes Category = "includes"
)

var SupportedCategories = []Category{
	Docs, Blog, Career, Community,
	Example, Partner, Release, Training,
	CommonResources, Includes,
}

// IsSupportedCategory checks if the given category is supported.
func IsSupportedCategory(category string) bool {
	return slices.Contains(SupportedCategories, Category(category))
}
