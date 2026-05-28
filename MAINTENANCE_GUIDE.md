# Maintenance Request Management System - Implementation Guide

## Overview
The maintenance request management page provides a comprehensive interface for managing maintenance requests in the Vinhomes property management system. This guide covers all features, usage patterns, and technical implementation details.

## File Location
- **Path**: `/pages/maintenance.html`
- **Module Dependencies**: 
  - `js/modules/auth-manager.js` - Authentication & authorization
  - `js/api/maintenance.js` - Maintenance API endpoints
  - `js/api/buildings.js` - Building & apartment data
  - `js/api/residents.js` - Resident & staff data
  - `js/utils/helpers.js` - Utility functions
  - `js/utils/constants.js` - Constants & enums

## Features

### 1. Dashboard Overview
**KPI Cards** display key metrics:
- **Yêu cầu mới (New Requests)**: Count of unassigned requests
- **Đang xử lý (In Progress)**: Count of requests being worked on
- **Hoàn thành (Completed)**: Count of resolved requests
- **Khẩn cấp (Urgent)**: Count of urgent priority requests

### 2. Request Display
The main table displays all maintenance requests with columns:
- **Yêu cầu #**: Unique request identifier
- **Căn hộ**: Apartment number where issue occurred
- **Loại vấn đề**: Issue category (plumbing, electrical, structural, appliances, cleaning, other)
- **Mức độ**: Priority level (Low, Medium, High, Urgent) with color coding
- **Trạng thái**: Current status (New, Assigned, In Progress, Completed, Cancelled) with color coding
- **Phân công cho**: Name of assigned staff member
- **Ngày tạo**: Request creation date
- **Thao tác**: Action buttons

### 3. Filtering & Search
**Search**: Search by apartment number, issue type, or description
**Status Filter**: Filter by request status
**Priority Filter**: Filter by priority level

All filters work together and update KPI metrics in real-time.

### 4. Create Maintenance Request
**Modal Dialog**: Opens when clicking "Tạo yêu cầu" button

**Form Fields**:
- **Căn hộ** (Required): Search/dropdown for apartment selection
- **Loại vấn đề** (Required): Select issue type
- **Mức độ ưu tiên** (Required): Radio buttons for priority
- **Mô tả chi tiết** (Required): Text area for detailed description
- **Phân công cho**: Dropdown to assign to staff member
- **Tải ảnh**: File upload for photos (optional, multiple)

**Validation**:
- All required fields must be filled
- Apartment must be selected from dropdown
- At least one character in description

**Photo Upload**:
- Supports multiple images
- Displays thumbnail previews
- Remove button for each photo
- Files are converted to base64 data URIs

### 5. Request Actions

#### Assign Request (Status: New)
- Open assign modal from action buttons
- Select staff member from dropdown
- Add optional assignment notes
- Updates status to "Assigned"

#### Update Status (Status: Assigned/In Progress)
- Three status options available:
  1. Assign: Move to "Assigned" status
  2. In Progress: Move to "In Progress" status
  3. Complete: Move to "Completed" status

#### Complete Request (Status: Assigned/In Progress)
**Complete Modal Opens** with fields:
- **Mô tả giải pháp** (Required): How the issue was resolved
- **Ghi chú hoàn thành**: Additional completion notes
- **Chi phí**: Cost of maintenance work (optional)
- **Người thực hiện**: Name of technician (pre-filled if assigned)
- **Ảnh hoàn thành công việc**: Photo uploads of completed work

#### Cancel Request (Status: New/Assigned)
- Confirmation dialog before cancellation
- Request status changes to "Cancelled"
- Cannot be reopened

### 6. View Details
- Click "Chi tiết" button to open details modal
- Shows all request information including:
  - Request ID
  - Apartment details
  - Issue type and description
  - Priority and status
  - Assignment information
  - Resolution details (if completed)
  - Cost (if available)
  - Completion date (if completed)

## Status Workflow

```
NEW
  ├─→ [Assign] → ASSIGNED
  │              ├─→ [Start Work] → IN_PROGRESS
  │              │                  ├─→ [Complete] → COMPLETED
  │              │                  └─→ [Cancel] → CANCELLED
  │              └─→ [Cancel] → CANCELLED
  └─→ [Cancel] → CANCELLED
```

## Data Structure

