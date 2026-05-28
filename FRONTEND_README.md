# Vinhomes Property Management System - Frontend

## 📋 Overview

This is a comprehensive frontend for the **Vinhomes Property Management System**, a modern GIS 3D-based solution for managing residential properties, contracts, invoices, maintenance, and residents.

The frontend integrates with a Go + Gin backend and provides four distinct user roles:
- **Admin** - System administration
- **Manager** - Property management operations
- **Staff** - Data entry and daily operations
- **Resident** - Portal for residents

## ✨ Features Implemented

### ✅ Phase 1: Foundation (Complete)
- [x] **API Client** - JWT token management, error handling, request/response interceptors
- [x] **Authentication** - Login, registration, session management
- [x] **Constants** - All app constants in one place (roles, statuses, colors, storage keys)
- [x] **Utilities** - Date formatting, currency formatting, validation, helpers
- [x] **API Wrappers** - 8 modular API client modules covering 40+ endpoints

### ✅ Phase 2: Core UI (In Progress)
- [x] **Login Page** - Professional login/register with validation
- [x] **Dashboard** - KPI cards, recent activities, overdue invoices, maintenance requests
- [ ] **Resident Management** - CRUD for residents, dependents, vehicles
- [ ] **Contract Management** - Create, approve, extend, terminate contracts
- [ ] **Invoice Management** - Create, record payments, financial reports
- [ ] **Maintenance Management** - Request tracking, assignment, completion

### ⏳ Phase 3: GIS Integration (Planned)
- [ ] Connect 3D model to real building data
- [ ] Color-code apartments by status
- [ ] Interactive apartment selection
- [ ] Real-time 3D model updates

### ⏳ Phase 4: Advanced Features (Planned)
- [ ] Notification system
- [ ] Excel/PDF reports
- [ ] Audit logs
- [ ] Equipment management

## 📁 Project Structure

```
frontend/
├── js/
│   ├── api/                    # API endpoint wrappers
│   │   ├── client.js           # Base HTTP client with JWT
│   │   ├── auth.js             # Authentication endpoints
│   │   ├── buildings.js        # Building/floor/apartment endpoints
│   │   ├── residents.js        # Resident management
│   │   ├── contracts.js        # Contract management
│   │   ├── invoices.js         # Invoice management
│   │   ├── maintenance.js      # Maintenance & equipment
│   │   ├── notifications.js    # Notification system
│   │   └── reports.js          # Reports & dashboard
│   ├── modules/
│   │   ├── auth-manager.js     # Authentication business logic
│   │   ├── dashboard.js        # Dashboard logic (planned)
│   │   ├── user-manager.js     # User management (planned)
│   │   └── ... (more modules)
│   ├── utils/
│   │   ├── constants.js        # App constants
│   │   ├── helpers.js          # Utility functions
│   │   └── validators.js       # Input validation
│   └── ui/
│       ├── components.js       # Reusable UI components (planned)
│       └── modals.js           # Modal dialogs (planned)
├── pages/
│   ├── login.html              # Login page
│   ├── dashboard.html          # Dashboard
│   ├── residents.html          # (Planned)
│   ├── contracts.html          # (Planned)
│   ├── invoices.html           # (Planned)
│   └── maintenance.html        # (Planned)
├── css/
│   ├── style.css               # Existing 3D styles
│   ├── dashboard.css           # (Planned)
│   └── responsive.css          # (Planned)
└── index.html                  # Main 3D page (existing)
```

## 🚀 Quick Start

### 1. Ensure Backend is Running
```bash
cd backend
go run ./cmd/main.go
# Server running on http://localhost:8080
```

### 2. Open Frontend in Browser
```bash
# Simply open index.html or serve with a local server
# Using Python:
python3 -m http.server 8000
# Or using Node.js:
npx http-server
```

Then visit: `http://localhost:8000/pages/login.html`

### 3. Login with Demo Credentials
- **Admin:** admin@vinhomes.com / admin123
- **Manager:** manager@vinhomes.com / manager123
- **Staff:** staff@vinhomes.com / staff123
- **Resident:** resident@vinhomes.com / resident123

## 🔐 Authentication

