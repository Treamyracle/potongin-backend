package handlers

import (
	"fmt"
	"net/http"
	"potongin/config"
	"potongin/database"
	"potongin/models"
	"potongin/utils"

	// Import time
	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2" // Import Resend
	"gorm.io/gorm"
)

type AuthInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Signup(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hash password"})
		return
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate API key"})
		return
	}

	user := models.User{
		Email:      input.Email,
		Password:   hashedPassword,
		APIKey:     apiKey,
		IsVerified: false,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Buat token verifikasi
	verificationToken, err := utils.GenerateVerificationToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token verifikasi"})
		return
	}

	// --- PERUBAHAN KRUSIAL DI SINI ---
	// Kirim Email Verifikasi SECARA SINKRON (HAPUS 'go')
	// Sekarang kita menunggu email selesai dikirim
	err = sendVerificationEmail(user.Email, verificationToken)
	if err != nil {
		// Jika email gagal dikirim, beri tahu frontend
		// Error sudah di-log di dalam sendVerificationEmail
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Berhasil mendaftar, tetapi gagal mengirim email verifikasi. Coba lagi nanti."})
		return
	}
	// --- SELESAI PERUBAHAN ---

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil. Silakan cek email Anda untuk verifikasi."})
}

func Login(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun Anda belum diverifikasi. Silakan cek email Anda."})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// HANDLER BARU
func VerifyEmail(c *gin.Context) {
	tokenString := c.Param("token")

	userID, err := utils.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token verifikasi tidak valid atau kedaluwarsa."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terverifikasi sebelumnya."})
		return
	}

	database.DB.Model(&user).Update("is_verified", true)

	c.Redirect(http.StatusFound, "https://app.potong.in?verification=success")
}

// FUNGSI HELPER BARU
// --- UBAH FUNGSI INI AGAR MENGEMBALIKAN ERROR ---
func sendVerificationEmail(toEmail, token string) error {
	client := resend.NewClient(config.ResendAPIKey)

	verificationLink := fmt.Sprintf("https://potong.in/verify-email/%s", token)

	subject := "Verifikasi Email Anda - Potong.in"
	htmlBody := fmt.Sprintf(`
		<h1>Selamat datang di Potong.in!</h1>
		<p>Terima kasih telah mendaftar. Silakan klik link di bawah ini untuk memverifikasi email Anda:</p>
		<p><a href="%s" target="_blank">Verifikasi Email Saya</a></p>
		<p>Link ini hanya valid selama 15 menit.</p>
	`, verificationLink)

	params := &resend.SendEmailRequest{
		From:    config.EmailSender,
		To:      []string{toEmail},
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		// Log error ini di Vercel
		fmt.Printf("Gagal mengirim email verifikasi ke %s: %s\n", toEmail, err.Error())
		return err // <-- KEMBALIKAN ERROR
	}

	return nil // <-- KEMBALIKAN SUKSES
}
