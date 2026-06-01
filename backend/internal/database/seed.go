package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"github.com/uit/vinhomes-management/internal/models"
)

type demoUserSpec struct {
	username, email, password, role, fullName string
	gender                                    string
}

// SeedDemoData tạo dữ liệu mẫu: hạ tầng (tòa/căn) + nghiệp vụ (HĐ, HĐơn, bảo trì, TB...)
func SeedDemoData() {
	seedBaseInfrastructure()
	seedBusinessDemoData()
}

func seedBaseInfrastructure() {
	var aptCount int64
	DB.Model(&models.Apartment{}).Count(&aptCount)
	if aptCount > 0 {
		return
	}

	log.Println("Seeding base infrastructure (users, towers, apartments)...")

	users := ensureDemoUsers([]demoUserSpec{
		{"admin", "admin@vinhomes.com", "admin123", "admin", "Quản Trị Viên", "M"},
		{"manager", "manager@vinhomes.com", "manager123", "manager", "Trần Văn Quản Lý", "M"},
		{"staff", "staff@vinhomes.com", "staff123", "staff", "Phạm Thị Nhân Viên", "F"},
		{"resident", "resident@vinhomes.com", "resident123", "resident", "Nguyễn Văn Cư Dân", "M"},
	})
	_ = users

	towers := []struct {
		id, name string
		floors   int
	}{
		{"W1", "Tòa West 1", 39},
		{"W2", "Tòa West 2", 39},
		{"W3", "Tòa West 3", 35},
	}

	buildingIDs := map[string]string{}
	for _, t := range towers {
		bid := uuid.New().String()
		buildingIDs[t.id] = bid
		DB.Create(&models.Building{
			ID: bid, Name: t.name, Address: "Phạm Hùng, Mễ Trì, Nam Từ Liêm, Hà Nội",
			YearBuilt: 2020, TotalFloors: t.floors,
			Description: "Vinhomes West Point - " + t.id,
		})
	}

	statuses := []string{"empty", "rented", "rented", "maintenance", "empty"}

	for _, t := range towers {
		bid := buildingIDs[t.id]
		sampleFloors := []int{5, 8, 12, 18, 25}
		for fi, floorNum := range sampleFloors {
			fid := uuid.New().String()
			DB.Create(&models.Floor{
				ID: fid, BuildingID: bid, FloorNumber: floorNum,
				FloorHeight: 3.3, FloorArea: 800, Purpose: "residential",
			})
			for u := 1; u <= 4; u++ {
				status := statuses[(fi+u)%len(statuses)]
				code := t.id + "." + pad2(floorNum) + "." + pad2(u)
				bodyID := t.id + "-F" + pad2(floorNum) + "-U" + pad2(u)
				DB.Create(&models.Apartment{
					ID: uuid.New().String(), FloorID: fid, ApartmentCode: code, BodyID: bodyID,
					Area: 75 + float64(u*10), Bedrooms: u%3 + 1,
					Direction: "Đông", ListingPrice: 15000000,
					CurrentStatus: status,
				})
			}
		}
	}

	log.Println("✓ Base infrastructure seeded (4 users, 3 towers, 60 apartments)")
}

func seedBusinessDemoData() {
	var demoInvCount int64
	DB.Model(&models.Invoice{}).Where("invoice_code LIKE ?", "INV-DEMO-%").Count(&demoInvCount)
	if demoInvCount > 0 {
		return
	}

	log.Println("Seeding business demo data (contracts, invoices, maintenance, notifications...)...")

	users := ensureDemoUsers([]demoUserSpec{
		{"admin", "admin@vinhomes.com", "admin123", "admin", "Quản Trị Viên", "M"},
		{"manager", "manager@vinhomes.com", "manager123", "manager", "Trần Văn Quản Lý", "M"},
		{"staff", "staff@vinhomes.com", "staff123", "staff", "Phạm Thị Nhân Viên", "F"},
		{"resident", "resident@vinhomes.com", "resident123", "resident", "Nguyễn Văn Cư Dân", "M"},
		{"resident2", "resident2@vinhomes.com", "resident123", "resident", "Trần Thị Lan", "F"},
		{"resident3", "resident3@vinhomes.com", "resident123", "resident", "Lê Văn Minh", "M"},
	})

	residents := ensureExtraResidents(users)
	apartments := loadApartmentsForSeed()
	if len(apartments) < 6 {
		log.Println("⚠ Skip business seed: not enough apartments")
		return
	}

	now := time.Now()
	managerID := users["manager"]
	staffUserID := users["staff"]

	contracts := seedDemoContracts(residents, apartments, managerID, now)
	seedDemoInvoices(contracts, staffUserID, now)
	seedDemoMaintenance(residents, apartments, staffUserID, now)
	seedDemoNotifications(users, residents, contracts, now)
	seedDemoEquipment(now)
	seedDemoResidentsExtras(residents, contracts)

	log.Println("✓ Business demo data seeded (contracts, invoices, maintenance, notifications, equipment)")
}

