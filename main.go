package main

import (
	"fmt"
	"log"

	"chefskiss-backend/config"
	"chefskiss-backend/models"
	"chefskiss-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Memulai server Chef's Kiss...")

	// 1. Koneksi Database
	config.ConnectDatabase()

	// 2. Jalankan Migrasi Database
	fmt.Println("Menjalankan migrasi database...")
	err := config.DB.AutoMigrate(
		&models.Menu{},
		&models.Order{},
		&models.OrderItem{},
	)
	
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}
	fmt.Println("Migrasi tabel sukses!")

	// 3. Inisialisasi Router Gin
	r := gin.Default()

	// 4. Daftarkan semua routes
	routes.SetupRoutes(r)

	// 5. Jalankan Server di port 8080
	fmt.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}