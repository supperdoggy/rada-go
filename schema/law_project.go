package schema

// LawProjectDetails is the normalized app-facing law project detail payload.
type LawProjectDetails struct {
	ID         string          `json:"id"`
	Title      string          `json:"title"`
	Status     string          `json:"status"`
	Dates      []DateEntry     `json:"dates"`
	Initiators []Initiator     `json:"initiators"`
	Committees []Committee     `json:"committees"`
	Documents  []Document      `json:"documents"`
	Timeline   []TimelineEvent `json:"timeline"`
	SourceURL  string          `json:"sourceURL"`
}

// DateEntry captures labeled date metadata.
type DateEntry struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Initiator captures a legislative initiator.
type Initiator struct {
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
}

// Committee captures committee metadata.
type Committee struct {
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
}

// Document captures related document metadata.
type Document struct {
	Title string `json:"title"`
	Type  string `json:"type,omitempty"`
	Date  string `json:"date,omitempty"`
	URL   string `json:"url"`
}

// TimelineEvent captures a chronological event in law project history.
type TimelineEvent struct {
	Date   string `json:"date,omitempty"`
	Status string `json:"status"`
	Note   string `json:"note,omitempty"`
}
