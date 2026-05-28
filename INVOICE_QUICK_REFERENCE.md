# Invoice Management - Quick Reference Guide

## 🚀 Getting Started

### Access the Page
```
http://localhost:3000/pages/invoices.html
```

### Required APIs
Ensure the following endpoints are available:
- `GET /api/v1/invoices` - List all invoices
- `POST /api/v1/invoices` - Create new invoice
- `POST /api/v1/invoices/{id}/payments` - Record payment
- `GET /api/v1/residents` - List residents
- `GET /api/v1/apartments` - List apartments

---

## 📋 Key Components

### State Variables
```javascript
let allInvoices = [];        // All invoices from API
let allResidents = [];       // All residents
let allApartments = [];      // All apartments
let filteredInvoices = [];   // Filtered display
let editingInvoiceId = null; // Current edit ID
let lineItemsData = [];      // Line items temp storage
```

### Main Functions
| Function | Purpose |
|----------|---------|
| `init()` | Initialize page, load data, setup listeners |
| `loadInvoices()` | Fetch invoices from API |
| `loadResidents()` | Fetch residents from API |
| `loadApartments()` | Fetch apartments from API |
| `renderInvoices()` | Render invoice table |
| `filterInvoices()` | Apply search/filter |
| `updateKPIs()` | Calculate KPI values |
| `handleCreateInvoice()` | Submit new invoice |
| `handleRecordPayment()` | Submit payment |

---

## 🎯 Common Tasks

### Create New Invoice
1. User clicks "Tạo hóa đơn"
2. Modal opens with form
3. Select resident → Apartment auto-populates
4. Set dates → Add line items → Save

### Record Payment
1. Find invoice in table
2. Click "Thanh toán" button
3. Modal opens with form
4. Select payment method → Enter amount → Save

### Search & Filter
- **Search**: Type invoice number or resident name (real-time)
- **Filter**: Select status from dropdown (Draft/Issued/Overdue/Paid/Cancelled)

---

## 🎨 CSS Classes

### Layout Classes
- `.navbar` - Navigation bar
- `.container` - Main content area
- `.section` - Content block
- `.kpi-grid` - KPI cards grid
- `.kpi-card` - Individual KPI

### Table Classes
- `.table-responsive` - Responsive wrapper
- `.status-badge` - Status indicator
- `.action-btns` - Action button container
- `.btn-sm` - Small button

### Modal Classes
- `.modal` - Modal overlay
- `.modal.show` - Visible modal
- `.modal-content` - Modal box
- `.form-group` - Form field

### Status Classes
- `.status-draft` - Draft status (blue)
- `.status-issued` - Issued status (purple)
- `.status-overdue` - Overdue status (red)
- `.status-paid` - Paid status (green)
- `.status-cancelled` - Cancelled status (gray)
- `.tr.overdue` - Overdue row highlight

---

## 🔧 Form Validation Rules

### Invoice Form
```javascript
Resident:    Required
Apartment:   Required
Invoice Date: Required, valid date
Due Date:    Required, after invoice date
Line Items:  At least 1 item, amount > 0
```

### Payment Form
```javascript
Payment Method: Required
Payment Date:   Required, valid date
Amount Paid:    Required, > 0
```

---

## 🌐 API Request Examples

### Create Invoice
```javascript
POST /api/v1/invoices
{
  resident_id: "res_123",
  apartment_id: "apt_456",
  invoice_date: "2024-05-27",
  due_date: "2024-06-03",
  total_amount: 5000000,
  line_items: [
    { description: "Rent", quantity: 1, price: 5000000, total: 5000000 }
  ],
  description: "Monthly rent",
  status: "issued"
}
```

### Record Payment
```javascript
POST /api/v1/invoices/{id}/payments
{
  payment_method: "bank_transfer",
  payment_date: "2024-05-27",
  amount_paid: 5000000,
  reference_number: "TXN123456",
  notes: "Bank transfer received"
}
```

---

## 🐛 Debugging Tips

### Check Console
```javascript
// View all invoices
console.log(allInvoices);

// Check filtered results
console.log(filteredInvoices);

// View API calls
// Open DevTools → Network tab
```

### Common Issues
| Issue | Solution |
|-------|----------|
| Page blank | Check console (F12), clear cache |
| No data loads | Verify API endpoints, check network |
| Calculation wrong | Verify `updateLineItemTotals()` called |
| Modal won't close | Check modal element visibility |

---

## 📱 Responsive Breakpoints

```css
Desktop:  1400px+  /* Full layout */
Tablet:   769-1399px /* Adjusted */
Mobile:   <768px   /* Stacked */
```

---

## 🚨 Error Messages

### User-Facing Errors
- "Vui lòng điền đầy đủ thông tin bắt buộc" - Missing required fields
- "Hạn thanh toán phải sau ngày tạo hóa đơn" - Invalid date
- "Vui lòng thêm ít nhất một hạng mục" - No line items
- "Số tiền phải > 0" - Invalid amount

---

## 📊 KPI Calculations

```javascript
Total Revenue = SUM(all invoice amounts)
Overdue Amount = SUM(unpaid invoices where due_date < today)
Paid Amount = SUM(invoices where status = 'paid')
Total Invoices = COUNT(all invoices)
```

---

## 🔒 Security Notes

- All API calls use authentication token
- User role validated via `authManager`
- Form inputs sanitized
- XSS prevention through DOM methods
- CSRF tokens handled by API client

---

## 📞 Support

### For Issues
1. Check browser console (F12)
2. Verify API responses in Network tab
3. Check INVOICE_MANAGEMENT_GUIDE.md
4. Contact development team

### For Enhancement Requests
1. Document requirement
2. Create GitHub issue
3. Tag with `enhancement` label
4. Discuss with team

---

**Version**: 1.0  
**Last Updated**: May 27, 2024  
**Status**: Production Ready
