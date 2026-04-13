package schema

// SearchQuery defines the public request contract for law-project searches.
type SearchQuery struct {
	Term    string            `json:"term,omitempty"`
	Page    int               `json:"page,omitempty"`
	PerPage int               `json:"perPage,omitempty"`
	Filters map[string]string `json:"filters,omitempty"`
}

// SearchResponse is the normalized app-facing search payload.
type SearchResponse struct {
	Items []SearchItem `json:"items"`
	Count int          `json:"count"`
}

// SearchItem is a normalized search result item.
type SearchItem struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Status           string `json:"status"`
	RegistrationDate string `json:"registrationDate"`
	Subject          string `json:"subject"`
	URL              string `json:"url"`
}
