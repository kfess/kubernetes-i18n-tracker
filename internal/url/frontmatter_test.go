package url

import "testing"

func TestFrontMatter_IsPublic(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrontMatter
		expected bool
	}{
		{
			name:     "Build is nil",
			fm:       &FrontMatter{},
			expected: true,
		},
		{
			name: "Build.Render is nil",
			fm: &FrontMatter{
				Build: &BuildSettings{},
			},
			expected: true,
		},
		{
			name: "Build.Render is 'never' string",
			fm: &FrontMatter{
				Build: &BuildSettings{
					Render: "never",
				},
			},
			expected: false,
		},
		{
			name: "Build.Render is 'always' string",
			fm: &FrontMatter{
				Build: &BuildSettings{
					Render: "always",
				},
			},
			expected: true,
		},
		{
			name: "Build.Render is true bool",
			fm: &FrontMatter{
				Build: &BuildSettings{
					Render: true,
				},
			},
			expected: true,
		},
		{
			name: "Build.Render is false bool",
			fm: &FrontMatter{
				Build: &BuildSettings{
					Render: false,
				},
			},
			expected: false,
		},
		{
			name: "Build.Render is other type (int)",
			fm: &FrontMatter{
				Build: &BuildSettings{
					Render: 123,
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.IsPublic()
			if got != tt.expected {
				t.Errorf("IsPublic() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFrontMatter_HasSlug(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrontMatter
		expected bool
	}{
		{
			name:     "Empty slug",
			fm:       &FrontMatter{Slug: ""},
			expected: false,
		},
		{
			name:     "Non-empty slug",
			fm:       &FrontMatter{Slug: "my-post"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.HasSlug()
			if got != tt.expected {
				t.Errorf("HasSlug() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFrontMatter_HasDate(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrontMatter
		expected bool
	}{
		{
			name:     "Empty date",
			fm:       &FrontMatter{Date: ""},
			expected: false,
		},
		{
			name:     "Non-empty date",
			fm:       &FrontMatter{Date: "2025-11-01"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.HasDate()
			if got != tt.expected {
				t.Errorf("HasDate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFrontMatter_HasExplicitUrl(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrontMatter
		expected bool
	}{
		{
			name:     "Empty url",
			fm:       &FrontMatter{Url: ""},
			expected: false,
		},
		{
			name:     "Non-empty url",
			fm:       &FrontMatter{Url: "/docs/concepts/overview"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.HasExplicitUrl()
			if got != tt.expected {
				t.Errorf("HasExplicitUrl() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFrontMatter_HasFullLink(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrontMatter
		expected bool
	}{
		{
			name:     "Empty full_link",
			fm:       &FrontMatter{FullLink: ""},
			expected: false,
		},
		{
			name:     "Non-empty full_link",
			fm:       &FrontMatter{FullLink: "https://example.com"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.HasFullLink()
			if got != tt.expected {
				t.Errorf("HasFullLink() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFrontMatter_HasTitle(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrontMatter
		expected bool
	}{
		{
			name:     "Empty title",
			fm:       &FrontMatter{Title: ""},
			expected: false,
		},
		{
			name:     "Non-empty title",
			fm:       &FrontMatter{Title: "My Blog Post"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.HasTitle()
			if got != tt.expected {
				t.Errorf("HasTitle() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasTags(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrontMatter
		expected bool
	}{
		{
			name:     "No tags",
			fm:       &FrontMatter{Tags: []string{}},
			expected: false,
		},
		{
			name:     "With tags",
			fm:       &FrontMatter{Tags: []string{"kubernetes", "docs"}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.HasTags()
			if got != tt.expected {
				t.Errorf("HasTags() = %v, want %v", got, tt.expected)
			}
		})
	}
}
