package profiles

type lawProjectV1Profile struct {
	selectors LawProjectSelectorSet
}

func newLawProjectV1Profile() LawProjectProfile {
	return lawProjectV1Profile{
		selectors: lawProjectV1Selectors{},
	}
}

func (p lawProjectV1Profile) ID() string {
	return "law_project/v1"
}

func (p lawProjectV1Profile) Version() string {
	return "v1"
}

func (p lawProjectV1Profile) Selectors() LawProjectSelectorSet {
	return p.selectors
}

type lawProjectV1Selectors struct{}

func (lawProjectV1Selectors) RootContainer() string { return ".wrap-body" }
func (lawProjectV1Selectors) Title() string         { return ".around-create.card-zakonoproektu h2.h2.text" }
func (lawProjectV1Selectors) SummaryRows() string {
	return ".around-create.card-zakonoproektu .info > .row"
}
func (lawProjectV1Selectors) SummaryLabel() string { return ".font-weight-bold" }
func (lawProjectV1Selectors) SummaryValue() string { return ".col" }
func (lawProjectV1Selectors) IDNode() string       { return "#qr-code" }
func (lawProjectV1Selectors) TimelineHeaderStatus() string {
	return "#nav-tab1 thead th:nth-child(2)"
}
func (lawProjectV1Selectors) TimelineRows() string            { return "#nav-tab1 tbody tr" }
func (lawProjectV1Selectors) TimelineDate() string            { return "td:nth-child(1)" }
func (lawProjectV1Selectors) TimelineStatus() string          { return "td:nth-child(2)" }
func (lawProjectV1Selectors) TimelineNote() string            { return "" }
func (lawProjectV1Selectors) RelatedRows() string             { return "#nav-tab3 tbody tr" }
func (lawProjectV1Selectors) RelatedNumberLink() string       { return "td:nth-child(1) a" }
func (lawProjectV1Selectors) RelatedRegistrationDate() string { return "td:nth-child(2)" }
func (lawProjectV1Selectors) RelatedTitle() string            { return "td:nth-child(3)" }
func (lawProjectV1Selectors) LinkURLAttr() string             { return "href" }
func (lawProjectV1Selectors) ChronologyEmbed() string         { return "#nav-tab5 embed" }
func (lawProjectV1Selectors) VotingResultsEmbed() string      { return "#nav-tab6 embed" }
