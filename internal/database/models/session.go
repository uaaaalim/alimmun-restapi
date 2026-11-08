package models

import (
	"time"
	"uuid"
)

type Session struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_sessions_user_created,priority:1"`

	UserAgent string `gorm:"type:text"`
	IP        string `gorm:"type:inet;not null"`

	TokenHash string `gorm:"type:char(64);not null;uniqueIndex"`

	CreatedAt time.Time `gorm:"index:idx_sessions_user_created,priority:2"`
	UpdatedAt time.Time
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
}
