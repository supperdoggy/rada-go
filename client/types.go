package client

import "github.com/supperdoggy/vr_api/schema"

type (
	SearchParams       = schema.SearchQuery
	SearchQuery        = schema.SearchQuery
	SearchResponse     = schema.SearchResponse
	SearchItem         = schema.SearchItem
	LawProjectDetails  = schema.LawProjectDetails
	DateEntry          = schema.DateEntry
	Registration       = schema.Registration
	Act                = schema.Act
	TextLink           = schema.TextLink
	Initiator          = schema.Initiator
	Committee          = schema.Committee
	Document           = schema.Document
	RelatedBill        = schema.RelatedBill
	TimelineEvent      = schema.TimelineEvent
	VotingResults      = schema.VotingResults
	VoteSummaryItem    = schema.VoteSummaryItem
	BillVotingResults  = schema.BillVotingResults
	BillVoteEvent      = schema.BillVoteEvent
	FactionVoteSummary = schema.FactionVoteSummary
	PersonVote         = schema.PersonVote
)
