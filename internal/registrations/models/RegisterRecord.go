package models

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRecord struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MemberID   uuid.UUID  `gorm:"type:uuid;notNull;index"`
	RecordDate *time.Time `gorm:"type:timestamp"`
	EventName  string     `gorm:"type:varchar(255);not null"`
	Rating     string     `gorm:"type:varchar(50)"`
	Suggestion string     `gorm:"type:text"`

	// Belongs To relationship (GORM will automatically handle the foreign key association)
	Member Member `gorm:"foreignKey:MemberID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName overrides the default table name if needed
func (RegisterRecord) TableName() string {
	return "savara_umosan.register_records"
}
