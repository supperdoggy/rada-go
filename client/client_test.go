package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	c := New("")
	if c.baseURL != defaultBaseURL {
		t.Fatalf("expected default base url %q, got %q", defaultBaseURL, c.baseURL)
	}
	if c.searchProfileVersion != "" {
		t.Fatalf("expected empty search profile version override, got %q", c.searchProfileVersion)
	}
	if c.lawProjectProfileVersion != "" {
		t.Fatalf("expected empty law profile version override, got %q", c.lawProjectProfileVersion)
	}
}

func TestNew_UsesProvidedBaseURL(t *testing.T) {
	c := New("https://example.com", WithSearchProfileVersion("v2"), WithLawProjectProfileVersion("v3"))
	if c.baseURL != "https://example.com" {
		t.Fatalf("unexpected baseURL: %q", c.baseURL)
	}
	if c.searchProfileVersion != "v2" {
		t.Fatalf("unexpected search version: %q", c.searchProfileVersion)
	}
	if c.lawProjectProfileVersion != "v3" {
		t.Fatalf("unexpected law version: %q", c.lawProjectProfileVersion)
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	c := NewClient(
		WithBaseURL("https://example.com"),
		WithSearchProfileVersion("v2"),
		WithLawProjectProfileVersion("v3"),
	)
	if c.baseURL != "https://example.com" {
		t.Fatalf("unexpected baseURL: %q", c.baseURL)
	}
	if c.searchProfileVersion != "v2" {
		t.Fatalf("unexpected search version: %q", c.searchProfileVersion)
	}
	if c.lawProjectProfileVersion != "v3" {
		t.Fatalf("unexpected law version: %q", c.lawProjectProfileVersion)
	}
}

func TestSearch_FetchesAndParsesSearchResults(t *testing.T) {
	searchHTML := readFixtureFile(t, "search_v1_fixture.html")

	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(searchHTML)
	}))
	defer srv.Close()

	c := New(srv.URL)
	resp, err := c.Search(SearchParams{
		Term:    "budget",
		Page:    2,
		PerPage: 25,
		Filters: map[string]string{
			"status":  "registered",
			"session": "9",
		},
	})
	if err != nil {
		t.Fatalf("search: %v", err)
	}

	if gotQuery.Get("term") != "budget" {
		t.Fatalf("unexpected term query: %q", gotQuery.Get("term"))
	}
	if gotQuery.Get("page") != "2" {
		t.Fatalf("unexpected page query: %q", gotQuery.Get("page"))
	}
	if gotQuery.Get("perPage") != "25" {
		t.Fatalf("unexpected perPage query: %q", gotQuery.Get("perPage"))
	}
	if gotQuery.Get("status") != "registered" {
		t.Fatalf("unexpected status query: %q", gotQuery.Get("status"))
	}
	if gotQuery.Get("session") != "9" {
		t.Fatalf("unexpected session query: %q", gotQuery.Get("session"))
	}
	if resp.Count != 2 {
		t.Fatalf("expected 2 items, got %d", resp.Count)
	}
	if resp.Items[0].URL != srv.URL+"/bill/123" {
		t.Fatalf("unexpected normalized url: %q", resp.Items[0].URL)
	}
}

func TestSearchContext_UsesProvidedContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := New("https://example.com")
	_, err := c.SearchContext(ctx, SearchParams{Term: "budget"})
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestGet_FetchesAndParsesLawProject(t *testing.T) {
	lawHTML := withLocalVotingResults(readRepoFixtureFile(t, "57707.html"), "/vote/57707")
	votingHTML := readFixtureFile(t, "voting_results_v1_fixture.html")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/billinfo/Bills/Card/57707":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(lawHTML)
		case "/vote/57707":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(votingHTML)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	c := New(srv.URL)
	resp, err := c.Get("57707")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if resp.ID != "57707" {
		t.Fatalf("unexpected id: %q", resp.ID)
	}
	if resp.Title == "" {
		t.Fatal("expected parsed title")
	}
	if resp.SourceURL != srv.URL+"/billinfo/Bills/Card/57707" {
		t.Fatalf("unexpected source url: %q", resp.SourceURL)
	}
	if len(resp.Documents) != 1 || resp.Documents[0].URL != srv.URL+"/billinfo/Bills/pubFile/3129875" {
		t.Fatalf("unexpected document url: %+v", resp.Documents)
	}
	if resp.VotingResults == nil || resp.VotingResults.URL != srv.URL+"/vote/57707" {
		t.Fatalf("unexpected voting results: %+v", resp.VotingResults)
	}
	if len(resp.VotingResults.Summary) != 3 {
		t.Fatalf("unexpected voting summary: %+v", resp.VotingResults.Summary)
	}
}

func TestLawProjectDetails_BackwardCompatible(t *testing.T) {
	lawHTML := withLocalVotingResults(readRepoFixtureFile(t, "57707.html"), "/vote/57707")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/billinfo/Bills/Card/57707":
			_, _ = w.Write(lawHTML)
		case "/vote/57707":
			_, _ = w.Write([]byte("<html><body></body></html>"))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	c := New(srv.URL)
	resp, err := c.LawProjectDetails(context.Background(), "57707")
	if err != nil {
		t.Fatalf("law project details: %v", err)
	}
	if resp.ID != "57707" {
		t.Fatalf("unexpected id: %q", resp.ID)
	}
	if resp.VotingResults == nil {
		t.Fatal("expected voting results metadata")
	}
	if len(resp.VotingResults.Summary) != 0 {
		t.Fatalf("expected empty voting summary fallback, got %+v", resp.VotingResults.Summary)
	}
}

func TestGet_ReturnsHTTPStatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadGateway)
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.Get("123")
	if err == nil {
		t.Fatal("expected an error for non-2xx response")
	}
}

func readFixtureFile(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join("..", "testdata", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}

	return content
}

func readRepoFixtureFile(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join("..", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}

	return content
}

func withLocalVotingResults(raw []byte, localPath string) []byte {
	return []byte(strings.Replace(
		string(raw),
		"https://w2.rada.gov.ua/pls/radan_gs09/ns_zakon_gol_dep_wohf?zn=1231%2F%D0%9F",
		localPath,
		1,
	))
}
