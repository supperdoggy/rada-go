package parser

import (
	"context"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/supperdoggy/vr_api/schema"
)

var voteSummaryPattern = regexp.MustCompile(`(?i)(За|Проти|Утрималися|Утрималось|Не голосували|Не голосувало|Всього|Рішення)\s*[:\-]?\s*(\d+|Прийнято|Не прийнято)`)

type VotingResultsHTMLParser struct {
	loader DocumentLoader
}

func NewVotingResultsHTMLParser(loader DocumentLoader) *VotingResultsHTMLParser {
	return &VotingResultsHTMLParser{loader: loader}
}

func (p *VotingResultsHTMLParser) Parse(ctx context.Context, rawHTML []byte) ([]schema.VoteSummaryItem, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	doc, err := p.loader.Load(rawHTML)
	if err != nil {
		return nil, err
	}

	items := p.parseRows(doc)
	if len(items) > 0 {
		return items, nil
	}

	return parseVoteSummaryFromText(normalizeWhitespace(doc.Text())), nil
}

func (p *VotingResultsHTMLParser) parseRows(doc *goquery.Document) []schema.VoteSummaryItem {
	items := make([]schema.VoteSummaryItem, 0)
	seen := make(map[string]struct{})

	doc.Find("tr").Each(func(_ int, rowSel *goquery.Selection) {
		cells := rowSel.Find("th,td")
		if cells.Length() < 2 {
			return
		}

		label := normalizeVoteLabel(cleanSelectionText(cells.First()))
		value := normalizeWhitespace(cleanSelectionText(cells.Eq(1)))
		if label == "" || value == "" {
			return
		}
		if !voteLabelAllowed(label) {
			return
		}

		if _, ok := seen[label]; ok {
			return
		}
		seen[label] = struct{}{}
		items = append(items, schema.VoteSummaryItem{
			Label: label,
			Value: value,
		})
	})

	return items
}

func parseVoteSummaryFromText(text string) []schema.VoteSummaryItem {
	if text == "" {
		return nil
	}

	matches := voteSummaryPattern.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil
	}

	items := make([]schema.VoteSummaryItem, 0, len(matches))
	seen := make(map[string]struct{})
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		label := normalizeVoteLabel(match[1])
		if !voteLabelAllowed(label) {
			continue
		}
		if _, ok := seen[label]; ok {
			continue
		}

		seen[label] = struct{}{}
		items = append(items, schema.VoteSummaryItem{
			Label: label,
			Value: normalizeWhitespace(match[2]),
		})
	}
	return items
}

func normalizeVoteLabel(label string) string {
	switch strings.ToLower(normalizeWhitespace(label)) {
	case "за":
		return "За"
	case "проти":
		return "Проти"
	case "утримались", "утрималися":
		return "Утрималися"
	case "утрималось":
		return "Утрималось"
	case "не голосувало":
		return "Не голосувало"
	case "не голосували":
		return "Не голосували"
	case "всього":
		return "Всього"
	case "рішення":
		return "Рішення"
	default:
		return ""
	}
}

func voteLabelAllowed(label string) bool {
	return label != ""
}
