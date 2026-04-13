package parser

import (
	"context"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/supperdoggy/vr_api/internal/apperr"
	"github.com/supperdoggy/vr_api/internal/profiles"
	"github.com/supperdoggy/vr_api/schema"
)

var (
	registrationPattern = regexp.MustCompile(`^\s*(.+?)\s+від\s+(\d{2}\.\d{2}\.\d{4})(?:\s+\(([^)]+)\))?\s*$`)
	dateSuffixPattern   = regexp.MustCompile(`^(.*?)\s*\((\d{2}\.\d{2}\.\d{4})\)\s*$`)
	actDatePattern      = regexp.MustCompile(`від\s+(\d{2}\.\d{2}\.\d{4})`)
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

	title := cleanText(root, p.selectors.Title())
	if title == "" {
		title, _ = doc.Find(`meta[name="description"]`).First().Attr("content")
		title = normalizeWhitespace(title)
	}

	id, _ := root.Find(p.selectors.IDNode()).First().Attr("data-id")
	summaryRows := p.summaryRows(root)

	registration := parseRegistration(summaryRows["Номер, дата реєстрації"])
	act := parseAct(summaryRows["Номер, дата акту"])
	initiators := parseInitiators(summaryRows["Ініціатор(и) законопроекту"])
	mainCommittee := parseTextLink(summaryRows["Головний комітет"])
	otherCommittees := parseTextLinks(summaryRows["Інші комітети"])
	committees := buildCompatibilityCommittees(mainCommittee, otherCommittees)
	documents := parseDocuments(summaryRows["Текст законопроекту та супровідні документи"])
	timeline := p.parseTimeline(root)
	relatedBills := p.parseRelatedBills(root)

	status := cleanText(root, p.selectors.TimelineHeaderStatus())
	if status == "" && len(timeline) > 0 {
		status = timeline[0].Status
	}

	dates := buildCompatibilityDates(registration, act)
	chronologyURL := extractAttr(root, p.selectors.ChronologyEmbed(), "src")
	votingResultsURL := extractAttr(root, p.selectors.VotingResultsEmbed(), "src")

	var votingResults *schema.VotingResults
	if votingResultsURL != "" {
		votingResults = &schema.VotingResults{URL: strings.TrimSpace(votingResultsURL)}
	}

	return schema.LawProjectDetails{
		ID:                  strings.TrimSpace(id),
		Title:               title,
		Status:              status,
		Dates:               dates,
		Registration:        registration,
		Act:                 act,
		RegistrationSession: rowText(summaryRows["Сесія реєстрації"]),
		Category:            rowText(summaryRows["Рубрика законопроекту"]),
		InitiativeSubject:   rowText(summaryRows["Суб'єкт права законодавчої ініціативи"]),
		Initiators:          initiators,
		MainCommittee:       mainCommittee,
		OtherCommittees:     otherCommittees,
		Committees:          committees,
		Documents:           documents,
		RelatedBills:        relatedBills,
		Timeline:            timeline,
		ChronologyURL:       strings.TrimSpace(chronologyURL),
		VotingResults:       votingResults,
		SourceURL:           "",
	}, nil
}

func (p *LawProjectHTMLParser) summaryRows(root *goquery.Selection) map[string]*goquery.Selection {
	rows := make(map[string]*goquery.Selection)
	root.Find(p.selectors.SummaryRows()).Each(func(_ int, rowSel *goquery.Selection) {
		label := normalizeSummaryLabel(cleanText(rowSel, p.selectors.SummaryLabel()))
		if label == "" {
			return
		}

		valueSel := rowSel.Find(p.selectors.SummaryValue()).Last()
		if valueSel.Length() == 0 {
			return
		}
		rows[label] = valueSel
	})
	return rows
}

func (p *LawProjectHTMLParser) parseTimeline(root *goquery.Selection) []schema.TimelineEvent {
	timeline := make([]schema.TimelineEvent, 0)
	root.Find(p.selectors.TimelineRows()).Each(func(_ int, rowSel *goquery.Selection) {
		status := cleanText(rowSel, p.selectors.TimelineStatus())
		if status == "" {
			return
		}

		event := schema.TimelineEvent{
			Date:   cleanText(rowSel, p.selectors.TimelineDate()),
			Status: status,
		}
		if noteSelector := p.selectors.TimelineNote(); noteSelector != "" {
			event.Note = cleanText(rowSel, noteSelector)
		}
		timeline = append(timeline, event)
	})
	return timeline
}

func (p *LawProjectHTMLParser) parseRelatedBills(root *goquery.Selection) []schema.RelatedBill {
	relatedBills := make([]schema.RelatedBill, 0)
	root.Find(p.selectors.RelatedRows()).Each(func(_ int, rowSel *goquery.Selection) {
		numberLink := rowSel.Find(p.selectors.RelatedNumberLink()).First()
		title := cleanText(rowSel, p.selectors.RelatedTitle())
		if title == "" && numberLink.Length() == 0 {
			return
		}

		url, _ := numberLink.Attr(p.selectors.LinkURLAttr())
		relatedBills = append(relatedBills, schema.RelatedBill{
			Number:           cleanSelectionText(numberLink),
			RegistrationDate: cleanText(rowSel, p.selectors.RelatedRegistrationDate()),
			Title:            title,
			URL:              strings.TrimSpace(url),
		})
	})
	return relatedBills
}

func parseRegistration(selection *goquery.Selection) *schema.Registration {
	text := rowText(selection)
	if text == "" {
		return nil
	}

	matches := registrationPattern.FindStringSubmatch(text)
	if len(matches) == 0 {
		return &schema.Registration{Number: text}
	}

	return &schema.Registration{
		Number:   normalizeWhitespace(matches[1]),
		Date:     matches[2],
		Archived: strings.Contains(matches[3], "Архів"),
	}
}

