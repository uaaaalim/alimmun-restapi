package database

import "gorm.io/gorm"

type DatabaseServices struct {
}

func NewDatabaseServices(db *gorm.DB) *DatabaseServices {
	return &DatabaseServices{}
}
