package database

import (
	"log"
	"potongin/config"
	"potongin/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect terhubung ke database dan melakukan auto-migrasi
func Connect() {
	var err error

	// Load config dulu
	config.LoadConfig()

	DB, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	log.Println("Koneksi database berhasil.")

	// Auto-migrasi skema
	err = DB.AutoMigrate(&models.User{}, &models.Shortener{})
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}
	log.Println("Migrasi database berhasil.")
}
