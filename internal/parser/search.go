package parser

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/supperdoggy/vr_api/internal/apperr"
	"github.com/supperdoggy/vr_api/internal/profiles"
	"github.com/supperdoggy/vr_api/schema"
)

type SearchHTMLParser struct {
	loader    DocumentLoader
	selectors profiles.SearchSelectorSet
}

func NewSearchHTMLParser(loader DocumentLoader, selectors profiles.SearchSelectorSet) *SearchHTMLParser {
	return &SearchHTMLParser{
		loader:    loader,
		selectors: selectors,
	}
}

func (p *SearchHTMLParser) Parse(ctx context.Context, rawHTML []byte) (schema.SearchResponse, error) {
	if err := ctx.Err(); err != nil {
		return schema.SearchResponse{}, err
	}

	doc, err := p.loader.Load(rawHTML)
	if err != nil {
		return schema.SearchResponse{}, err
	}

	container := doc.Find(p.selectors.ResultsContainer()).First()
	if container.Length() == 0 {
		return schema.SearchResponse{}, apperr.NewLayoutChangedError(
			"parser.search",
			"results container selector not found",
			nil,
		)
	}

	items := make([]schema.SearchItem, 0)
	container.Find(p.selectors.ResultRows()).Each(func(_ int, rowSel *goquery.Selection) {
		linkSel := rowSel.Find(p.selectors.TitleLink()).First()
		url, _ := linkSel.Attr(p.selectors.LinkURLAttr())
		items = append(items, schema.SearchItem{
			ID:               cleanText(rowSel, p.selectors.ID()),
			Title:            cleanText(rowSel, p.selectors.Title()),
			Status:           cleanText(rowSel, p.selectors.Status()),
			RegistrationDate: cleanText(rowSel, p.selectors.RegistrationDate()),
			Subject:          cleanText(rowSel, p.selectors.Subject()),
			URL:              strings.TrimSpace(url),
		})
	})

	return schema.SearchResponse{
		Items: items,
		Count: len(items),
	}, nil
}