func parseAct(selection *goquery.Selection) *schema.Act {
	if selection == nil || selection.Length() == 0 {
		return nil
	}

	text := rowText(selection)
	if text == "" {
		return nil
	}

	linkSel := selection.Find("a[href]").First()
	url, _ := linkSel.Attr("href")
	act := &schema.Act{
		Text: cleanSelectionText(linkSel),
		URL:  strings.TrimSpace(url),
	}

	if act.Text == "" {
		act.Text = strings.TrimSpace(actDatePattern.ReplaceAllString(text, ""))
	}

	if matches := actDatePattern.FindStringSubmatch(text); len(matches) > 0 {
		act.Date = matches[1]
	}

	if act.Text == "" && act.Date == "" && act.URL == "" {
		return nil
	}
	return act
}

func parseInitiators(selection *goquery.Selection) []schema.Initiator {
	if selection == nil || selection.Length() == 0 {
		return nil
	}

	initiators := make([]schema.Initiator, 0)
	selection.Find("a[href]").Each(func(_ int, linkSel *goquery.Selection) {
		url, _ := linkSel.Attr("href")
		text := cleanSelectionText(linkSel)
		if text == "" {
			return
		}
		initiators = append(initiators, schema.Initiator{
			Name: text,
			URL:  strings.TrimSpace(url),
		})
	})
	if len(initiators) > 0 {
		return initiators
	}

	text := rowText(selection)
	if text == "" {
		return nil
	}
	return []schema.Initiator{{Name: text}}
}

func parseTextLink(selection *goquery.Selection) *schema.TextLink {
	text := rowText(selection)
	if text == "" {
		return nil
	}

	linkSel := selection.Find("a[href]").First()
	url, _ := linkSel.Attr("href")
	return &schema.TextLink{
		Text: text,
		URL:  strings.TrimSpace(url),
	}
}

func parseTextLinks(selection *goquery.Selection) []schema.TextLink {
	if selection == nil || selection.Length() == 0 {
		return nil
	}

	links := make([]schema.TextLink, 0)
	selection.Find("a[href]").Each(func(_ int, linkSel *goquery.Selection) {
		text := cleanSelectionText(linkSel)
		if text == "" {
			return
		}
		url, _ := linkSel.Attr("href")
		links = append(links, schema.TextLink{
			Text: text,
			URL:  strings.TrimSpace(url),
		})
	})
	if len(links) > 0 {
		return links
	}

	text := rowText(selection)
	if text == "" {
		return nil
	}
	return []schema.TextLink{{Text: text}}
}

func parseDocuments(selection *goquery.Selection) []schema.Document {
	if selection == nil || selection.Length() == 0 {
		return nil
	}

	documents := make([]schema.Document, 0)
	selection.Find("a[href]").Each(func(_ int, linkSel *goquery.Selection) {
		url, _ := linkSel.Attr("href")
		title, date := splitTitleAndDate(cleanSelectionText(linkSel))
		if title == "" {
			return
		}

		documents = append(documents, schema.Document{
			Title: title,
			Date:  date,
			URL:   strings.TrimSpace(url),
		})
	})
	return documents
}

func buildCompatibilityCommittees(mainCommittee *schema.TextLink, otherCommittees []schema.TextLink) []schema.Committee {
	committees := make([]schema.Committee, 0, len(otherCommittees)+1)
	if mainCommittee != nil && mainCommittee.Text != "" {
		committees = append(committees, schema.Committee{
			Name: mainCommittee.Text,
			Role: "Головний комітет",
			URL:  mainCommittee.URL,
		})
	}
	for _, committee := range otherCommittees {
		if committee.Text == "" {
			continue
		}
		committees = append(committees, schema.Committee{
			Name: committee.Text,
			Role: "Інший комітет",
			URL:  committee.URL,
		})
	}
	return committees
}

func buildCompatibilityDates(registration *schema.Registration, act *schema.Act) []schema.DateEntry {
	dates := make([]schema.DateEntry, 0, 2)
	if registration != nil && (registration.Number != "" || registration.Date != "") {
		value := strings.TrimSpace(registration.Number)
		if registration.Date != "" {
			value = strings.TrimSpace(strings.Join([]string{value, "від", registration.Date}, " "))
		}
		dates = append(dates, schema.DateEntry{
			Label: "Номер, дата реєстрації",
			Value: value,
		})
	}
	if act != nil && (act.Text != "" || act.Date != "") {
		value := act.Text
		if act.Date != "" {
			value = strings.TrimSpace(strings.Join([]string{act.Text, "від", act.Date}, " "))
		}
		dates = append(dates, schema.DateEntry{
			Label: "Номер, дата акту",
			Value: value,
		})
	}
	if len(dates) == 0 {
		return nil
	}
	return dates
}

func splitTitleAndDate(text string) (string, string) {
	if text == "" {
		return "", ""
	}

	matches := dateSuffixPattern.FindStringSubmatch(text)
	if len(matches) == 0 {
		return text, ""
	}
	return normalizeWhitespace(matches[1]), matches[2]
}

func normalizeSummaryLabel(label string) string {
	label = strings.TrimSpace(strings.TrimSuffix(label, ":"))
	return normalizeWhitespace(label)
}

func rowText(selection *goquery.Selection) string {
	return cleanSelectionText(selection)
}

func extractAttr(selection *goquery.Selection, selector, attr string) string {
	if selector == "" {
		return ""
	}

	value, _ := selection.Find(selector).First().Attr(attr)
	return strings.TrimSpace(value)
}
