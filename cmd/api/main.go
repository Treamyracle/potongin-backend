package main

import (
	"potongin/config"
	"potongin/database"
	"potongin/routes"
)

func main() {
	// 1. Load Config dulu
	config.LoadConfig()

	// 2. Konek DB
	database.Connect()

	// 3. Setup Router
	r := routes.SetupRouter()

	r.Run(":8080")
}