### Maintenance Request Object
```javascript
{
  id: string,                      // Unique identifier
  request_number: string,          // Display number
  apartment_id: string,            // Foreign key to apartment
  apartment_number: string,        // Apartment display number
  issue_type: string,              // plumbing|electrical|structural|appliances|cleaning|other
  description: string,             // Detailed issue description
  priority: string,                // low|medium|high|urgent
  status: string,                  // new|assigned|in_progress|completed|cancelled
  assigned_to: string,             // Staff member ID (nullable)
  assigned_to_name: string,        // Staff member name
  resolution: string,              // How issue was resolved (nullable)
  cost: number,                    // Cost in VND (nullable)
  technician_name: string,         // Name of technician (nullable)
  photos: array,                   // Array of photo data URIs
  created_date: ISO8601,           // Request creation timestamp
  completed_date: ISO8601,         // Request completion timestamp (nullable)
}
```

## API Endpoints Used

### Maintenance API
```javascript
// Get all maintenance requests
maintenanceAPI.getMaintenanceRequests()

// Create new request
maintenanceAPI.createMaintenanceRequest(data)

// Assign request to staff
maintenanceAPI.assignMaintenanceRequest(id, data)

// Update request status
maintenanceAPI.updateMaintenanceStatus(id, status)

// Complete request
maintenanceAPI.completeMaintenanceRequest(id, data)
```

### Building API
```javascript
// Get all apartments
buildingAPI.getApartments()
```

### Resident API
```javascript
// Get residents/staff
residentAPI.getResidents()
```

## Color Coding

### Priority Levels
- **Thấp (Low)**: Green (#2e7d32)
- **Trung bình (Medium)**: Orange (#e65100)
- **Cao (High)**: Orange-Red (#d84315)
- **Khẩn cấp (Urgent)**: Red (#c62828)

### Status Colors
- **Mới (New)**: Blue (#1976d2)
- **Đã phân công (Assigned)**: Purple (#7b1fa2)
- **Đang xử lý (In Progress)**: Orange (#e65100)
- **Hoàn thành (Completed)**: Green (#2e7d32)
- **Hủy bỏ (Cancelled)**: Gray (#666)

## Responsive Design
- Mobile-friendly layout (tested on 768px and below)
- Table scrolling on small screens
- Touch-friendly button sizes
- Modal adjusts to screen size

## Error Handling
- Form validation with user-friendly error messages
- API error handling with toast notifications
- Loading states during async operations
- Empty states when no data available

## User Permissions
- **Admin/Manager**: Full access to all features
- **Staff**: Can view and update assigned requests
- **Resident**: Can view own requests (if implemented)

## Navigation
- Navbar links to Dashboard, Residents, Invoices, Maintenance, Contracts
- Active page indicator
- User profile section with name and role
- Logout button

## Keyboard Support
- Tab navigation through form fields
- Enter to submit forms
- Escape to close modals (if implemented)

## Performance Features
- Debounced search (300ms delay)
- Efficient data filtering
- Lazy loading of staff/apartment data
- Optimized rendering

## Future Enhancements
1. Real-time updates using WebSockets
2. Calendar view for scheduled maintenance
3. Equipment tracking and maintenance history
4. Service provider integration
5. Advanced reporting and analytics
6. Mobile app integration
7. Estimated time of arrival (ETA) tracking
8. Before/after photo comparison
9. Maintenance templates for common issues
10. Automated notification system

## Testing Checklist
- [ ] Create maintenance request with all required fields
- [ ] Search filters work correctly
- [ ] Priority and status color coding displays properly
- [ ] Assign request to staff member
- [ ] Complete request with photos and cost
- [ ] Cancel request confirmation works
- [ ] KPI metrics update correctly
- [ ] Responsive layout on mobile devices
- [ ] Form validation prevents empty submissions
- [ ] Error messages display clearly
- [ ] Empty state displays when no requests exist
- [ ] Apartment search dropdown filters correctly
- [ ] Photo upload and preview works
- [ ] Photo removal works correctly
- [ ] Details modal shows all information
- [ ] Status workflow follows expected transitions

## Troubleshooting

### Request not appearing in list
- Check API response in browser console
- Verify appointment is not filtered out by status/priority
- Check user has correct permissions

### Apartment dropdown not showing
- Ensure apartments are loaded from API
- Check apartment data has apartment_number field
- Verify search term matches apartment name/number

### Photos not uploading
- Check file size is reasonable
- Verify file is image format
- Check browser console for errors

### Assigned staff not showing
- Verify staff members have correct role (staff/manager)
- Check residentAPI returns staff data correctly
- Verify staff names are populated in database

## Support
For issues or feature requests, contact the development team.
