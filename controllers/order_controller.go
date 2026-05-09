package controllers

import (
	"net/http"
	"time"

	"chefskiss-backend/config"
	"chefskiss-backend/models"

	"github.com/gin-gonic/gin"
)

func isValidPickupDate(date time.Time) bool {
	weekday := date.Weekday()

	if weekday == time.Tuesday || weekday == time.Wednesday || weekday == time.Thursday {
		return true
	}

	return false
}

func CreateOrder(c *gin.Context) {
	var input models.Order

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data yang dikirim tidak lengkap"})
		return
	}

	// Cek hari untuk order
	if !isValidPickupDate(input.PickupDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Hari pengambilan tidak valid! Pre-order hanya tersedia untuk hari Selasa, Rabu, dan Kamis.",
		})
		return
	}

	var total float64

	// Lakukan perulangan untuk mengecek setiap item pesanan
	for i, item := range input.Items {
		var trueMenu models.Menu

		// Cari data menu di database berdasarkan ID yang dikirim client (misal: "m1")
		if err := config.DB.First(&trueMenu, "id = ?", item.MenuID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Menu dengan ID " + item.MenuID + " tidak ditemukan"})
			return
		}

		// OVERRIDE: Timpa data dari client dengan data asli dari database
		input.Items[i].Price = trueMenu.Price
		input.Items[i].MenuName = trueMenu.Name

		// Hitung subtotal menggunakan harga ASLI
		total += trueMenu.Price * float64(item.Quantity)
	}

	input.TotalPrice = total

	// Simpan ke database
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesanan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pesanan berhasil dicatat dengan aman!",
		"data":    input,
	})
}