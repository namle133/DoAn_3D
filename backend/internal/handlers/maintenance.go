package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/services"
)

type CreateMaintenanceReqPayload struct {
	ApartmentID string   `json:"apartment_id" binding:"required"`
	ResidentID  string   `json:"resident_id"`
	IssueType   string   `json:"issue_type" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Priority    string   `json:"priority" binding:"required"`
	Location    string   `json:"location"`
	Photos      []string `json:"photos,omitempty"`
	AssignedTo  string   `json:"assigned_to,omitempty"`
}

type AssignMaintenanceReqPayload struct {
	AssignedTo string `json:"assigned_to" binding:"required"`
	Notes      string `json:"notes,omitempty"`
}

type CompleteMaintenanceReqPayload struct {
	CompletionNotes string  `json:"completion_notes"`
	Resolution      string  `json:"resolution"`
	Notes           string  `json:"notes"`
	Cost            float64 `json:"cost"`
	TechnicianName  string  `json:"technician_name"`
}

type RegisterEquipmentRequest struct {
	BuildingID        string    `json:"building_id" binding:"required"`
	EquipmentName     string    `json:"equipment_name" binding:"required"`
	EquipmentType     string    `json:"equipment_type" binding:"required"`
	Manufacturer      string    `json:"manufacturer"`
	InstallationDate  time.Time `json:"installation_date" binding:"required"`
	LocationNodeID    string    `json:"location_node_id,omitempty"`
	MaintenancePeriod int       `json:"maintenance_period" binding:"required"`
}

func mapPriorityToBackend(p string) string {
	switch strings.ToLower(p) {
	case "urgent", "high":
		return "urgent"
	case "medium":
		return "normal"
	default:
		return "low"
	}
}

func resolveResidentID(c *gin.Context, apartmentID, residentID string) (string, error) {
	if residentID != "" {
		return residentID, nil
	}
	if role, ok := c.Get("role"); ok && role.(string) == "resident" {
		userID, _ := c.Get("user_id")
		var resident models.Resident
		if err := database.DB.Where("user_id = ?", userID).First(&resident).Error; err != nil {
			return "", fmt.Errorf("Không tìm thấy hồ sơ cư dân")
		}
		return resident.ID, nil
	}
	var contract models.Contract
	if err := database.DB.Where("apartment_id = ? AND status IN ?", apartmentID, []string{"active", "extended"}).
		First(&contract).Error; err != nil {
		return "", fmt.Errorf("Không tìm thấy cư dân có HĐ active cho căn hộ này")
	}
	return contract.ResidentID, nil
}

// ========== Maintenance Request Handlers ==========

func CreateMaintenanceRequest(c *gin.Context) {
	var req CreateMaintenanceReqPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	residentID, err := resolveResidentID(c, req.ApartmentID, req.ResidentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := "new"
	var assignedTo *string
	if req.AssignedTo != "" {
		status = "assigned"
		assignedTo = &req.AssignedTo
	}

	now := time.Now()
	maintenanceReq := models.MaintenanceRequest{
		ID:          uuid.New().String(),
		ApartmentID: req.ApartmentID,
		ResidentID:  residentID,
		RequestCode: "MR-" + time.Now().Format("20060102150405"),
		IssueType:   req.IssueType,
		Description: req.Description,
		Priority:    mapPriorityToBackend(req.Priority),
		Location:    req.Location,
		Photos:      pq.StringArray(req.Photos),
		Status:      status,
		AssignedTo:  assignedTo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if status == "assigned" {
		maintenanceReq.AssignedAt = &now
	}

	if err := database.DB.Create(&maintenanceReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance request"})
		return
	}

	ns := services.NewNotificationService()
	_ = ns.NotifyMaintenanceRequestCreated(&maintenanceReq)
	if maintenanceReq.AssignedTo != nil {
		_ = ns.NotifyMaintenanceAssigned(&maintenanceReq)
	}

	c.JSON(http.StatusCreated, gin.H{"data": maintenanceReq, "success": true})
}

type MaintenanceListItem struct {
	ID               string     `json:"id"`
	RequestCode      string     `json:"request_code"`
	RequestNumber    string     `json:"request_number"`
	ApartmentID      string     `json:"apartment_id"`
	ApartmentCode    string     `json:"apartment_code"`
	ApartmentNumber  string     `json:"apartment_number"`
	ResidentID       string     `json:"resident_id"`
	ResidentName     string     `json:"resident_name"`
	IssueType        string     `json:"issue_type"`
	Description      string     `json:"description"`
	Priority         string     `json:"priority"`
	Location         string     `json:"location"`
	Status           string     `json:"status"`
	AssignedTo       *string    `json:"assigned_to,omitempty"`
	AssignedToName   string     `json:"assigned_to_name"`
	AssignedAt       *time.Time `json:"assigned_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CompletedDate    *time.Time `json:"completed_date,omitempty"`
	CompletionNotes  string     `json:"completion_notes,omitempty"`
	Resolution       string     `json:"resolution,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedDate      time.Time  `json:"created_date"`
	Photos           []string   `json:"photos,omitempty"`
}

func buildMaintenanceListItem(req models.MaintenanceRequest) MaintenanceListItem {
	item := MaintenanceListItem{
		ID: req.ID, RequestCode: req.RequestCode, RequestNumber: req.RequestCode,
		ApartmentID: req.ApartmentID, ResidentID: req.ResidentID,
		IssueType: req.IssueType, Description: req.Description,
		Priority: req.Priority, Location: req.Location, Status: req.Status,
		AssignedTo: req.AssignedTo, AssignedAt: req.AssignedAt,
		CompletedAt: req.CompletedAt, CompletedDate: req.CompletedAt,
		CompletionNotes: req.CompletionNotes, CreatedAt: req.CreatedAt, CreatedDate: req.CreatedAt,
		Photos: []string(req.Photos),
	}
	if req.Apartment.ApartmentCode != "" {
		item.ApartmentCode = req.Apartment.ApartmentCode
		item.ApartmentNumber = req.Apartment.ApartmentCode
	}
	if req.Resident.FullName != "" {
		item.ResidentName = req.Resident.FullName
	}
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		var user models.User
		if database.DB.First(&user, "id = ?", *req.AssignedTo).Error == nil {
			item.AssignedToName = user.Username
		}
	}
	if req.CompletionNotes != "" {
		item.Resolution = req.CompletionNotes
	}
	return item
}

func GetMaintenanceRequests(c *gin.Context) {
	var requests []models.MaintenanceRequest

	query := database.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", mapPriorityToBackend(priority))
	}
	if apartmentID := c.Query("apartment_id"); apartmentID != "" {
		query = query.Where("apartment_id = ?", apartmentID)
	}

	if role, ok := c.Get("role"); ok && role.(string) == "resident" {
		userID, _ := c.Get("user_id")
		var resident models.Resident
		if err := database.DB.Where("user_id = ?", userID).First(&resident).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"data": []MaintenanceListItem{}})
			return
		}
		query = query.Where("resident_id = ?", resident.ID)
	}

	if err := query.Preload("Apartment").Preload("Resident").
		Order("created_at DESC").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance requests"})
		return
	}

	result := make([]MaintenanceListItem, 0, len(requests))
	for _, req := range requests {
		result = append(result, buildMaintenanceListItem(req))
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetMaintenanceRequestByID(c *gin.Context) {
	id := c.Param("id")
	var request models.MaintenanceRequest

	if err := database.DB.Preload("Apartment").Preload("Resident").
		First(&request, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance request not found"})
		return
	}

	if role, ok := c.Get("role"); ok && role.(string) == "resident" {
		userID, _ := c.Get("user_id")
		var resident models.Resident
		if err := database.DB.Where("user_id = ?", userID).First(&resident).Error; err != nil ||
			request.ResidentID != resident.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": buildMaintenanceListItem(request)})
}

func AssignMaintenanceRequest(c *gin.Context) {
	id := c.Param("id")
	var req AssignMaintenanceReqPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ? AND role IN ?", req.AssignedTo, []string{"admin", "manager", "staff"}).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nhân viên không hợp lệ"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"assigned_to": req.AssignedTo,
		"assigned_at": now,
		"status":      "assigned",
		"updated_at":  now,
	}
	if req.Notes != "" {
		var mr models.MaintenanceRequest
		if database.DB.First(&mr, "id = ?", id).Error == nil {
			updates["description"] = mr.Description + "\n[Ghi chú phân công]: " + req.Notes
		}
	}

	if err := database.DB.Model(&models.MaintenanceRequest{}).Where("id = ?", id).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign maintenance request"})
		return
	}

	var mr models.MaintenanceRequest
	if database.DB.First(&mr, "id = ?", id).Error == nil {
		mr.AssignedTo = &req.AssignedTo
		_ = services.NewNotificationService().NotifyMaintenanceAssigned(&mr)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance request assigned successfully", "success": true})
}

func UpdateMaintenanceStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=new assigned in_progress completed cancelled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.MaintenanceRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": req.Status, "updated_at": time.Now()}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully", "success": true})
}

func CompleteMaintenanceRequest(c *gin.Context) {
	id := c.Param("id")
	var req CompleteMaintenanceReqPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notes := req.CompletionNotes
	if notes == "" {
		if req.Resolution != "" {
			notes = "Giải pháp: " + req.Resolution
		}
		if req.Notes != "" {
			if notes != "" {
				notes += "\n"
			}
			notes += "Ghi chú: " + req.Notes
		}
		if req.Cost > 0 {
			if notes != "" {
				notes += "\n"
			}
			notes += fmt.Sprintf("Chi phí: %.0f VND", req.Cost)
		}
		if req.TechnicianName != "" {
			if notes != "" {
				notes += "\n"
			}
			notes += "Người thực hiện: " + req.TechnicianName
		}
	}
	if notes == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng mô tả giải pháp"})
		return
	}

	now := time.Now()
	if err := database.DB.Model(&models.MaintenanceRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           "completed",
			"completed_at":     now,
			"completion_notes": notes,
			"updated_at":       now,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance request completed successfully", "success": true})
}

// ========== Equipment Handlers ==========

func RegisterEquipment(c *gin.Context) {
	var req RegisterEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nextMaintenance := req.InstallationDate.AddDate(0, req.MaintenancePeriod, 0)

	equipment := models.Equipment{
		ID:                  uuid.New().String(),
		BuildingID:          req.BuildingID,
		EquipmentName:       req.EquipmentName,
		EquipmentType:       req.EquipmentType,
		Manufacturer:        req.Manufacturer,
		InstallationDate:    req.InstallationDate,
		LocationNodeID:      req.LocationNodeID,
		MaintenancePeriod:   req.MaintenancePeriod,
		LastMaintenanceDate: &req.InstallationDate,
		NextMaintenanceDate: &nextMaintenance,
		Status:              "operational",
	}

	if err := database.DB.Create(&equipment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register equipment"})
		return
	}

	c.JSON(http.StatusCreated, equipment)
}

func GetEquipmentByBuilding(c *gin.Context) {
	buildingID := c.Param("id")
	var equipment []models.Equipment

	if err := database.DB.Where("building_id = ?", buildingID).Find(&equipment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch equipment"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func GetEquipmentDueForMaintenance(c *gin.Context) {
	var equipment []models.Equipment

	if err := database.DB.Where("next_maintenance_date <= NOW()").Find(&equipment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch equipment"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func UpdateEquipmentMaintenanceRecord(c *gin.Context) {
	id := c.Param("id")

	now := time.Now()
	var equipment models.Equipment

	if err := database.DB.First(&equipment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	nextMaintenance := now.AddDate(0, equipment.MaintenancePeriod, 0)

	if err := database.DB.Model(&equipment).Updates(map[string]interface{}{
		"last_maintenance_date": now,
		"next_maintenance_date": nextMaintenance,
		"status":                "operational",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance record updated successfully"})
}
