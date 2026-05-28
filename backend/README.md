# Vinhomes Property Management Backend

## Project Setup

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (optional)
- Git

### Installation

1. **Clone and Setup**
```bash
cd backend
cp config/.env.example config/.env
```

2. **Edit `.env` file** with your configuration

3. **Install dependencies**
```bash
go mod download
go mod tidy
```

4. **Run Database (Docker)**
```bash
docker-compose up -d postgres
```

Or install PostgreSQL locally with PostGIS extension.

5. **Run Application**
```bash
go run ./cmd/main.go
```

Or build and run:
```bash
go build -o vinhomes ./cmd
./vinhomes
```

6. **Docker Compose (Full Stack)**
```bash
docker-compose up -d
```

## API Documentation

### Authentication Endpoints

**POST** `/api/v1/auth/login`
- Login and receive JWT token
- Request: `{ "username": "...", "password": "..." }`

**POST** `/api/v1/auth/register`
- Create new user account
- Request: `{ "username": "...", "email": "...", "password": "...", "role": "..." }`

**GET** `/api/v1/profile`
- Get current user profile (requires authentication)

### Building Management

**POST** `/api/v1/buildings` (admin, manager)
- Create new building
**GET** `/api/v1/buildings`
- List all buildings
**GET** `/api/v1/buildings/:id`
- Get building details
**PUT** `/api/v1/buildings/:id` (admin, manager)
- Update building
**DELETE** `/api/v1/buildings/:id` (admin)
- Delete building

### Floor Management

**POST** `/api/v1/floors` (admin, manager)
- Create floor
**GET** `/api/v1/buildings/:building_id/floors`
- List floors in building

### Apartment Management

**POST** `/api/v1/apartments` (admin, manager)
- Create apartment
**GET** `/api/v1/apartments`
- List apartments (supports `floor_id` filter)
**GET** `/api/v1/apartments/:id`
- Get apartment details
**PUT** `/api/v1/apartments/:id/status` (admin, manager)
- Update apartment status

### Contract Management

**POST** `/api/v1/contracts` (admin, manager)
- Create contract
**GET** `/api/v1/contracts`
- List contracts (supports `status`, `resident_id` filters)
**GET** `/api/v1/contracts/:id`
- Get contract details
**POST** `/api/v1/contracts/:id/approve` (admin, manager)
- Approve contract
**POST** `/api/v1/contracts/:id/extend` (admin, manager)
- Extend contract
**POST** `/api/v1/contracts/:id/terminate` (admin, manager)
- Terminate contract

### Resident Management

**POST** `/api/v1/residents` (admin, manager)
- Create resident
**GET** `/api/v1/residents`
- List residents
**GET** `/api/v1/residents/:id`
- Get resident details
**PUT** `/api/v1/residents/:id`
- Update resident profile
**POST** `/api/v1/residents/:resident_id/dependents`
- Add dependent
**DELETE** `/api/v1/dependents/:id`
- Remove dependent

### Vehicle Management

**POST** `/api/v1/residents/:resident_id/vehicles`
- Register vehicle
**GET** `/api/v1/residents/:resident_id/vehicles`
- List resident vehicles
**DELETE** `/api/v1/vehicles/:id`
- Unregister vehicle

### Invoice Management

**POST** `/api/v1/invoices` (admin, manager)
- Create invoice
**GET** `/api/v1/invoices`
- List invoices (supports `status`, `contract_id` filters)
**GET** `/api/v1/invoices/:id`
- Get invoice details
**POST** `/api/v1/invoices/:invoice_id/payments` (staff, manager)
- Record payment
**GET** `/api/v1/invoices/overdue`
- Get overdue invoices
**POST** `/api/v1/reports/financial`
- Generate financial report

### Maintenance Management

**POST** `/api/v1/maintenance-requests`
- Create maintenance request
**GET** `/api/v1/maintenance-requests`
- List requests (supports `status`, `priority`, `apartment_id` filters)
**GET** `/api/v1/maintenance-requests/:id`
- Get request details
**POST** `/api/v1/maintenance-requests/:id/assign` (admin, manager)
- Assign request
**PUT** `/api/v1/maintenance-requests/:id/status`
- Update status
**POST** `/api/v1/maintenance-requests/:id/complete`
- Complete request

