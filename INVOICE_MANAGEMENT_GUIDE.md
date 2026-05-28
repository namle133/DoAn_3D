# Invoice Management Page - User & Developer Guide

## 📋 Table of Contents
1. [User Guide](#user-guide)
2. [Developer Guide](#developer-guide)
3. [API Integration](#api-integration)
4. [Customization](#customization)
5. [Troubleshooting](#troubleshooting)

---

## User Guide

### Accessing the Invoice Management Page
1. Navigate to `pages/invoices.html` from the Vinhomes system
2. The page requires authentication (will redirect to login if not authenticated)
3. Your user role and name will be displayed in the top-right navbar

### Dashboard Overview
The page displays financial metrics at the top:
- **Total Revenue**: Sum of all invoice amounts
- **Overdue Amount**: Sum of unpaid invoices past their due date (highlighted in red)
- **Paid Amount**: Sum of completed payments
- **Total Invoices**: Total count of all invoices

### Searching and Filtering

#### Search Bar
- Search for invoices by **Invoice Number** or **Resident Name**
- Search is real-time and updates as you type
- Case-insensitive matching

#### Status Filter
Select from dropdown to filter invoices:
- **All statuses** (default)
- **Draft**: Not yet issued invoices
- **Issued**: Active unpaid invoices
- **Overdue**: Unpaid invoices past due date
- **Paid**: Completed invoices
- **Cancelled**: Cancelled invoices

### Managing Invoices

#### Creating a New Invoice
1. Click the **"+ Tạo hóa đơn"** (Create Invoice) button
2. Select a **Resident** from the dropdown
3. Select the **Apartment** (auto-filtered for selected resident)
4. Set the **Invoice Date** (default: today)
5. Set the **Due Date** (default: 7 days from invoice date)
6. Add **Line Items**:
   - Click "+ Thêm hạng mục" to add items
   - Enter description (e.g., "Rent", "Utilities")
   - Enter quantity
   - Enter unit price
   - Total is calculated automatically
   - Remove items with the "Xóa" button
7. (Optional) Add a **Note** describing the invoice
8. Click **"Lưu hóa đơn"** to save

#### Recording a Payment
1. Find the invoice in the table
2. Click the **"Thanh toán"** (Payment) button
3. Enter **Payment Method**:
   - Cash
   - Bank Transfer
   - Check
   - Online
4. Enter **Payment Date** (default: today)
5. Enter **Amount Paid** (default: full invoice amount)
6. (Optional) Enter **Reference Number** (e.g., transaction ID)
7. (Optional) Add payment **Notes**
8. Click **"Ghi nhận thanh toán"** to record

#### Viewing Invoice Details
1. Click the **"Xem"** (View) button on an invoice
2. A popup displays:
   - Invoice Number
   - Resident Name
   - Apartment Number
   - Invoice Amount
   - Current Status

#### Editing an Invoice
- Click **"Sửa"** (Edit) button
- Available only for Draft and Issued invoices
- *(Foundation laid for future development)*

#### Deleting an Invoice
- Click **"Xóa"** (Delete) button
- Only available for Draft and Issued invoices
- *(Foundation laid for future development)*

### Understanding Invoice Status

| Status | Color | Description |
|--------|-------|-------------|
| Draft | Blue | Not yet issued, can be edited/deleted |
| Issued | Purple | Active invoice awaiting payment |
| Overdue | Red | Payment deadline has passed |
| Paid | Green | Payment received in full |
| Cancelled | Gray | Voided invoice |

**Overdue Highlighting**: Rows with overdue unpaid invoices appear with a light red background for quick identification.

---

## Developer Guide

### File Location
```
/Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/pages/invoices.html
```

### Architecture Overview

#### State Management
```javascript
// Global state variables
let currentUser = null;           // Current authenticated user
let currentRole = null;           // User role (admin, manager, staff, resident)
let allInvoices = [];            // All invoices from API
let allResidents = [];           // All residents from API
let allApartments = [];          // All apartments from API
let filteredInvoices = [];       // Filtered invoices for display
let editingInvoiceId = null;     // Currently editing invoice ID
let lineItemsData = [];          // Line items in form
```

#### Module Imports
```javascript
import authManager from '../js/modules/auth-manager.js';
import invoiceAPI from '../js/api/invoices.js';
import residentAPI from '../js/api/residents.js';
import buildingAPI from '../js/api/buildings.js';
import { formatDate, formatCurrency, showToast, debounce } from '../js/utils/helpers.js';
import { INVOICE_STATUS, ROLES } from '../js/utils/constants.js';
```

### Key Functions

#### Initialization
```javascript
async function init()
```
- Guards authentication
- Sets user information in navbar
- Loads data from APIs
- Sets up event listeners
- Initializes UI

#### Data Loading
```javascript
async function loadInvoices()           // Fetch all invoices
async function loadResidents()          // Fetch all residents
async function loadApartments()         // Fetch all apartments
```

#### UI Rendering
```javascript
function renderInvoices(invoices)       // Render invoice table
function populateResidentDropdown()     // Populate resident dropdown
function populateApartmentDropdown()    // Populate apartment dropdown
function updateKPIs()                   // Calculate and update KPI cards
```

#### Event Handlers
```javascript
function filterInvoices()               // Handle search and filter
async function handleCreateInvoice()    // Handle form submission
async function handleRecordPayment()    // Handle payment recording
function addLineItem()                  // Add line item to form
function removeLineItem(itemId)         // Remove line item
function updateLineItemTotals()         // Recalculate line item totals
```

#### Modal Management
```javascript
function openAddInvoiceModal()           // Open create invoice modal
function closeInvoiceModal()             // Close create invoice modal
function openPaymentModal(invoiceId)    // Open payment recording modal
function closePaymentModal()             // Close payment modal
```

### CSS Classes

#### Layout
```css
.navbar                                 /* Navigation bar */
.container                              /* Main content container */
.page-header                            /* Page title area */
.kpi-grid                               /* KPI cards grid */
.kpi-card                               /* Individual KPI card */
.filters-bar                            /* Search and filter bar */
.section                                /* Content section */
.table-responsive                       /* Responsive table wrapper */
```

#### Modals
```css
.modal                                  /* Modal overlay */
.modal.show                             /* Visible modal */
.modal-content                          /* Modal content box */
.modal-header                           /* Modal title area */
.form-group                             /* Form field wrapper */
.form-row                               /* Two-column form layout */
.form-buttons                           /* Form action buttons */
```

#### Status Badges
```css
.status-draft                           /* Draft status */
.status-issued                          /* Issued status */
.status-overdue                         /* Overdue status */
.status-paid                            /* Paid status */
.status-cancelled                       /* Cancelled status */
.tr.overdue                             /* Overdue row highlight */
```

#### Line Items
```css
.line-items                             /* Line items container */
.line-item                              /* Individual line item */
.line-item-total                        /* Line item total display */
```

### Form Validation

#### Invoice Form Validation
```javascript
- Resident: Required, must be selected
- Apartment: Required, must be selected
- Invoice Date: Required, valid date
- Due Date: Required, valid date, must be after invoice date
- Line Items: At least 1 item with positive amount
- Total Amount: Must be > 0
```

#### Payment Form Validation
```javascript
- Payment Method: Required, must be selected
- Payment Date: Required, valid date
- Amount Paid: Required, must be > 0
```

### Error Handling

#### Try-Catch Pattern
```javascript
try {
  // API call
  const response = await invoiceAPI.createInvoice(data);
  // Handle success
  showToast('Success message', 'success');
} catch (error) {
  console.error('Error context:', error);
  showFormError(elementId, error.message);
}
```

#### User-Facing Errors
- Form validation errors displayed inline
- API errors shown in alert/modal
- Toast notifications for quick feedback
- Console logging for debugging

### Customization Guide

#### Adding New Line Item Types
Modify `addLineItem()` to include custom default descriptions:
```javascript
function addLineItem() {
  // ... existing code ...
  lineItem.innerHTML = `
    <input type="text" placeholder="..." class="item-description" 
           value="${lineItemsData.length === 0 ? 'Rent' : ''}" />
    // ... rest of input ...
  `;
}
```

#### Changing Status Labels
Modify the `getStatusLabel()` function:
```javascript
function getStatusLabel(status) {
  const labels = {
    [INVOICE_STATUS.DRAFT]: 'Draft Status',
    // ... add more ...
  };
  return labels[status] || status;
}
```

#### Adding New Payment Methods
Update the payment method dropdown:
```html
<select id="paymentMethod" name="paymentMethod">
  <option value="">-- Select Method --</option>
  <option value="cash">Cash</option>
  <option value="bank_transfer">Bank Transfer</option>
  <option value="crypto">Cryptocurrency</option>
</select>
```

#### Styling Customization
All styles are in `<style>` block at the top of HTML:
- Color schemes easily adjustable
- Responsive breakpoints at line ~650
- CSS variables can be added for easier theming

---

## API Integration

### Required API Endpoints

#### Invoice API (`js/api/invoices.js`)
```javascript
// GET /invoices - Get all invoices
invoiceAPI.getInvoices()

// GET /invoices/{id} - Get specific invoice
invoiceAPI.getInvoiceByID(id)

// POST /invoices - Create new invoice
invoiceAPI.createInvoice(data)

// POST /invoices/{invoiceID}/payments - Record payment
invoiceAPI.recordPayment(invoiceID, data)

// GET /invoices/overdue - Get overdue invoices
invoiceAPI.getOverdueInvoices()
```

#### Resident API (`js/api/residents.js`)
```javascript
// GET /residents - Get all residents
residentAPI.getResidents()
```

#### Building API (`js/api/buildings.js`)
```javascript
// GET /apartments - Get all apartments
buildingAPI.getApartments()
```

### API Request/Response Format

#### Create Invoice Request
```javascript
{
  resident_id: "res_123",
  apartment_id: "apt_456",
  invoice_date: "2024-05-27",
  due_date: "2024-06-03",
  total_amount: 5000000,
  line_items: [
    {
      description: "Rent",
      quantity: 1,
      price: 5000000,
      total: 5000000
    }
  ],
  description: "Monthly rent invoice",
  status: "issued"
}
```

#### Record Payment Request
```javascript
{
  payment_method: "bank_transfer",
  payment_date: "2024-05-27",
  amount_paid: 5000000,
  reference_number: "TXN123456",
  notes: "Payment received via bank transfer"
}
```

---

## Troubleshooting

### Page Not Loading
**Problem**: Page shows blank or doesn't render
**Solutions**:
1. Check browser console for errors (F12)
2. Verify all imported modules exist in correct paths
3. Check if running on HTTPS (some browsers require this)
4. Clear browser cache and reload

### Authentication Redirect
**Problem**: Redirected to login page
**Solutions**:
1. Log in with valid credentials
2. Check if session token is expired
3. Verify localStorage has auth token: `localStorage.getItem('vinhomes_auth_token')`

### Data Not Loading
**Problem**: Tables show "Đang tải dữ liệu..." indefinitely
**Solutions**:
1. Check API endpoint URLs in network tab (F12 → Network)
2. Verify backend API is running
3. Check API response status codes (200 = success)
4. Look for CORS errors in console
5. Check if API returns correct data structure

### Form Submission Fails
**Problem**: Cannot create invoice or record payment
**Solutions**:
1. Check all required fields are filled (marked with *)
2. Verify date inputs are valid
3. Check line items have positive amounts
4. Look at browser console for specific errors
5. Verify resident and apartment are properly selected

### Calculation Errors
**Problem**: Invoice totals incorrect
**Solutions**:
1. Check line item inputs (quantity × price)
2. Ensure all line items are removed except active ones
3. Verify `updateLineItemTotals()` is called
4. Check browser console for JavaScript errors

### UI Styling Issues
**Problem**: Layout broken or styles not applied
**Solutions**:
1. Hard refresh browser (Ctrl+Shift+R or Cmd+Shift+R)
2. Check if CSS is loading in Network tab
3. Verify no browser extensions interfering
4. Try in incognito/private window
5. Check browser console for CSS warnings

### API Errors
**Problem**: API returns 400, 401, 403, 404, or 500
**Solutions**:
- **400**: Invalid request format - check data being sent
- **401**: Unauthorized - check authentication token
- **403**: Forbidden - check user role permissions
- **404**: Endpoint not found - verify API URL
- **500**: Server error - check backend logs

### Performance Issues
**Problem**: Page slow or unresponsive
**Solutions**:
1. Reduce number of invoices displayed (pagination)
2. Optimize API response data
3. Check browser Performance tab (F12)
4. Reduce DOM manipulation
5. Clear browser cache

---

## Support & Maintenance

### Regular Tasks
- Monitor API performance
- Check for deprecated APIs
- Update dependencies
- Test with various data sizes
- Review error logs

### Future Enhancements
- [ ] Invoice PDF export
- [ ] Email notifications
- [ ] Batch operations
- [ ] Invoice templates
- [ ] Payment history
- [ ] Advanced reporting
- [ ] Multi-currency support

### Version History
- v1.0 (2024-05-27): Initial release with full CRUD operations

---

**Last Updated**: May 27, 2024
**Status**: Production Ready
**Maintainer**: Development Team
