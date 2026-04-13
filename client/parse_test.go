package client

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseSearchHTML_Golden(t *testing.T) {
	c := NewClient()
	html := readFixture(t, "search_v1_fixture.html")

	got, err := c.ParseSearchHTML(context.Background(), html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertGoldenJSON(t, got, "search_v1_fixture.golden.json")
}

func TestParseLawProjectHTML_Golden(t *testing.T) {
	c := NewClient()
	html := readFixture(t, "law_v1_fixture.html")

	got, err := c.ParseLawProjectHTML(context.Background(), html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertGoldenJSON(t, got, "law_v1_fixture.golden.json")
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "testdata", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	return content
}

func assertGoldenJSON(t *testing.T, payload any, goldenName string) {
	t.Helper()

	got, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	goldenPath := filepath.Join("..", "testdata", goldenName)
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v", goldenPath, err)
	}

	if !bytes.Equal(bytes.TrimSpace(got), bytes.TrimSpace(want)) {
		t.Fatalf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", goldenPath, got, want)
	}
}
