package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedOrigins = []string{
	"http://localhost:5500",
	"http://127.0.0.1:5500",
	"http://localhost:3000",
	"http://127.0.0.1:3000",
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	// Cho phép mọi localhost trong môi trường dev
	if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") {
		return true
	}
	return false
}

// CORSMiddleware xử lý preflight OPTIONS và header CORS cho frontend (:5500 → :8080)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
