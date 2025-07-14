package main

import (
	"log"
	"url-crawler-backend/internal/api"
	"url-crawler-backend/internal/db"
	"url-crawler-backend/internal/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	db.Connect()

	e := echo.New()

	e.Use(middleware.JWTMiddleware())
	api.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
