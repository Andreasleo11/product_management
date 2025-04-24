package routes

import (
	"backend_prodman/controllers"
	"backend_prodman/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// ✅ CORS config yang bener (dari dokumentasi resmi)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"}, // HTML-mu dari sini kan?
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Router untuk products
	router.GET("/products", middlewares.AuthAdmin(false), controllers.GetProducts)
	router.GET("/products/:id", middlewares.AuthAdmin(true), controllers.GetProductByID)
	router.GET("/products/category/:category_id", controllers.GetProductsByCategoryID)
	router.GET("/products/search", controllers.SearchProducts)
	router.POST("/products", middlewares.AuthAdmin(true), controllers.CreateProduct)
	router.PUT("/products/:id", controllers.UpdateProduct)
	router.DELETE("/products/:id", middlewares.AuthAdmin(true), controllers.DeleteProduct)

	//Router untuk categories
	router.GET("/categories", controllers.GetCategories)
	router.GET("/categories/:id", controllers.GetCategoryByID)
	router.POST("/categories", controllers.CreateCategory)
	router.PUT("/categories/:id", controllers.UpdateCategory)
	router.DELETE("/categories/:id", controllers.DeleteCategory)

	//Router untuk users
	router.POST("/users/register", controllers.RegisterUser)
	router.GET("/users", controllers.GetUsers)
	router.POST("/login", controllers.LoginUser)

	//Router untuk stocklog
	router.GET("/stock-logs", controllers.GetAllStockLogs) // Endpoint baru untuk mengambil semua log stok
	router.PUT("/products/:id/stock", middlewares.AuthAdmin(true), controllers.UpdateStock)

	return router
}
