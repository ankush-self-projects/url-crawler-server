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
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB: %v", err))
	}

	// 🔑 Auto-migrate your tables here!
	if err := connection.AutoMigrate(&model.URL{}); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	DB = connection
}
