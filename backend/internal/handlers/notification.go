package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type CreateNotificationRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Type      string `json:"type" binding:"required"`
	RelatedID string `json:"related_id,omitempty"`
}

// ========== Notification Handlers ==========

func CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := models.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Title:     req.Title,
		Message:   req.Message,
		Type:      req.Type,
		IsRead:    false,
		RelatedID: req.RelatedID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

func GetUserNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var notifications []models.Notification

	// Get unread or recent notifications
	if err := database.DB.Where("user_id = ?", userID.(string)).
		Order("created_at DESC").
		Limit(50).
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func GetUnreadNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var notifications []models.Notification

	if err := database.DB.Where("user_id = ? AND is_read = ?", userID.(string), false).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func MarkNotificationAsRead(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()

	if err := database.DB.Model(&models.Notification{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	now := time.Now()

	if err := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID.(string)).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func DeleteNotification(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Notification{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}

// ========== Dashboard/KPI Handlers ==========

func GetDashboardKPIs(c *gin.Context) {
	type KPIData struct {
		TotalApartments        int64   `json:"total_apartments"`
		RentedApartments       int64   `json:"rented_apartments"`
		EmptyApartments        int64   `json:"empty_apartments"`
		MaintenanceApartments  int64   `json:"maintenance_apartments"`
		OccupancyRate          float64 `json:"occupancy_rate"`
		ActiveContracts        int64   `json:"active_contracts"`
		ExpiringContractsCount int64   `json:"expiring_contracts_count"`
		TotalResidents         int64   `json:"total_residents"`
		OutstandingInvoices    float64 `json:"outstanding_invoices"`
		CollectionRate         float64 `json:"collection_rate"`
		PendingMaintenanceReqs int64   `json:"pending_maintenance_requests"`
		AvgResponseTime        float64 `json:"avg_response_time_hours"`
		MonthlyRevenue         float64 `json:"monthly_revenue"`
		OverdueInvoices        int64   `json:"overdue_invoices"`
		OverdueAmount          float64 `json:"overdue_amount"`
	}

	kpi := KPIData{}

	// Total apartments
	database.DB.Model(&models.Apartment{}).Count(&kpi.TotalApartments)

	// Rented apartments
	database.DB.Model(&models.Apartment{}).Where("current_status = ?", "rented").Count(&kpi.RentedApartments)

	// Empty apartments
	database.DB.Model(&models.Apartment{}).Where("current_status = ?", "empty").Count(&kpi.EmptyApartments)

	// Maintenance apartments
	database.DB.Model(&models.Apartment{}).Where("current_status = ?", "maintenance").Count(&kpi.MaintenanceApartments)

	// Occupancy rate
	if kpi.TotalApartments > 0 {
		kpi.OccupancyRate = float64(kpi.RentedApartments) / float64(kpi.TotalApartments) * 100
	}

	// Active contracts
	database.DB.Model(&models.Contract{}).Where("status = ? AND end_date >= NOW()", "active").Count(&kpi.ActiveContracts)

	// Contracts expiring soon (within 30 days)
	database.DB.Model(&models.Contract{}).
		Where("status = ? AND end_date BETWEEN NOW() AND NOW() + INTERVAL '30 days'", "active").
		Count(&kpi.ExpiringContractsCount)

	// Total residents
	database.DB.Model(&models.Resident{}).Count(&kpi.TotalResidents)

	// Outstanding invoices (sum of unpaid amounts)
	database.DB.Model(&models.Invoice{}).Where("status IN ?", []string{"pending", "overdue"}).
		Select("SUM(total_amount - paid_amount)").Row().Scan(&kpi.OutstandingInvoices)

	// Overdue invoices count and amount
	database.DB.Model(&models.Invoice{}).Where("status = ?", "overdue").Count(&kpi.OverdueInvoices)
	database.DB.Model(&models.Invoice{}).Where("status = ?", "overdue").
		Select("SUM(total_amount - paid_amount)").Row().Scan(&kpi.OverdueAmount)

	// Collection rate (paid amount / total billed)
	var totalBilled, totalPaid float64
	database.DB.Model(&models.Invoice{}).
		Select("SUM(total_amount) as total_billed, SUM(paid_amount) as total_paid").
		Row().Scan(&totalBilled, &totalPaid)
	if totalBilled > 0 {
		kpi.CollectionRate = (totalPaid / totalBilled) * 100
	}

	// Monthly revenue (sum of paid invoices this month)
	database.DB.Model(&models.Invoice{}).
		Where("status = ? AND EXTRACT(YEAR FROM updated_at) = EXTRACT(YEAR FROM NOW()) AND EXTRACT(MONTH FROM updated_at) = EXTRACT(MONTH FROM NOW())", "paid").
		Select("SUM(paid_amount)").Row().Scan(&kpi.MonthlyRevenue)

	// Pending maintenance requests
	database.DB.Model(&models.MaintenanceRequest{}).Where("status IN ?", []string{"new", "assigned"}).Count(&kpi.PendingMaintenanceReqs)

	c.JSON(http.StatusOK, kpi)
}

func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog

	query := database.DB
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if entity := c.Query("entity"); entity != "" {
		query = query.Where("entity = ?", entity)
	}

	if err := query.Order("timestamp DESC").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
