package handlers

import (
	"net/http"
	"potongin/database"
	"potongin/models"

	"github.com/gin-gonic/gin"
)

// Redirect akan mengarahkan short_code ke URL asli
func Redirect(c *gin.Context) {
	shortCode := c.Param("code")

	var link models.Shortener
	if err := database.DB.Where("short_code = ?", shortCode).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link tidak ditemukan"})
		return
	}

	// Update jumlah klik (bisa dijalankan di goroutine agar tidak memblokir)
	go func() {
		database.DB.Model(&link).UpdateColumn("clicks", link.Clicks+1)
	}()

	// Redirect ke URL asli
	c.Redirect(http.StatusFound, link.OriginalURL)
}
