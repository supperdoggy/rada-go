package parser

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/supperdoggy/vr_api/internal/apperr"
)

type DocumentLoader interface {
	Load(rawHTML []byte) (*goquery.Document, error)
}

type GoqueryLoader struct{}

func NewGoqueryLoader() GoqueryLoader {
	return GoqueryLoader{}
}

func (GoqueryLoader) Load(rawHTML []byte) (*goquery.Document, error) {
	if len(bytes.TrimSpace(rawHTML)) == 0 {
		return nil, apperr.NewMissingFieldError("parser.loader", "html payload is empty", nil)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return nil, apperr.NewLayoutChangedError("parser.loader", "unable to parse html document", err)
	}
	return doc, nil
}
