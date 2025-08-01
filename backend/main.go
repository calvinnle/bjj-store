package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/calvinnle/bjj-store/backend/config"
	"github.com/calvinnle/bjj-store/backend/handlers"
	"github.com/calvinnle/bjj-store/backend/middleware"
	"github.com/calvinnle/bjj-store/backend/models"
	"github.com/calvinnle/bjj-store/backend/services"
)

// @title BJJ Store API
// @version 1.0
// @description A complete e-commerce API for Brazilian Jiu-Jitsu gear and equipment
// @termsOfService http://swagger.io/terms/

// @contact.name Calvin Le
// @contact.url http://www.bjjstore.com/support
// @contact.email admin@bjjstore.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration with Viper
	config.LoadConfig()

	// Connect to database
	models.ConnectDatabase()

	// Auto migrate database
	models.AutoMigrate()

	// Initialize MinIO
	if err := services.InitMinIO(); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Setup Gin router
	if config.AppConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// API documentation endpoint (placeholder)
	r.GET("/docs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "BJJ Store API Documentation",
			"version": "1.0",
			"endpoints": gin.H{
				"health": "GET /api/health",
				"products": "GET /api/products",
				"orders": "POST /api/orders",
				"admin": "POST /api/admin/auth/login",
			},
		})
	})

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "BJJ Store API is running",
			"env":     config.AppConfig.Server.Environment,
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public product routes (no authentication needed)
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)

		// Public order routes (for customers)
		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders/track/:orderNumber", handlers.TrackOrder)
		api.GET("/orders/email/:email", handlers.GetOrdersByEmail)

		// Payment routes (public)
		api.POST("/payment/process", handlers.ProcessPayment)

		// Admin authentication routes (no auth required for login)
		adminAuth := api.Group("/admin/auth")
		{
			adminAuth.POST("/login", handlers.AdminLogin)
		}

		// Protected admin routes (JWT required)
		adminAPI := api.Group("/admin", middleware.JWTAuthMiddleware())
		{
			// Admin profile and session management
			adminAPI.POST("/logout", handlers.AdminLogout)
			adminAPI.GET("/profile", handlers.GetAdminProfile)
			adminAPI.GET("/stats", handlers.GetAdminStats)

			// Product management (require product permissions)
			adminAPI.POST("/products", middleware.RequirePermission("create_products"), handlers.CreateProduct)
			adminAPI.PUT("/products/:id", middleware.RequirePermission("update_products"), handlers.UpdateProduct)
			adminAPI.DELETE("/products/:id", middleware.RequirePermission("delete_products"), handlers.DeleteProduct)

			// Order management (require order permissions)
			adminAPI.GET("/orders", middleware.RequirePermission("view_orders"), handlers.GetAllOrders)
			adminAPI.PUT("/orders/:id/status", middleware.RequirePermission("update_orders"), handlers.UpdateOrderStatus)

			// Image upload (require product permissions)
			adminAPI.POST("/upload/image", middleware.RequirePermission("create_products"), handlers.UploadImage)
			adminAPI.DELETE("/upload/image", middleware.RequirePermission("delete_products"), handlers.DeleteImage)
		}
	}

	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	log.Printf("Environment: %s", config.AppConfig.Server.Environment)
	log.Printf("Default admin: %s", config.AppConfig.Admin.DefaultEmail)
	log.Printf("Swagger documentation: http://localhost:%s/swagger/index.html", config.AppConfig.Server.Port)
	r.Run(":" + config.AppConfig.Server.Port)
}
