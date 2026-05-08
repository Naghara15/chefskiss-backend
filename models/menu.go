package models

// menus
type Menu struct {
	ID    		string  `gorm:"primaryKey;size:50" json:"id"`
	Name  		string  `gorm:"size:255;not null" json:"name"`
	Desc		string `json:"desc"`
	OldPrice 	float64 `json:"oldPrice"`
	Price 		float64 `gorm:"not null" json:"price"`
	Badge 		string `json:"badge"`
	Image 		string `json:"image"`

}