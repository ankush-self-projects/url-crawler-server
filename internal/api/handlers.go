package api

import (
	"net/http"
	"net/url"
	"time"

	"url-crawler-backend/internal/crawler"
	"url-crawler-backend/internal/db"
	"url-crawler-backend/internal/model"

	"github.com/labstack/echo/v4"
)

type AddURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func AddURL(c echo.Context) error {
	req := new(AddURLRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if req.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL is required")
	}

	parsed, err := url.ParseRequestURI(req.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please include http:// or https:// in your URL.")
	}

	urlRecord := model.URL{
		URL:    req.URL,
		Status: "queued",
	}

	if err := db.DB.Create(&urlRecord).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save URL")
	}

	return c.JSON(http.StatusCreated, urlRecord)
}

func GetURLs(c echo.Context) error {
	var urls []model.URL

	if err := db.DB.Find(&urls).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch URLs")
	}

	return c.JSON(http.StatusOK, urls)
}

func StartBulkCrawl(c echo.Context) error {
	var req DeleteURLsRequest // reuse the struct with IDs []uint
	if err := c.Bind(&req); err != nil || len(req.IDs) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: must provide non-empty 'ids' array")
	}

	var notFound []uint
	for _, id := range req.IDs {
		var urlRecord model.URL
		if err := db.DB.First(&urlRecord, id).Error; err != nil {
			notFound = append(notFound, id)
			continue
		}
		urlRecord.Status = "running"
		db.DB.Save(&urlRecord)
		go func(urlModel model.URL) {
			err := crawler.CrawlURL(&urlModel)
			if err != nil {
				urlModel.Status = "error"
			} else {
				urlModel.Status = "done"
			}
			urlModel.UpdatedAt = time.Now()
			db.DB.Save(&urlModel)
		}(urlRecord)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Bulk crawl started",
		"not_found": notFound,
	})
}

type DeleteURLsRequest struct {
	IDs []uint `json:"ids"`
}

func DeleteURLs(c echo.Context) error {
	var req DeleteURLsRequest
	if err := c.Bind(&req); err != nil || len(req.IDs) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: must provide non-empty 'ids' array")
	}
	if err := db.DB.Delete(&model.URL{}, req.IDs).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete URLs")
	}
	return c.NoContent(http.StatusNoContent)
}
