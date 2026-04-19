package schema

// SearchQuery defines the public request contract for law-project searches.
type SearchQuery struct {
	Session                            int    `json:"session,omitempty"`
	RegistrationNumberCompareOperation int    `json:"registrationNumberCompareOperation,omitempty"`
	RegistrationNumber                 string `json:"registrationNumber,omitempty"`
	RegistrationRangeStart             string `json:"registrationRangeStart,omitempty"`
	RegistrationRangeEnd               string `json:"registrationRangeEnd,omitempty"`
	Name                               string `json:"name,omitempty"`
	DetailView                         bool   `json:"detailView,omitempty"`
	Page                               int    `json:"page,omitempty"`
	PerPage                            int    `json:"perPage,omitempty"`
}

// SearchResponse is the normalized app-facing search payload.
type SearchResponse struct {
	Items      []SearchItem `json:"items"`
	Count      int          `json:"count"`
	Page       int          `json:"page,omitempty"`
	PerPage    int          `json:"perPage,omitempty"`
	TotalPages int          `json:"totalPages,omitempty"`
}

// SearchItem is a normalized search result item.
type SearchItem struct {
	ID                 string `json:"id"`
	RegistrationNumber string `json:"registrationNumber,omitempty"`
	Title              string `json:"title"`
	Status             string `json:"status"`
	RegistrationDate   string `json:"registrationDate"`
	InitiativeSubject  string `json:"initiativeSubject,omitempty"`
	Subject            string `json:"subject"`
	URL                string `json:"url"`
}
