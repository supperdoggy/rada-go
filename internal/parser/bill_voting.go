package parser

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/supperdoggy/vr_api/schema"
	"golang.org/x/net/html"
)

var (
	detailedVoteSummaryPattern = regexp.MustCompile(`(?i)(За|Проти|Утримались|Утрималися|Утрималось|Не голосували|Не голосувало|Всього|Відсутні|Відсутній)\s*[:\-]?\s*(\d+)`)
	factionMembersPattern      = regexp.MustCompile(`Кількість депутатів:\s*(\d+)`)
)

type ChronologySitting struct {
	Date string
	NomS string
}

type PlenaryVoteRef struct {
	GID   string
	Title string
	URL   string
}

type ChronologyVotingHTMLParser struct {
	loader DocumentLoader
}

func NewChronologyVotingHTMLParser(loader DocumentLoader) *ChronologyVotingHTMLParser {
	return &ChronologyVotingHTMLParser{loader: loader}
}

func (p *ChronologyVotingHTMLParser) Parse(ctx context.Context, rawHTML []byte) ([]ChronologySitting, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	doc, err := p.loader.Load(rawHTML)
	if err != nil {
		return nil, err
	}

	items := make([]ChronologySitting, 0)
	seen := make(map[string]struct{})

	doc.Find(`a[href*="ns_pd2?"]`).Each(func(_ int, linkSel *goquery.Selection) {
		href, ok := linkSel.Attr("href")
		if !ok {
			return
		}

		day := queryValueFromURL(href, "day_")
		month := queryValueFromURL(href, "month_")
		year := queryValueFromURL(href, "year")
		nomS := queryValueFromURL(href, "nom_s")
		if day == "" || month == "" || year == "" || nomS == "" {
			return
		}

		date := fmt.Sprintf("%02s%02s%s", day, month, year)
		key := date + "|" + nomS
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		items = append(items, ChronologySitting{
			Date: date,
			NomS: nomS,
		})
	})

	return items, nil
}

type PlenaryBillVotingHTMLParser struct {
	loader DocumentLoader
}

func NewPlenaryBillVotingHTMLParser(loader DocumentLoader) *PlenaryBillVotingHTMLParser {
	return &PlenaryBillVotingHTMLParser{loader: loader}
}

func (p *PlenaryBillVotingHTMLParser) Parse(ctx context.Context, registrationNumber string, rawHTML []byte) ([]PlenaryVoteRef, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	doc, err := p.loader.Load(rawHTML)
	if err != nil {
		return nil, err
	}

	marker := "(№" + normalizeWhitespace(registrationNumber) + ")"
	items := make([]PlenaryVoteRef, 0)
	seen := make(map[string]struct{})

	doc.Find(`a[href*="ns_golos?g_id="]`).Each(func(_ int, linkSel *goquery.Selection) {
		href, ok := linkSel.Attr("href")
		if !ok {
			return
		}

		gid := queryValueFromURL(href, "g_id")
		if gid == "" {
			return
		}
		if _, ok := seen[gid]; ok {
			return
		}

		rowSel := linkSel.Closest("tr")
		title := ""
		if rowSel.Length() > 0 {
			cells := rowSel.Find("td")
			if cells.Length() > 1 {
				title = normalizeWhitespace(cleanSelectionText(cells.Eq(1)))
			}
		}
		if title == "" || !strings.Contains(title, marker) {
			return
		}

		seen[gid] = struct{}{}
		items = append(items, PlenaryVoteRef{
			GID:   gid,
			Title: title,
			URL:   strings.TrimSpace(href),
		})
	})

	return items, nil
}

type DetailedVotingResultsHTMLParser struct {
	loader DocumentLoader
}

func NewDetailedVotingResultsHTMLParser(loader DocumentLoader) *DetailedVotingResultsHTMLParser {
	return &DetailedVotingResultsHTMLParser{loader: loader}
}

