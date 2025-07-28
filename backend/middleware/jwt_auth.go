package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/calvinnle/bjj-store/backend/config"
    "github.com/calvinnle/bjj-store/backend/services"
    "github.com/calvinnle/bjj-store/backend/models"
)

func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
                "hint": "Add 'Authorization: Bearer <token>' header",
            })
            c.Abort()
            return
        }
        
        // Extract token (remove "Bearer " prefix)
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization format",
                "hint": "Use 'Bearer <token>' format",
            })
            c.Abort()
            return
        }
        
        // Create JWT service with config
        jwtService := services.NewJWTService(config.AppConfig.JWT.Secret)
        
        // Validate token
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
                "details": err.Error(),
            })
            c.Abort()
            return
        }
        
        // Check if token is revoked
        if jwtService.IsTokenRevoked(tokenString) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token has been revoked",
            })
            c.Abort()
            return
        }
        
        // Store admin info in context
        c.Set("admin_id", claims.AdminID)
        c.Set("admin_email", claims.Email)
        c.Set("admin_role", claims.Role)
        
        c.Next()
    }
}

func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        adminRole, exists := c.Get("admin_role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Admin role not found in context",
            })
            c.Abort()
            return
        }
        
        role := adminRole.(models.AdminRole)
        
        // Create temporary admin user to check permissions
        admin := models.AdminUser{Role: role}
        
        if !admin.HasPermission(permission) {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Insufficient permissions",
                "required_permission": permission,
                "your_role": role,
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
