package adapters

// Common interfaces
type SearchProvider interface {
	Search(query string, options map[string]string) (*SearchResponse, error)
}

type ExtractProvider interface {
	Extract(urls []string) (*ExtractResponse, error)
}

// Response types
type SearchResult struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

type ExtractResult struct {
	URL        string `json:"url"`
	RawContent string `json:"raw_content"`
}

type FailedResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

type ExtractResponse struct {
	Results       []ExtractResult `json:"results"`
	FailedResults []FailedResult  `json:"failed_results"`
	ResponseTime  float64         `json:"response_time"`
}
