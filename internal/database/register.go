package database

import "gorm.io/gorm"

type Service struct {
}

func NewDatabaseServices(db *gorm.DB) *Service {
	return &Service{}
}
