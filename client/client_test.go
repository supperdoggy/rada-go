package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
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
	lawHTML := readFixtureFile(t, "law_v1_fixture.html")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bill/123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(lawHTML)
	}))
	defer srv.Close()

	c := New(srv.URL)
	resp, err := c.Get("123")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if resp.ID != "123" {
		t.Fatalf("unexpected id: %q", resp.ID)
	}
	if resp.SourceURL != srv.URL+"/bill/123" {
		t.Fatalf("unexpected source url: %q", resp.SourceURL)
	}
	if len(resp.Documents) != 1 || resp.Documents[0].URL != srv.URL+"/docs/123/table.pdf" {
		t.Fatalf("unexpected document url: %+v", resp.Documents)
	}
}

func TestLawProjectDetails_BackwardCompatible(t *testing.T) {
	lawHTML := readFixtureFile(t, "law_v1_fixture.html")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(lawHTML)
	}))
	defer srv.Close()

	c := New(srv.URL)
	resp, err := c.LawProjectDetails(context.Background(), "123")
	if err != nil {
		t.Fatalf("law project details: %v", err)
	}
	if resp.ID != "123" {
		t.Fatalf("unexpected id: %q", resp.ID)
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
