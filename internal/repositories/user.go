package repositories

import (
	"context"
	"errors"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

var (
	ErrEmailTaken = errors.New("email is already registered")
)

func (repo *UserRepository) Create(ctx context.Context, user *models.User) error {
	err := repo.db.WithContext(ctx).Create(user).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrEmailTaken
	}

	return err
}
