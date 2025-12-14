package handlers

import (
	"fmt"
	"net/http"
	"potongin/database"
	"potongin/models"
	"potongin/utils"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type CreateLinkInput struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomCode  string `json:"custom_code,omitempty"` // opsional
}

type UpdateLinkInput struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomCode  string `json:"custom_code"` // Field baru (opsional)
}

// CreateShortLink membuat link pendek baru
func CreateShortLink(c *gin.Context) {
	var input CreateLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	var shortCode string
	if input.CustomCode != "" {
		// Cek apakah custom code sudah dipakai
		var existing models.Shortener
		if err := database.DB.Where("short_code = ?", input.CustomCode).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Custom link ini sudah dipakai"})
			return
		}
		shortCode = input.CustomCode
	} else {
		// Generate random code
		var err error
		for i := 0; i < 5; i++ { // Coba 5x jika terjadi kolisi
			shortCode, err = utils.GenerateShortCode()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate code"})
				return
			}
			// Pastikan code belum ada
			var existing models.Shortener
			if err := database.DB.Where("short_code = ?", shortCode).First(&existing).Error; err == gorm.ErrRecordNotFound {
				break // Code unik, bisa dipakai
			}
		}
	}

	link := models.Shortener{
		OriginalURL: input.OriginalURL,
		ShortCode:   shortCode,
		UserID:      userID.(uint),
	}

	if err := database.DB.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan link"})
		return
	}

	// Domain Anda (bisa juga diambil dari config)
	domain := "potong.in"

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Link berhasil dibuat",
		"short_link": fmt.Sprintf("https://%s/%s", domain, link.ShortCode),
		"data":       link,
	})
}

// GetMyLinks mendapatkan semua link milik user yang login
func GetMyLinks(c *gin.Context) {
	userID, _ := c.Get("userID")

	var links []models.Shortener
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data links"})
		return
	}

	c.JSON(http.StatusOK, links)
}

type UpdateLinkInput struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

// UpdateLink mengubah URL asli dari link pendek
func UpdateLink(c *gin.Context) {
	userID, _ := c.Get("userID")
	linkID := c.Param("id")

	// Cari link di database
	var link models.Shortener
	if err := database.DB.Where("id = ? AND user_id = ?", linkID, userID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link tidak ditemukan atau Anda tidak punya akses"})
		return
	}

	// Bind JSON input
	var input UpdateLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Selalu update URL asli
	link.OriginalURL = input.OriginalURL

	// Cek apakah user ingin mengubah Custom URL (Short Code)
	if input.CustomCode != "" && input.CustomCode != link.ShortCode {
		// Cek apakah custom code baru sudah dipakai orang lain?
		var existing models.Shortener
		if err := database.DB.Where("short_code = ?", input.CustomCode).First(&existing).Error; err == nil {
			// Jika err == nil, berarti ketemu (sudah ada yang pakai)
			c.JSON(http.StatusConflict, gin.H{"error": "Custom link ini sudah dipakai, cari yang lain"})
			return
		}
		
		// Jika aman, update short code
		link.ShortCode = input.CustomCode
	}

	// Simpan perubahan ke DB
	if err := database.DB.Save(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link berhasil diupdate", "data": link})
}

// DeleteLink menghapus link pendek
func DeleteLink(c *gin.Context) {
	userID, _ := c.Get("userID")
	linkID := c.Param("id")

	var link models.Shortener
	if err := database.DB.Where("id = ? AND user_id = ?", linkID, userID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link tidak ditemukan atau Anda tidak punya akses"})
		return
	}

	database.DB.Delete(&link)
	c.JSON(http.StatusOK, gin.H{"message": "Link berhasil dihapus"})
}

// GetQRCode membuat dan mengembalikan QR Code sebagai gambar PNG
func GetQRCode(c *gin.Context) {
	shortCode := c.Param("code")

	var link models.Shortener
	if err := database.DB.Where("short_code = ?", shortCode).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link tidak ditemukan"})
		return
	}

	// URL lengkap untuk QR Code
	// Kita ambil host dari request, misal "potong.in" atau "localhost:8080"
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s/%s", scheme, c.Request.Host, link.ShortCode)

	// Generate QR code
	png, err := qrcode.Encode(fullURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate QR code"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "image/png", png)
}
