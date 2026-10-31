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
	ID    int64   `gorm:"primaryKey" json:"id"`
	Name  string  `gorm:"type:varchar(100);not null;unique" json:"name"`
	Slug  string  `gorm:"type:varchar(100);not null;unique" json:"slug"`
	Image *string `gorm:"type:varchar(2048)" json:"image"`
}

type Project struct {
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`

	Name             string        `gorm:"type:varchar(100);not null" json:"name"`
	ShortDescription string        `gorm:"type:varchar(150);not null" json:"short_description"`
	FullDescription  string        `gorm:"type:text;not null" json:"full_description"`
	Status           ProjectStatus `gorm:"type:project_status;not null" json:"status"`
	StartYear        int           `gorm:"not null" json:"start_year"`
	CloseYear        *int          `json:"close_year"`

	Technologies []ProjectTechnology `gorm:"many2many:projects_technologies;joinForeignKey:ProjectID;joinReferences:TechnologyID" json:"technologies"`

	SourceCodeLink *string `gorm:"type:varchar(255)" json:"source_code_link"`
	ProductionLink *string `gorm:"type:varchar(255)" json:"production_link"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
