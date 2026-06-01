package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type InvoiceLineItem struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

type UpdateInvoiceRequest struct {
	DueDate       time.Time         `json:"due_date" binding:"required"`
	BillingPeriod string            `json:"billing_period"`
	TotalAmount   float64           `json:"total_amount" binding:"required"`
	LineItems     []InvoiceLineItem `json:"line_items"`
	Description   string            `json:"description"`
	Status        string            `json:"status"`
}

type InvoiceDetailResponse struct {
	ID                string             `json:"id"`
	InvoiceCode       string             `json:"invoice_code"`
	InvoiceNumber     string             `json:"invoice_number"`
	ContractID        string             `json:"contract_id"`
	ApartmentID       string             `json:"apartment_id"`
	ResidentID        string             `json:"resident_id"`
	ResidentName      string             `json:"resident_name"`
	ApartmentCode     string             `json:"apartment_code"`
	BillingPeriod     string             `json:"billing_period"`
	DueDate           time.Time          `json:"due_date"`
	TotalAmount       float64            `json:"total_amount"`
	PaidAmount        float64            `json:"paid_amount"`
	OutstandingAmount float64            `json:"outstanding_amount"`
	Status            string             `json:"status"`
	Description       string             `json:"description"`
	LineItems         []InvoiceLineItem  `json:"line_items"`
	CreatedAt         time.Time          `json:"created_at"`
	Payments          []models.InvoicePayment `json:"payments,omitempty"`
}

func parseInvoiceDetailsJSON(details datatypes.JSON) ([]InvoiceLineItem, string) {
	var payload map[string]interface{}
	if len(details) == 0 {
		return nil, ""
	}
	if err := json.Unmarshal(details, &payload); err != nil {
		return nil, ""
	}
	desc, _ := payload["description"].(string)
	var items []InvoiceLineItem
	if raw, ok := payload["line_items"].([]interface{}); ok {
		for _, r := range raw {
			m, ok := r.(map[string]interface{})
			if !ok {
				continue
			}
			item := InvoiceLineItem{
				Description: fmt.Sprint(m["description"]),
			}
			if q, ok := m["quantity"].(float64); ok {
				item.Quantity = q
			}
			if p, ok := m["price"].(float64); ok {
				item.Price = p
			}
			if t, ok := m["total"].(float64); ok {
				item.Total = t
			} else {
				item.Total = item.Quantity * item.Price
			}
			items = append(items, item)
		}
	}
	return items, desc
}

func buildInvoiceDetail(inv models.Invoice) InvoiceDetailResponse {
	items, desc := parseInvoiceDetailsJSON(inv.Details)
	outstanding := inv.TotalAmount - inv.PaidAmount
	if outstanding < 0 {
		outstanding = 0
	}
	resp := InvoiceDetailResponse{
		ID: inv.ID, InvoiceCode: inv.InvoiceCode, InvoiceNumber: inv.InvoiceCode,
		ContractID: inv.ContractID, ApartmentID: inv.ApartmentID,
		BillingPeriod: inv.BillingPeriod, DueDate: inv.DueDate,
		TotalAmount: inv.TotalAmount, PaidAmount: inv.PaidAmount,
		OutstandingAmount: outstanding, Status: inv.Status,
		Description: desc, LineItems: items, CreatedAt: inv.CreatedAt,
		Payments: inv.Payments,
	}
	return resp
}

type CreateInvoiceRequest struct {
	ContractID    string            `json:"contract_id"`
	ApartmentID   string            `json:"apartment_id" binding:"required"`
	ResidentID    string            `json:"resident_id"`
	BillingPeriod string            `json:"billing_period"`
	DueDate       time.Time         `json:"due_date" binding:"required"`
	TotalAmount   float64           `json:"total_amount" binding:"required"`
	LineItems     []InvoiceLineItem `json:"line_items"`
	Description   string            `json:"description"`
	Status        string            `json:"status"`
}

