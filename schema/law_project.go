package schema

// LawProjectDetails is the normalized app-facing law project detail payload.
type LawProjectDetails struct {
	ID                  string          `json:"id"`
	Title               string          `json:"title"`
	Status              string          `json:"status"`
	Dates               []DateEntry     `json:"dates"`
	Registration        *Registration   `json:"registration,omitempty"`
	Act                 *Act            `json:"act,omitempty"`
	RegistrationSession string          `json:"registrationSession,omitempty"`
	Category            string          `json:"category,omitempty"`
	InitiativeSubject   string          `json:"initiativeSubject,omitempty"`
	Initiators          []Initiator     `json:"initiators"`
	MainCommittee       *TextLink       `json:"mainCommittee,omitempty"`
	OtherCommittees     []TextLink      `json:"otherCommittees,omitempty"`
	Committees          []Committee     `json:"committees"`
	Documents           []Document      `json:"documents"`
	RelatedBills        []RelatedBill   `json:"relatedBills,omitempty"`
	Timeline            []TimelineEvent `json:"timeline"`
	ChronologyURL       string          `json:"chronologyURL,omitempty"`
	VotingResults       *VotingResults  `json:"votingResults,omitempty"`
	SourceURL           string          `json:"sourceURL"`
}

// DateEntry captures labeled date metadata.
type DateEntry struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Registration captures registration metadata.
type Registration struct {
	Number   string `json:"number"`
	Date     string `json:"date"`
	Archived bool   `json:"archived"`
}

// Act captures adopted-act metadata.
type Act struct {
	Text string `json:"text"`
	Date string `json:"date,omitempty"`
	URL  string `json:"url,omitempty"`
}

// TextLink captures display text and hyperlink metadata.
type TextLink struct {
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
}

// Initiator captures a legislative initiator.
type Initiator struct {
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Committee captures committee metadata.
type Committee struct {
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Document captures related document metadata.
type Document struct {
	Title string `json:"title"`
	Type  string `json:"type,omitempty"`
	Date  string `json:"date,omitempty"`
	URL   string `json:"url"`
}

// RelatedBill captures linked law-project metadata.
type RelatedBill struct {
	Number           string `json:"number"`
	RegistrationDate string `json:"registrationDate,omitempty"`
	Title            string `json:"title"`
	URL              string `json:"url,omitempty"`
}

// TimelineEvent captures a chronological event in law project history.
type TimelineEvent struct {
	Date   string `json:"date,omitempty"`
	Status string `json:"status"`
	Note   string `json:"note,omitempty"`
}

// VotingResults captures the results page link and parsed numeric summary.
type VotingResults struct {
	URL     string            `json:"url"`
	Summary []VoteSummaryItem `json:"summary,omitempty"`
}

// VoteSummaryItem captures one labeled vote-result item.
type VoteSummaryItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
