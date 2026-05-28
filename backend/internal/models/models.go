package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// ========== M1: Authentication & Account Management ==========

type User struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	Username  string     `gorm:"uniqueIndex" json:"username"`
	Email     string     `gorm:"uniqueIndex" json:"email"`
	Password  string     `json:"-"`      // Hashed password
	Role      string     `json:"role"`   // admin, manager, staff, resident
	Status    string     `json:"status"` // active, inactive, locked
	ResetToken       string     `json:"-"`
	ResetTokenExpiry *time.Time `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Relationships
	Admin    *Admin    `json:"admin,omitempty"`
	Manager  *Manager  `json:"manager,omitempty"`
	Staff    *Staff    `json:"staff,omitempty"`
	Resident *Resident `json:"resident,omitempty"`
}

// ========== M2: GIS 3D Visualization ==========

type Building struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	YearBuilt   int            `json:"year_built"`
	TotalFloors int            `json:"total_floors"`
	Description string         `json:"description"`
	Geometry    datatypes.JSON `json:"geometry"` // 3D geometry for GIS
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relationships
	Floors []Floor `json:"floors,omitempty"`
}

type Floor struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	BuildingID  string         `gorm:"index" json:"building_id"`
	FloorNumber int            `json:"floor_number"`
	FloorHeight float64        `json:"floor_height"`
	FloorArea   float64        `json:"floor_area"`
	Purpose     string         `json:"purpose"`  // residential, technical, service
	Geometry    datatypes.JSON `json:"geometry"` // 3D geometry
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	Building   Building    `gorm:"foreignKey:BuildingID" json:"-"`
	Apartments []Apartment `json:"apartments,omitempty"`
}

type Apartment struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	FloorID       string         `gorm:"index" json:"floor_id"`
	ApartmentCode string         `gorm:"uniqueIndex" json:"apartment_code"` // e.g., A01.01
	Area          float64        `json:"area"`
	Bedrooms      int            `json:"bedrooms"`
	Direction     string         `json:"direction"` // N, S, E, W, NE, SW, etc.
	ListingPrice  float64        `json:"listing_price"`
	CurrentStatus string         `json:"current_status"` // empty, rented, maintenance
	BodyID        string         `json:"body_id"`        // Reference to 3D model
	Geometry      datatypes.JSON `json:"geometry"`       // 3D geometry
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`

	Floor         Floor                    `gorm:"foreignKey:FloorID" json:"-"`
	StatusHistory []ApartmentStatusHistory `json:"status_history,omitempty"`
	Contracts     []Contract               `json:"contracts,omitempty"`
	MainRequests  []MaintenanceRequest     `json:"maintenance_requests,omitempty"`
}

type ApartmentStatusHistory struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	ApartmentID string    `gorm:"index" json:"apartment_id"`
	OldStatus   string    `json:"old_status"`
	NewStatus   string    `json:"new_status"`
	Reason      string    `json:"reason"`
	ChangedBy   string    `json:"changed_by"` // User ID
	ChangedAt   time.Time `json:"changed_at"`

	Apartment Apartment `gorm:"foreignKey:ApartmentID" json:"-"`
}

// ========== M3: Building-Floor-Apartment Management ==========
// Already covered above in GIS 3D section

// ========== M4: Contract Management ==========

