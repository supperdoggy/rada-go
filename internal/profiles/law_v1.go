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

func (lawProjectV1Selectors) RootContainer() string { return "#LawProjectContainer" }
func (lawProjectV1Selectors) ID() string            { return ".law-id" }
func (lawProjectV1Selectors) Title() string         { return ".law-title" }
func (lawProjectV1Selectors) Status() string        { return ".law-status" }
func (lawProjectV1Selectors) DatesRows() string     { return ".law-dates .date-row" }
func (lawProjectV1Selectors) DateLabel() string     { return ".date-label" }
func (lawProjectV1Selectors) DateValue() string     { return ".date-value" }
func (lawProjectV1Selectors) InitiatorRows() string { return ".law-initiators .initiator-row" }
func (lawProjectV1Selectors) InitiatorName() string { return ".initiator-name" }
func (lawProjectV1Selectors) InitiatorRole() string { return ".initiator-role" }
func (lawProjectV1Selectors) CommitteeRows() string { return ".law-committees .committee-row" }
func (lawProjectV1Selectors) CommitteeName() string { return ".committee-name" }
func (lawProjectV1Selectors) CommitteeRole() string { return ".committee-role" }
func (lawProjectV1Selectors) DocumentRows() string  { return ".law-documents .document-row" }
func (lawProjectV1Selectors) DocumentTitle() string { return ".document-title" }
func (lawProjectV1Selectors) DocumentType() string  { return ".document-type" }
func (lawProjectV1Selectors) DocumentDate() string  { return ".document-date" }
func (lawProjectV1Selectors) DocumentLink() string  { return ".document-link" }
func (lawProjectV1Selectors) DocumentURLAttr() string {
	return "href"
}
func (lawProjectV1Selectors) TimelineRows() string {
	return ".law-timeline .timeline-row"
}
func (lawProjectV1Selectors) TimelineDate() string { return ".timeline-date" }
func (lawProjectV1Selectors) TimelineStatus() string {
	return ".timeline-status"
}
func (lawProjectV1Selectors) TimelineNote() string { return ".timeline-note" }
func (lawProjectV1Selectors) SourceURLAttr() string {
	return "data-source-url"
}
