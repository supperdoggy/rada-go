package client

import (
	"encoding/json"
	"testing"
)

func TestSearchContractJSONFields(t *testing.T) {
	payload := SearchResponse{
		Items: []SearchItem{
			{
				ID:               "123",
				Title:            "Draft Law",
				Status:           "Registered",
				RegistrationDate: "2026-04-13",
				Subject:          "Public Safety",
				URL:              "https://itd.rada.gov.ua/bill/123",
			},
		},
		Count: 1,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"items":[{"id":"123","title":"Draft Law","status":"Registered","registrationDate":"2026-04-13","subject":"Public Safety","url":"https://itd.rada.gov.ua/bill/123"}],"count":1}`
	if string(raw) != want {
		t.Fatalf("unexpected json contract\n got: %s\nwant: %s", raw, want)
	}
}

func TestLawProjectContractJSONFields(t *testing.T) {
	payload := LawProjectDetails{
		ID:        "123",
		Title:     "Draft Law 123",
		Status:    "In committee",
		SourceURL: "https://itd.rada.gov.ua/bill/123",
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"id":"123","title":"Draft Law 123","status":"In committee","dates":null,"initiators":null,"committees":null,"documents":null,"timeline":null,"sourceURL":"https://itd.rada.gov.ua/bill/123"}`
	if string(raw) != want {
		t.Fatalf("unexpected json contract\n got: %s\nwant: %s", raw, want)
	}
}

func TestLawProjectContractJSONFields_Expanded(t *testing.T) {
	payload := LawProjectDetails{
		ID:     "57707",
		Title:  "Проект Постанови",
		Status: "Постанову підписано",
		Registration: &Registration{
			Number:   "1231/П",
			Date:     "19.09.2019",
			Archived: true,
		},
		Act: &Act{
			Text: "126-IX",
			Date: "20.09.2019",
			URL:  "https://zakon.rada.gov.ua/go/126-IX",
		},
		RegistrationSession: "2 сесія IX скликання",
		Category:            "Соціальна політика",
		InitiativeSubject:   "Народний депутат України",
		MainCommittee: &TextLink{
			Text: "Комітет",
		},
		RelatedBills: []RelatedBill{
			{
				Number:           "1231",
				RegistrationDate: "02.09.2019",
				Title:            "Пов'язаний законопроект",
				URL:              "https://itd.rada.gov.ua/billinfo/Bills/Card/57708",
			},
		},
		VotingResults: &VotingResults{
			URL: "https://w2.rada.gov.ua/vote",
			Summary: []VoteSummaryItem{
				{Label: "За", Value: "252"},
			},
		},
		SourceURL: "https://itd.rada.gov.ua/billinfo/Bills/Card/57707",
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"id":"57707","title":"Проект Постанови","status":"Постанову підписано","dates":null,"registration":{"number":"1231/П","date":"19.09.2019","archived":true},"act":{"text":"126-IX","date":"20.09.2019","url":"https://zakon.rada.gov.ua/go/126-IX"},"registrationSession":"2 сесія IX скликання","category":"Соціальна політика","initiativeSubject":"Народний депутат України","initiators":null,"mainCommittee":{"text":"Комітет"},"committees":null,"documents":null,"relatedBills":[{"number":"1231","registrationDate":"02.09.2019","title":"Пов'язаний законопроект","url":"https://itd.rada.gov.ua/billinfo/Bills/Card/57708"}],"timeline":null,"votingResults":{"url":"https://w2.rada.gov.ua/vote","summary":[{"label":"За","value":"252"}]},"sourceURL":"https://itd.rada.gov.ua/billinfo/Bills/Card/57707"}`
	if string(raw) != want {
		t.Fatalf("unexpected json contract\n got: %s\nwant: %s", raw, want)
	}
}
