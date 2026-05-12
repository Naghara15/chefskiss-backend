package controllers

import (
	"net/http"

	"chefskiss-backend/config"
	"chefskiss-backend/models"

	"github.com/gin-gonic/gin"
)

// SubmitRating menangani proses penyimpanan review dari user
func SubmitRating(c *gin.Context) {
	var input models.MenuRating

	// 1. Tangkap JSON dari Frontend
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data review tidak valid. Pastikan review berisi angka 1-5."})
		return
	}

	// 2. Validasi apakah MenuID benar-benar ada di database kita
	var menu models.Menu
	if err := config.DB.First(&menu, "id = ?", input.MenuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu dengan ID " + input.MenuID + " tidak ditemukan"})
		return
	}

	// 3. Simpan Review ke Database
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan review ke database"})
		return
	}

	// 4. Berikan Response Sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Terima kasih! Review kamu berhasil diunggah.",
		"data":    input,
	})
}