type RecordPaymentRequest struct {
	Amount        float64 `json:"amount"`
	AmountPaid    float64 `json:"amount_paid"`
	Method        string  `json:"method"`
	PaymentMethod string  `json:"payment_method"`
	Reference     string  `json:"reference"`
	ReferenceNum  string  `json:"reference_number"`
}

type FinancialReportRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Status    string `json:"status,omitempty"`
}

func parseReportDate(s string, endOfDay bool) (time.Time, error) {
	var t time.Time
	var err error
	if t, err = time.Parse("2006-01-02", s); err != nil {
		t, err = time.Parse(time.RFC3339, s)
	}
	if err != nil {
		return t, err
	}
	if endOfDay {
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local), nil
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local), nil
}

// overdueInvoiceQuery — HĐ quá hạn: status overdue HOẶC pending đã quá due_date, còn nợ > 0
func overdueInvoiceQuery() *gorm.DB {
	now := time.Now()
	return database.DB.Model(&models.Invoice{}).
		Where("(status = ? OR (status = ? AND due_date < ?))", "overdue", "pending", now).
		Where("(total_amount - paid_amount) > 0")
}

func applyInvoiceStatusFilter(query *gorm.DB, status string) *gorm.DB {
	if status == "" {
		return query
	}
	now := time.Now()
	switch status {
	case "overdue":
		return query.Where("(status = ? OR (status = ? AND due_date < ?))", "overdue", "pending", now)
	case "pending":
		return query.Where("status = ? AND due_date >= ?", "pending", now)
	default:
		return query.Where("status = ?", status)
	}
}

// ========== Invoice Handlers ==========

func CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tìm HĐ active nếu chưa gửi contract_id
	contractID := req.ContractID
	if contractID == "" {
		var contract models.Contract
		q := database.DB.Where("apartment_id = ? AND status IN ?", req.ApartmentID, []string{"active", "extended"})
		if req.ResidentID != "" {
			q = q.Where("resident_id = ?", req.ResidentID)
		}
		if err := q.First(&contract).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy hợp đồng active cho căn hộ/cư dân này. Vui lòng duyệt hợp đồng trước."})
			return
		}
		contractID = contract.ID
		if req.ResidentID == "" {
			req.ResidentID = contract.ResidentID
		}
	}

	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", contractID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	billingPeriod := req.BillingPeriod
	if billingPeriod == "" {
		billingPeriod = time.Now().Format("2006-01")
	}

	detailsMap := map[string]interface{}{
		"line_items":  req.LineItems,
		"description": req.Description,
		"resident_id": req.ResidentID,
	}
	detailsJSON, _ := json.Marshal(detailsMap)

	status := req.Status
	if status == "" || status == "issued" || status == "draft" {
		status = "pending"
	}

	invoice := models.Invoice{
		ID:            uuid.New().String(),
		ContractID:    contractID,
		ApartmentID:   req.ApartmentID,
		InvoiceCode:   fmt.Sprintf("INV-%s-%s", billingPeriod, uuid.New().String()[:8]),
		BillingPeriod: billingPeriod,
		DueDate:       req.DueDate,
		TotalAmount:   req.TotalAmount,
		PaidAmount:    0,
		Status:        status,
		Details:       datatypes.JSON(detailsJSON),
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": invoice, "success": true})
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

	// Cư dân chỉ xem HĐ của mình
	if role, ok := c.Get("role"); ok && role.(string) == "resident" {
		userID, _ := c.Get("user_id")
		var resident models.Resident
		if err := database.DB.Where("user_id = ?", userID).First(&resident).Error; err == nil {
			query = query.Joins("JOIN contracts ON contracts.id = invoices.contract_id").
				Where("contracts.resident_id = ?", resident.ID)
		}
	}

	if err := query.Preload("Apartment").Preload("Contract.Resident").Preload("Payments").
		Order("created_at DESC").Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
		return
	}

	type InvoiceListItem struct {
		ID              string    `json:"id"`
		InvoiceCode     string    `json:"invoice_code"`
		InvoiceNumber   string    `json:"invoice_number"`
		ContractID      string    `json:"contract_id"`
		ApartmentID     string    `json:"apartment_id"`
		ResidentID        string    `json:"resident_id"`
		ResidentName      string    `json:"resident_name"`
		ApartmentNumber   string    `json:"apartment_number"`
		ApartmentCode     string    `json:"apartment_code"`
		BillingPeriod   string    `json:"billing_period"`
		TotalAmount     float64   `json:"total_amount"`
		PaidAmount      float64   `json:"paid_amount"`
		DueDate         time.Time `json:"due_date"`
		Status          string    `json:"status"`
		CreatedAt       time.Time `json:"created_at"`
	}

	result := make([]InvoiceListItem, 0, len(invoices))
	for _, inv := range invoices {
		item := InvoiceListItem{
			ID: inv.ID, InvoiceCode: inv.InvoiceCode, InvoiceNumber: inv.InvoiceCode,
			ContractID: inv.ContractID, ApartmentID: inv.ApartmentID,
			BillingPeriod: inv.BillingPeriod, TotalAmount: inv.TotalAmount,
			PaidAmount: inv.PaidAmount, DueDate: inv.DueDate, Status: inv.Status,
			CreatedAt: inv.CreatedAt,
		}
		if inv.Apartment.ApartmentCode != "" {
			item.ApartmentNumber = inv.Apartment.ApartmentCode
			item.ApartmentCode = inv.Apartment.ApartmentCode
		}
		if inv.Contract.Resident.ID != "" {
			item.ResidentID = inv.Contract.Resident.ID
			item.ResidentName = inv.Contract.Resident.FullName
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice

	query := database.DB.Preload("Payments").Preload("Apartment").Preload("Contract.Resident")
	if role, ok := c.Get("role"); ok && role.(string) == "resident" {
		userID, _ := c.Get("user_id")
		var resident models.Resident
		if err := database.DB.Where("user_id = ?", userID).First(&resident).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		query = query.Joins("JOIN contracts ON contracts.id = invoices.contract_id").
			Where("invoices.id = ? AND contracts.resident_id = ?", id, resident.ID)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	resp := buildInvoiceDetail(invoice)
	if invoice.Apartment.ApartmentCode != "" {
		resp.ApartmentCode = invoice.Apartment.ApartmentCode
	}
	if invoice.Contract.Resident.FullName != "" {
		resp.ResidentName = invoice.Contract.Resident.FullName
		resp.ResidentID = invoice.Contract.Resident.ID
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var req UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var invoice models.Invoice
	if err := database.DB.First(&invoice, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	if invoice.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể sửa hóa đơn đã thanh toán đủ"})
		return
	}
	if invoice.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể sửa hóa đơn đã hủy"})
		return
	}
	if req.TotalAmount < invoice.PaidAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tổng tiền không được nhỏ hơn số đã thu"})
		return
	}

	_, descOld := parseInvoiceDetailsJSON(invoice.Details)
	desc := req.Description
	if desc == "" {
		desc = descOld
	}

	billingPeriod := req.BillingPeriod
	if billingPeriod == "" {
		billingPeriod = invoice.BillingPeriod
	}

	status := req.Status
	if status == "" || status == "issued" || status == "draft" {
		status = invoice.Status
	}
	if status != "pending" && status != "overdue" && status != "cancelled" {
		status = "pending"
	}
	if status == "pending" && req.DueDate.Before(time.Now()) {
		status = "overdue"
	}

	detailsMap := map[string]interface{}{
		"line_items":  req.LineItems,
		"description": desc,
	}
	var oldPayload map[string]interface{}
	if json.Unmarshal(invoice.Details, &oldPayload) == nil {
		if rid, ok := oldPayload["resident_id"].(string); ok {
			detailsMap["resident_id"] = rid
		}
	}
	detailsJSON, _ := json.Marshal(detailsMap)

	updates := map[string]interface{}{
		"billing_period": billingPeriod,
		"due_date":       req.DueDate,
		"total_amount":   req.TotalAmount,
		"status":         status,
		"details":        datatypes.JSON(detailsJSON),
		"updated_at":     time.Now(),
	}

	if err := database.DB.Model(&invoice).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}

	database.DB.Preload("Payments").Preload("Apartment").Preload("Contract.Resident").First(&invoice, "id = ?", id)
	resp := buildInvoiceDetail(invoice)
	if invoice.Apartment.ApartmentCode != "" {
		resp.ApartmentCode = invoice.Apartment.ApartmentCode
	}
	if invoice.Contract.Resident.FullName != "" {
		resp.ResidentName = invoice.Contract.Resident.FullName
	}

	c.JSON(http.StatusOK, gin.H{"data": resp, "success": true})
}

