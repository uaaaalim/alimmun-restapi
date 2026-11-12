package repositories

import (
	"context"
	"errors"
	"uuid"

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
	ErrEmailTaken   = errors.New("email is already registered")
	ErrUserNotFound = errors.New("user not found")
)

func (repo *UserRepository) Create(ctx context.Context, user *models.User) error {
	err := repo.db.WithContext(ctx).Create(user).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrEmailTaken
	}

	return err
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := repo.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	err := repo.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
