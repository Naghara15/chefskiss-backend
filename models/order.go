package models

import "time"

// orders
type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Customer   string      `gorm:"size:255;not null" json:"customer_name" binding:"required"`
	PickupDate string      `gorm:"size:50;not null" json:"pickup_date" binding:"required"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  time.Time   `json:"created_at"`
	
	// 1 orders punya banyak orderitems
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items" binding:"required"` 
}

// order_items
type OrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	OrderID  uint    `json:"order_id"` // Foreign key orders.ID
	MenuID   string  `gorm:"size:50;not null" json:"menu_id" binding:"required"`
	MenuName string  `gorm:"size:255;not null" json:"menu_name"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Price    float64 `json:"price" binding:"required"`
}