func RecordPayment(c *gin.Context) {
	invoiceID := c.Param("id")
	var req RecordPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount := req.Amount
	if amount == 0 && req.AmountPaid > 0 {
		amount = req.AmountPaid
	}
	method := req.Method
	if method == "" {
		method = req.PaymentMethod
	}
	if method == "bank_transfer" {
		method = "transfer"
	}
	reference := req.Reference
	if reference == "" {
		reference = req.ReferenceNum
	}
	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Số tiền thanh toán không hợp lệ"})
		return
	}
	if method == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn phương thức thanh toán"})
		return
	}

	var invoice models.Invoice
	if err := database.DB.First(&invoice, "id = ?", invoiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	userID, _ := c.Get("user_id")

	payment := models.InvoicePayment{
		ID:         uuid.New().String(),
		InvoiceID:  invoiceID,
		Amount:     amount,
		Method:     method,
		Reference:  reference,
		PaidAt:     time.Now(),
		RecordedBy: userID.(string),
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment"})
		return
	}

	invoice.PaidAmount += amount
	if invoice.PaidAmount >= invoice.TotalAmount {
		invoice.Status = "paid"
	}

	database.DB.Save(&invoice)

	c.JSON(http.StatusCreated, gin.H{"data": payment, "success": true})
}

