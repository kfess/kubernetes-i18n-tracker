package translation

import (
	"testing"
)

func TestSeverity(t *testing.T) {
	tests := []struct {
		name       string
		totalLines int
		want       Severity
	}{
		{
			name:       "No changes",
			totalLines: 0,
			want:       SeverityCurrent,
		},
		{
			name:       "Minor changes",
			totalLines: 30,
			want:       SeverityMinor,
		},
		{
			name:       "Moderate changes",
			totalLines: 150,
			want:       SeverityModerate,
		},
		{
			name:       "Significant changes",
			totalLines: 300,
			want:       SeveritySignificant,
		},
		{
			name:       "Critical changes",
			totalLines: 600,
			want:       SeverityCritical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateSeverity(tt.totalLines)
			if got != tt.want {
				t.Errorf("calculateSeverity(%d) = %v; want %v", tt.totalLines, got, tt.want)
			}
		})
	}
}