### Token Management
- Tokens are stored in **localStorage**
- JWT token sent in `Authorization: Bearer <token>` header
- Automatic token refresh on 401 response (planned)
- Auto-logout on token expiration

### User Session
- User profile stored in localStorage
- User role checked for RBAC
- Session cleared on logout

## 📡 API Integration

### Available API Endpoints (40+)

#### Authentication
- `POST /auth/login` - User login
- `POST /auth/register` - User registration

#### Buildings
- `GET /buildings` - Get all buildings
- `GET /buildings/:id` - Get building details
- `GET /buildings/:id/floors` - Get building floors
- `POST /buildings` - Create building (Admin/Manager)
- `PUT /buildings/:id` - Update building (Admin/Manager)
- `DELETE /buildings/:id` - Delete building (Admin)

#### Apartments
- `GET /apartments` - Get all apartments
- `GET /apartments/:id` - Get apartment details
- `POST /apartments` - Create apartment (Admin/Manager)
- `PUT /apartments/:id/status` - Update apartment status

#### Residents
- `GET /residents` - Get all residents
- `GET /residents/:id` - Get resident profile
- `POST /residents` - Create resident
- `PUT /residents/:id` - Update resident profile
- `POST /residents/:id/dependents` - Add dependent
- `DELETE /dependents/:id` - Remove dependent
- `POST /residents/:id/vehicles` - Register vehicle
- `GET /residents/:id/vehicles` - Get resident vehicles
- `DELETE /vehicles/:id` - Unregister vehicle

#### Contracts
- `GET /contracts` - Get all contracts
- `GET /contracts/:id` - Get contract details
- `POST /contracts` - Create contract
- `POST /contracts/:id/approve` - Approve contract
- `POST /contracts/:id/extend` - Extend contract
- `POST /contracts/:id/terminate` - Terminate contract

#### Invoices
- `GET /invoices` - Get all invoices
- `GET /invoices/:id` - Get invoice details
- `POST /invoices` - Create invoice (Admin/Manager)
- `POST /invoices/:id/payments` - Record payment
- `GET /invoices/overdue` - Get overdue invoices
- `POST /reports/financial` - Get financial report

#### Maintenance
- `GET /maintenance-requests` - Get all requests
- `GET /maintenance-requests/:id` - Get request details
- `POST /maintenance-requests` - Create request
- `POST /maintenance-requests/:id/assign` - Assign request
- `PUT /maintenance-requests/:id/status` - Update status
- `POST /maintenance-requests/:id/complete` - Complete request

#### Equipment
- `POST /equipment` - Register equipment (Admin/Manager)
- `GET /buildings/:id/equipment` - Get building equipment
- `GET /equipment/due-for-maintenance` - Get due equipment
- `PUT /equipment/:id/maintenance` - Update maintenance record

#### Notifications
- `GET /notifications` - Get all notifications
- `GET /notifications/unread` - Get unread notifications
- `POST /notifications` - Create notification (Admin/Manager)
- `PUT /notifications/:id/read` - Mark as read
- `PUT /notifications/read-all` - Mark all as read
- `DELETE /notifications/:id` - Delete notification

#### Reports
- `GET /dashboard/kpis` - Get dashboard KPIs
- `GET /audit-logs` - Get audit logs (Admin)

## 🎨 UI Components (Planned)

### Forms
- Resident form
- Contract form
- Invoice form
- Maintenance request form
- Equipment registration form

### Modals
- Confirm dialogs
- Detail views
- Status update dialogs
- Payment recording

### Tables
- Data tables with sorting/filtering
- Pagination
- Bulk actions
- Export to CSV/PDF

## 🔄 Real-time Updates

### Current Implementation
- Data refreshed via polling (30 seconds)
- Manual refresh buttons available

### Planned Improvements
- WebSocket support for real-time updates
- Event-driven notifications
- Live status synchronization

## 🎯 GIS 3D Integration (Planned)

### Features
- Apartments colored by status (empty, rented, maintenance, reserved)
- Click apartment for details
- Highlight maintenance issues
- Highlight overdue invoices
- Spatial queries (filter by floor, building)

### Color Mapping
- **Gray (0xcccccc)** - Empty
- **Blue (0x0099ff)** - Rented
- **Yellow (0xffcc00)** - Maintenance
- **Orange (0xff9900)** - Reserved

