package routes

import (
	"chefskiss-backend/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes berfungsi mendaftarkan semua endpoint API
func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Fitur order
		api.POST("/orders", controllers.CreateOrder)

		// Rating
		api.POST("/ratings", controllers.SubmitRating)
		api.GET("/ratings/average/:menuId", controllers.GetAverageRating)
	}
}