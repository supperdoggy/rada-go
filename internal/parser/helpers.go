package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func cleanText(selection *goquery.Selection, selector string) string {
	return strings.TrimSpace(selection.Find(selector).First().Text())
}
