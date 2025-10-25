package language

import "testing"

func TestIsSupportedLanguage(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want bool
	}{
		{"Supported English", "en", true},
		{"Supported Chinese", "zh-cn", true},
		{"Supported Spanish", "ja", true},
		{"Unsupported Arabic", "ar", false},
		{"Deprecated Norwegian", "no", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupportedLanguage(tt.lang); got != tt.want {
				t.Errorf("IsSupportedLanguage(%q) = %v; want %v", tt.lang, got, tt.want)
			}
		})
	}
}
