package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
)

// TestHelper provides utilities for testing
type TestHelper struct {
	t *testing.T
}

func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// SetupTestDB initializes a test database
func (h *TestHelper) SetupTestDB() *config.Config {
	// Use test database
	cfg := &config.Config{
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "vinhomes_user",
		DBPassword:     "vinhomes_pass",
		DBName:         "vinhomes_test",
		ServerPort:     "8080",
		GinMode:        "test",
		JWTSecret:      "test-secret",
		JWTExpiryHours: 24,
	}

	if err := database.InitDatabase(cfg); err != nil {
		h.t.Fatalf("Failed to initialize test database: %v", err)
	}

	return cfg
}

// SetupTestRouter creates a test Gin router
func (h *TestHelper) SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TeardownTestDB cleans up test database
func (h *TestHelper) TeardownTestDB() {
	database.CloseDatabase()
}
