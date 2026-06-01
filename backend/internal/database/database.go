package database

import (
	"fmt"
	"log"

	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) error {
	dsn := cfg.GetDSN()
	// PreferSimpleProtocol né bug "insufficient arguments" của pgx/v5 với GORM
	// AutoMigrate khi gặp các query không có tham số.
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("✓ Database connected successfully")

	// Run migrations
	if err := runMigrations(); err != nil {
		return err
	}

	SeedDemoData()
	RepairOrphanResidents()
	return nil
}

func runMigrations() error {
	log.Println("Running database migrations...")

	// AutoMigrate will create tables and columns if they don't exist
	err := DB.AutoMigrate(
		// Authentication
		&models.User{},
		&models.Admin{},
		&models.Manager{},
		&models.Staff{},

		// GIS 3D & Building Management
		&models.Building{},
		&models.Floor{},
		&models.Apartment{},
		&models.ApartmentStatusHistory{},

		// Contract Management
		&models.Contract{},

		// Resident Management
		&models.Resident{},
		&models.ResidentDependent{},
		&models.Vehicle{},
		&models.ResidenceHistory{},

		// Financial Management
		&models.Invoice{},
		&models.InvoicePayment{},

		// Maintenance Management
		&models.MaintenanceRequest{},
		&models.Equipment{},

		// Notifications & Audit
		&models.Notification{},
		&models.AuditLog{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✓ All migrations completed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDatabase() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