type userMap map[string]string // role or key -> user id
type residentMap map[string]string // key -> resident id

func ensureDemoUsers(specs []demoUserSpec) userMap {
	result := userMap{}
	for _, s := range specs {
		var existing models.User
		if err := DB.Where("email = ?", s.email).First(&existing).Error; err == nil {
			result[s.username] = existing.ID
			result[s.role+"_"+s.username] = existing.ID
			if s.role == "resident" {
				var res models.Resident
				if DB.Where("user_id = ?", existing.ID).First(&res).Error == nil {
					result["resident_"+s.username] = res.ID
				}
			}
			continue
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte(s.password), bcrypt.DefaultCost)
		id := uuid.New().String()
		DB.Create(&models.User{
			ID: id, Username: s.username, Email: s.email,
			Password: string(hashed), Role: s.role, Status: "active",
		})
		result[s.username] = id
		result[s.role] = id

		switch s.role {
		case "admin":
			DB.Create(&models.Admin{ID: uuid.New().String(), UserID: id})
		case "manager":
			DB.Create(&models.Manager{ID: uuid.New().String(), UserID: id})
		case "staff":
			DB.Create(&models.Staff{ID: uuid.New().String(), UserID: id})
		case "resident":
			resID := uuid.New().String()
			DB.Create(&models.Resident{
				ID: resID, UserID: id, FullName: s.fullName,
				IDCard: uniqueIDCard(s.username), PhoneNumber: phoneFor(s.username),
				Email: s.email, Gender: s.gender,
				DateOfBirth: time.Date(1990, 3, 10, 0, 0, 0, 0, time.UTC),
				PermanentAddress: "Hà Nội, Việt Nam",
			})
			result["resident_"+s.username] = resID
			result["resident_profile"] = resID
		}
	}
	return result
}

func ensureExtraResidents(users userMap) residentMap {
	m := residentMap{}
	keys := []string{"resident", "resident2", "resident3"}
	for _, k := range keys {
		if id, ok := users["resident_"+k]; ok {
			m[k] = id
		}
	}
	if _, ok := m["resident"]; !ok {
		var res models.Resident
		if DB.Where("full_name != ''").First(&res).Error == nil {
			m["resident"] = res.ID
		}
	}
	return m
}

func loadApartmentsForSeed() []models.Apartment {
	var apts []models.Apartment
	DB.Order("apartment_code ASC").Find(&apts)
	return apts
}

type contractSeed struct {
	model   models.Contract
	aptCode string
}

func seedDemoContracts(residents residentMap, apartments []models.Apartment, managerID string, now time.Time) []contractSeed {
	aptByCode := map[string]models.Apartment{}
	for _, a := range apartments {
		aptByCode[a.ApartmentCode] = a
	}

	specs := []struct {
		code, aptSuffix, resKey, status string
		monthsAgoStart, monthsToEnd   int
		rent                          float64
	}{
		{"HD-DEMO-001", "W1.05.01", "resident", "active", 12, -1, 15000000},
		{"HD-DEMO-002", "W1.08.02", "resident", "active", 6, 1, 16500000},
		{"HD-DEMO-003", "W2.12.01", "resident2", "active", 8, 0, 18000000},
		{"HD-DEMO-004", "W2.18.03", "resident3", "active", 3, 8, 17500000},
		{"HD-DEMO-005", "W1.25.02", "resident2", "pending", 0, 12, 16000000},
		{"HD-DEMO-006", "W3.05.03", "resident", "active", 10, 0, 14000000}, // sắp hết hạn ~20 ngày
	}

	var created []contractSeed
	for _, sp := range specs {
		apt, ok := aptByCode[sp.aptSuffix]
		if !ok && len(apartments) > 0 {
			apt = apartments[len(created)%len(apartments)]
		}
		resID := residents[sp.resKey]
		if resID == "" {
			resID = residents["resident"]
		}
		if resID == "" {
			continue
		}

		start := now.AddDate(0, -sp.monthsAgoStart, 0)
		end := now.AddDate(0, sp.monthsToEnd, 20)
		if sp.code == "HD-DEMO-006" {
			end = now.AddDate(0, 0, 18)
		}

		cid := uuid.New().String()
		approvedAt := now.AddDate(0, -sp.monthsAgoStart, -2)
		c := models.Contract{
			ID: cid, ApartmentID: apt.ID, ResidentID: resID,
			ContractCode: sp.code, StartDate: start, EndDate: end,
			MonthlyRent: sp.rent, PaymentPeriod: 1, Deposit: sp.rent * 2,
			Status: sp.status, CreatedAt: start,
		}
		if sp.status == "active" {
			c.ApprovedBy = &managerID
			c.ApprovedAt = &approvedAt
			DB.Model(&apt).Update("current_status", "rented")
		}
		DB.Create(&c)
		created = append(created, contractSeed{model: c, aptCode: apt.ApartmentCode})
	}
	return created
}

