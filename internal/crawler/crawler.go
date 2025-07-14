package crawler

import (
	"fmt"
	"url-crawler-backend/internal/model"
)

func CrawlURL(u *model.URL) error {
	bodyStr, err := fetchHTML(u.URL)
	if err != nil {
		u.Status = "error"
		return fmt.Errorf("fetch error: %w", err)
	}

	u.HTMLVersion = detectHTMLVersion(bodyStr)

	doc, err := parseDocument(bodyStr)
	if err != nil {
		u.Status = "error"
		return fmt.Errorf("parse error: %w", err)
	}

	u.PageTitle = extractTitle(doc)
	u.Headings = extractHeadingSummary(doc)

	baseURL := u.URL
	internal, external, broken := analyzeLinks(doc, baseURL)
	u.InternalLinks = internal
	u.ExternalLinks = external
	u.BrokenLinks = broken

	u.HasLoginForm = detectLoginForm(doc)

	u.Status = "done"
	return nil
}
