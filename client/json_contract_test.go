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