func seedDemoInvoices(contracts []contractSeed, staffUserID string, now time.Time) {
	for ci, cs := range contracts {
		if cs.model.Status != "active" {
			continue
		}
		for m := 0; m < 4; m++ {
			period := now.AddDate(0, -m, 0).Format("2006-01")
			invCode := fmt.Sprintf("INV-DEMO-%s-%02d", cs.model.ContractCode, m+1)
			due := now.AddDate(0, -m+1, 5)
			if m == 0 {
				due = now.AddDate(0, 0, 7)
			}

			status := "pending"
			paid := 0.0
			if m >= 2 {
				status = "paid"
				paid = cs.model.MonthlyRent
			} else if m == 1 {
				status = "pending"
				due = now.AddDate(0, 0, -15) // quá hạn
			}

			details, _ := json.Marshal(map[string]interface{}{
				"line_items": []map[string]interface{}{
					{"description": "Tiền thuê căn hộ", "quantity": 1, "price": cs.model.MonthlyRent, "total": cs.model.MonthlyRent},
					{"description": "Phí quản lý", "quantity": 1, "price": 500000, "total": 500000},
				},
				"description": "Hóa đơn tháng " + period,
			})
			total := cs.model.MonthlyRent + 500000

			invID := uuid.New().String()
			inv := models.Invoice{
				ID: invID, ContractID: cs.model.ID, ApartmentID: cs.model.ApartmentID,
				InvoiceCode: invCode, BillingPeriod: period,
				DueDate: due, TotalAmount: total, PaidAmount: paid,
				Status: status, Details: datatypes.JSON(details),
				CreatedAt: now.AddDate(0, -m, 0),
			}
			if status == "pending" && due.Before(now) {
				inv.Status = "overdue"
			}
			DB.Create(&inv)

			if status == "paid" {
				DB.Create(&models.InvoicePayment{
					ID: uuid.New().String(), InvoiceID: invID,
					Amount: paid, Method: "transfer",
					Reference: fmt.Sprintf("TT-DEMO-%d-%d", ci, m),
					PaidAt: now.AddDate(0, -m, 3), RecordedBy: staffUserID,
				})
			}
		}
	}
}

func seedDemoMaintenance(residents residentMap, apartments []models.Apartment, staffUserID string, now time.Time) {
	resID := residents["resident"]
	if resID == "" {
		return
	}
	assigned := staffUserID
	specs := []struct {
		code, issue, priority, status, desc string
		aptIdx                              int
		daysAgo                             int
		assign                              bool
	}{
		{"MR-DEMO-001", "plumbing", "urgent", "new", "Rò rỉ ống nước nhà bếp", 0, 1, false},
		{"MR-DEMO-002", "electrical", "normal", "assigned", "Đèn phòng khách không sáng", 1, 3, true},
		{"MR-DEMO-003", "appliances", "normal", "in_progress", "Điều hòa không lạnh", 2, 5, true},
		{"MR-DEMO-004", "structural", "low", "completed", "Vết nứt nhỏ trần nhà", 3, 10, true},
		{"MR-DEMO-005", "cleaning", "low", "cancelled", "Yêu cầu vệ sinh thảm — đã tự xử lý", 4, 2, false},
	}

	for i, sp := range specs {
		apt := apartments[i%len(apartments)]
		mr := models.MaintenanceRequest{
			ID: uuid.New().String(), ApartmentID: apt.ID, ResidentID: resID,
			RequestCode: sp.code, IssueType: sp.issue,
			Description: sp.desc, Priority: sp.priority, Location: "Phòng khách",
			Status: sp.status, CreatedAt: now.AddDate(0, 0, -sp.daysAgo),
		}
		if sp.assign {
			mr.AssignedTo = &assigned
			t := now.AddDate(0, 0, -sp.daysAgo+1)
			mr.AssignedAt = &t
		}
		if sp.status == "completed" {
			t := now.AddDate(0, 0, -1)
			mr.CompletedAt = &t
			mr.CompletionNotes = "Giải pháp: Đã sửa xong. Chi phí: 500000 VND"
		}
		DB.Create(&mr)
	}
}

