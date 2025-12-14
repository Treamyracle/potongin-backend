package models

import "gorm.io/gorm"

// User sesuai dengan tabel 'users'
type User struct {
	gorm.Model
	Email      string      `gorm:"unique;not null;index" json:"email"` // GANTI DARI USERNAME
	Password   string      `gorm:"not null" json:"-"`
	APIKey     string      `gorm:"unique;index" json:"api_key"`
	IsVerified bool        `gorm:"default:false" json:"is_verified"` // TAMBAHKAN INI
	Shorteners []Shortener `gorm:"foreignKey:UserID" json:"shorteners,omitempty"`
}

// Shortener sesuai dengan tabel 'shorteners'
type Shortener struct {
	gorm.Model
	OriginalURL string `gorm:"not null" json:"original_url"`
	ShortCode   string `gorm:"unique;not null;index" json:"short_code"`
	UserID      uint   `gorm:"not null" json:"user_id"`
	Clicks      uint   `gorm:"default:0" json:"clicks"`
}

// Catatan: Kolom 'qr' tidak disimpan di DB.
// QR code akan di-generate secara on-the-fly berdasarkan ShortCode.
