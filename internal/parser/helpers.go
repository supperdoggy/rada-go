package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func cleanText(selection *goquery.Selection, selector string) string {
	return normalizeWhitespace(selection.Find(selector).First().Text())
}

func cleanSelectionText(selection *goquery.Selection) string {
	if selection == nil {
		return ""
	}

	return normalizeWhitespace(selection.First().Text())
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
