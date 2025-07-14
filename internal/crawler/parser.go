package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func detectHTMLVersion(htmlStr string) string {
	if strings.Contains(strings.ToLower(htmlStr), "<!doctype html>") {
		return "HTML5"
	}
	return "Older HTML"
}

func parseDocument(htmlStr string) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
}

func extractTitle(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("title").Text())
}

func extractHeadingSummary(doc *goquery.Document) string {
	h1Count := doc.Find("h1").Length()
	h2Count := doc.Find("h2").Length()
	h3Count := doc.Find("h3").Length()
	return fmt.Sprintf("H1: %d, H2: %d, H3: %d", h1Count, h2Count, h3Count)
}

func detectLoginForm(doc *goquery.Document) bool {
	hasLogin := false
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		if s.Find(`input[type="password"]`).Length() > 0 {
			hasLogin = true
		}
	})
	return hasLogin
}
