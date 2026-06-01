package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type CreateNotificationRequest struct {
	UserID    string `json:"user_id"`
	Target    string `json:"target"` // single, all_users, residents, staff, managers
	Title     string `json:"title" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Type      string `json:"type" binding:"required"`
	RelatedID string `json:"related_id,omitempty"`
}

type NotificationListItem struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Username  string     `json:"username,omitempty"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"`
	IsRead    bool       `json:"is_read"`
	RelatedID string     `json:"related_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

func buildNotificationListItem(n models.Notification) NotificationListItem {
	item := NotificationListItem{
		ID: n.ID, UserID: n.UserID, Title: n.Title, Message: n.Message,
		Type: n.Type, IsRead: n.IsRead, RelatedID: n.RelatedID,
		CreatedAt: n.CreatedAt, ReadAt: n.ReadAt,
	}
	if n.User.Username != "" {
		item.Username = n.User.Username
	}
	return item
}

func resolveNotificationRecipients(req CreateNotificationRequest) ([]string, error) {
	target := req.Target
	if target == "" || target == "single" {
		if req.UserID == "" {
			return nil, fmt.Errorf("Vui lòng chọn người nhận hoặc đối tượng gửi")
		}
		return []string{req.UserID}, nil
	}

	var users []models.User
	q := database.DB.Where("status = ?", "active")
	switch target {
	case "all_users":
		// all active users
	case "residents":
		q = q.Where("role = ?", "resident")
	case "staff":
		q = q.Where("role = ?", "staff")
	case "managers":
		q = q.Where("role IN ?", []string{"manager", "admin"})
	default:
		return nil, fmt.Errorf("Đối tượng gửi không hợp lệ")
	}
	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("Không tìm thấy người nhận phù hợp")
	}
	return ids, nil
}

// ========== Notification Handlers ==========

func CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipientIDs, err := resolveNotificationRecipients(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	created := make([]models.Notification, 0, len(recipientIDs))
	for _, uid := range recipientIDs {
		n := models.Notification{
			ID: uuid.New().String(), UserID: uid,
			Title: req.Title, Message: req.Message, Type: req.Type,
			IsRead: false, RelatedID: req.RelatedID, CreatedAt: now,
		}
		if err := database.DB.Create(&n).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
			return
		}
		created = append(created, n)
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    created,
		"count":   len(created),
		"success": true,
		"message": fmt.Sprintf("Đã gửi %d thông báo", len(created)),
	})
}

func GetUserNotifications(c *gin.Context) {
	currentUserID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	scopeAll := c.Query("scope") == "all" && (role.(string) == "admin" || role.(string) == "manager")

	var notifications []models.Notification
	query := database.DB.Preload("User")
	if scopeAll {
		query = query.Order("created_at DESC").Limit(100)
	} else {
		query = query.Where("user_id = ?", currentUserID.(string)).
			Order("created_at DESC").Limit(100)
	}
	if err := query.Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	var unreadCount int64
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", currentUserID.(string), false).
		Count(&unreadCount)

	items := make([]NotificationListItem, 0, len(notifications))
	for _, n := range notifications {
		items = append(items, buildNotificationListItem(n))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":          items,
		"unread_count":  unreadCount,
		"scope":         map[bool]string{true: "all", false: "mine"}[scopeAll],
	})
}

func GetUnreadNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var notifications []models.Notification

	if err := database.DB.Preload("User").Where("user_id = ? AND is_read = ?", userID.(string), false).
		Order("created_at DESC").Limit(50).
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	items := make([]NotificationListItem, 0, len(notifications))
	for _, n := range notifications {
		items = append(items, buildNotificationListItem(n))
	}

	c.JSON(http.StatusOK, gin.H{"data": items, "unread_count": len(items)})
}

func MarkNotificationAsRead(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	now := time.Now()

	res := database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID.(string)).
		Updates(map[string]interface{}{"is_read": true, "read_at": now})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy thông báo"})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read", "success": true})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	now := time.Now()

	if err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID.(string), false).
		Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read", "success": true})
}

func DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	q := database.DB.Where("id = ?", id)
	if role.(string) != "admin" {
		q = q.Where("user_id = ?", userID.(string))
	}
	res := q.Delete(&models.Notification{})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy thông báo"})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted", "success": true})
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
		MonthlyCollected       float64 `json:"monthly_collected"`
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

	// Hóa đơn quá hạn — cùng logic với GET /invoices/overdue
	now := time.Now()
	overdueQ := database.DB.Model(&models.Invoice{}).
		Where("(status = ? OR (status = ? AND due_date < ?))", "overdue", "pending", now).
		Where("(total_amount - paid_amount) > 0")
	overdueQ.Count(&kpi.OverdueInvoices)
	overdueQ.Select("COALESCE(SUM(total_amount - paid_amount), 0)").Row().Scan(&kpi.OverdueAmount)

	// Collection rate (paid amount / total billed)
	var totalBilled, totalPaid float64
	database.DB.Model(&models.Invoice{}).
		Select("SUM(total_amount) as total_billed, SUM(paid_amount) as total_paid").
		Row().Scan(&totalBilled, &totalPaid)
	if totalBilled > 0 {
		kpi.CollectionRate = (totalPaid / totalBilled) * 100
	}

	// Doanh thu tháng = tổng tiền thuê từ HĐ đang active trong tháng hiện tại
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
	database.DB.Model(&models.Contract{}).
		Where("status = ? AND start_date <= ? AND end_date >= ?", "active", monthEnd, monthStart).
		Select("COALESCE(SUM(monthly_rent), 0)").Row().Scan(&kpi.MonthlyRevenue)

	// Tiền đã thu thực tế trong tháng (hóa đơn status=paid)
	database.DB.Model(&models.Invoice{}).
		Where("status = ? AND EXTRACT(YEAR FROM updated_at) = EXTRACT(YEAR FROM NOW()) AND EXTRACT(MONTH FROM updated_at) = EXTRACT(MONTH FROM NOW())", "paid").
		Select("COALESCE(SUM(paid_amount), 0)").Row().Scan(&kpi.MonthlyCollected)

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
