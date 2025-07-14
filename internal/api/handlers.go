package api

import (
	"net/http"
	"time"

	"url-crawler-backend/internal/crawler"
	"url-crawler-backend/internal/db"
	"url-crawler-backend/internal/model"

	"github.com/labstack/echo/v4"
)

// Request body for adding a new URL
type AddURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// AddURL handles POST /api/urls
func AddURL(c echo.Context) error {
	req := new(AddURLRequest)

	// Bind and validate input
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if req.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL is required")
	}

	// Create new URL record
	urlRecord := model.URL{
		URL:    req.URL,
		Status: "queued",
	}

	if err := db.DB.Create(&urlRecord).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save URL")
	}

	return c.JSON(http.StatusCreated, urlRecord)
}

// GetURLs handles GET /api/urls
func GetURLs(c echo.Context) error {
	var urls []model.URL

	if err := db.DB.Find(&urls).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch URLs")
	}

	return c.JSON(http.StatusOK, urls)
}

// StartCrawl handles POST /api/urls/:id/start
func StartCrawl(c echo.Context) error {
	id := c.Param("id")
	var urlRecord model.URL

	if err := db.DB.First(&urlRecord, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "URL not found")
	}

	// Update status to running
	urlRecord.Status = "running"
	db.DB.Save(&urlRecord)

	// Run the crawl in a goroutine so we don't block the API
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

	return c.JSON(http.StatusOK, map[string]string{"message": "Crawl started"})
}
