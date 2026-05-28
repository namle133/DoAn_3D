package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type CreateResidentRequest struct {
	FullName         string    `json:"full_name" binding:"required"`
	IDCard           string    `json:"id_card" binding:"required"`
	DateOfBirth      time.Time `json:"date_of_birth" binding:"required"`
	Gender           string    `json:"gender" binding:"required,oneof=M F"`
	PhoneNumber      string    `json:"phone_number" binding:"required"`
	Email            string    `json:"email" binding:"required,email"`
	PermanentAddress string    `json:"permanent_address" binding:"required"`
	Username         string    `json:"username" binding:"required"`
	Password         string    `json:"password" binding:"required,min=8"`
}

type AddDependentRequest struct {
	FullName     string    `json:"full_name" binding:"required"`
	IDCard       string    `json:"id_card"`
	DateOfBirth  time.Time `json:"date_of_birth" binding:"required"`
	Relationship string    `json:"relationship" binding:"required"`
	ContractID   string    `json:"contract_id" binding:"required"`
}

type RegisterVehicleRequest struct {
	LicensePlate string `json:"license_plate" binding:"required"`
	VehicleType  string `json:"vehicle_type" binding:"required,oneof=car motorcycle bicycle"`
	Color        string `json:"color,omitempty"`
}

// ========== Resident Handlers ==========

func CreateResident(c *gin.Context) {
	var req CreateResidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user account first
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "resident",
		Status:   "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create resident profile
	resident := models.Resident{
		ID:               uuid.New().String(),
		UserID:           user.ID,
		FullName:         req.FullName,
		IDCard:           req.IDCard,
		DateOfBirth:      req.DateOfBirth,
		Gender:           req.Gender,
		PhoneNumber:      req.PhoneNumber,
		Email:            req.Email,
		PermanentAddress: req.PermanentAddress,
	}

	if err := database.DB.Create(&resident).Error; err != nil {
		database.DB.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resident profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Resident created successfully",
		"resident_id": resident.ID,
		"user_id":     user.ID,
	})
}

func GetResidents(c *gin.Context) {
	var residents []models.Resident

	if err := database.DB.Preload("Contracts").Find(&residents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch residents"})
		return
	}

	c.JSON(http.StatusOK, residents)
}

func GetResidentByID(c *gin.Context) {
	id := c.Param("id")
	var resident models.Resident

	if err := database.DB.Preload("Contracts").Preload("Dependents").
		Preload("Vehicles").Preload("ResidenceHistory").
		First(&resident, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		return
	}

	c.JSON(http.StatusOK, resident)
}

func UpdateResidentProfile(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		FullName         string `json:"full_name,omitempty"`
		PhoneNumber      string `json:"phone_number,omitempty"`
		Email            string `json:"email,omitempty"`
		PermanentAddress string `json:"permanent_address,omitempty"`
		Avatar           string `json:"avatar,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.Resident{}).Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ========== Dependent Handlers ==========

func AddDependent(c *gin.Context) {
	residentID := c.Param("id")
	var req AddDependentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dependent := models.ResidentDependent{
		ID:           uuid.New().String(),
		ResidentID:   residentID,
		ContractID:   req.ContractID,
		FullName:     req.FullName,
		IDCard:       req.IDCard,
		DateOfBirth:  req.DateOfBirth,
		Relationship: req.Relationship,
		StartDate:    time.Now(),
	}

	if err := database.DB.Create(&dependent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add dependent"})
		return
	}

	c.JSON(http.StatusCreated, dependent)
}

func RemoveDependent(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()

	if err := database.DB.Model(&models.ResidentDependent{}).
		Where("id = ?", id).Update("end_date", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove dependent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dependent removed successfully"})
}

// ========== Vehicle Handlers ==========

func RegisterVehicle(c *gin.Context) {
	residentID := c.Param("id")
	var req RegisterVehicleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle := models.Vehicle{
		ID:           uuid.New().String(),
		ResidentID:   residentID,
		LicensePlate: req.LicensePlate,
		VehicleType:  req.VehicleType,
		Color:        req.Color,
		RegisteredAt: time.Now(),
	}

	if err := database.DB.Create(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register vehicle"})
		return
	}

	c.JSON(http.StatusCreated, vehicle)
}

func GetResidentVehicles(c *gin.Context) {
	residentID := c.Param("id")
	var vehicles []models.Vehicle

	if err := database.DB.Where("resident_id = ?", residentID).Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

func UnregisterVehicle(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Vehicle{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unregister vehicle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle unregistered successfully"})
}
