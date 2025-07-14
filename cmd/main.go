package main

import (
	"url-crawler-backend/internal/api"
	"url-crawler-backend/internal/db"
	"url-crawler-backend/internal/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	db.Connect()

	e := echo.New()

	e.Use(middleware.JWTMiddleware())
	api.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
