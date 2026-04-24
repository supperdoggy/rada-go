package parser

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestChronologyVotingHTMLParser_Parse(t *testing.T) {
	raw := readParserFixture(t, "chronology_voting_fixture.html")

	got, err := NewChronologyVotingHTMLParser(NewGoqueryLoader()).Parse(context.Background(), raw)
	if err != nil {
		t.Fatalf("parse chronology voting html: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 unique sittings, got %d", len(got))
	}
	if got[0].Date != "29102019" || got[0].NomS != "3" {
		t.Fatalf("unexpected first sitting: %+v", got[0])
	}
	if got[1].Date != "30102019" || got[1].NomS != "2" {
		t.Fatalf("unexpected second sitting: %+v", got[1])
	}
}

func TestPlenaryBillVotingHTMLParser_Parse(t *testing.T) {
	raw := readParserFixture(t, "plenary_votes_fixture.html")

	got, err := NewPlenaryBillVotingHTMLParser(NewGoqueryLoader()).Parse(context.Background(), "0001", raw)
	if err != nil {
		t.Fatalf("parse plenary vote html: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 bill vote refs, got %d", len(got))
	}
	if got[0].GID != "1314" {
		t.Fatalf("unexpected first g_id: %+v", got[0])
	}
	if got[1].GID != "1315" {
		t.Fatalf("unexpected second g_id: %+v", got[1])
	}
}

func TestDetailedVotingResultsHTMLParser_Parse(t *testing.T) {
	raw := readParserFixture(t, "bill_vote_details_fixture.html")

	got, err := NewDetailedVotingResultsHTMLParser(NewGoqueryLoader()).Parse(context.Background(), raw)
	if err != nil {
		t.Fatalf("parse detailed vote html: %v", err)
	}

	if got.Title == "" || got.DateTime != "29.10.2019 16:35" {
		t.Fatalf("unexpected vote header: %+v", got)
	}
	if got.Decision != "Рішення прийнято" {
		t.Fatalf("unexpected decision: %q", got.Decision)
	}
	if len(got.Summary) != 5 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
	if len(got.FactionSummary) != 2 {
		t.Fatalf("unexpected faction summary: %+v", got.FactionSummary)
	}
	if len(got.People) != 3 {
		t.Fatalf("unexpected people count: %+v", got.People)
	}
	if got.People[0].Name != "Новинський Вадим Владиславович" {
		t.Fatalf("expected full name enrichment, got %+v", got.People[0])
	}
	if got.People[0].DeputyID != "40" {
		t.Fatalf("expected deputy id join, got %+v", got.People[0])
	}
	if got.People[0].Faction != `Фракція політичної партії "ОПОЗИЦІЙНА ПЛАТФОРМА-ЗА ЖИТТЯ"` {
		t.Fatalf("unexpected faction mapping: %+v", got.People[0])
	}
	if got.People[2].Status != "Не голосував" || got.People[2].RawStatus != "Не голосувала" {
		t.Fatalf("unexpected status normalization: %+v", got.People[2])
	}
}

func readParserFixture(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join("..", "..", "testdata", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	return content
}
