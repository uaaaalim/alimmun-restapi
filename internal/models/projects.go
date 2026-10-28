package models

import "time"

type ProjectStatus string

const (
	Active        ProjectStatus = "active"
	InDevelopment ProjectStatus = "in_development"
	Suspended     ProjectStatus = "suspended"
	Abandoned     ProjectStatus = "abandoned"
)

type ProjectTechnology struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null;unique"`
	Slug string `gorm:"type:varchar(100);not null;unique"`
}

type Project struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	Name             string        `gorm:"type:varchar(100);not null"`
	ShortDescription string        `gorm:"type:varchar(150);not null"`
	FullDescription  string        `gorm:"type:text;not null"`
	Status           ProjectStatus `gorm:"type:project_status;not null"`
	StartYear        int           `gorm:"not null"`
	CloseYear        *int

	Technologies []ProjectTechnology `gorm:"many2many:project_technologies;"`

	SourceCodeLink *string `gorm:"type:varchar(255)"`
	ProductionLink *string `gorm:"type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
