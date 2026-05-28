package services

import (
	"time"

	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type MaintenanceService struct{}

func NewMaintenanceService() *MaintenanceService {
	return &MaintenanceService{}
}

// GetUrgentMaintenanceRequests returns all urgent maintenance requests
func (s *MaintenanceService) GetUrgentMaintenanceRequests() ([]models.MaintenanceRequest, error) {
	var requests []models.MaintenanceRequest

	if err := database.DB.Where("priority = ? AND status IN ?", "urgent", []string{"new", "assigned"}).
		Preload("Apartment").
		Preload("Resident").
		Order("created_at ASC").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

// GetAverageResolutionTime returns average time to resolve maintenance requests (in hours)
func (s *MaintenanceService) GetAverageResolutionTime() (float64, error) {
	var avgTime float64

	err := database.DB.Model(&models.MaintenanceRequest{}).
		Where("status = ? AND completed_at IS NOT NULL", "completed").
		Select("AVG(EXTRACT(EPOCH FROM (completed_at - created_at))/3600)").
		Row().Scan(&avgTime)

	return avgTime, err
}

// GetEquipmentMaintenanceSchedule returns all equipment with upcoming maintenance
func (s *MaintenanceService) GetEquipmentMaintenanceSchedule(daysAhead int) ([]models.Equipment, error) {
	var equipment []models.Equipment

	futureDate := time.Now().AddDate(0, 0, daysAhead)
	err := database.DB.Where("next_maintenance_date BETWEEN NOW() AND ?", futureDate).
		Order("next_maintenance_date ASC").
		Find(&equipment).Error

	return equipment, err
}

// GetMaintenanceStatistics returns statistics about maintenance requests
func (s *MaintenanceService) GetMaintenanceStatistics(startDate, endDate time.Time) (map[string]interface{}, error) {
	type Stats struct {
		TotalRequests      int64
		CompletedCount     int64
		PendingCount       int64
		AvgResolutionHours float64
		MostCommonIssue    string
		CriticalCount      int64
	}

	var stats Stats

	database.DB.Model(&models.MaintenanceRequest{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalRequests)

	database.DB.Model(&models.MaintenanceRequest{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").
		Count(&stats.CompletedCount)

	database.DB.Model(&models.MaintenanceRequest{}).
		Where("created_at BETWEEN ? AND ? AND status IN ?", startDate, endDate, []string{"new", "assigned"}).
		Count(&stats.PendingCount)

	database.DB.Model(&models.MaintenanceRequest{}).
		Where("created_at BETWEEN ? AND ? AND status = ? AND completed_at IS NOT NULL", startDate, endDate, "completed").
		Select("AVG(EXTRACT(EPOCH FROM (completed_at - created_at))/3600)").
		Row().Scan(&stats.AvgResolutionHours)

	database.DB.Model(&models.MaintenanceRequest{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("issue_type").
		Group("issue_type").
		Order("COUNT(*) DESC").
		Limit(1).
		Row().Scan(&stats.MostCommonIssue)

	database.DB.Model(&models.MaintenanceRequest{}).
		Where("created_at BETWEEN ? AND ? AND priority = ?", startDate, endDate, "urgent").
		Count(&stats.CriticalCount)

	return map[string]interface{}{
		"total_requests":       stats.TotalRequests,
		"completed_count":      stats.CompletedCount,
		"pending_count":        stats.PendingCount,
		"avg_resolution_hours": stats.AvgResolutionHours,
		"most_common_issue":    stats.MostCommonIssue,
		"critical_count":       stats.CriticalCount,
		"completion_rate":      float64(stats.CompletedCount) / float64(stats.TotalRequests) * 100,
	}, nil
}
