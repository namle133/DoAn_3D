package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		log.Printf("[%d] %s %s - %v", statusCode, method, path, duration)
	}
}

func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()

		// Only log write operations
		if method != "GET" && method != "HEAD" {
			auditLog := models.AuditLog{
				ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
				UserID:    userID.(string),
				Action:    method,
				Entity:    path, // In real app, extract entity type from path
				IPAddress: ip,
				Timestamp: time.Now(),
			}

			database.DB.Create(&auditLog)
		}

		c.Next()
	}
}