type Contract struct {
	ID                string     `gorm:"primaryKey" json:"id"`
	ApartmentID       string     `gorm:"index" json:"apartment_id"`
	ResidentID        string     `gorm:"index" json:"resident_id"`
	ContractCode      string     `gorm:"uniqueIndex" json:"contract_code"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	MonthlyRent       float64    `json:"monthly_rent"`
	PaymentPeriod     int        `json:"payment_period"` // 1=monthly, 3=quarterly, etc
	Deposit           float64    `json:"deposit"`
	Status            string     `json:"status"` // pending, active, expired, terminated
	TerminationReason string     `json:"termination_reason,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ApprovedBy        *string    `json:"approved_by,omitempty"` // Manager ID
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`

	Apartment  Apartment           `gorm:"foreignKey:ApartmentID" json:"-"`
	Resident   Resident            `gorm:"foreignKey:ResidentID" json:"-"`
	Invoices   []Invoice           `json:"invoices,omitempty"`
	Dependents []ResidentDependent `json:"dependents,omitempty"`
}

// ========== M5: Resident Management ==========

type Resident struct {
	ID               string    `gorm:"primaryKey" json:"id"`
	UserID           string    `gorm:"index" json:"user_id"`
	FullName         string    `json:"full_name"`
	IDCard           string    `gorm:"uniqueIndex" json:"id_card"` // CCCD
	DateOfBirth      time.Time `json:"date_of_birth"`
	Gender           string    `json:"gender"` // M, F
	PhoneNumber      string    `json:"phone_number"`
	Email            string    `json:"email"`
	PermanentAddress string    `json:"permanent_address"`
	Avatar           string    `json:"avatar,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	User             User                 `gorm:"foreignKey:UserID" json:"-"`
	Contracts        []Contract           `json:"contracts,omitempty"`
	Dependents       []ResidentDependent  `json:"dependents,omitempty"`
	Vehicles         []Vehicle            `json:"vehicles,omitempty"`
	MainRequests     []MaintenanceRequest `json:"maintenance_requests,omitempty"`
	ResidenceHistory []ResidenceHistory   `json:"residence_history,omitempty"`
}

type ResidentDependent struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	ResidentID   string     `gorm:"index" json:"resident_id"`
	ContractID   string     `gorm:"index" json:"contract_id"`
	FullName     string     `json:"full_name"`
	IDCard       string     `json:"id_card,omitempty"`
	DateOfBirth  time.Time  `json:"date_of_birth"`
	Relationship string     `json:"relationship"` // spouse, child, parent, etc
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`

	Resident Resident `gorm:"foreignKey:ResidentID" json:"-"`
	Contract Contract `gorm:"foreignKey:ContractID" json:"-"`
}

type Vehicle struct {
	ID                string    `gorm:"primaryKey" json:"id"`
	ResidentID        string    `gorm:"index" json:"resident_id"`
	LicensePlate      string    `gorm:"uniqueIndex" json:"license_plate"`
	VehicleType       string    `json:"vehicle_type"` // car, motorcycle, bicycle
	Color             string    `json:"color,omitempty"`
	ParkingSlotNodeID string    `json:"parking_slot_node_id,omitempty"` // Reference to parking slot in GIS
	RegisteredAt      time.Time `json:"registered_at"`

	Resident Resident `gorm:"foreignKey:ResidentID" json:"-"`
}

type ResidenceHistory struct {
	ID          string     `gorm:"primaryKey" json:"id"`
	ResidentID  string     `gorm:"index" json:"resident_id"`
	ApartmentID string     `gorm:"index" json:"apartment_id"`
	ContractID  string     `json:"contract_id"`
	MoveInDate  time.Time  `json:"move_in_date"`
	MoveOutDate *time.Time `json:"move_out_date,omitempty"`

	Resident  Resident  `gorm:"foreignKey:ResidentID" json:"-"`
	Apartment Apartment `gorm:"foreignKey:ApartmentID" json:"-"`
}

// ========== M6: Financial Management & Invoices ==========

type Invoice struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	ContractID    string         `gorm:"index" json:"contract_id"`
	ApartmentID   string         `gorm:"index" json:"apartment_id"`
	InvoiceCode   string         `gorm:"uniqueIndex" json:"invoice_code"`
	BillingPeriod string         `json:"billing_period"` // YYYY-MM
	DueDate       time.Time      `json:"due_date"`
	TotalAmount   float64        `json:"total_amount"`
	PaidAmount    float64        `json:"paid_amount"`
	Status        string         `json:"status"`  // pending, paid, overdue, cancelled
	Details       datatypes.JSON `json:"details"` // Line items: management fee, utilities, parking, etc
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`

	Contract  Contract         `gorm:"foreignKey:ContractID" json:"-"`
	Apartment Apartment        `gorm:"foreignKey:ApartmentID" json:"-"`
	Payments  []InvoicePayment `json:"payments,omitempty"`
}

