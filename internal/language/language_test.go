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
			if got := IsSupported(tt.lang); got != tt.want {
				t.Errorf("IsSupportedLanguage(%q) = %v; want %v", tt.lang, got, tt.want)
			}
		})
	}
}

func TestIsDeprecatedLanguage(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want bool
	}{
		{"Deprecated Chinese", "cn", true},
		{"Deprecated Norwegian", "no", true},
		{"Deprecated Portuguese", "pt", true},
		{"Supported English", "en", false},
		{"Unsupported Arabic", "ar", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDeprecated(tt.lang); got != tt.want {
				t.Errorf("IsDeprecated(%q) = %v; want %v", tt.lang, got, tt.want)
			}
		})
	}
}

func TestIsNotSupportedYetLanguage(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want bool
	}{
		{"Not Supported Yet Arabic", "ar", true},
		{"Supported English", "en", false},
		{"Deprecated Chinese", "cn", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotSupportedYet(tt.lang); got != tt.want {
				t.Errorf("IsNotSupportedYet(%q) = %v; want %v", tt.lang, got, tt.want)
			}
		})
	}
}
