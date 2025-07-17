package api

import (
	"url-crawler-backend/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/login", Login)

	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware())

	api.POST("/urls", AddURL)
	api.GET("/urls", GetURLs)
	api.POST("/urls/:id/start", StartCrawl)
	api.DELETE("/urls", DeleteURLs)
}
