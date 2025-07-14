package api

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api")
	api.POST("/urls", AddURL)
	api.GET("/urls", GetURLs)
	api.POST("/urls/:id/start", StartCrawl)
}
