package models

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RecordDate *time.Time `gorm:"type:timestamp"`
	Name       string     `gorm:"type:varchar(255);not null"`
	Email      string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone      *string    `gorm:"type:varchar(50)"`

	// Optional: Inverse relationship back to records if you want to preload them
	Records []RegisterRecord `gorm:"foreignKey:MemberID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName overrides the default table name if needed
func (Member) TableName() string {
	return "savara_umosan.members"
}