## 📊 Data Types & Constants

### Apartment Status
- `empty` - Available for rent
- `rented` - Currently rented
- `maintenance` - Under maintenance
- `reserved` - Reserved for future rent

### Contract Status
- `pending` - Awaiting approval
- `active` - Currently active
- `extended` - Extension approved
- `terminated` - Contract ended
- `expired` - Contract expired

### Invoice Status
- `draft` - In preparation
- `issued` - Sent to resident
- `overdue` - Payment overdue
- `paid` - Payment received
- `cancelled` - Invoice cancelled

### Maintenance Priority
- `low` - Non-urgent
- `medium` - Normal
- `high` - Urgent
- `urgent` - Critical, immediate attention needed

## 🔒 Security Features

- [x] JWT authentication
- [x] Role-based access control (RBAC)
- [x] Protected routes with role checking
- [x] Secure token storage
- [x] CORS headers (configured on backend)
- [ ] Rate limiting (planned)
- [ ] Audit logging (planned)

## 📱 Responsive Design

- ✅ Mobile-first approach
- ✅ Desktop optimized (1400px max-width)
- ✅ Tablet support
- ✅ Touch-friendly buttons

## ⚙️ Configuration

### API Base URL
Edit `js/utils/constants.js`:
```javascript
const API_CONFIG = {
  BASE_URL: 'http://localhost:8080/api/v1', // Change if needed
  TIMEOUT: 10000,
};
```

### Role Labels
Customize role display names in each page script.

## 🧪 Testing Workflow

1. **Login** with different roles
2. **Dashboard** loads KPIs and data
3. **Data fetching** - Verify API calls in Network tab
4. **Sorting/Filtering** - Test table features
5. **Form submission** - Create/update operations
6. **Error handling** - Test with invalid data
7. **GIS 3D** - Verify apartment updates

## 🐛 Troubleshooting

### "Unauthorized" Error
- Check login credentials
- Verify backend is running
- Clear browser storage: `localStorage.clear()`
- Check JWT token in browser console: `localStorage.getItem('vinhomes_auth_token')`

### "Cannot GET /api/v1/..." Error
- Backend not running
- Wrong API base URL
- Backend port not 8080

### CORS Errors
- Backend CORS middleware not configured
- Check backend headers

### Module Import Errors
- Use `.js` extension in imports
- Ensure server supports ES6 modules
- Use `type="module"` in script tags

## 📚 Code Examples

### Login
```javascript
import authManager from './js/modules/auth-manager.js';

await authManager.login('user@example.com', 'password123');
authManager.redirectToDashboard();
```

### API Call
```javascript
import buildingAPI from './js/api/buildings.js';

const buildings = await buildingAPI.getBuildings();
console.log(buildings.data);
```

### RBAC Check
```javascript
import authManager from './js/modules/auth-manager.js';

if (authManager.isManager()) {
  // Show manager-only features
}
```

### Format Currency
```javascript
import { formatCurrency } from './js/utils/helpers.js';

const price = formatCurrency(1000000); // 1,000,000 ₫
```

## 🚧 Roadmap

### Week 1 (Current)
- [x] API client & authentication
- [x] Login page
- [x] Dashboard skeleton
- [ ] Complete dashboard data loading

### Week 2
- [ ] Resident CRUD pages
- [ ] Contract management
- [ ] Invoice management

### Week 3
- [ ] Maintenance system
- [ ] Notifications
- [ ] Reports & export

### Week 4
- [ ] GIS 3D integration
- [ ] Real-time updates
- [ ] Mobile optimization

## 🤝 Contributing

Please follow these guidelines:
1. Use modular structure (API → Module → Page)
2. Add error handling with `showToast()`
3. Validate all inputs
4. Add comments for complex logic
5. Keep components reusable

## 📄 License

All rights reserved. Vinhomes Property Management System.

## ✉️ Support

For issues or questions:
1. Check troubleshooting section
2. Review browser console for errors
3. Check Network tab for API issues
4. Contact development team

---

**Status:** In Active Development 🚀

Last Updated: 2026-05-27
Version: 0.2.0 - Alpha
