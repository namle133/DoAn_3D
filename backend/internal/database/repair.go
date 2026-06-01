package database

import (
	"log"

	"github.com/uit/vinhomes-management/internal/models"
)

// RepairOrphanResidents xóa bản ghi Resident rỗng (do bug createRoleRecord cũ) nếu đã có hồ sơ đầy đủ
func RepairOrphanResidents() {
	var emptyProfiles []models.Resident
	DB.Where("full_name = '' OR full_name IS NULL").Find(&emptyProfiles)

	fixed := 0
	for _, empty := range emptyProfiles {
		var full models.Resident
		if err := DB.Where("user_id = ? AND full_name != '' AND full_name IS NOT NULL", empty.UserID).
			First(&full).Error; err == nil {
			DB.Delete(&empty)
			fixed++
			continue
		}

		// Chỉ còn bản ghi rỗng — thử lấy thông tin từ User liên kết (username → full_name tạm)
		var user models.User
		if err := DB.First(&user, "id = ?", empty.UserID).Error; err == nil && user.Role == "resident" {
			DB.Model(&empty).Updates(map[string]interface{}{
				"full_name": user.Username,
				"email":     user.Email,
			})
			fixed++
		}
	}

	if fixed > 0 {
		log.Printf("✓ Repaired %d orphan resident profile(s)", fixed)
	}
}
