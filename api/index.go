package handler

import (
	"net/http"
	"sync" // 1. Import package sync

	"potongin/config"
	"potongin/database"
	"potongin/routes"

	"github.com/gin-gonic/gin"
)

var (
	once   sync.Once      // 2. Variabel untuk memastikan init hanya jalan sekali
	router *gin.Engine
)

// initApp menggantikan fungsi init() bawaan Go.
// Fungsi ini tidak akan jalan otomatis, tapi dipanggil manual lewat once.Do()
func initApp() {
	// 1. Load Config
	config.LoadConfig()

	// 2. Konek database
	database.Connect()

	// 3. Setup router
	router = routes.SetupRouter()
}

// Handler adalah entry point Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// 3. Panggil initApp menggunakan sync.Once
	// Ini menjamin database & router hanya di-load 1 kali saja per instance,
	// aman dari race condition (rebutan proses) saat banyak request masuk bersamaan.
	once.Do(initApp)

	router.ServeHTTP(w, r)
}