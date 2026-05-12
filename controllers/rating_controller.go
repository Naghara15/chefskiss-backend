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

func GetAverageRating(c *gin.Context) {
	// Ambil menu_id dari URL parameter (contoh: /api/ratings/average/m1)
	menuID := c.Param("menuId")

	// Kita buat struct sementara untuk menampung hasil perhitungan dari database
	var result struct {
		Average float64 `json:"average"`
		Count   int64   `json:"total_reviews"`
	}

	// Gunakan fungsi agregasi bawaan database (AVG dan COUNT) untuk performa maksimal
	err := config.DB.Model(&models.MenuRating{}).
		Select("COALESCE(AVG(review), 0) as average, COUNT(id) as count").
		Where("menu_id = ?", menuID).
		Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung rata-rata rating"})
		return
	}

	// Berikan response lengkap
	c.JSON(http.StatusOK, gin.H{
		"menu_id": menuID,
		"data":    result,
	})
}

func GetTopRatings(c *gin.Context) {
	var ratings []models.MenuRating

	// 1. Query ke database dengan filter dan randomizer
	err := config.DB.Preload("Menu"). // Ambil detail menu terkait
		Where("review = ?", 5).        // Hanya bintang 5
		Where("note <> '' AND note IS NOT NULL"). // Note tidak boleh kosong atau null
		Order("RANDOM()").             // Fungsi random bawaan PostgreSQL
		Limit(5).                      // Batasi maksimal 5 data
		Find(&ratings).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data testimoni"})
		return
	}

	// 2. Berikan response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(ratings),
		"data":   ratings,
	})
}