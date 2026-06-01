package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/handlers"
	"github.com/uit/vinhomes-management/internal/middleware"
	"github.com/uit/vinhomes-management/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Println("✓ Configuration loaded")
	fmt.Println(cfg)
	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Start background scheduler
	scheduler := services.GetScheduler()
	scheduler.StartScheduler()
	defer scheduler.StopScheduler()

	// Setup Gin engine
	gin.SetMode(cfg.GinMode)
	router := gin.New()

	// Global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.AuditLogMiddleware())

	// ========== Public Routes (No Authentication) ==========
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", handlers.Login(cfg))
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/forgot-password", handlers.ForgotPassword)
		public.POST("/auth/reset-password", handlers.ResetPassword)
		public.GET("/buildings", handlers.GetBuildings)
		public.GET("/buildings/:id", handlers.GetBuildingByID)
		public.GET("/gis/apartment-status", handlers.GetGISApartmentStatus)
	}

	// ========== Protected Routes (Requires Authentication) ==========
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// ========== Profile Management ==========
		protected.GET("/profile", handlers.GetProfile)
		protected.PUT("/profile", handlers.UpdateProfile)

		// ========== User Management (Admin) ==========
		protected.GET("/users", middleware.RequireRole("admin"), handlers.GetUsers)
		protected.POST("/users", middleware.RequireRole("admin"), handlers.CreateUser)
		protected.PUT("/users/:id", middleware.RequireRole("admin"), handlers.UpdateUser)
		protected.DELETE("/users/:id", middleware.RequireRole("admin"), handlers.DeleteUser)

		// ========== Building Management ==========
		protected.POST("/buildings", middleware.RequireRole("admin", "manager"), handlers.CreateBuilding)
		protected.PUT("/buildings/:id", middleware.RequireRole("admin", "manager"), handlers.UpdateBuilding)
		protected.DELETE("/buildings/:id", middleware.RequireRole("admin"), handlers.DeleteBuilding)

		// ========== Floor Management ==========
		protected.POST("/floors", middleware.RequireRole("admin", "manager"), handlers.CreateFloor)
		protected.GET("/buildings/:id/floors", handlers.GetFloorsByBuilding)

		// ========== Apartment Management ==========
		protected.POST("/apartments", middleware.RequireRole("admin", "manager"), handlers.CreateApartment)
		protected.GET("/apartments", handlers.GetApartments)
		protected.GET("/apartments/:id", handlers.GetApartmentByID)
		protected.PUT("/apartments/:id/status", middleware.RequireRole("admin", "manager"), handlers.UpdateApartmentStatus)

		// ========== Contract Management ==========
		protected.POST("/contracts", middleware.RequireRole("admin", "manager"), handlers.CreateContract)
		protected.GET("/contracts", handlers.GetContracts)
		protected.GET("/contracts/:id", handlers.GetContractByID)
		protected.POST("/contracts/:id/approve", middleware.RequireRole("admin", "manager"), handlers.ApproveContract)
		protected.POST("/contracts/:id/extend", middleware.RequireRole("admin", "manager"), handlers.ExtendContract)
		protected.POST("/contracts/:id/terminate", middleware.RequireRole("admin", "manager"), handlers.TerminateContract)

		// ========== Resident Management ==========
		protected.POST("/residents", middleware.RequireRole("admin", "manager"), handlers.CreateResident)
		protected.GET("/residents", handlers.GetResidents)
		protected.GET("/residents/:id", handlers.GetResidentByID)
		protected.PUT("/residents/:id", handlers.UpdateResidentProfile)

		// ========== Dependent Management ==========
		protected.POST("/residents/:id/dependents", handlers.AddDependent)
		protected.DELETE("/dependents/:id", handlers.RemoveDependent)

		// ========== Vehicle Management ==========
		protected.POST("/residents/:id/vehicles", handlers.RegisterVehicle)
		protected.GET("/residents/:id/vehicles", handlers.GetResidentVehicles)
		protected.DELETE("/vehicles/:id", handlers.UnregisterVehicle)

		// ========== Invoice Management ==========
		protected.POST("/invoices", middleware.RequireRole("admin", "manager"), handlers.CreateInvoice)
		protected.GET("/invoices", handlers.GetInvoices)
		protected.GET("/invoices/:id", handlers.GetInvoiceByID)
		protected.PUT("/invoices/:id", middleware.RequireRole("admin", "manager"), handlers.UpdateInvoice)
		protected.POST("/invoices/:id/payments", middleware.RequireRole("staff", "manager", "admin"), handlers.RecordPayment)
		protected.GET("/invoices/overdue", handlers.GetOverdueInvoices)
		protected.POST("/reports/financial", handlers.GetFinancialReport)

		// ========== Maintenance Management ==========
		protected.POST("/maintenance-requests", handlers.CreateMaintenanceRequest)
		protected.GET("/maintenance-requests", handlers.GetMaintenanceRequests)
		protected.GET("/maintenance-requests/:id", handlers.GetMaintenanceRequestByID)
		protected.POST("/maintenance-requests/:id/assign", middleware.RequireRole("admin", "manager"), handlers.AssignMaintenanceRequest)
		protected.PUT("/maintenance-requests/:id/status", handlers.UpdateMaintenanceStatus)
		protected.POST("/maintenance-requests/:id/complete", handlers.CompleteMaintenanceRequest)

		// ========== Equipment Management ==========
		protected.POST("/equipment", middleware.RequireRole("admin", "manager"), handlers.RegisterEquipment)
		protected.GET("/buildings/:id/equipment", handlers.GetEquipmentByBuilding)
		protected.GET("/equipment/due-for-maintenance", handlers.GetEquipmentDueForMaintenance)
		protected.PUT("/equipment/:id/maintenance", middleware.RequireRole("staff", "manager"), handlers.UpdateEquipmentMaintenanceRecord)

		// ========== Notification Management ==========
		protected.POST("/notifications", middleware.RequireRole("admin", "manager"), handlers.CreateNotification)
		protected.GET("/notifications", handlers.GetUserNotifications)
		protected.GET("/notifications/unread", handlers.GetUnreadNotifications)
		protected.PUT("/notifications/:id/read", handlers.MarkNotificationAsRead)
		protected.PUT("/notifications/read-all", handlers.MarkAllNotificationsAsRead)
		protected.DELETE("/notifications/:id", handlers.DeleteNotification)

		// ========== Dashboard & Reports ==========
		protected.GET("/dashboard/kpis", handlers.GetDashboardKPIs)
		protected.GET("/audit-logs", middleware.RequireRole("admin"), handlers.GetAuditLogs)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"app":    "Vinhomes Property Management",
		})
	})

	// Start server in goroutine
	go func() {
		address := fmt.Sprintf(":%s", cfg.ServerPort)
		log.Printf("✓ Starting server on %s\n", address)
		if err := router.Run(address); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("✓ Graceful shutdown initiated")
}
