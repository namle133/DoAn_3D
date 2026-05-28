package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type CreateMaintenanceReqPayload struct {
	ApartmentID string   `json:"apartment_id" binding:"required"`
	ResidentID  string   `json:"resident_id" binding:"required"`
	IssueType   string   `json:"issue_type" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Priority    string   `json:"priority" binding:"required,oneof=urgent normal low"`
	Location    string   `json:"location"`
	Photos      []string `json:"photos,omitempty"`
}

type AssignMaintenanceReqPayload struct {
	AssignedTo string `json:"assigned_to" binding:"required"`
}

type CompleteMaintenanceReqPayload struct {
	CompletionNotes string `json:"completion_notes"`
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

// ========== Maintenance Request Handlers ==========

func CreateMaintenanceRequest(c *gin.Context) {
	var req CreateMaintenanceReqPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	maintenanceReq := models.MaintenanceRequest{
		ID:          uuid.New().String(),
		ApartmentID: req.ApartmentID,
		ResidentID:  req.ResidentID,
		RequestCode: "MR-" + time.Now().Format("20060102150405"),
		IssueType:   req.IssueType,
		Description: req.Description,
		Priority:    req.Priority,
		Location:    req.Location,
		Photos:      pq.StringArray(req.Photos),
		Status:      "new",
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&maintenanceReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance request"})
		return
	}

	c.JSON(http.StatusCreated, maintenanceReq)
}

func GetMaintenanceRequests(c *gin.Context) {
	var requests []models.MaintenanceRequest

	query := database.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if apartmentID := c.Query("apartment_id"); apartmentID != "" {
		query = query.Where("apartment_id = ?", apartmentID)
	}

	if err := query.Preload("Apartment").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance requests"})
		return
	}

	// Build response with apartment_code
	type MaintenanceResponse struct {
		ID          string    `json:"id"`
		RequestCode string    `json:"request_code"`
		ApartmentID string    `json:"apartment_id"`
		IssueType   string    `json:"issue_type"`
		Priority    string    `json:"priority"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var response []MaintenanceResponse
	for _, req := range requests {
		response = append(response, MaintenanceResponse{
			ID:          req.ID,
			RequestCode: req.RequestCode,
			ApartmentID: req.ApartmentID,
			IssueType:   req.IssueType,
			Priority:    req.Priority,
			Status:      req.Status,
			CreatedAt:   req.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetMaintenanceRequestByID(c *gin.Context) {
	id := c.Param("id")
	var request models.MaintenanceRequest

	if err := database.DB.Preload("Apartment").Preload("Resident").
		First(&request, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance request not found"})
		return
	}

	c.JSON(http.StatusOK, request)
}

func AssignMaintenanceRequest(c *gin.Context) {
	id := c.Param("id")
	var req AssignMaintenanceReqPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	if err := database.DB.Model(&models.MaintenanceRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"assigned_to": req.AssignedTo,
			"assigned_at": now,
			"status":      "assigned",
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign maintenance request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance request assigned successfully"})
}

func UpdateMaintenanceStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=new assigned in_progress completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.MaintenanceRequest{}).Where("id = ?", id).
		Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func CompleteMaintenanceRequest(c *gin.Context) {
	id := c.Param("id")
	var req CompleteMaintenanceReqPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	if err := database.DB.Model(&models.MaintenanceRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           "completed",
			"completed_at":     now,
			"completion_notes": req.CompletionNotes,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance request completed successfully"})
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

	// Find equipment where next maintenance date is before or equal to today
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
