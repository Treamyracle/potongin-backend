package handlers

import (
	"net/http"
	"potongin/database"
	"potongin/models"
	"potongin/utils"

	"github.com/gin-gonic/gin"
)

// GetMe mendapatkan info user yang sedang login
func GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Kirim balik data user (tanpa password)
	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"email":       user.Email, // <-- PERBAIKAN DI SINI
		"api_key":     user.APIKey,
		"created_at":  user.CreatedAt,
		"is_verified": user.IsVerified, // <-- TAMBAHAN YANG BAGUS
	})
}

// GenerateAPIKey membuat atau me-regenerasi API key untuk user
// (Fungsi ini tidak perlu diubah, sudah benar)
func GenerateAPIKey(c *gin.Context) {
	userID, _ := c.Get("userID")

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate API key"})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("api_key", apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API Key berhasil di-generate",
		"api_key": apiKey,
	})
}
