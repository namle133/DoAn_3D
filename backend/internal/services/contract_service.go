package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type ContractService struct{}

func NewContractService() *ContractService {
	return &ContractService{}
}

// AutoTerminateExpiredContracts terminates contracts that have reached end date
func (s *ContractService) AutoTerminateExpiredContracts() error {
	var contracts []models.Contract

	if err := database.DB.Where("status = ? AND end_date < ?", "active", time.Now()).
		Find(&contracts).Error; err != nil {
		return err
	}

	for _, contract := range contracts {
		contract.Status = "expired"
		contract.TerminationReason = "Contract automatically expired"

		if err := database.DB.Save(&contract).Error; err != nil {
			return err
		}

		// Update apartment status to empty
		if err := database.DB.Model(&models.Apartment{}).Where("id = ?", contract.ApartmentID).
			Update("current_status", "empty").Error; err != nil {
			return err
		}

		// Create status history
		database.DB.Create(&models.ApartmentStatusHistory{
			ID:          uuid.New().String(),
			ApartmentID: contract.ApartmentID,
			OldStatus:   "rented",
			NewStatus:   "empty",
			Reason:      "Contract expired",
			ChangedAt:   time.Now(),
		})

		// Notify resident
		notificationService := NewNotificationService()
		notificationService.NotifyContractExpiring()
	}

	return nil
}

// GetContractsByResident returns all contracts for a resident
func (s *ContractService) GetContractsByResident(residentID string) ([]models.Contract, error) {
	var contracts []models.Contract

	if err := database.DB.Where("resident_id = ?", residentID).
		Preload("Apartment").
		Preload("Invoices").
		Order("start_date DESC").
		Find(&contracts).Error; err != nil {
		return nil, err
	}

	return contracts, nil
}

// GetActiveContracts returns all currently active contracts
func (s *ContractService) GetActiveContracts() ([]models.Contract, error) {
	var contracts []models.Contract

	if err := database.DB.Where("status = ? AND end_date >= ?", "active", time.Now()).
		Preload("Apartment").
		Preload("Resident").
		Find(&contracts).Error; err != nil {
		return nil, err
	}

	return contracts, nil
}

// GetContractRevenue returns total revenue from contracts for a period
func (s *ContractService) GetContractRevenue(startDate, endDate time.Time) (float64, error) {
	var totalRevenue float64

	err := database.DB.Model(&models.Contract{}).
		Where("status = ? AND start_date <= ? AND end_date >= ?", "active", endDate, startDate).
		Select("SUM(monthly_rent)").
		Row().Scan(&totalRevenue)

	return totalRevenue, err
}
