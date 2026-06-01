package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type InvoiceService struct{}

func NewInvoiceService() *InvoiceService {
	return &InvoiceService{}
}

// CreateInvoiceForContractIfNotExists tạo hóa đơn tháng đầu khi duyệt HĐ (nếu chưa có)
func (s *InvoiceService) CreateInvoiceForContractIfNotExists(contract *models.Contract) error {
	billingPeriod := time.Now().Format("2006-01")

	var existing models.Invoice
	if err := database.DB.Where("contract_id = ? AND billing_period = ?", contract.ID, billingPeriod).
		First(&existing).Error; err == nil {
		return nil // đã có hóa đơn kỳ này
	}

	dueDate := time.Now().AddDate(0, 0, 15) // hạn TT sau 15 ngày
	if contract.StartDate.After(time.Now()) {
		dueDate = contract.StartDate.AddDate(0, 0, 15)
	}

	invoice := models.Invoice{
		ID:            uuid.New().String(),
		ContractID:    contract.ID,
		ApartmentID:   contract.ApartmentID,
		InvoiceCode:   fmt.Sprintf("INV-%s-%s", billingPeriod, uuid.New().String()[:8]),
		BillingPeriod: billingPeriod,
		DueDate:       dueDate,
		TotalAmount:   contract.MonthlyRent,
		PaidAmount:    0,
		Status:        "pending",
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		return err
	}

	NewNotificationService().NotifyContractBilling(&invoice, contract)
	return nil
}

// GenerateMonthlyInvoices creates invoices for all active contracts
func (s *InvoiceService) GenerateMonthlyInvoices(billingPeriod string) error {
	var contracts []models.Contract

	// Get all active contracts
	if err := database.DB.Where("status = ? AND end_date >= ?", "active", time.Now()).
		Find(&contracts).Error; err != nil {
		return err
	}

	for _, contract := range contracts {
		// Check if invoice already exists for this period
		var existingInvoice models.Invoice
		if err := database.DB.Where("contract_id = ? AND billing_period = ?", contract.ID, billingPeriod).
			First(&existingInvoice).Error; err == nil {
			continue // Invoice already exists
		}

		// Create new invoice
		dueDate := time.Now().AddDate(0, 1, 0)
		invoice := models.Invoice{
			ID:            uuid.New().String(),
			ContractID:    contract.ID,
			ApartmentID:   contract.ApartmentID,
			InvoiceCode:   fmt.Sprintf("INV-%s-%s", billingPeriod, uuid.New().String()[:8]),
			BillingPeriod: billingPeriod,
			DueDate:       dueDate,
			TotalAmount:   contract.MonthlyRent,
			PaidAmount:    0,
			Status:        "pending",
		}

		if err := database.DB.Create(&invoice).Error; err != nil {
			return err
		}

		// Create notification for resident
		notificationService := NewNotificationService()
		notificationService.NotifyContractBilling(&invoice, &contract)
	}

	return nil
}

// CheckOverdueInvoices finds overdue invoices and marks them
func (s *InvoiceService) CheckOverdueInvoices() error {
	var invoices []models.Invoice

	if err := database.DB.Where("status = ? AND due_date < ?", "pending", time.Now()).
		Find(&invoices).Error; err != nil {
		return err
	}

	for _, invoice := range invoices {
		if err := database.DB.Model(&invoice).Update("status", "overdue").Error; err != nil {
			return err
		}

		// Notify resident about overdue payment
		notificationService := NewNotificationService()
		notificationService.NotifyOverduePayment(&invoice)
	}

	return nil
}

// GetDebtorReport returns residents with outstanding payments
func (s *InvoiceService) GetDebtorReport() ([]map[string]interface{}, error) {
	type DebtorInfo struct {
		ResidentID        string
		ResidentName      string
		OutstandingAmount float64
		OverdueCount      int64
		LatestDueDate     time.Time
	}

	var debtors []DebtorInfo
	err := database.DB.Raw(`
		SELECT 
			r.id as resident_id,
			r.full_name as resident_name,
			SUM(i.total_amount - i.paid_amount) as outstanding_amount,
			COUNT(CASE WHEN i.status = 'overdue' THEN 1 END) as overdue_count,
			MAX(i.due_date) as latest_due_date
		FROM residents r
		JOIN contracts c ON r.id = c.resident_id
		JOIN invoices i ON c.id = i.contract_id
		WHERE i.status IN ('pending', 'overdue')
		GROUP BY r.id, r.full_name
		ORDER BY outstanding_amount DESC
	`).Scan(&debtors).Error

	if err != nil {
		return nil, err
	}

	// Convert to map for consistency
	result := make([]map[string]interface{}, len(debtors))
	for i, debtor := range debtors {
		result[i] = map[string]interface{}{
			"resident_id":        debtor.ResidentID,
			"resident_name":      debtor.ResidentName,
			"outstanding_amount": debtor.OutstandingAmount,
			"overdue_count":      debtor.OverdueCount,
			"latest_due_date":    debtor.LatestDueDate,
		}
	}

	return result, nil
}