### Equipment Management

**POST** `/api/v1/equipment` (admin, manager)
- Register equipment
**GET** `/api/v1/buildings/:building_id/equipment`
- List building equipment
**GET** `/api/v1/equipment/due-for-maintenance`
- Get equipment needing maintenance
**PUT** `/api/v1/equipment/:id/maintenance` (staff, manager)
- Update maintenance record

### Notification Management

**POST** `/api/v1/notifications` (admin, manager)
- Create notification
**GET** `/api/v1/notifications`
- Get user notifications
**GET** `/api/v1/notifications/unread`
- Get unread notifications
**PUT** `/api/v1/notifications/:id/read`
- Mark as read
**PUT** `/api/v1/notifications/read-all`
- Mark all as read
**DELETE** `/api/v1/notifications/:id`
- Delete notification

### Dashboard & Reports

**GET** `/api/v1/dashboard/kpis`
- Get KPI dashboard
**GET** `/api/v1/audit-logs` (admin only)
- Get audit logs

## Project Structure

```
backend/
├── cmd/
│   └── main.go                    # Application entry point
├── config/
│   ├── config.go                  # Configuration management
│   └── .env.example               # Environment variables template
├── internal/
│   ├── database/
│   │   └── database.go            # Database connection & migrations
│   ├── models/
│   │   └── models.go              # Data models for all entities
│   ├── handlers/
│   │   ├── auth.go                # Authentication handlers
│   │   ├── building.go            # Building/Floor/Apartment handlers
│   │   ├── contract.go            # Contract handlers
│   │   ├── resident.go            # Resident/Dependent/Vehicle handlers
│   │   ├── invoice.go             # Invoice & Financial handlers
│   │   ├── maintenance.go         # Maintenance & Equipment handlers
│   │   └── notification.go        # Notification & Dashboard handlers
│   └── middleware/
│       ├── auth.go                # JWT authentication & authorization
│       └── logger.go              # Request logging & audit
├── go.mod                         # Go module definition
├── go.sum                         # Dependency checksums
├── Dockerfile                     # Container build configuration
├── docker-compose.yml             # Multi-container orchestration
└── README.md                      # This file

## Technology Stack

- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL + PostGIS (spatial database)
- **ORM**: GORM (database abstraction)
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Containerization**: Docker & Docker Compose

## Features Implemented

### M1: Authentication & Account Management
- User registration with role assignment
- JWT-based login authentication
- Password hashing with bcrypt
- User profile management
- Role-based access control (admin, manager, staff, resident)

### M2: GIS 3D Visualization
- Building with geometry (3D coordinates)
- Floor with spatial data
- Apartment with BodyID reference for 3D model
- Status history tracking

### M3: Building/Floor/Apartment Management
- CRUD operations for all entities
- Status tracking with history
- Hierarchical relationships

### M4: Contract Management
- Contract creation and approval workflow
- Contract extension and termination
- Status tracking (pending, active, expired, terminated)
- Payment terms management

### M5: Resident Management
- Resident profile creation
- Dependent management (family members)
- Vehicle registration
- Residence history tracking

### M6: Financial Management
- Invoice generation
- Payment recording
- Overdue tracking
- Financial reporting with KPIs
- Collection rate calculation

### M7: Maintenance Management
- Maintenance request creation
- Priority-based assignment
- Status workflow (new → assigned → in-progress → completed)
- Equipment registration and tracking
- Scheduled maintenance management

### M8: Notifications & Dashboard
- Real-time notifications
- KPI dashboard
- Audit logging
- Unread notification tracking

## Security Features

- JWT token-based authentication
- Role-based access control (RBAC)
- Password hashing with bcrypt
- Audit logging for all operations
- Request validation
- SQL injection prevention (GORM ORM)

## Next Steps

1. **GIS Integration**: Connect with ArcGIS for spatial queries and visualization
2. **Email Notifications**: Integrate email service for contract/payment notifications
3. **Real-time Updates**: Implement WebSocket for live notifications
4. **Advanced Reporting**: Add PDF generation for invoices and contracts
5. **Testing**: Add unit and integration tests
6. **Logging**: Implement structured logging with rotation
7. **Caching**: Add Redis for caching frequently accessed data
