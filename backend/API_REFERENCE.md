# Vinhomes Property Management - API Reference

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

All protected endpoints require a JWT token in the `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

## Response Format

All responses are in JSON format:

### Success Response
```json
{
  "data": {},
  "message": "Success"
}
```

### Error Response
```json
{
  "error": "Error message description"
}
```

## Authentication Endpoints

### 1. User Registration
**POST** `/auth/register`

Request:
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password_123",
  "role": "resident"
}
```

Response:
```json
{
  "message": "User created successfully",
  "user_id": "uuid"
}
```

Roles: `admin`, `manager`, `staff`, `resident`

---

### 2. User Login
**POST** `/auth/login`

Request:
```json
{
  "username": "john_doe",
  "password": "secure_password_123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": "uuid",
  "role": "resident",
  "username": "john_doe"
}
```

---

### 3. Get Current User Profile
**GET** `/profile`

Headers:
```
Authorization: Bearer <token>
```

Response:
```json
{
  "id": "uuid",
  "username": "john_doe",
  "email": "john@example.com",
  "role": "resident",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

## Building Management

### 1. Create Building
**POST** `/buildings` (Requires: admin, manager)

Request:
```json
{
  "name": "Vinhomes West Point A",
  "address": "Tay Ho, Hanoi",
  "year_built": 2023,
  "total_floors": 40,
  "description": "Luxury apartment complex"
}
```

Response:
```json
{
  "id": "uuid",
  "name": "Vinhomes West Point A",
  "total_floors": 40,
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

### 2. List Buildings
**GET** `/buildings`

Response:
```json
[
  {
    "id": "uuid",
    "name": "Vinhomes West Point A",
    "address": "Tay Ho, Hanoi",
    "total_floors": 40,
    "floors": []
  }
]
```

---

### 3. Get Building Details
**GET** `/buildings/:id`

Response:
```json
{
  "id": "uuid",
  "name": "Vinhomes West Point A",
  "floors": [
    {
      "id": "uuid",
      "floor_number": 1,
      "apartments": []
    }
  ]
}
```

---

### 4. Update Building
**PUT** `/buildings/:id` (Requires: admin, manager)

Request:
```json
{
  "name": "Updated Building Name"
}
```

---

### 5. Delete Building
**DELETE** `/buildings/:id` (Requires: admin)

---

## Floor Management

### 1. Create Floor
**POST** `/floors` (Requires: admin, manager)

Request:
```json
{
  "building_id": "uuid",
  "floor_number": 5,
  "floor_height": 3.5,
  "floor_area": 1200,
  "purpose": "residential"
}
```

---

### 2. List Building Floors
**GET** `/buildings/:building_id/floors`

Response:
```json
[
  {
    "id": "uuid",
    "floor_number": 5,
    "floor_area": 1200,
    "apartments": []
  }
]
```

---

## Apartment Management

### 1. Create Apartment
**POST** `/apartments` (Requires: admin, manager)

Request:
```json
{
  "floor_id": "uuid",
  "apartment_code": "A05.01",
  "area": 85.5,
  "bedrooms": 2,
  "direction": "NE",
  "listing_price": 800000000,
  "current_status": "empty",
  "body_id": "apartment_mesh_001"
}
```

---

### 2. List Apartments
**GET** `/apartments?floor_id=uuid`

Query Parameters:
- `floor_id` (optional): Filter by floor

---

### 3. Get Apartment Details
**GET** `/apartments/:id`

---

### 4. Update Apartment Status
**PUT** `/apartments/:id/status` (Requires: admin, manager)

Request:
```json
{
  "status": "rented",
  "reason": "Lease signed"
}
```

Statuses: `empty`, `rented`, `maintenance`

---

## Contract Management

### 1. Create Contract
**POST** `/contracts` (Requires: admin, manager)

Request:
```json
{
  "apartment_id": "uuid",
  "resident_id": "uuid",
  "start_date": "2024-02-01T00:00:00Z",
  "end_date": "2025-02-01T00:00:00Z",
  "monthly_rent": 8000000,
  "payment_period": 1,
  "deposit": 16000000
}
```

---

### 2. List Contracts
**GET** `/contracts?status=active&resident_id=uuid`

Query Parameters:
- `status` (optional): `pending`, `active`, `expired`, `terminated`
- `resident_id` (optional): Filter by resident

---

### 3. Get Contract Details
**GET** `/contracts/:id`

---

### 4. Approve Contract
**POST** `/contracts/:id/approve` (Requires: admin, manager)

Request:
```json
{
  "approved": true,
  "notes": "All documents verified"
}
```

---

### 5. Extend Contract
**POST** `/contracts/:id/extend` (Requires: admin, manager)

Request:
```json
{
  "new_end_date": "2026-02-01T00:00:00Z"
}
```

---

### 6. Terminate Contract
**POST** `/contracts/:id/terminate` (Requires: admin, manager)

Request:
```json
{
  "reason": "Resident requested early termination"
}
```

---

## Resident Management

### 1. Create Resident
**POST** `/residents` (Requires: admin, manager)

Request:
```json
{
  "full_name": "Nguyen Van A",
  "id_card": "001234567890",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "gender": "M",
  "phone_number": "0987654321",
  "email": "nguyenvana@example.com",
  "permanent_address": "123 Nguyen Hue St, Hanoi",
  "username": "nguyenvana",
  "password": "secure_password_123"
}
```

---

### 2. List Residents
**GET** `/residents`

---

### 3. Get Resident Details
**GET** `/residents/:id`

---

### 4. Update Resident Profile
**PUT** `/residents/:id`

Request:
```json
{
  "phone_number": "0988888888",
  "email": "newemail@example.com"
}
```

---

### 5. Add Dependent
**POST** `/residents/:resident_id/dependents`

Request:
```json
{
  "full_name": "Nguyen Van B",
  "id_card": "001234567891",
  "date_of_birth": "1992-05-20T00:00:00Z",
  "relationship": "spouse",
  "contract_id": "uuid"
}
```

---

### 6. Remove Dependent
**DELETE** `/dependents/:id`

---

### 7. Register Vehicle
**POST** `/residents/:resident_id/vehicles`

Request:
```json
{
  "license_plate": "29A-12345",
  "vehicle_type": "car",
  "color": "black"
}
```

Vehicle Types: `car`, `motorcycle`, `bicycle`

---

### 8. Get Resident Vehicles
**GET** `/residents/:resident_id/vehicles`

---

### 9. Unregister Vehicle
**DELETE** `/vehicles/:id`

---

## Invoice Management

### 1. Create Invoice
**POST** `/invoices` (Requires: admin, manager)

Request:
```json
{
  "contract_id": "uuid",
  "apartment_id": "uuid",
  "billing_period": "2024-01",
  "due_date": "2024-02-05T00:00:00Z",
  "total_amount": 8000000,
  "details": {
    "management_fee": 500000,
    "water_fee": 100000,
    "electricity_fee": 200000,
    "parking_fee": 150000
  }
}
```

---

### 2. List Invoices
**GET** `/invoices?status=pending&contract_id=uuid`

Query Parameters:
- `status` (optional): `pending`, `paid`, `overdue`, `cancelled`
- `contract_id` (optional): Filter by contract

---

### 3. Get Invoice Details
**GET** `/invoices/:id`

---

### 4. Record Payment
**POST** `/invoices/:invoice_id/payments` (Requires: staff, manager)

Request:
```json
{
  "amount": 8000000,
  "method": "transfer",
  "reference": "BANK_TRANSFER_12345"
}
```

Payment Methods: `cash`, `transfer`, `card`

---

### 5. Get Overdue Invoices
**GET** `/invoices/overdue`

---

### 6. Generate Financial Report
**POST** `/reports/financial`

Request:
```json
{
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-31T23:59:59Z"
}
```

Response:
```json
{
  "total_invoices": 150,
  "total_amount": 1200000000,
  "total_paid": 1100000000,
  "total_outstanding": 100000000,
  "collection_rate": 91.67
}
```

---

## Maintenance Management

### 1. Create Maintenance Request
**POST** `/maintenance-requests`

Request:
```json
{
  "apartment_id": "uuid",
  "resident_id": "uuid",
  "issue_type": "electrical",
  "description": "Light switches not working",
  "priority": "normal",
  "location": "living room"
}
```

Priority: `urgent`, `normal`, `low`

---

### 2. List Maintenance Requests
**GET** `/maintenance-requests?status=new&priority=urgent`

Query Parameters:
- `status` (optional): `new`, `assigned`, `in_progress`, `completed`
- `priority` (optional): `urgent`, `normal`, `low`
- `apartment_id` (optional): Filter by apartment

---

### 3. Get Request Details
**GET** `/maintenance-requests/:id`

---

### 4. Assign Request
**POST** `/maintenance-requests/:id/assign` (Requires: admin, manager)

Request:
```json
{
  "assigned_to": "staff_uuid"
}
```

---

### 5. Update Status
**PUT** `/maintenance-requests/:id/status`

Request:
```json
{
  "status": "in_progress"
}
```

---

### 6. Complete Request
**POST** `/maintenance-requests/:id/complete`

Request:
```json
{
  "completion_notes": "Replaced light switches, problem resolved"
}
```

---

## Equipment Management

### 1. Register Equipment
**POST** `/equipment` (Requires: admin, manager)

Request:
```json
{
  "building_id": "uuid",
  "equipment_name": "Main Water Pump",
  "equipment_type": "pump",
  "manufacturer": "Grundfos",
  "installation_date": "2023-01-15T00:00:00Z",
  "maintenance_period": 12
}
```

---

### 2. Get Building Equipment
**GET** `/buildings/:building_id/equipment`

---

### 3. Get Equipment Due for Maintenance
**GET** `/equipment/due-for-maintenance`

---

### 4. Update Maintenance Record
**PUT** `/equipment/:id/maintenance` (Requires: staff, manager)

---

## Notification Management

### 1. Create Notification
**POST** `/notifications` (Requires: admin, manager)

Request:
```json
{
  "user_id": "uuid",
  "title": "Lease Expiration Notice",
  "message": "Your lease will expire in 30 days",
  "type": "contract_expiring"
}
```

---

### 2. Get User Notifications
**GET** `/notifications`

---

### 3. Get Unread Notifications
**GET** `/notifications/unread`

---

### 4. Mark as Read
**PUT** `/notifications/:id/read`

---

### 5. Mark All as Read
**PUT** `/notifications/read-all`

---

### 6. Delete Notification
**DELETE** `/notifications/:id`

---

## Dashboard

### 1. Get KPI Dashboard
**GET** `/dashboard/kpis`

Response:
```json
{
  "total_apartments": 500,
  "occupied_apartments": 450,
  "vacant_apartments": 50,
  "occupancy_rate": 90.0,
  "active_contracts": 450,
  "expiring_contracts_count": 15,
  "total_residents": 1200,
  "outstanding_invoices": 50000000,
  "collection_rate": 95.5,
  "pending_maintenance_requests": 8,
  "avg_response_time": 4.5
}
```

---

### 2. Get Audit Logs
**GET** `/audit-logs` (Requires: admin)

Query Parameters:
- `user_id` (optional): Filter by user
- `entity` (optional): Filter by entity type

---

## Error Handling

Common HTTP Status Codes:
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error Response Example:
```json
{
  "error": "Invalid request: apartment not found"
}
```

---

## Rate Limiting

Currently no rate limiting. To be implemented in future versions.

---

## Pagination

Currently no pagination. Add query parameters `page` and `limit` in future versions.

---

## Filtering & Search

Most list endpoints support query parameters for filtering.

Example:
```
GET /contracts?status=active&resident_id=uuid
```
