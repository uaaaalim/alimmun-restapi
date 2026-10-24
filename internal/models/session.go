package models

import (
	"time"
	"uuid"
)

type Session struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID int64     `gorm:"type:bigint;not null;index"`

	UserAgent  string `gorm:"not null;default:''"`
	DeviceName string `gorm:"not null;default:''"`
	IPAddress  string `gorm:"type:inet;not null"`

	TokenHash string `gorm:"not null;unique" json:"-"`

	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	LastSeenAt time.Time `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	RevokedAt  *time.Time
}
