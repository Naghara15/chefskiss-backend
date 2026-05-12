package models

import "time"

type MenuRating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MenuID    string    `gorm:"size:50;not null" json:"menu_id" binding:"required"`
	Menu     Menu    `gorm:"foreignKey:MenuID;references:ID" json:"-"` 
	// Angka rating (misal: 1 sampai 5)
	Review    int       `gorm:"not null" json:"review" binding:"required,min=1,max=5"`
	// Catatan/Komentar dari user (bisa kosong, jadi tidak perlu not null atau required)
	Note      string    `gorm:"type:text" json:"note"`
	CreatedAt time.Time `json:"created_at"`
}