func (p *DetailedVotingResultsHTMLParser) Parse(ctx context.Context, rawHTML []byte) (schema.BillVoteEvent, error) {
	if err := ctx.Err(); err != nil {
		return schema.BillVoteEvent{}, err
	}

	doc, err := p.loader.Load(rawHTML)
	if err != nil {
		return schema.BillVoteEvent{}, err
	}

	event := schema.BillVoteEvent{}
	headLines := selectionLines(doc.Find(".head_gol").First())
	if len(headLines) > 0 {
		event.Title = headLines[0]
	}
	if len(headLines) > 1 {
		event.DateTime = headLines[1]
	}
	if len(headLines) > 2 {
		event.Summary = parseDetailedVoteSummary(headLines[2])
	}
	if len(headLines) > 3 {
		event.Decision = headLines[3]
	}

	if rtfURL, ok := doc.Find(`a[href*="ns_golos_rtf?g_id="]`).First().Attr("href"); ok {
		event.RTFURL = strings.TrimSpace(rtfURL)
	}
	if printURL, ok := doc.Find(`a[href*="ns_golos_print?g_id="]`).First().Attr("href"); ok {
		event.PrintURL = strings.TrimSpace(printURL)
	}

	factionSummaries, factionByShortName := parseFactionSummaries(doc)
	event.FactionSummary = factionSummaries
	event.People = parseDetailedPeople(doc, factionByShortName)

	return event, nil
}

func parseFactionSummaries(doc *goquery.Document) ([]schema.FactionVoteSummary, map[string]string) {
	items := make([]schema.FactionVoteSummary, 0)
	factionByShortName := make(map[string]string)

	doc.Find(`li[id^="0idf"]`).Each(func(_ int, factionSel *goquery.Selection) {
		headerSel := factionSel.ChildrenFiltered(".frn").First()
		if headerSel.Length() == 0 {
			return
		}

		name := cleanSelectionText(headerSel.Find("b").First())
		if name == "" {
			name = cleanSelectionText(headerSel)
		}
		if name == "" {
			return
		}

		headerText := normalizeWhitespace(cleanSelectionText(headerSel))
		members := ""
		if matches := factionMembersPattern.FindStringSubmatch(headerText); len(matches) > 1 {
			members = matches[1]
		}

		summaryText := cleanSelectionText(headerSel.Find("p").First())
		items = append(items, schema.FactionVoteSummary{
			Name:    name,
			Members: members,
			Summary: parseDetailedVoteSummary(summaryText),
		})

		factionSel.ChildrenFiltered("ul.frd").First().Find("li").Each(func(_ int, personSel *goquery.Selection) {
			shortName := cleanText(personSel, ".dep")
			if shortName == "" {
				return
			}
			if _, exists := factionByShortName[shortName]; !exists {
				factionByShortName[shortName] = name
			}
		})
	})

	return items, factionByShortName
}

func parseDetailedPeople(doc *goquery.Document, factionByShortName map[string]string) []schema.PersonVote {
	type personAccumulator struct {
		shortName string
		fullName  string
		deputyID  string
		status    string
		rawStatus string
		faction   string
	}

	peopleByKey := make(map[string]*personAccumulator)
	order := make([]string, 0)

	doc.Find(`li[id^="0idd"]`).Each(func(_ int, personSel *goquery.Selection) {
		id, ok := personSel.Attr("id")
		if !ok {
			return
		}

		key := strings.TrimPrefix(strings.TrimSpace(id), "0")
		if key == "" {
			return
		}

		shortName := cleanText(personSel, ".dep")
		rawStatus := extractStatusText(cleanText(personSel, ".golos"))
		if shortName == "" || rawStatus == "" {
			return
		}

		acc := &personAccumulator{
			shortName: shortName,
			deputyID:  strings.TrimPrefix(key, "idd"),
			status:    normalizeDetailedVoteStatus(rawStatus),
			rawStatus: rawStatus,
		}
		if faction := factionByShortName[shortName]; faction != "" {
			acc.faction = faction
		}

		peopleByKey[key] = acc
		order = append(order, key)
	})

	doc.Find(`a[id^="1idd"]`).Each(func(_ int, linkSel *goquery.Selection) {
		id, ok := linkSel.Attr("id")
		if !ok {
			return
		}

		key := strings.TrimPrefix(strings.TrimSpace(id), "1")
		if key == "" {
			return
		}

		title, ok := linkSel.Attr("title")
		if !ok {
			return
		}

		rawStatus, fullName := parseHallMapTitle(title)
		if rawStatus == "" && fullName == "" {
			return
		}

		acc, ok := peopleByKey[key]
		if !ok {
			acc = &personAccumulator{
				deputyID: strings.TrimPrefix(key, "idd"),
			}
			peopleByKey[key] = acc
			order = append(order, key)
		}
		if fullName != "" {
			acc.fullName = fullName
		}
		if acc.deputyID == "" {
			acc.deputyID = strings.TrimPrefix(key, "idd")
		}
		if acc.rawStatus == "" && rawStatus != "" {
			acc.rawStatus = rawStatus
			acc.status = normalizeDetailedVoteStatus(rawStatus)
		}
	})

	if len(order) == 0 {
		fallback := make([]schema.PersonVote, 0)
		doc.Find(`li[id^="0idf"] ul.frd li`).Each(func(_ int, personSel *goquery.Selection) {
			shortName := cleanText(personSel, ".dep")
			rawStatus := extractStatusText(cleanText(personSel, ".golos"))
			if shortName == "" || rawStatus == "" {
				return
			}
			fallback = append(fallback, schema.PersonVote{
				Name:      shortName,
				Faction:   factionByShortName[shortName],
				Status:    normalizeDetailedVoteStatus(rawStatus),
				RawStatus: rawStatus,
			})
		})
		return fallback
	}

	people := make([]schema.PersonVote, 0, len(order))
	seen := make(map[string]struct{})
	for _, key := range order {
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		acc := peopleByKey[key]
		if acc == nil {
			continue
		}

		name := acc.fullName
		if name == "" {
			name = acc.shortName
		}
		if name == "" || acc.status == "" {
			continue
		}
		if acc.faction == "" && acc.shortName != "" {
			acc.faction = factionByShortName[acc.shortName]
		}

		people = append(people, schema.PersonVote{
			Name:      name,
			DeputyID:  acc.deputyID,
			Faction:   acc.faction,
			Status:    acc.status,
			RawStatus: acc.rawStatus,
		})
	}

	return people
}

