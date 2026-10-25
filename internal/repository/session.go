package repository

import (
	"context"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (repo *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	return repo.DB.WithContext(ctx).Create(session).Error
}