type InvoicePayment struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	InvoiceID  string    `gorm:"index" json:"invoice_id"`
	Amount     float64   `json:"amount"`
	Method     string    `json:"method"`              // cash, transfer, card
	Reference  string    `json:"reference,omitempty"` // Transaction reference
	PaidAt     time.Time `json:"paid_at"`
	RecordedBy string    `json:"recorded_by"` // Staff ID

	Invoice Invoice `gorm:"foreignKey:InvoiceID" json:"-"`
}

// ========== M7: Maintenance Management ==========

type MaintenanceRequest struct {
	ID              string         `gorm:"primaryKey" json:"id"`
	ApartmentID     string         `gorm:"index" json:"apartment_id"`
	ResidentID      string         `gorm:"index" json:"resident_id"`
	RequestCode     string         `gorm:"uniqueIndex" json:"request_code"`
	IssueType       string         `json:"issue_type"` // electrical, water, elevator, kitchen, etc
	Description     string         `json:"description"`
	Priority        string         `json:"priority"`                            // urgent, normal, low
	Location        string         `json:"location"`                            // living room, bedroom, kitchen, etc
	Photos          pq.StringArray `gorm:"type:text[]" json:"photos,omitempty"` // URLs
	Status          string         `json:"status"`                              // new, assigned, in_progress, completed
	AssignedTo      *string        `json:"assigned_to,omitempty"`               // Staff/Technician ID
	AssignedAt      *time.Time     `json:"assigned_at,omitempty"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	CompletionNotes string         `json:"completion_notes,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	Apartment Apartment `gorm:"foreignKey:ApartmentID" json:"-"`
	Resident  Resident  `gorm:"foreignKey:ResidentID" json:"-"`
}

type Equipment struct {
	ID                  string     `gorm:"primaryKey" json:"id"`
	BuildingID          string     `gorm:"index" json:"building_id"`
	EquipmentName       string     `json:"equipment_name"`
	EquipmentType       string     `json:"equipment_type"` // elevator, pump, generator, etc
	Manufacturer        string     `json:"manufacturer"`
	InstallationDate    time.Time  `json:"installation_date"`
	LocationNodeID      string     `json:"location_node_id"`   // Reference to GIS node
	MaintenancePeriod   int        `json:"maintenance_period"` // months
	LastMaintenanceDate *time.Time `json:"last_maintenance_date,omitempty"`
	NextMaintenanceDate *time.Time `json:"next_maintenance_date,omitempty"`
	Status              string     `json:"status"` // operational, maintenance, faulty

	Building Building `gorm:"foreignKey:BuildingID" json:"-"`
}

// ========== M8: Notifications & Dashboard ==========

type Notification struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	UserID    string     `gorm:"index" json:"user_id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"` // contract_expiring, invoice_overdue, maintenance_request, etc
	IsRead    bool       `json:"is_read"`
	RelatedID string     `json:"related_id,omitempty"` // Contract/Invoice/Request ID
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type AuditLog struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"index" json:"user_id"`
	Action    string         `json:"action"` // CREATE, UPDATE, DELETE
	Entity    string         `json:"entity"` // User, Contract, Invoice, etc
	EntityID  string         `json:"entity_id"`
	OldValues datatypes.JSON `json:"old_values,omitempty"`
	NewValues datatypes.JSON `json:"new_values,omitempty"`
	IPAddress string         `json:"ip_address,omitempty"`
	Timestamp time.Time      `json:"timestamp"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// ========== Role-specific models ==========

type Admin struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type Manager struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	UserID         string    `gorm:"index" json:"user_id"`
	DepartmentName string    `json:"department_name,omitempty"`
	PhoneNumber    string    `json:"phone_number,omitempty"`
	CreatedAt      time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type Staff struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"index" json:"user_id"`
	Position    string    `json:"position"` // operations, maintenance, etc
	Department  string    `json:"department"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