func seedDemoNotifications(users userMap, residents residentMap, contracts []contractSeed, now time.Time) {
	type notifSpec struct {
		userKey, title, msg, ntype string
		daysAgo                    int
		read                       bool
		relatedID                  string
	}

	specs := []notifSpec{
		{"resident", "Hợp đồng sắp hết hạn", "Hợp đồng thuê căn hộ của bạn sẽ hết hạn trong 18 ngày. Vui lòng liên hệ BQL.", "contract_expiring", 2, false, ""},
		{"resident", "Hóa đơn quá hạn", "Hóa đơn INV-DEMO-HD-DEMO-001-02 đã quá hạn thanh toán.", "invoice_overdue", 5, false, ""},
		{"resident", "Hóa đơn mới", "Hóa đơn tháng hiện tại đã phát hành. Hạn thanh toán trong 7 ngày.", "invoice_generated", 1, true, ""},
		{"manager", "Yêu cầu bảo trì mới", "Cư dân gửi yêu cầu sửa ống nước — ưu tiên khẩn cấp.", "maintenance_request", 1, false, ""},
		{"staff", "Được phân công bảo trì", "Bạn được phân công xử lý yêu cầu MR-DEMO-002.", "maintenance_assigned", 3, false, ""},
		{"manager", "Thông báo hệ thống", "Báo cáo tài chính tháng đã sẵn sàng trên portal.", "general", 0, false, ""},
		{"resident2", "Chào mừng Vinhomes West Point", "Tài khoản cư dân của bạn đã được kích hoạt.", "general", 7, true, ""},
	}

	for _, sp := range specs {
		uid := users[sp.userKey]
		if uid == "" {
			continue
		}
		n := models.Notification{
			ID: uuid.New().String(), UserID: uid,
			Title: sp.title, Message: sp.msg, Type: sp.ntype,
			IsRead: sp.read, RelatedID: sp.relatedID,
			CreatedAt: now.AddDate(0, 0, -sp.daysAgo),
		}
		if sp.read {
			t := now.AddDate(0, 0, -sp.daysAgo+1)
			n.ReadAt = &t
		}
		DB.Create(&n)
	}
}

func seedDemoEquipment(now time.Time) {
	var b models.Building
	if DB.First(&b).Error != nil {
		return
	}
	equip := []struct {
		name, etype, mfr string
		monthsAgoInstall int
		period           int
		overdue          bool
	}{
		{"Thang máy tòa W1-A", "elevator", "Otis", 24, 6, false},
		{"Máy bơm nước W1", "pump", "Grundfos", 18, 3, true},
		{"Máy phát điện", "generator", "Cummins", 36, 12, false},
	}
	for _, e := range equip {
		install := now.AddDate(0, -e.monthsAgoInstall, 0)
		next := install.AddDate(0, e.period, 0)
		if e.overdue {
			next = now.AddDate(0, -1, 0)
		}
		last := install
		DB.Create(&models.Equipment{
			ID: uuid.New().String(), BuildingID: b.ID,
			EquipmentName: e.name, EquipmentType: e.etype, Manufacturer: e.mfr,
			InstallationDate: install, MaintenancePeriod: e.period,
			LastMaintenanceDate: &last, NextMaintenanceDate: &next,
			Status: "operational",
		})
	}
}

func seedDemoResidentsExtras(residents residentMap, contracts []contractSeed) {
	resID := residents["resident"]
	if resID == "" || len(contracts) == 0 {
		return
	}
	cid := contracts[0].model.ID
	aptID := contracts[0].model.ApartmentID

	DB.Create(&models.Vehicle{
		ID: uuid.New().String(), ResidentID: resID,
		LicensePlate: "30A-DEMO-001", VehicleType: "car", Color: "Trắng",
		RegisteredAt: time.Now().AddDate(-1, 0, 0),
	})
	DB.Create(&models.ResidentDependent{
		ID: uuid.New().String(), ResidentID: resID, ContractID: cid,
		FullName: "Nguyễn Thị Hoa", Relationship: "spouse",
		DateOfBirth: time.Date(1992, 6, 15, 0, 0, 0, 0, time.UTC),
		StartDate: time.Now().AddDate(-1, 0, 0),
	})
	DB.Create(&models.ResidenceHistory{
		ID: uuid.New().String(), ResidentID: resID,
		ApartmentID: aptID, ContractID: cid,
		MoveInDate: time.Now().AddDate(-1, 0, 0),
	})
}

func uniqueIDCard(username string) string {
	h := 0
	for _, c := range username {
		h = h*31 + int(c)
	}
	return fmt.Sprintf("%012d", 100000000000+h%899999999999)
}

func phoneFor(username string) string {
	return fmt.Sprintf("09%08d", hashStr(username)%100000000)
}

func hashStr(s string) int {
	h := 0
	for _, c := range s {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}

func pad2(n int) string {
	return fmt.Sprintf("%02d", n)
}
