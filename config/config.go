package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DatabaseURL  string
	JWTSecret    string
	ResendAPIKey string
	EmailSender  string
)

func LoadConfig() {
	// 1. Coba load .env (untuk Local Development)
	_ = godotenv.Load() 

	// 2. Baca dari Environment Variables (Vercel & Local yang sudah di-load)
	DatabaseURL = os.Getenv("DATABASE_NEON_DATABASE_URL")
	JWTSecret = os.Getenv("JWT_SECRET")        
	ResendAPIKey = os.Getenv("RESEND_API_KEY") 
	EmailSender = os.Getenv("EMAIL_SENDER")

// --- DEBUGGING BLOCK (Hapus nanti setelah fix) ---
	if DatabaseURL == "" {
		fmt.Println("=== DAFTAR ENV VARIABLE YANG DITERIMA SERVER ===")
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			// Hanya print KEY-nya saja agar aman, jangan print Value
			fmt.Println(pair[0]) 
		}
		fmt.Println("================================================")
		
		log.Fatal("Environment Variable DATABASE_URL tidak ditemukan (Cek logs di atas)")
	}
	// ------------------------------------------------
	if JWTSecret == "" {
		log.Fatal("Environment Variable JWT_SECRET tidak ditemukan")
	}
	// ... validasi lainnya
	
	log.Println("Konfigurasi berhasil dimuat dari Environment.")
}