package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lenguyenkhang/bjj-store/config"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if the origin is allowed
		allowed := false
		for _, allowedOrigin := range config.AppConfig.CORS.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(config.AppConfig.CORS.AllowedOrigins) > 0 {
			// Default to first allowed origin if request origin is not allowed
			c.Header("Access-Control-Allow-Origin", config.AppConfig.CORS.AllowedOrigins[0])
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
