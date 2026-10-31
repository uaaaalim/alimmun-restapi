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

func (repo *SessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*models.Session, error) {
	var session models.Session
	err := repo.DB.WithContext(ctx).Where("token_hash = ?", tokenHash).Take(&session).Error
	return &session, err
}
