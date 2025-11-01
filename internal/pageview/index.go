package pageview

// Index provides efficient lookup of page view statistics by URL.
type Index struct {
	// Map from URL to page view statistics
	byURL map[string]*PageViewStats
}

// NewIndex creates an Index from page view statistics.
func NewIndex(stats map[string]*PageViewStats) *Index {
	return &Index{
		byURL: stats,
	}
}

// GetPageViewStats returns page view statistics for the given URL.
// Returns nil if no statistics are found.
func (idx *Index) GetPageViewStats(url string) *PageViewStats {
	return idx.byURL[url]
}

// HasStats returns true if there are statistics for the given URL.
func (idx *Index) HasStats(url string) bool {
	return idx.byURL[url] != nil
}

// TotalURLs returns the number of unique URLs with statistics.
func (idx *Index) TotalURLs() int {
	return len(idx.byURL)
}
