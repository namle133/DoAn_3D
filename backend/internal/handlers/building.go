package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/datatypes"
)

type CreateBuildingRequest struct {
	Name        string         `json:"name" binding:"required"`
	Address     string         `json:"address" binding:"required"`
	YearBuilt   int            `json:"year_built"`
	TotalFloors int            `json:"total_floors" binding:"required"`
	Description string         `json:"description"`
	Geometry    datatypes.JSON `json:"geometry,omitempty"` // 3D geometry
}

type CreateFloorRequest struct {
	BuildingID  string         `json:"building_id" binding:"required"`
	FloorNumber int            `json:"floor_number" binding:"required"`
	FloorHeight float64        `json:"floor_height"`
	FloorArea   float64        `json:"floor_area"`
	Purpose     string         `json:"purpose"` // residential, technical, service
	Geometry    datatypes.JSON `json:"geometry,omitempty"`
}

type CreateApartmentRequest struct {
	FloorID       string         `json:"floor_id" binding:"required"`
	ApartmentCode string         `json:"apartment_code" binding:"required"`
	Area          float64        `json:"area" binding:"required"`
	Bedrooms      int            `json:"bedrooms"`
	Direction     string         `json:"direction"`
	ListingPrice  float64        `json:"listing_price" binding:"required"`
	CurrentStatus string         `json:"current_status"`
	BodyID        string         `json:"body_id,omitempty"` // 3D model reference
	Geometry      datatypes.JSON `json:"geometry,omitempty"`
}

// ========== Building Handlers ==========

func CreateBuilding(c *gin.Context) {
	var req CreateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	building := models.Building{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Address:     req.Address,
		YearBuilt:   req.YearBuilt,
		TotalFloors: req.TotalFloors,
		Description: req.Description,
		Geometry:    req.Geometry,
	}

	if err := database.DB.Create(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create building"})
		return
	}

	c.JSON(http.StatusCreated, building)
}

func GetBuildings(c *gin.Context) {
	var buildings []models.Building

	if err := database.DB.Preload("Floors").Find(&buildings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch buildings"})
		return
	}

	c.JSON(http.StatusOK, buildings)
}

func GetBuildingByID(c *gin.Context) {
	id := c.Param("id")
	var building models.Building

	if err := database.DB.Preload("Floors.Apartments").First(&building, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
		return
	}

	c.JSON(http.StatusOK, building)
}

func UpdateBuilding(c *gin.Context) {
	id := c.Param("id")
	var req CreateBuildingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.Building{}).Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update building"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Building updated successfully"})
}

func DeleteBuilding(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Building{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete building"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Building deleted successfully"})
}

// ========== Floor Handlers ==========

func CreateFloor(c *gin.Context) {
	var req CreateFloorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	floor := models.Floor{
		ID:          uuid.New().String(),
		BuildingID:  req.BuildingID,
		FloorNumber: req.FloorNumber,
		FloorHeight: req.FloorHeight,
		FloorArea:   req.FloorArea,
		Purpose:     req.Purpose,
		Geometry:    req.Geometry,
	}

	if err := database.DB.Create(&floor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create floor"})
		return
	}

	c.JSON(http.StatusCreated, floor)
}

func GetFloorsByBuilding(c *gin.Context) {
	buildingID := c.Param("id")
	var floors []models.Floor

	if err := database.DB.Where("building_id = ?", buildingID).Preload("Apartments").Find(&floors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch floors"})
		return
	}

	c.JSON(http.StatusOK, floors)
}

// ========== Apartment Handlers ==========

func CreateApartment(c *gin.Context) {
	var req CreateApartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apartment := models.Apartment{
		ID:            uuid.New().String(),
		FloorID:       req.FloorID,
		ApartmentCode: req.ApartmentCode,
		Area:          req.Area,
		Bedrooms:      req.Bedrooms,
		Direction:     req.Direction,
		ListingPrice:  req.ListingPrice,
		CurrentStatus: req.CurrentStatus,
		BodyID:        req.BodyID,
		Geometry:      req.Geometry,
	}

	if err := database.DB.Create(&apartment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create apartment"})
		return
	}

	c.JSON(http.StatusCreated, apartment)
}

func GetApartments(c *gin.Context) {
	var apartments []models.Apartment

	// Filter by floor_id if provided
	query := database.DB
	if floorID := c.Query("floor_id"); floorID != "" {
		query = query.Where("floor_id = ?", floorID)
	}

	if err := query.Find(&apartments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch apartments"})
		return
	}

	c.JSON(http.StatusOK, apartments)
}

func GetApartmentByID(c *gin.Context) {
	id := c.Param("id")
	var apartment models.Apartment

	if err := database.DB.Preload("Contracts").Preload("StatusHistory").
		First(&apartment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Apartment not found"})
		return
	}

	c.JSON(http.StatusOK, apartment)
}

func UpdateApartmentStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
		Reason string `json:"reason,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var apartment models.Apartment
	if err := database.DB.First(&apartment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Apartment not found"})
		return
	}

	oldStatus := apartment.CurrentStatus

	// Update apartment status
	if err := database.DB.Model(&apartment).Update("current_status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update apartment"})
		return
	}

	// Record status history
	userID, _ := c.Get("user_id")
	history := models.ApartmentStatusHistory{
		ID:          uuid.New().String(),
		ApartmentID: id,
		OldStatus:   oldStatus,
		NewStatus:   req.Status,
		Reason:      req.Reason,
		ChangedBy:   userID.(string),
		ChangedAt:   database.DB.NowFunc(),
	}

	database.DB.Create(&history)

	c.JSON(http.StatusOK, gin.H{"message": "Apartment status updated successfully"})
}
