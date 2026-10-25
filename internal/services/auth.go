package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/uaaaalim/alimmun-restapi/internal/common/security"
	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"github.com/uaaaalim/alimmun-restapi/internal/repository"
	"gorm.io/gorm"
)

var AuthInvalidCredentials = errors.New("invalid credentials: email or password")
var AuthUserAlreadyExists = errors.New("user is already created")

type SessionInfo struct {
	IPAddress string
	UserAgent string
}

type AuthService struct {
	UsersRepo   *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

func NewAuthService(usersRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{UsersRepo: usersRepo, SessionRepo: sessionRepo}
}

func (s *AuthService) Login(
	ctx context.Context,
	email, password string,
) (*models.User, error) {
	user, err := s.UsersRepo.GetByEmail(ctx, email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, AuthInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := security.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, AuthInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) Register(
	ctx context.Context,
	name, email, password string,
) (*models.User, error) {
	_, err := s.UsersRepo.GetByEmail(ctx, email)

	if err == nil {
		return nil, AuthUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check user: %w", err)
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}

	if err := s.UsersRepo.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, AuthUserAlreadyExists
		}

		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
func (s *AuthService) CreateSession(ctx context.Context, session *models.Session) error {
	if err := s.SessionRepo.Create(ctx, session); err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	return nil
}
