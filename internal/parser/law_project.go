package parser

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/supperdoggy/vr_api/internal/apperr"
	"github.com/supperdoggy/vr_api/internal/profiles"
	"github.com/supperdoggy/vr_api/schema"
)

type LawProjectHTMLParser struct {
	loader    DocumentLoader
	selectors profiles.LawProjectSelectorSet
}

func NewLawProjectHTMLParser(loader DocumentLoader, selectors profiles.LawProjectSelectorSet) *LawProjectHTMLParser {
	return &LawProjectHTMLParser{
		loader:    loader,
		selectors: selectors,
	}
}

func (p *LawProjectHTMLParser) Parse(ctx context.Context, rawHTML []byte) (schema.LawProjectDetails, error) {
	if err := ctx.Err(); err != nil {
		return schema.LawProjectDetails{}, err
	}

	doc, err := p.loader.Load(rawHTML)
	if err != nil {
		return schema.LawProjectDetails{}, err
	}

	root := doc.Find(p.selectors.RootContainer()).First()
	if root.Length() == 0 {
		return schema.LawProjectDetails{}, apperr.NewLayoutChangedError(
			"parser.law_project",
			"root container selector not found",
			nil,
		)
	}

	dates := make([]schema.DateEntry, 0)
	root.Find(p.selectors.DatesRows()).Each(func(_ int, rowSel *goquery.Selection) {
		dates = append(dates, schema.DateEntry{
			Label: cleanText(rowSel, p.selectors.DateLabel()),
			Value: cleanText(rowSel, p.selectors.DateValue()),
		})
	})

	initiators := make([]schema.Initiator, 0)
	root.Find(p.selectors.InitiatorRows()).Each(func(_ int, rowSel *goquery.Selection) {
		initiators = append(initiators, schema.Initiator{
			Name: cleanText(rowSel, p.selectors.InitiatorName()),
			Role: cleanText(rowSel, p.selectors.InitiatorRole()),
		})
	})

	committees := make([]schema.Committee, 0)
	root.Find(p.selectors.CommitteeRows()).Each(func(_ int, rowSel *goquery.Selection) {
		committees = append(committees, schema.Committee{
			Name: cleanText(rowSel, p.selectors.CommitteeName()),
			Role: cleanText(rowSel, p.selectors.CommitteeRole()),
		})
	})

	documents := make([]schema.Document, 0)
	root.Find(p.selectors.DocumentRows()).Each(func(_ int, rowSel *goquery.Selection) {
		linkSel := rowSel.Find(p.selectors.DocumentLink()).First()
		url, _ := linkSel.Attr(p.selectors.DocumentURLAttr())
		documents = append(documents, schema.Document{
			Title: cleanText(rowSel, p.selectors.DocumentTitle()),
			Type:  cleanText(rowSel, p.selectors.DocumentType()),
			Date:  cleanText(rowSel, p.selectors.DocumentDate()),
			URL:   strings.TrimSpace(url),
		})
	})

	timeline := make([]schema.TimelineEvent, 0)
	root.Find(p.selectors.TimelineRows()).Each(func(_ int, rowSel *goquery.Selection) {
		timeline = append(timeline, schema.TimelineEvent{
			Date:   cleanText(rowSel, p.selectors.TimelineDate()),
			Status: cleanText(rowSel, p.selectors.TimelineStatus()),
			Note:   cleanText(rowSel, p.selectors.TimelineNote()),
		})
	})

	sourceURL, _ := root.Attr(p.selectors.SourceURLAttr())
	return schema.LawProjectDetails{
		ID:         cleanText(root, p.selectors.ID()),
		Title:      cleanText(root, p.selectors.Title()),
		Status:     cleanText(root, p.selectors.Status()),
		Dates:      dates,
		Initiators: initiators,
		Committees: committees,
		Documents:  documents,
		Timeline:   timeline,
		SourceURL:  strings.TrimSpace(sourceURL),
	}, nil
}
