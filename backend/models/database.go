package models

import (
	"fmt"
	"log"

	"github.com/calvinnle/bjj-store/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	log.Println("Database connected successfully!")
}

func AutoMigrate() {
	err := DB.AutoMigrate(&Product{}, &Order{}, &OrderItem{}, &AdminUser{}, &AdminSession{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed!")

	// Create default admin user if none exists
	createDefaultAdmin()
}

func createDefaultAdmin() {
	var count int64
	DB.Model(&AdminUser{}).Count(&count)

	if count == 0 {
		admin := AdminUser{
			Email:    config.AppConfig.Admin.DefaultEmail,
			Role:     RoleSuperAdmin,
			IsActive: true,
		}

		if err := admin.HashPassword(config.AppConfig.Admin.DefaultPassword); err != nil {
			log.Printf("Failed to hash default admin password: %v", err)
			return
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("Failed to create default admin: %v", err)
			return
		}

		log.Printf("Default admin created: %s", admin.Email)
	}
}
