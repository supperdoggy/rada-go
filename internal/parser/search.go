package parser

import (
	"context"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/supperdoggy/vr_api/internal/apperr"
	"github.com/supperdoggy/vr_api/internal/profiles"
	"github.com/supperdoggy/vr_api/schema"
)

var digitsPattern = regexp.MustCompile(`\d+`)

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
		linkSel := rowSel.Find(p.selectors.RegistrationNumberLink()).First()
		resultURL, _ := linkSel.Attr(p.selectors.LinkURLAttr())
		resultURL = strings.TrimSpace(resultURL)

		initiativeSubject := cleanText(rowSel, p.selectors.InitiativeSubject())
		if linkSel.Length() == 0 && cleanText(rowSel, p.selectors.Title()) == "" {
			return
		}

		items = append(items, schema.SearchItem{
			ID:                 deriveCardID(resultURL),
			RegistrationNumber: cleanSelectionText(linkSel),
			Title:              cleanText(rowSel, p.selectors.Title()),
			Status:             "",
			RegistrationDate:   cleanText(rowSel, p.selectors.RegistrationDate()),
			InitiativeSubject:  initiativeSubject,
			Subject:            initiativeSubject,
			URL:                resultURL,
		})
	})

	return schema.SearchResponse{
		Items:      items,
		Count:      parseFirstInt(cleanText(container, p.selectors.Count())),
		Page:       parseInputValue(container, p.selectors.CurrentPage()),
		PerPage:    parseSelectedOptionValue(container, p.selectors.PerPage()),
		TotalPages: parseTotalPages(container.Find(p.selectors.TotalPages()).First()),
	}, nil
}

func deriveCardID(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	id := path.Base(parsed.Path)
	if id == "." || id == "/" {
		return ""
	}
	return strings.TrimSpace(id)
}

func parseFirstInt(text string) int {
	match := digitsPattern.FindString(text)
	if match == "" {
		return 0
	}

	value, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return value
}

func parseInputValue(selection *goquery.Selection, selector string) int {
	value, _ := selection.Find(selector).First().Attr("value")
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return parsed
}

func parseSelectedOptionValue(selection *goquery.Selection, selector string) int {
	option := selection.Find(selector).First().Find("option[selected]").First()
	if option.Length() == 0 {
		option = selection.Find(selector).First().Find("option").First()
	}
	value, _ := option.Attr("value")
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return parsed
}

func parseTotalPages(selection *goquery.Selection) int {
	maxPage := 0

	selection.Find("[data-page]").Each(func(_ int, itemSel *goquery.Selection) {
		value, _ := itemSel.Attr("data-page")
		page, err := strconv.Atoi(strings.TrimSpace(value))
		if err == nil && page > maxPage {
			maxPage = page
		}
	})
	selection.Find(".active").Each(func(_ int, itemSel *goquery.Selection) {
		page := parseFirstInt(cleanSelectionText(itemSel))
		if page > maxPage {
			maxPage = page
		}
	})

	return maxPage
}
