package crawler

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func analyzeLinks(doc *goquery.Document, baseURL string) (int, int, int) {
	internalLinks := 0
	externalLinks := 0
	brokenLinks := 0

	base, _ := url.Parse(baseURL)

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			return
		}

		resolved := base.ResolveReference(linkURL)
		if resolved.Hostname() == base.Hostname() {
			internalLinks++
		} else {
			externalLinks++
		}

		linkResp, err := http.Head(resolved.String())
		if err != nil || linkResp.StatusCode >= 400 {
			brokenLinks++
		}
	})

	return internalLinks, externalLinks, brokenLinks
}
