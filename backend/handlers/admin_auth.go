package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/calvinnle/bjj-store/backend/config"
	"github.com/calvinnle/bjj-store/backend/models"
	"github.com/calvinnle/bjj-store/backend/services"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Success bool              `json:"success"`
	Token   string            `json:"token"`
	Admin   AdminUserResponse `json:"admin"`
	Message string            `json:"message"`
}

type AdminUserResponse struct {
	ID        uint             `json:"id"`
	Email     string           `json:"email"`
	Role      models.AdminRole `json:"role"`
	IsActive  bool             `json:"is_active"`
	LastLogin *time.Time       `json:"last_login"`
}

// AdminLogin godoc
// @Summary Admin login
// @Description Authenticate admin user and get JWT token
// @Tags admin,auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Admin login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/auth/login [post]
func AdminLogin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Find admin user
	var admin models.AdminUser
	if err := models.DB.Where("email = ? AND is_active = true", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	// Check password
	if !admin.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	// Create JWT service with config
	jwtService := services.NewJWTService(config.AppConfig.JWT.Secret)

	// Generate JWT token
	token, err := jwtService.GenerateToken(&admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate token",
		})
		return
	}

	// Store session
	if err := jwtService.StoreSession(admin.ID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create session",
		})
		return
	}

	// Update last login
	now := time.Now()
	admin.LastLogin = &now
	models.DB.Save(&admin)

	// Return response
	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Token:   token,
		Admin: AdminUserResponse{
			ID:        admin.ID,
			Email:     admin.Email,
			Role:      admin.Role,
			IsActive:  admin.IsActive,
			LastLogin: admin.LastLogin,
		},
		Message: "Login successful",
	})
}

// AdminLogout godoc
// @Summary Admin logout
// @Description Logout admin user and revoke JWT token
// @Tags admin,auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "No authorization header"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/logout [post]
func AdminLogout(c *gin.Context) {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No authorization header",
		})
		return
	}

	// Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Create JWT service
	jwtService := services.NewJWTService(config.AppConfig.JWT.Secret)

	// Revoke token
	if err := jwtService.RevokeToken(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// GetAdminProfile godoc
// @Summary Get admin profile
// @Description Get current admin user profile
// @Tags admin,auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Admin profile"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Admin not found"
// @Security BearerAuth
// @Router /admin/profile [get]
func GetAdminProfile(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var admin models.AdminUser
	if err := models.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Admin not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"admin": AdminUserResponse{
			ID:        admin.ID,
			Email:     admin.Email,
			Role:      admin.Role,
			IsActive:  admin.IsActive,
			LastLogin: admin.LastLogin,
		},
	})
}

// GetAdminStats godoc
// @Summary Get admin dashboard statistics
// @Description Get overview statistics for admin dashboard
// @Tags admin,dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Dashboard statistics"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Insufficient permissions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /admin/stats [get]
func GetAdminStats(c *gin.Context) {
	var stats struct {
		TotalProducts int64          `json:"total_products"`
		TotalOrders   int64          `json:"total_orders"`
		TotalRevenue  float64        `json:"total_revenue"`
		PendingOrders int64          `json:"pending_orders"`
		RecentOrders  []models.Order `json:"recent_orders"`
	}

	// Count products
	models.DB.Model(&models.Product{}).Count(&stats.TotalProducts)

	// Count orders
	models.DB.Model(&models.Order{}).Count(&stats.TotalOrders)

	// Calculate total revenue
	var result struct {
		Revenue float64
	}
	models.DB.Model(&models.Order{}).
		Where("status != ?", models.OrderStatusCancelled).
		Select("COALESCE(SUM(total_amount), 0) as revenue").
		Scan(&result)
	stats.TotalRevenue = result.Revenue

	// Count pending orders
	models.DB.Model(&models.Order{}).
		Where("status = ?", models.OrderStatusPending).
		Count(&stats.PendingOrders)

	// Get recent orders
	models.DB.Preload("Items.Product").
		Order("created_at desc").
		Limit(5).
		Find(&stats.RecentOrders)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats":   stats,
	})
}
