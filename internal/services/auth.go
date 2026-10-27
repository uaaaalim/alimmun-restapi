package services

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	UsersRepo             *repository.UserRepository
	SessionRepo           *repository.SessionRepository
	LoginAttemptRepo      *repository.LoginAttemptRepository
	RegistrationLimitRepo *repository.RegistrationLimitRepository
}

func NewAuthService(
	usersRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	loginAttemptRepo *repository.LoginAttemptRepository,
	registrationLimitRepo *repository.RegistrationLimitRepository,
) *AuthService {
	return &AuthService{
		UsersRepo:             usersRepo,
		SessionRepo:           sessionRepo,
		LoginAttemptRepo:      loginAttemptRepo,
		RegistrationLimitRepo: registrationLimitRepo,
	}
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

// rate limiting

func (s *AuthService) RecordLoginFailure(ctx context.Context, ip, email string) error {
	const (
		maxAttempts    = 5
		attemptsWindow = 15 * time.Minute
		banTime        = 30 * time.Minute
	)

	count, err := s.LoginAttemptRepo.IncrementLoginAttempt(ctx, ip, email)
	if err != nil {
		return fmt.Errorf("increment attempt: %w", err)
	}

	if count == 1 {
		if err := s.LoginAttemptRepo.ExpireAttempts(ctx, ip, email, attemptsWindow); err != nil {
			return fmt.Errorf("expire attempts: %w", err)
		}
	} else if count >= maxAttempts {
		if err := s.LoginAttemptRepo.BanAttempts(ctx, ip, email, banTime); err != nil {
			return fmt.Errorf("set ban with duration: %w", err)
		}

		if err := s.LoginAttemptRepo.ResetFailures(ctx, ip, email); err != nil {
			return fmt.Errorf("delete attempts: %w", err)
		}
	}

	return nil
}

func (s *AuthService) LoginBanRemaining(ctx context.Context, ip, email string) (time.Duration, error) {
	remaining, err := s.LoginAttemptRepo.BanRemaining(ctx, ip, email)
	if err != nil {
		return 0, fmt.Errorf("get login ban duration: %w", err)
	}

	// -2 unknown key, -1 no ttl
	switch remaining {
	case -2:
		return 0, nil
	case -1:
		return 0, fmt.Errorf("login ban has no expiration")
	}

	if remaining == 0 {
		return time.Millisecond, nil
	}

	return remaining, nil
}

func (s *AuthService) RegisterBanRemaining(ctx context.Context, ip string) (time.Duration, error) {
	return s.RegistrationLimitRepo.BanRemaining(ctx, ip)
}

func (s *AuthService) ResetLoginFailures(ctx context.Context, ip, email string) error {
	if err := s.LoginAttemptRepo.ResetFailures(ctx, ip, email); err != nil {
		return fmt.Errorf("reset login attempts: %w", err)
	}

	return nil
}
