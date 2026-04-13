package profiles

import (
	"errors"
	"testing"

	"github.com/supperdoggy/vr_api/internal/apperr"
)

func TestRegistry_DefaultVersions(t *testing.T) {
	r := NewRegistry()
	if got := r.DefaultSearchVersion(); got != "v1" {
		t.Fatalf("unexpected default search version: %q", got)
	}
	if got := r.DefaultLawProjectVersion(); got != "v1" {
		t.Fatalf("unexpected default law version: %q", got)
	}
}

func TestRegistry_UnsupportedVersion(t *testing.T) {
	r := NewRegistry()
	_, err := r.Search("missing")
	if err == nil {
		t.Fatal("expected unsupported profile error")
	}
	if !errors.Is(err, apperr.ErrUnsupportedProfile) {
		t.Fatalf("expected ErrUnsupportedProfile, got %v", err)
	}
}

func TestRegistry_CustomSearchProfile(t *testing.T) {
	r := NewRegistry()
	r.RegisterSearch(mockSearchProfile{
		version: "v2",
		id:      "search/v2",
		sel:     searchV1Selectors{},
	})

	profile, err := r.Search("v2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.ID() != "search/v2" {
		t.Fatalf("unexpected profile id: %q", profile.ID())
	}
}

type mockSearchProfile struct {
	version string
	id      string
	sel     SearchSelectorSet
}

func (m mockSearchProfile) ID() string {
	return m.id
}

func (m mockSearchProfile) Version() string {
	return m.version
}

func (m mockSearchProfile) Selectors() SearchSelectorSet {
	return m.sel
}
