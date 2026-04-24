package client

import (
	"encoding/json"
	"testing"
)

func TestSearchContractJSONFields(t *testing.T) {
	payload := SearchResponse{
		Items: []SearchItem{
			{
				ID:                 "57706",
				RegistrationNumber: "0001",
				Title:              "Проект Закону про ратифікацію Угоди",
				Status:             "",
				RegistrationDate:   "09.09.2019",
				InitiativeSubject:  "Кабінет Міністрів України",
				Subject:            "Кабінет Міністрів України",
				URL:                "https://itd.rada.gov.ua/billinfo/Bills/Card/57706",
			},
		},
		Count:      14891,
		Page:       1,
		PerPage:    30,
		TotalPages: 497,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"items":[{"id":"57706","registrationNumber":"0001","title":"Проект Закону про ратифікацію Угоди","status":"","registrationDate":"09.09.2019","initiativeSubject":"Кабінет Міністрів України","subject":"Кабінет Міністрів України","url":"https://itd.rada.gov.ua/billinfo/Bills/Card/57706"}],"count":14891,"page":1,"perPage":30,"totalPages":497}`
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

func TestBillVotingResultsContractJSONFields(t *testing.T) {
	payload := BillVotingResults{
		BillID:             "57706",
		RegistrationNumber: "0001",
		SourceURL:          "https://w2.rada.gov.ua/pls/radan_gs09/ns_zakon_gol_dep_wohf?zn=0001",
		Votes: []BillVoteEvent{
			{
				GID:      "1315",
				Title:    "Поіменне голосування про проект Закону (№0001) - в цілому",
				DateTime: "29.10.2019 16:35",
				URL:      "https://w1.c1.rada.gov.ua/pls/radan_gs09/ns_golos?g_id=1315",
				RTFURL:   "https://w2.rada.gov.ua/pls/radan_gs09/ns_golos_rtf?g_id=1315&vid=1",
				PrintURL: "https://w2.rada.gov.ua/pls/radan_gs09/ns_golos_print?g_id=1315&vid=1",
				Decision: "Рішення прийнято",
				Summary: []VoteSummaryItem{
					{Label: "За", Value: "310"},
				},
				FactionSummary: []FactionVoteSummary{
					{
						Name:    `Фракція політичної партії "СЛУГА НАРОДУ"`,
						Members: "252",
						Summary: []VoteSummaryItem{
							{Label: "За", Value: "214"},
						},
					},
				},
				People: []PersonVote{
					{
						Name:      "Новинський Вадим Владиславович",
						DeputyID:  "40",
						Faction:   `Фракція політичної партії "ОПОЗИЦІЙНА ПЛАТФОРМА-ЗА ЖИТТЯ"`,
						Status:    "Відсутній",
						RawStatus: "Відсутній",
					},
				},
			},
		},
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"billId":"57706","registrationNumber":"0001","sourceURL":"https://w2.rada.gov.ua/pls/radan_gs09/ns_zakon_gol_dep_wohf?zn=0001","votes":[{"gId":"1315","title":"Поіменне голосування про проект Закону (№0001) - в цілому","dateTime":"29.10.2019 16:35","url":"https://w1.c1.rada.gov.ua/pls/radan_gs09/ns_golos?g_id=1315","rtfURL":"https://w2.rada.gov.ua/pls/radan_gs09/ns_golos_rtf?g_id=1315\u0026vid=1","printURL":"https://w2.rada.gov.ua/pls/radan_gs09/ns_golos_print?g_id=1315\u0026vid=1","decision":"Рішення прийнято","summary":[{"label":"За","value":"310"}],"factionSummary":[{"name":"Фракція політичної партії \"СЛУГА НАРОДУ\"","members":"252","summary":[{"label":"За","value":"214"}]}],"people":[{"name":"Новинський Вадим Владиславович","deputyId":"40","faction":"Фракція політичної партії \"ОПОЗИЦІЙНА ПЛАТФОРМА-ЗА ЖИТТЯ\"","status":"Відсутній","rawStatus":"Відсутній"}]}]}`
	if string(raw) != want {
		t.Fatalf("unexpected json contract\n got: %s\nwant: %s", raw, want)
	}
}
