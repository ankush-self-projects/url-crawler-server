package db

import (
	"fmt"
	"os"
	"url-crawler-backend/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Get database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Set default port if not provided
	if dbPort == "" {
		dbPort = "3306"
	}

	// Build the DSN (Data Source Name) with separate host and port
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB: %v", err))
	}

	// 🔑 Auto-migrate your tables here!
	if err := connection.AutoMigrate(&model.URL{}, &model.User{}); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	DB = connection
}
