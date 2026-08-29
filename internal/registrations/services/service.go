package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wksn753/umosan-backend/internal/registrations/infrastructure"
	"github.com/wksn753/umosan-backend/internal/registrations/models"
	"gorm.io/gorm"
)

type RegistrationService struct {
	Repo *infrastructure.Infrastructure
}

func NewRegistrationService(repo *infrastructure.Infrastructure) *RegistrationService {
	return &RegistrationService{Repo: repo}
}

// RegisterOrCheckIn handles the core logic: check if member exists by email/phone,
// create if missing, then log the registration/attendance record.
func (s *RegistrationService) RegisterOrCheckIn(req RegisterMemberRequest) error {
	var member models.Member
	now := time.Now()

	// 1. Check if member exists by Email or Phone
	err := s.Repo.DB.Where("email = ? OR phone = ?", strings.TrimSpace(req.Email), strings.TrimSpace(req.Phone)).First(&member).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 2. If member does not exist, create a new member record
		member = models.Member{
			ID:         uuid.New(),
			Name:       req.Name,
			Email:      req.Email,
			Phone:      &req.Phone,
			RecordDate: &now,
		}
		if createErr := s.Repo.CreateMember(&member); createErr != nil {
			return createErr
		}
	} else if err != nil {
		// Return unexpected database errors
		return err
	}

	// 3. Create the event registration record linked to the member
	record := models.RegisterRecord{
		ID:         uuid.New(),
		MemberID:   member.ID,
		RecordDate: &now,
		EventName:  req.EventName,
		Rating:     req.Rating,
		Suggestion: req.Suggestion,
	}

	return s.Repo.CreateRecord(&record)
}

// GetAttendanceByDate retrieves all registration/attendance records for a specific calendar date (YYYY-MM-DD)
func (s *RegistrationService) GetAttendanceByDate(dateStr string) ([]models.RegisterRecord, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var records []models.RegisterRecord
	err = s.Repo.DB.Preload("Member").
		Where("record_date >= ? AND record_date < ?", startOfDay, endOfDay).
		Find(&records).Error

	return records, err
}

// GetMemberStats returns general aggregate statistics for the dashboard
func (s *RegistrationService) GetMemberStats() (map[string]int64, error) {
	var totalMembers int64
	var totalRegistrations int64

	if err := s.Repo.DB.Model(&models.Member{}).Count(&totalMembers).Error; err != nil {
		return nil, err
	}

	if err := s.Repo.DB.Model(&models.RegisterRecord{}).Count(&totalRegistrations).Error; err != nil {
		return nil, err
	}

	stats := map[string]int64{
		"total_members":       totalMembers,
		"total_registrations": totalRegistrations,
	}

	return stats, nil
}
