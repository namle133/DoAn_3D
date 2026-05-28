# Maintenance Request Management - Quick Reference

## Quick Links
- **Page**: `/pages/maintenance.html`
- **API Module**: `js/api/maintenance.js`
- **Constants**: `js/utils/constants.js` → `MAINTENANCE_STATUS`, `MAINTENANCE_PRIORITY`

## Key Features at a Glance

### Creating a Request
1. Click "Tạo yêu cầu" button
2. Fill form: Apartment (search), Issue type, Priority, Description, Optional: photos + assignment
3. Submit → Request created

### Assigning a Request
1. Click "Phân công" button on request
2. Select staff member from dropdown
3. Add optional notes
4. Submit → Request status changes to "Assigned"

### Completing a Request
1. Click "Hoàn thành" button on request
2. Fill: Resolution (required), Notes, Cost, Technician name, Photos
3. Submit → Request status changes to "Completed"

### Cancelling a Request
1. Click "Hủy" button on request
2. Confirm cancellation
3. Request status changes to "Cancelled"

## Priority Levels
```
low      → Green badge
medium   → Orange badge
high     → Orange-red badge
urgent   → Red badge
```

## Status States
```
new          → Blue (initial state)
assigned     → Purple (assigned to staff)
in_progress  → Orange (being worked on)
completed    → Green (finished)
cancelled    → Gray (cancelled)
```

## Common Tasks

### Get all requests
```javascript
maintenanceAPI.getMaintenanceRequests()
```

### Create request
```javascript
maintenanceAPI.createMaintenanceRequest({
  apartment_id: "apt-id",
  issue_type: "plumbing",
  priority: "high",
  description: "Water leak in bathroom",
  assigned_to: "staff-id" // optional
})
```

### Assign to staff
```javascript
maintenanceAPI.assignMaintenanceRequest(requestId, {
  assigned_to: "staff-id",
  assigned_to_name: "John Doe",
  notes: "Assign note"
})
```

### Complete request
```javascript
maintenanceAPI.completeMaintenanceRequest(requestId, {
  resolution: "Fixed the leak",
  notes: "Additional notes",
  cost: 500000,
  technician_name: "John Doe",
  photos: [base64_images]
})
```

### Update status
```javascript
maintenanceAPI.updateMaintenanceStatus(requestId, "in_progress")
```

## Issue Types
- `plumbing` - Cấp nước / Thoát nước
- `electrical` - Điện / Điện tử
- `structural` - Kết cấu / Vách tường
- `appliances` - Thiết bị / Nội thất
- `cleaning` - Vệ sinh / Dọn dẹp
- `other` - Khác

## Data Structure (Request Object)
```javascript
{
  id: "uuid",
  request_number: "REQ-001",
  apartment_id: "apt-id",
  apartment_number: "A101",
  issue_type: "plumbing",
  priority: "high",
  status: "assigned",
  description: "Issue details...",
  assigned_to: "staff-id",
  assigned_to_name: "Staff Name",
  resolution: "How it was fixed...",
  cost: 500000,
  technician_name: "Tech Name",
  photos: ["data:image/..."],
  created_date: "2024-05-27T...",
  completed_date: "2024-05-28T..."
}
```

## KPI Dashboard
- **Yêu cầu mới**: Count where status = "new"
- **Đang xử lý**: Count where status = "in_progress"
- **Hoàn thành**: Count where status = "completed"
- **Khẩn cấp**: Count where priority = "urgent"

## Filters Available
1. **Search**: Apartment number, issue type, description
2. **Status**: new, assigned, in_progress, completed, cancelled
3. **Priority**: low, medium, high, urgent

## Public API Methods
```javascript
window.maintenanceApp = {
  selectApartment(apartmentId, apartmentNumber),
  removePhoto(previewId, index),
  viewDetails(requestId),
  openAssignModal(requestId),
  openCompleteModal(requestId),
  cancelRequest(requestId)
}
```

## Required Permissions
- **Admin**: Create, assign, complete, cancel any request
- **Manager**: Create, assign, complete, cancel any request
- **Staff**: View and update assigned requests
- **Resident**: View own requests (if implemented)

## Form Validations
- **Apartment**: Required, must select from dropdown
- **Issue Type**: Required
- **Priority**: Required (defaults to "low")
- **Description**: Required, at least 1 character
- **Photos**: Optional, multiple image files

## Error Handling
- Form validation errors display in red alert below form
- API errors show as toast notifications
- Failed submissions prevent modal close
- User can correct and resubmit

## Responsive Breakpoints
- Desktop: Full table display
- Tablet (768px): Compressed columns, stacked buttons
- Mobile (< 480px): Scrollable table, button stacking

## Local State Variables
```javascript
allRequests        // All maintenance requests
filteredRequests   // Filtered by search/filters
allApartments      // All apartments
allStaff          // All staff members
currentUser       // Logged-in user
currentRole       // User role
selectedRequestId // Current modal's request
createPhotos      // Photos for new request
completePhotos    // Photos for completed request
```

## Event Listeners Setup
- Search input: Debounced 300ms
- Status/Priority filters: Immediate
- Form submissions: Validation then API call
- Modal close buttons: Clear forms and modals
- Outside click on modal: Closes modal
- File uploads: Convert to base64

## Performance Tips
- Search is debounced (300ms) to reduce API calls
- Photos converted to base64 for easy transmission
- Staff/apartments cached in memory
- KPI calculations done client-side
- No pagination (consider for large datasets)

## Debugging Tips
1. Check browser console for errors
2. Check API response in Network tab
3. Verify user has correct role permissions
4. Check apartment data has apartment_number field
5. Verify staff members have staff/manager role
6. Check photo file size and format

## Browser Support
- Chrome/Edge: Full support
- Firefox: Full support
- Safari: Full support (including iOS)
- IE 11: Not supported (uses ES6 modules)

## Accessibility
- Color-blind friendly badges (text + background color)
- Semantic HTML with proper labels
- Keyboard navigation support
- ARIA labels on buttons
- Error messages displayed clearly

---
**Last Updated**: May 27, 2024
**Status**: Production Ready
