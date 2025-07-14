package main

import (
	"log"

	"url-crawler-backend/internal/db"
	"url-crawler-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	db.Connect()

	var existing model.User
	result := db.DB.Where("username = ?", "admin").First(&existing)
	if result.Error == nil {
		log.Println("User 'admin' already exists. Skipping seed.")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Create the user
	user := model.User{
		Username: "admin",
		Password: string(hashed),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Fatal("Failed to create user:", err)
	}

	log.Println("✅ Seeded user: admin / testpassword")
}
