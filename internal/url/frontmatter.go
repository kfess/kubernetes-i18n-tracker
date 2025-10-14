package url

type FrontMatter struct {
	Url      string         `yaml:"url"`
	Slug     string         `yaml:"slug"`
	Date     string         `yaml:"date"`
	Title    string         `yaml:"title"`
	FullLink string         `yaml:"full_link"`
	Build    *BuildSettings `yaml:"_build"`
}

type BuildSettings struct {
	Render interface{} `yaml:"render"` // string "never" または bool false の両方に対応
}

// IsPublic determines if the content should be publicly accessible based on the Build settings.
// Returns false only if Build.Render is explicitly set to "never" or false.
// Default behavior is to consider content as public.
func (fm *FrontMatter) IsPublic() bool {
	if fm.Build != nil {
		return true
	}

	switch v := fm.Build.Render.(type) {
	case string:
		return v != "never"
	case bool:
		return v
	default:
		return true
	}
}

func (fm *FrontMatter) HasSlug() bool {
	return fm.Slug != ""
}

func (fm *FrontMatter) HasDate() bool {
	return fm.Date != ""
}

func (fm *FrontMatter) HasExplicitUrl() bool {
	return fm.Url != ""
}

func (fm *FrontMatter) HasFullLink() bool {
	return fm.FullLink != ""
}

func (fm *FrontMatter) HasTitle() bool {
	return fm.Title != ""
}
