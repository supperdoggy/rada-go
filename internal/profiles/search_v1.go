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

func (searchV1Selectors) ResultsContainer() string { return "#SearchResultContainer" }
func (searchV1Selectors) ResultRows() string       { return ".search-result-row" }
func (searchV1Selectors) ID() string               { return ".result-id" }
func (searchV1Selectors) Title() string            { return ".result-title" }
func (searchV1Selectors) TitleLink() string        { return ".result-title" }
func (searchV1Selectors) Status() string           { return ".result-status" }
func (searchV1Selectors) RegistrationDate() string { return ".result-registration-date" }
func (searchV1Selectors) Subject() string          { return ".result-subject" }
func (searchV1Selectors) LinkURLAttr() string      { return "href" }
