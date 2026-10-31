package repository

import (
	"context"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	if err := repo.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := repo.DB.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) Create(ctx context.Context, user *models.User) error {
	return repo.DB.WithContext(ctx).Create(user).Error
}
