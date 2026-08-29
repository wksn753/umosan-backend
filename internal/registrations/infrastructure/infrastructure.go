package infrastructure

import (
	"errors"

	"github.com/google/uuid"
	"github.com/wksn753/umosan-backend/internal/registrations/models" // Adjust to your actual module import path
	"gorm.io/gorm"
)

type Infrastructure struct {
	DB *gorm.DB
}

func NewInfrastructure(db *gorm.DB) *Infrastructure {
	return &Infrastructure{DB: db}
}

// ==================== MEMBER OPERATIONS ====================

func (r *Infrastructure) CreateMember(member *models.Member) error {
	return r.DB.Create(member).Error
}

func (r *Infrastructure) GetMemberByID(id uuid.UUID) (*models.Member, error) {
	var member models.Member
	err := r.DB.Preload("Records").First(&member, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *Infrastructure) UpdateMember(member *models.Member) error {
	return r.DB.Save(member).Error
}

func (r *Infrastructure) DeleteMember(id uuid.UUID) error {
	return r.DB.Delete(&models.Member{}, "id = ?", id).Error
}

// ==================== REGISTER RECORD OPERATIONS ====================

func (r *Infrastructure) CreateRecord(record *models.RegisterRecord) error {
	// Verify member exists first if needed
	var member models.Member
	if err := r.DB.First(&member, "id = ?", record.MemberID).Error; err != nil {
		return errors.New("associated member not found")
	}
	return r.DB.Create(record).Error
}

func (r *Infrastructure) GetRecordByID(id uuid.UUID) (*models.RegisterRecord, error) {
	var record models.RegisterRecord
	err := r.DB.Preload("Member").First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *Infrastructure) UpdateRecord(record *models.RegisterRecord) error {
	return r.DB.Save(record).Error
}

func (r *Infrastructure) DeleteRecord(id uuid.UUID) error {
	return r.DB.Delete(&models.RegisterRecord{}, "id = ?", id).Error
}
