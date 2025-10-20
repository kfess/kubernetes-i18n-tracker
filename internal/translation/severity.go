package translation

// Severity represents the severity level of translation issues.
type Severity string

const (
	SeverityCurrent     Severity = "current"
	SeverityMinor       Severity = "minor"
	SeverityModerate    Severity = "moderate"
	SeveritySignificant Severity = "significant"
	SeverityCritical    Severity = "critical"
)

// calculateSeverity determines severity based on total changed lines.
func calculateSeverity(totalLines int) Severity {
	switch {
	case totalLines == 0:
		return SeverityCurrent
	case totalLines < 50:
		return SeverityMinor
	case totalLines < 200:
		return SeverityModerate
	case totalLines < 500:
		return SeveritySignificant
	default:
		return SeverityCritical
	}
}
