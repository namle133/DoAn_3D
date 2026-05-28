package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type GISApartmentStatus struct {
	ID            string  `json:"id"`
	ApartmentCode string  `json:"apartment_code"`
	BodyID        string  `json:"body_id"`
	TowerID       string  `json:"tower_id"`
	FloorNumber   int     `json:"floor_number"`
	UnitIndex     int     `json:"unit_index"`
	CurrentStatus string  `json:"current_status"`
	DisplayStatus string  `json:"display_status"`
	Area          float64 `json:"area"`
	Bedrooms      int     `json:"bedrooms"`
	ListingPrice  float64 `json:"listing_price"`
	ResidentName  string  `json:"resident_name,omitempty"`
	ContractEnd   *string `json:"contract_end,omitempty"`
}

// GetGISApartmentStatus — public endpoint cho GIS 3D (UC 2.7.8–2.7.12)
func GetGISApartmentStatus(c *gin.Context) {
	towerFilter := c.Query("tower_id")
	statusFilter := c.Query("status")
	floorMin := c.DefaultQuery("floor_min", "0")
	floorMax := c.DefaultQuery("floor_max", "999")

	var apartments []models.Apartment
	query := database.DB
	if towerFilter != "" {
		query = query.Where("body_id LIKE ? OR apartment_code LIKE ?", towerFilter+"%", towerFilter+".%")
	}
	if err := query.Find(&apartments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch apartments"})
		return
	}

	now := time.Now()
	in30 := now.AddDate(0, 0, 30)

	type contractInfo struct {
		EndDate      time.Time
		ResidentName string
	}
	contractMap := map[string]contractInfo{}
	var contracts []models.Contract
	database.DB.Where("status = ? AND end_date > ?", "active", now).Find(&contracts)
	for _, ct := range contracts {
		name := ""
		var resident models.Resident
		if database.DB.First(&resident, "id = ?", ct.ResidentID).Error == nil {
			name = resident.FullName
		}
		contractMap[ct.ApartmentID] = contractInfo{EndDate: ct.EndDate, ResidentName: name}
	}

	result := make([]GISApartmentStatus, 0, len(apartments))
	for _, apt := range apartments {
		towerID, floorNum, unitIdx := parseApartmentLocation(apt.ApartmentCode, apt.BodyID)
		if towerID == "" {
			continue
		}
		if floorNum < parseIntOr(floorMin, 0) || floorNum > parseIntOr(floorMax, 999) {
			continue
		}

		displayStatus := apt.CurrentStatus
		var contractEnd *string
		var residentName string

		if ci, ok := contractMap[apt.ID]; ok {
			endStr := ci.EndDate.Format("2006-01-02")
			contractEnd = &endStr
			residentName = ci.ResidentName
			if apt.CurrentStatus == "rented" && ci.EndDate.After(now) && !ci.EndDate.After(in30) {
				displayStatus = "expiring"
			}
		}

		if statusFilter != "" && displayStatus != statusFilter && apt.CurrentStatus != statusFilter {
			continue
		}

		result = append(result, GISApartmentStatus{
			ID:            apt.ID,
			ApartmentCode: apt.ApartmentCode,
			BodyID:        apt.BodyID,
			TowerID:       towerID,
			FloorNumber:   floorNum,
			UnitIndex:     unitIdx,
			CurrentStatus: apt.CurrentStatus,
			DisplayStatus: displayStatus,
			Area:          apt.Area,
			Bedrooms:      apt.Bedrooms,
			ListingPrice:  apt.ListingPrice,
			ResidentName:  residentName,
			ContractEnd:   contractEnd,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func parseApartmentLocation(code, bodyID string) (towerID string, floorNum int, unitIdx int) {
	// body_id: W1-F12-U03  |  apartment_code: W1.12.03
	src := bodyID
	if src == "" {
		src = code
	}
	if src == "" {
		return "", 0, 0
	}

	parts := splitDelims(src, ".-_")
	if len(parts) >= 3 {
		towerID = parts[0]
		floorNum = parseIntOr(parts[1], 0)
		unitIdx = parseIntOr(parts[2], 0)
		return towerID, floorNum, unitIdx
	}
	return "", 0, 0
}

func splitDelims(s string, delims string) []string {
	var parts []string
	current := ""
	for _, ch := range s {
		isDelim := false
		for _, d := range delims {
			if ch == d {
				isDelim = true
				break
			}
		}
		if isDelim {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func parseIntOr(s string, def int) int {
	n := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			n = n*10 + int(ch-'0')
		}
	}
	if n == 0 && s != "0" {
		return def
	}
	return n
}
