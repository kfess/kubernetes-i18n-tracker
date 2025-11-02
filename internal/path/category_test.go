package path

import "testing"

func TestIsSupportedCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{"docs is supported", "docs", true},
		{"blog is supported", "blog", true},
		{"unknown is not supported", "unknown", false},
		{"empty string is not supported", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSupportedCategory(tt.category)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
