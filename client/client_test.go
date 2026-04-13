package client

import (
	"context"
	"errors"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient()
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

func TestFutureFetchMethods_NotImplemented(t *testing.T) {
	c := NewClient()

	_, err := c.Search(context.Background(), SearchQuery{Term: "budget"})
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented for Search, got %v", err)
	}

	_, err = c.LawProjectDetails(context.Background(), "123")
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented for LawProjectDetails, got %v", err)
	}
}