func selectionLines(selection *goquery.Selection) []string {
	if selection == nil || selection.Length() == 0 {
		return nil
	}

	lines := make([]string, 0)
	var buffer strings.Builder

	flush := func() {
		line := normalizeWhitespace(buffer.String())
		buffer.Reset()
		if line != "" {
			lines = append(lines, line)
		}
	}

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node == nil {
			return
		}

		if node.Type == html.TextNode {
			buffer.WriteString(node.Data)
		}
		if node.Type == html.ElementNode && node.Data == "br" {
			flush()
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	for _, node := range selection.Nodes {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	flush()

	return lines
}

func parseDetailedVoteSummary(text string) []schema.VoteSummaryItem {
	if text == "" {
		return nil
	}

	matches := detailedVoteSummaryPattern.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil
	}

	items := make([]schema.VoteSummaryItem, 0, len(matches))
	seen := make(map[string]struct{})
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		label := normalizeDetailedSummaryLabel(match[1])
		if label == "" {
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

func normalizeDetailedSummaryLabel(label string) string {
	switch strings.ToLower(normalizeWhitespace(label)) {
	case "за":
		return "За"
	case "проти":
		return "Проти"
	case "утримались", "утрималися", "утрималось":
		return "Утрималися"
	case "не голосували", "не голосувало":
		return "Не голосували"
	case "всього":
		return "Всього"
	case "відсутні", "відсутній":
		return "Відсутні"
	default:
		return ""
	}
}

func normalizeDetailedVoteStatus(status string) string {
	switch strings.ToLower(normalizeWhitespace(status)) {
	case "за":
		return "За"
	case "проти":
		return "Проти"
	case "утримався", "утрималась", "утримались", "утрималися", "утрималось":
		return "Утримався"
	case "не голосував", "не голосувала", "не голосували", "не голосувало":
		return "Не голосував"
	case "відсутній", "відсутня", "відсутні", "відсутнє":
		return "Відсутній"
	default:
		return normalizeWhitespace(status)
	}
}

func extractStatusText(text string) string {
	text = normalizeWhitespace(text)
	text = strings.TrimPrefix(text, `Голос - "`)
	text = strings.TrimPrefix(text, "Голос - ")
	text = strings.Trim(text, `"`)
	return normalizeWhitespace(text)
}

func parseHallMapTitle(title string) (string, string) {
	lines := strings.Split(strings.ReplaceAll(title, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return "", ""
	}

	status := extractStatusText(lines[0])
	fullName := ""
	if len(lines) > 1 {
		fullName = normalizeWhitespace(lines[1])
	}

	return status, fullName
}

func queryValueFromURL(rawURL, key string) string {
	idx := strings.Index(rawURL, "?")
	if idx < 0 || idx >= len(rawURL)-1 {
		return ""
	}

	for _, pair := range strings.Split(rawURL[idx+1:], "&") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 || parts[0] != key {
			continue
		}
		return strings.TrimSpace(parts[1])
	}
	return ""
}
