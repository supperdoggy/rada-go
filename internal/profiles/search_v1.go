package profiles

type searchV1Profile struct {
	selectors SearchSelectorSet
}

func newSearchV1Profile() SearchProfile {
	return searchV1Profile{
		selectors: searchV1Selectors{},
	}
}

func (p searchV1Profile) ID() string {
	return "search/v1"
}

func (p searchV1Profile) Version() string {
	return "v1"
}

func (p searchV1Profile) Selectors() SearchSelectorSet {
	return p.selectors
}

type searchV1Selectors struct{}

func (searchV1Selectors) ResultsContainer() string       { return "#searchResultContainer" }
func (searchV1Selectors) Count() string                  { return "p.condition" }
func (searchV1Selectors) ResultRows() string             { return "table.table tbody tr" }
func (searchV1Selectors) RegistrationNumberLink() string { return "td:nth-child(2) a" }
func (searchV1Selectors) Title() string                  { return "td:nth-child(5)" }
func (searchV1Selectors) RegistrationDate() string       { return "td:nth-child(3)" }
func (searchV1Selectors) InitiativeSubject() string      { return "td:nth-child(4)" }
func (searchV1Selectors) CurrentPage() string            { return "#pagingPage" }
func (searchV1Selectors) PerPage() string                { return "#Paging_per_page" }
func (searchV1Selectors) TotalPages() string             { return ".pager-container .pagination" }
func (searchV1Selectors) LinkURLAttr() string            { return "href" }
