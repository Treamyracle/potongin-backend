package middleware

import (
	"net/http"
	"potongin/database"
	"potongin/models"
	"potongin/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memeriksa JWT (Bearer Token)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header authorization kosong"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format header authorization salah"})
			return
		}

		tokenString := parts[1]
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		// Set userID di context agar bisa dipakai di handler
		c.Set("userID", userID)
		c.Next()
	}
}

// APIKeyMiddleware memeriksa X-API-Key
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Jika tidak ada API Key, coba fallback ke AuthMiddleware (JWT)
			AuthMiddleware()(c)
			return
		}

		var user models.User
		if err := database.DB.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API Key tidak valid"})
			return
		}

		// Set userID di context
		c.Set("userID", user.ID)
		c.Next()
	}
}
