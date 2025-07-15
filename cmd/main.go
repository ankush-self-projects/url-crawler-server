package main

import (
	"log"
	"url-crawler-backend/internal/api"
	"url-crawler-backend/internal/db"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	db.Connect()

	e := echo.New()

	api.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
