package services

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// NotifyContractExpiring sends notification to residents about expiring contracts
func (s *NotificationService) NotifyContractExpiring() error {
	var contracts []models.Contract

	// Get contracts expiring within 30 days
	expiringDate := time.Now().AddDate(0, 0, 30)
	if err := database.DB.Where("status = ? AND end_date BETWEEN NOW() AND ?", "active", expiringDate).
		Preload("Resident.User").Find(&contracts).Error; err != nil {
		return err
	}

	for _, contract := range contracts {
		notification := models.Notification{
			ID:        uuid.New().String(),
			UserID:    contract.Resident.UserID,
			Title:     "Contract Expiring Soon",
			Message:   fmt.Sprintf("Your rental contract for apartment %s will expire on %s", contract.ApartmentID, contract.EndDate.Format("2006-01-02")),
			Type:      "contract_expiring",
			IsRead:    false,
			RelatedID: contract.ID,
			CreatedAt: time.Now(),
		}

		if err := database.DB.Create(&notification).Error; err != nil {
			log.Printf("Error creating notification: %v", err)
		}
	}

	return nil
}

// NotifyContractBilling sends invoice notification to resident
func (s *NotificationService) NotifyContractBilling(invoice *models.Invoice, contract *models.Contract) {
	notification := models.Notification{
		ID:        uuid.New().String(),
		UserID:    contract.Resident.UserID,
		Title:     "New Invoice Generated",
		Message:   fmt.Sprintf("Invoice %s for amount %.2f VND is now due on %s", invoice.InvoiceCode, invoice.TotalAmount, invoice.DueDate.Format("2006-01-02")),
		Type:      "invoice_generated",
		IsRead:    false,
		RelatedID: invoice.ID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification: %v", err)
	}
}

// NotifyOverduePayment sends overdue payment notification
func (s *NotificationService) NotifyOverduePayment(invoice *models.Invoice) {
	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", invoice.ContractID).Error; err != nil {
		return
	}

	var resident models.Resident
	if err := database.DB.First(&resident, "id = ?", contract.ResidentID).Error; err != nil {
		return
	}

	notification := models.Notification{
		ID:        uuid.New().String(),
		UserID:    resident.UserID,
		Title:     "Overdue Payment Alert",
		Message:   fmt.Sprintf("Your invoice %s is now overdue. Outstanding amount: %.2f VND", invoice.InvoiceCode, invoice.TotalAmount-invoice.PaidAmount),
		Type:      "invoice_overdue",
		IsRead:    false,
		RelatedID: invoice.ID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification: %v", err)
	}
}

// NotifyMaintenanceRequestCreated sends notification to managers
func (s *NotificationService) NotifyMaintenanceRequestCreated(request *models.MaintenanceRequest) error {
	var managers []models.Manager

	if err := database.DB.Preload("User").Find(&managers).Error; err != nil {
		return err
	}

	for _, manager := range managers {
		notification := models.Notification{
			ID:        uuid.New().String(),
			UserID:    manager.UserID,
			Title:     "New Maintenance Request",
			Message:   fmt.Sprintf("New %s maintenance request from apartment %s: %s", request.Priority, request.ApartmentID, request.Description),
			Type:      "maintenance_request",
			IsRead:    false,
			RelatedID: request.ID,
			CreatedAt: time.Now(),
		}

		if err := database.DB.Create(&notification).Error; err != nil {
			return err
		}
	}

	return nil
}

// NotifyMaintenanceAssigned sends notification to assigned staff (AssignedTo = user id)
func (s *NotificationService) NotifyMaintenanceAssigned(request *models.MaintenanceRequest) error {
	if request.AssignedTo == nil || *request.AssignedTo == "" {
		return nil
	}
	var user models.User
	if err := database.DB.Where("id = ? AND role IN ?", *request.AssignedTo, []string{"staff", "manager", "admin"}).
		First(&user).Error; err != nil {
		return err
	}

	notification := models.Notification{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Title:     "Maintenance Request Assigned",
		Message:   fmt.Sprintf("Maintenance request for apartment %s has been assigned to you. Issue: %s", request.ApartmentID, request.IssueType),
		Type:      "maintenance_assigned",
		IsRead:    false,
		RelatedID: request.ID,
		CreatedAt: time.Now(),
	}

	return database.DB.Create(&notification).Error
}

// NotifyEquipmentMaintenanceDue sends notification about equipment maintenance
func (s *NotificationService) NotifyEquipmentMaintenanceDue() error {
	var equipment []models.Equipment

	if err := database.DB.Where("next_maintenance_date <= ?", time.Now()).Find(&equipment).Error; err != nil {
		return err
	}

	// Notify all managers
	var managers []models.Manager
	if err := database.DB.Preload("User").Find(&managers).Error; err != nil {
		return err
	}

	for _, eq := range equipment {
		for _, manager := range managers {
			notification := models.Notification{
				ID:        uuid.New().String(),
				UserID:    manager.UserID,
				Title:     "Equipment Maintenance Due",
				Message:   fmt.Sprintf("Equipment '%s' (%s) is due for maintenance", eq.EquipmentName, eq.EquipmentType),
				Type:      "equipment_maintenance",
				IsRead:    false,
				RelatedID: eq.ID,
				CreatedAt: time.Now(),
			}

			if err := database.DB.Create(&notification).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanupOldNotifications removes notifications older than 30 days
func (s *NotificationService) CleanupOldNotifications() error {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	if err := database.DB.Where("created_at < ? AND is_read = ?", thirtyDaysAgo, true).
		Delete(&models.Notification{}).Error; err != nil {
		return err
	}

	return nil
}
