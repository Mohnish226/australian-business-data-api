package models

// SearchResponse represents the API response structure
type SearchResponse struct {
	Help    string                 `json:"help"`
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result"`
}

// SearchParams represents the parameters for a search request
type SearchParams struct {
	BaseURL    string
	ResourceID string
	Query      string
	Filter     map[string]string
	Limit      int
	UseCache   bool
}
