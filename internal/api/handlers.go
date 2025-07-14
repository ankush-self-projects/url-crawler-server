package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddURL(c echo.Context) error {
	return c.JSON(http.StatusOK, "Add URL placeholder")
}

func GetURLs(c echo.Context) error {
	return c.JSON(http.StatusOK, "Get URLs placeholder")
}

func StartCrawl(c echo.Context) error {
	return c.JSON(http.StatusOK, "Start Crawl placeholder")
}
