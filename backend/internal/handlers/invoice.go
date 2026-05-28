package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/datatypes"
)

type CreateInvoiceRequest struct {
	ContractID    string         `json:"contract_id" binding:"required"`
	ApartmentID   string         `json:"apartment_id" binding:"required"`
	BillingPeriod string         `json:"billing_period" binding:"required"`
	DueDate       time.Time      `json:"due_date" binding:"required"`
	TotalAmount   float64        `json:"total_amount" binding:"required"`
	Details       datatypes.JSON `json:"details"`
}

type RecordPaymentRequest struct {
	Amount    float64 `json:"amount" binding:"required"`
	Method    string  `json:"method" binding:"required,oneof=cash transfer card"`
	Reference string  `json:"reference,omitempty"`
}

type FinancialReportRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Status    string    `json:"status,omitempty"`
}

// ========== Invoice Handlers ==========

func CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify contract and apartment exist
	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", req.ContractID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	invoice := models.Invoice{
		ID:            uuid.New().String(),
		ContractID:    req.ContractID,
		ApartmentID:   req.ApartmentID,
		InvoiceCode:   fmt.Sprintf("INV-%s-%s", req.BillingPeriod, time.Now().Format("150405")),
		BillingPeriod: req.BillingPeriod,
		DueDate:       req.DueDate,
		TotalAmount:   req.TotalAmount,
		PaidAmount:    0,
		Status:        "pending",
		Details:       req.Details,
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

func GetInvoices(c *gin.Context) {
	var invoices []models.Invoice

	query := database.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if contractID := c.Query("contract_id"); contractID != "" {
		query = query.Where("contract_id = ?", contractID)
	}

	if err := query.Preload("Payments").Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

func GetInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice

	if err := database.DB.Preload("Payments").First(&invoice, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func RecordPayment(c *gin.Context) {
	invoiceID := c.Param("id")
	var req RecordPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var invoice models.Invoice
	if err := database.DB.First(&invoice, "id = ?", invoiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	userID, _ := c.Get("user_id")

	// Record payment
	payment := models.InvoicePayment{
		ID:         uuid.New().String(),
		InvoiceID:  invoiceID,
		Amount:     req.Amount,
		Method:     req.Method,
		Reference:  req.Reference,
		PaidAt:     time.Now(),
		RecordedBy: userID.(string),
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment"})
		return
	}

	// Update invoice
	invoice.PaidAmount += req.Amount
	if invoice.PaidAmount >= invoice.TotalAmount {
		invoice.Status = "paid"
	}

	database.DB.Save(&invoice)

	c.JSON(http.StatusCreated, payment)
}

func GetOverdueInvoices(c *gin.Context) {
	var invoices []models.Invoice

	// Find invoices past due date with pending status
	if err := database.DB.Where("status = ? AND due_date < ?", "pending", time.Now()).
		Preload("Apartment").Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch overdue invoices"})
		return
	}

	// Build response with apartment_code
	type InvoiceResponse struct {
		ID            string    `json:"id"`
		InvoiceCode   string    `json:"invoice_code"`
		ApartmentCode string    `json:"apartment_code"`
		TotalAmount   float64   `json:"total_amount"`
		DueDate       time.Time `json:"due_date"`
		Status        string    `json:"status"`
	}

	var response []InvoiceResponse
	for _, inv := range invoices {
		apartmentCode := "-"
		if inv.Apartment.ApartmentCode != "" {
			apartmentCode = inv.Apartment.ApartmentCode
		}
		response = append(response, InvoiceResponse{
			ID:            inv.ID,
			InvoiceCode:   inv.InvoiceCode,
			ApartmentCode: apartmentCode,
			TotalAmount:   inv.TotalAmount,
			DueDate:       inv.DueDate,
			Status:        inv.Status,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetFinancialReport(c *gin.Context) {
	var req FinancialReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type FinancialSummary struct {
		TotalInvoices    int64   `json:"total_invoices"`
		TotalAmount      float64 `json:"total_amount"`
		TotalPaid        float64 `json:"total_paid"`
		TotalOutstanding float64 `json:"total_outstanding"`
		PaidCount        int64   `json:"paid_count"`
		PendingCount     int64   `json:"pending_count"`
		OverdueCount     int64   `json:"overdue_count"`
		CollectionRate   float64 `json:"collection_rate"`
	}

	var invoices []models.Invoice
	query := database.DB.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Find(&invoices)

	var summary FinancialSummary
	database.DB.Model(&models.Invoice{}).
		Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate).
		Select("COUNT(*) as total_invoices, SUM(total_amount) as total_amount, SUM(paid_amount) as total_paid").
		Scan(&summary)

	summary.TotalOutstanding = summary.TotalAmount - summary.TotalPaid
	database.DB.Model(&models.Invoice{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", req.StartDate, req.EndDate, "paid").
		Count(&summary.PaidCount)
	database.DB.Model(&models.Invoice{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", req.StartDate, req.EndDate, "pending").
		Count(&summary.PendingCount)
	database.DB.Model(&models.Invoice{}).
		Where("created_at BETWEEN ? AND ? AND status = ? AND due_date < NOW()", req.StartDate, req.EndDate, "pending").
		Count(&summary.OverdueCount)

	if summary.TotalAmount > 0 {
		summary.CollectionRate = (summary.TotalPaid / summary.TotalAmount) * 100
	}

	c.JSON(http.StatusOK, summary)
}
