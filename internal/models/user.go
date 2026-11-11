package models

import (
	"time"
	"uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"type:varchar(120);not null"`
	Email        string    `gorm:"type:varchar(254);not null"`
	PasswordHash string    `gorm:"type:text;not null"`

	IsActivated bool    `gorm:"not null;default:true"`
	IsBanned    bool    `gorm:"not null;default:false"`
	BanReason   *string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Sessions []Session `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
