package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/supperdoggy/vr_api/internal/apperr"
)

func TestParserErrorTyping(t *testing.T) {
	err := apperr.NewLayoutChangedError("parser.search", "selector missing", nil)
	if !errors.Is(err, ErrLayoutChanged) {
		t.Fatalf("expected ErrLayoutChanged match, got %v", err)
	}
	if errors.Is(err, ErrMissingField) {
		t.Fatalf("expected no ErrMissingField match, got %v", err)
	}
}

func TestParserErrorMessage(t *testing.T) {
	err := apperr.NewMissingFieldError("parser.loader", "html payload is empty", nil)
	msg := err.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
	if !strings.Contains(msg, "missing_field") {
		t.Fatalf("expected code in message, got %q", msg)
	}
}
