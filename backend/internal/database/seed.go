package database

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/uit/vinhomes-management/internal/models"
)

// SeedDemoData tạo dữ liệu mẫu nếu DB trống (users, buildings, apartments cho GIS)
func SeedDemoData() {
	var aptCount int64
	DB.Model(&models.Apartment{}).Count(&aptCount)
	if aptCount > 0 {
		return
	}

	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	needUsers := userCount == 0

	log.Println("Seeding demo data...")

	userIDs := map[string]string{}
	if needUsers {
		type demoUser struct {
			username, email, password, role string
		}
		demoUsers := []demoUser{
			{"admin", "admin@vinhomes.com", "admin123", "admin"},
			{"manager", "manager@vinhomes.com", "manager123", "manager"},
			{"staff", "staff@vinhomes.com", "staff123", "staff"},
			{"resident", "resident@vinhomes.com", "resident123", "resident"},
		}

		for _, du := range demoUsers {
			hashed, _ := bcrypt.GenerateFromPassword([]byte(du.password), bcrypt.DefaultCost)
			id := uuid.New().String()
			user := models.User{
				ID: id, Username: du.username, Email: du.email,
				Password: string(hashed), Role: du.role, Status: "active",
			}
			DB.Create(&user)
			userIDs[du.role] = id
			switch du.role {
			case "admin":
				DB.Create(&models.Admin{ID: uuid.New().String(), UserID: id})
			case "manager":
				DB.Create(&models.Manager{ID: uuid.New().String(), UserID: id})
			case "staff":
				DB.Create(&models.Staff{ID: uuid.New().String(), UserID: id})
			case "resident":
				resID := uuid.New().String()
				DB.Create(&models.Resident{
					ID: resID, UserID: id, FullName: "Nguyễn Văn Cư Dân",
					IDCard: "001234567890", PhoneNumber: "0901234567",
					Email: du.email, Gender: "M",
					DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
					PermanentAddress: "Hà Nội",
				})
				userIDs["resident_profile"] = resID
			}
		}
	} else {
		// Lấy resident profile nếu đã có user
		var resident models.Resident
		if DB.First(&resident).Error == nil {
			userIDs["resident_profile"] = resident.ID
		}
	}

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
	apartmentIDs := []string{}

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
				aptID := uuid.New().String()
				DB.Create(&models.Apartment{
					ID: aptID, FloorID: fid, ApartmentCode: code, BodyID: bodyID,
					Area: 75 + float64(u*10), Bedrooms: u % 3 + 1,
					Direction: "Đông", ListingPrice: 15000000,
					CurrentStatus: status,
				})
				apartmentIDs = append(apartmentIDs, aptID)
			}
		}
	}

	// Hợp đồng sắp hết hạn (expiring - yellow on GIS)
	if len(apartmentIDs) > 2 {
		resID := userIDs["resident_profile"]
		DB.Create(&models.Contract{
			ID: uuid.New().String(), ApartmentID: apartmentIDs[1],
			ResidentID: resID, ContractCode: "HD-2024-001",
			StartDate: time.Now().AddDate(-1, 0, 0),
			EndDate:   time.Now().AddDate(0, 0, 20),
			MonthlyRent: 15000000, PaymentPeriod: 1, Deposit: 30000000,
			Status: "active",
		})
		DB.Model(&models.Apartment{}).Where("id = ?", apartmentIDs[1]).Update("current_status", "rented")

		DB.Create(&models.Contract{
			ID: uuid.New().String(), ApartmentID: apartmentIDs[5],
			ResidentID: resID, ContractCode: "HD-2024-002",
			StartDate: time.Now().AddDate(0, -6, 0),
			EndDate:   time.Now().AddDate(0, 6, 0),
			MonthlyRent: 18000000, PaymentPeriod: 1, Deposit: 36000000,
			Status: "active",
		})
	}

	log.Println("✓ Demo data seeded (4 users, 3 towers, 60 apartments)")
}

func pad2(n int) string {
	return fmt.Sprintf("%02d", n)
}
