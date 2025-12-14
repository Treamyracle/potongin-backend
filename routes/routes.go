package routes

import (
	"potongin/handlers"
	"potongin/middleware"
	"time" // <-- 1. IMPORT TAMBAHAN

	"github.com/gin-contrib/cors" // <-- 2. IMPORT TAMBAHAN
	"github.com/gin-gonic/gin"
)

// SetupRouter mengkonfigurasi semua rute aplikasi
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.RedirectTrailingSlash = false

	// --- 3. BLOK CORS DITAMBAHKAN DI SINI ---
	// Konfigurasi CORS agar frontend Anda (app.potong.in) diizinkan
	r.Use(cors.New(cors.Config{
		// Ganti "https://app.potong.in" dengan domain frontend Anda
		// "http://localhost:3000" adalah untuk testing frontend di komputer lokal
		AllowOrigins:     []string{"https://app.potong.in", "http://localhost:3000", "https://potongin-frontend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --- SELESAI BLOK CORS ---

	// Rute publik
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.GET("/verify-email/:token", handlers.VerifyEmail)

	// Rute redirect utama
	r.GET("/:code", handlers.Redirect)
	// Rute QR code publik
	r.GET("/qr/:code", handlers.GetQRCode)
	// Grup API v1
	apiV1 := r.Group("/api/v1")

	// Rute yang dilindungi (bisa pakai JWT atau API Key)
	authRoutes := apiV1.Group("/")
	authRoutes.Use(middleware.APIKeyMiddleware()) // Middleware ini cek API Key, jika tidak ada, fallback ke JWT
	{
		// Rute Link
		authRoutes.POST("/links", handlers.CreateShortLink)
		authRoutes.GET("/links", handlers.GetMyLinks)
		authRoutes.PUT("/links/:id", handlers.UpdateLink)
		authRoutes.DELETE("/links/:id", handlers.DeleteLink)

		// Rute User
		authRoutes.GET("/me", handlers.GetMe)
		authRoutes.POST("/me/apikey", handlers.GenerateAPIKey)
	}

	return r
}