func GetOverdueInvoices(c *gin.Context) {
	var invoices []models.Invoice
	now := time.Now()

	query := overdueInvoiceQuery().
		Preload("Apartment").
		Preload("Contract.Resident").
		Order("due_date ASC")

	// Cư dân chỉ xem HĐ quá hạn của mình
	if role, ok := c.Get("role"); ok && role.(string) == "resident" {
		userID, _ := c.Get("user_id")
		var resident models.Resident
		if err := database.DB.Where("user_id = ?", userID).First(&resident).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
			return
		}
		query = query.Joins("JOIN contracts ON contracts.id = invoices.contract_id").
			Where("contracts.resident_id = ?", resident.ID)
	}

	if err := query.Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch overdue invoices"})
		return
	}

	type InvoiceResponse struct {
		ID                 string    `json:"id"`
		InvoiceCode        string    `json:"invoice_code"`
		ApartmentCode      string    `json:"apartment_code"`
		ResidentName       string    `json:"resident_name"`
		BillingPeriod      string    `json:"billing_period"`
		TotalAmount        float64   `json:"total_amount"`
		PaidAmount         float64   `json:"paid_amount"`
		OutstandingAmount  float64   `json:"outstanding_amount"`
		DueDate            time.Time `json:"due_date"`
		Status             string    `json:"status"`
		DaysOverdue        int       `json:"days_overdue"`
	}

	response := make([]InvoiceResponse, 0, len(invoices))
	for _, inv := range invoices {
		outstanding := inv.TotalAmount - inv.PaidAmount
		if outstanding <= 0 {
			continue
		}
		daysOverdue := int(now.Sub(inv.DueDate).Hours() / 24)
		if daysOverdue < 0 {
			daysOverdue = 0
		}
		status := inv.Status
		if status == "pending" && inv.DueDate.Before(now) {
			status = "overdue"
		}
		row := InvoiceResponse{
			ID: inv.ID, InvoiceCode: inv.InvoiceCode, BillingPeriod: inv.BillingPeriod,
			TotalAmount: inv.TotalAmount, PaidAmount: inv.PaidAmount,
			OutstandingAmount: outstanding, DueDate: inv.DueDate, Status: status,
			DaysOverdue: daysOverdue,
		}
		if inv.Apartment.ApartmentCode != "" {
			row.ApartmentCode = inv.Apartment.ApartmentCode
		}
		if inv.Contract.Resident.FullName != "" {
			row.ResidentName = inv.Contract.Resident.FullName
		}
		response = append(response, row)
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func GetFinancialReport(c *gin.Context) {
	var req FinancialReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := parseReportDate(req.StartDate, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date không hợp lệ (YYYY-MM-DD)"})
		return
	}
	endDate, err := parseReportDate(req.EndDate, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date không hợp lệ (YYYY-MM-DD)"})
		return
	}
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date phải sau start_date"})
		return
	}

	now := time.Now()
	baseQuery := database.DB.Model(&models.Invoice{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate)
	filteredQuery := applyInvoiceStatusFilter(baseQuery, req.Status)

	type FinancialSummary struct {
		TotalInvoices       int64   `json:"total_invoices"`
		TotalAmount         float64 `json:"total_amount"`
		TotalPaid           float64 `json:"total_paid"`
		TotalOutstanding    float64 `json:"total_outstanding"`
		PaidCount           int64   `json:"paid_count"`
		PendingCount        int64   `json:"pending_count"`
		OverdueCount        int64   `json:"overdue_count"`
		CollectionRate      float64 `json:"collection_rate"`
		ExpectedRentRevenue float64 `json:"expected_rent_revenue"`
		TotalDebtors        int64   `json:"total_debtors"`
		TotalDebtAmount     float64 `json:"total_debt_amount"`
	}

	var summary FinancialSummary
	filteredQuery.Select("COUNT(*) as total_invoices, COALESCE(SUM(total_amount),0) as total_amount, COALESCE(SUM(paid_amount),0) as total_paid").
		Scan(&summary)
	summary.TotalOutstanding = summary.TotalAmount - summary.TotalPaid

	periodBase := database.DB.Model(&models.Invoice{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)
	if req.Status == "" {
		periodBase.Where("status = ?", "paid").Count(&summary.PaidCount)
		database.DB.Model(&models.Invoice{}).
			Where("created_at BETWEEN ? AND ? AND status = ? AND due_date >= ?", startDate, endDate, "pending", now).
			Count(&summary.PendingCount)
		database.DB.Model(&models.Invoice{}).
			Where("created_at BETWEEN ? AND ?", startDate, endDate).
			Where("(status = ? OR (status = ? AND due_date < ?))", "overdue", "pending", now).
			Count(&summary.OverdueCount)
	} else if req.Status == "paid" {
		summary.PaidCount = summary.TotalInvoices
	} else if req.Status == "pending" {
		summary.PendingCount = summary.TotalInvoices
	} else if req.Status == "overdue" {
		summary.OverdueCount = summary.TotalInvoices
	}

	if summary.TotalAmount > 0 {
		summary.CollectionRate = (summary.TotalPaid / summary.TotalAmount) * 100
	}

	database.DB.Model(&models.Contract{}).
		Where("status IN ? AND start_date <= ? AND end_date >= ?", []string{"active", "extended"}, endDate, startDate).
		Select("COALESCE(SUM(monthly_rent), 0)").Row().Scan(&summary.ExpectedRentRevenue)

	type DebtorRow struct {
		ResidentID        string    `json:"resident_id"`
		ResidentName      string    `json:"resident_name"`
		ApartmentCode     string    `json:"apartment_code"`
		OutstandingAmount float64   `json:"outstanding_amount"`
		OverdueCount      int64     `json:"overdue_count"`
		InvoiceCount      int64     `json:"invoice_count"`
		LatestDueDate     time.Time `json:"latest_due_date"`
	}
	var debtors []DebtorRow
	database.DB.Raw(`
		SELECT r.id as resident_id, r.full_name as resident_name,
			COALESCE(MAX(a.apartment_code), '') as apartment_code,
			COALESCE(SUM(i.total_amount - i.paid_amount), 0) as outstanding_amount,
			COUNT(CASE WHEN i.status = 'overdue' OR (i.status = 'pending' AND i.due_date < ?) THEN 1 END) as overdue_count,
			COUNT(i.id) as invoice_count,
			MAX(i.due_date) as latest_due_date
		FROM residents r
		JOIN contracts c ON r.id = c.resident_id
		JOIN invoices i ON c.id = i.contract_id
		LEFT JOIN apartments a ON i.apartment_id = a.id
		WHERE i.status IN ('pending', 'overdue') AND (i.total_amount - i.paid_amount) > 0
		GROUP BY r.id, r.full_name
		ORDER BY outstanding_amount DESC
	`, now).Scan(&debtors)

	summary.TotalDebtors = int64(len(debtors))
	for _, d := range debtors {
		summary.TotalDebtAmount += d.OutstandingAmount
	}

	type MonthlyRow struct {
		BillingPeriod string  `json:"billing_period"`
		InvoiceCount  int64   `json:"invoice_count"`
		Issued        float64 `json:"issued"`
		Collected     float64 `json:"collected"`
		Outstanding   float64 `json:"outstanding"`
	}
	var monthly []MonthlyRow
	mq := database.DB.Model(&models.Invoice{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate)
	mq = applyInvoiceStatusFilter(mq, req.Status)
	mq.Select(`billing_period,
		COUNT(*) as invoice_count,
		COALESCE(SUM(total_amount), 0) as issued,
		COALESCE(SUM(paid_amount), 0) as collected,
		COALESCE(SUM(total_amount - paid_amount), 0) as outstanding`).
		Group("billing_period").Order("billing_period ASC").Scan(&monthly)

	type InvoiceReportRow struct {
		ID            string    `json:"id"`
		InvoiceCode   string    `json:"invoice_code"`
		BillingPeriod string    `json:"billing_period"`
		ApartmentCode string    `json:"apartment_code"`
		ResidentName  string    `json:"resident_name"`
		TotalAmount   float64   `json:"total_amount"`
		PaidAmount    float64   `json:"paid_amount"`
		Status        string    `json:"status"`
		DueDate       time.Time `json:"due_date"`
		CreatedAt     time.Time `json:"created_at"`
	}
	var invoices []models.Invoice
	invQ := database.DB.Preload("Apartment").Preload("Contract.Resident").
		Where("created_at BETWEEN ? AND ?", startDate, endDate)
	invQ = applyInvoiceStatusFilter(invQ, req.Status)
	invQ.Order("created_at DESC").Limit(100).Find(&invoices)

	invoiceRows := make([]InvoiceReportRow, 0, len(invoices))
	for _, inv := range invoices {
		row := InvoiceReportRow{
			ID: inv.ID, InvoiceCode: inv.InvoiceCode, BillingPeriod: inv.BillingPeriod,
			TotalAmount: inv.TotalAmount, PaidAmount: inv.PaidAmount,
			Status: inv.Status, DueDate: inv.DueDate, CreatedAt: inv.CreatedAt,
		}
		if inv.Apartment.ApartmentCode != "" {
			row.ApartmentCode = inv.Apartment.ApartmentCode
		}
		if inv.Contract.Resident.FullName != "" {
			row.ResidentName = inv.Contract.Resident.FullName
		}
		if row.Status == "pending" && inv.DueDate.Before(now) {
			row.Status = "overdue"
		}
		invoiceRows = append(invoiceRows, row)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary":           summary,
			"debtors":           debtors,
			"monthly_breakdown": monthly,
			"invoices":          invoiceRows,
			"period": gin.H{
				"start_date": startDate,
				"end_date":   endDate,
			},
		},
		"success": true,
	})
}
