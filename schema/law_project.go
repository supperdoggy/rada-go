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

// BillVotingResults captures detailed bill-based vote events.
type BillVotingResults struct {
	BillID             string          `json:"billId"`
	RegistrationNumber string          `json:"registrationNumber,omitempty"`
	SourceURL          string          `json:"sourceURL"`
	Votes              []BillVoteEvent `json:"votes"`
}

// BillVoteEvent captures one discovered vote event for a bill.
type BillVoteEvent struct {
	GID            string               `json:"gId"`
	Title          string               `json:"title"`
	DateTime       string               `json:"dateTime,omitempty"`
	URL            string               `json:"url"`
	RTFURL         string               `json:"rtfURL,omitempty"`
	PrintURL       string               `json:"printURL,omitempty"`
	Decision       string               `json:"decision,omitempty"`
	Summary        []VoteSummaryItem    `json:"summary,omitempty"`
	FactionSummary []FactionVoteSummary `json:"factionSummary,omitempty"`
	People         []PersonVote         `json:"people,omitempty"`
}

// FactionVoteSummary captures one faction-level vote summary.
type FactionVoteSummary struct {
	Name    string            `json:"name"`
	Members string            `json:"members,omitempty"`
	Summary []VoteSummaryItem `json:"summary,omitempty"`
}

// PersonVote captures one deputy's vote in a vote event.
type PersonVote struct {
	Name      string `json:"name"`
	DeputyID  string `json:"deputyId,omitempty"`
	Faction   string `json:"faction,omitempty"`
	Status    string `json:"status"`
	RawStatus string `json:"rawStatus,omitempty"`
}
