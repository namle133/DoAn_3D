package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	ServerPort string
	GinMode    string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string

	// App
	AppEnv string
}

func LoadConfig() *Config {
	loadEnvFiles()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	return &Config{
		// Mặc định: PostgreSQL trong Docker (docker-compose port 5433)
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "5433"),
		DBUser:         getEnv("DB_USER", "vinhomes_user"),
		DBPassword:     getEnv("DB_PASSWORD", "vinhomes_pass"),
		DBName:         getEnv("DB_NAME", "vinhomes_db"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHours: jwtExpiry,
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		AppEnv:         getEnv("APP_ENV", "development"),
	}
}

// loadEnvFiles ưu tiên config/.env rồi .env ở thư mục backend.
func loadEnvFiles() {
	for _, path := range []string{"config/.env", ".env"} {
		if err := godotenv.Load(path); err == nil {
			return
		}
	}
	_ = godotenv.Load()
}

func (c *Config) GetDSN() string {
	host := c.DBHost
	if host == "localhost" {
		// Tránh kết nối IPv6 [::1] trùng instance Postgres khác trên máy
		host = "127.0.0.1"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
