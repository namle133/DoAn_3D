package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type CreateContractRequest struct {
	ApartmentID   string    `json:"apartment_id" binding:"required"`
	ResidentID    string    `json:"resident_id" binding:"required"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	MonthlyRent   float64   `json:"monthly_rent" binding:"required"`
	PaymentPeriod int       `json:"payment_period" binding:"required"`
	Deposit       float64   `json:"deposit"`
}

type ApproveContractRequest struct {
	Approved bool   `json:"approved" binding:"required"`
	Notes    string `json:"notes,omitempty"`
}

type ExtendContractRequest struct {
	NewEndDate time.Time `json:"new_end_date" binding:"required"`
}

// ========== Contract Handlers ==========

func CreateContract(c *gin.Context) {
	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify apartment and resident exist
	var apartment models.Apartment
	if err := database.DB.First(&apartment, "id = ?", req.ApartmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Apartment not found"})
		return
	}

	var resident models.Resident
	if err := database.DB.First(&resident, "id = ?", req.ResidentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		return
	}

	contract := models.Contract{
		ID:            uuid.New().String(),
		ApartmentID:   req.ApartmentID,
		ResidentID:    req.ResidentID,
		ContractCode:  "CONTRACT-" + time.Now().Format("20060102150405"),
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		MonthlyRent:   req.MonthlyRent,
		PaymentPeriod: req.PaymentPeriod,
		Deposit:       req.Deposit,
		Status:        "pending",
	}

	if err := database.DB.Create(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contract"})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

func GetContracts(c *gin.Context) {
	var contracts []models.Contract

	query := database.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if residentID := c.Query("resident_id"); residentID != "" {
		query = query.Where("resident_id = ?", residentID)
	}

	if err := query.Preload("Apartment").Preload("Resident").Find(&contracts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contracts"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

func GetContractByID(c *gin.Context) {
	id := c.Param("id")
	var contract models.Contract

	if err := database.DB.Preload("Apartment").Preload("Resident").Preload("Invoices").
		First(&contract, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func ApproveContract(c *gin.Context) {
	id := c.Param("id")
	var req ApproveContractRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	userID, _ := c.Get("user_id")
	now := time.Now()

	if req.Approved {
		contract.Status = "active"
		// Update apartment status to rented
		if err := database.DB.Model(&models.Apartment{}).
			Where("id = ?", contract.ApartmentID).
			Update("current_status", "rented").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update apartment status"})
			return
		}
	} else {
		contract.Status = "cancelled"
		// Update apartment status back to empty
		if err := database.DB.Model(&models.Apartment{}).
			Where("id = ?", contract.ApartmentID).
			Update("current_status", "empty").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update apartment status"})
			return
		}
	}

	contract.ApprovedBy = stringPtr(userID.(string))
	contract.ApprovedAt = &now

	if err := database.DB.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve contract"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func ExtendContract(c *gin.Context) {
	id := c.Param("id")
	var req ExtendContractRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	contract.EndDate = req.NewEndDate

	if err := database.DB.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend contract"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Contract extended successfully",
		"contract": contract,
	})
}

func TerminateContract(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	contract.Status = "terminated"
	contract.TerminationReason = req.Reason

	// Update apartment status back to empty
	if err := database.DB.Model(&models.Apartment{}).
		Where("id = ?", contract.ApartmentID).
		Update("current_status", "empty").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update apartment status"})
		return
	}

	if err := database.DB.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to terminate contract"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contract terminated successfully",
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
