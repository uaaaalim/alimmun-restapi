package services

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"
	"uuid"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"github.com/uaaaalim/alimmun-restapi/internal/repositories"
	"github.com/uaaaalim/alimmun-restapi/internal/security"
)

type AuthService struct {
	repos *repositories.Repositories
}

func NewAuthService(repos *repositories.Repositories) *AuthService {
	return &AuthService{repos: repos}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type rateLimitCheck struct {
	key    string
	window time.Duration
	limit  int
}

func (s *AuthService) Register(ctx context.Context, email, name, password string) (*models.User, error) {
	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		ID:           uuid.NewV4(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}
	if err := s.repos.Users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email, password, userAgent, ip, tokenHash string,
) (*models.Session, error) {
	user, err := s.repos.Users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	loggedIn, err := security.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("verify password: %w", err)
	}
	if !loggedIn {
		return nil, ErrInvalidCredentials
	}
	if user.IsBanned {
		return nil, ErrInvalidCredentials
	}

	session := &models.Session{
		ID:        uuid.NewV4(),
		UserID:    user.ID,
		UserAgent: userAgent,
		IP:        ip,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	if err := s.repos.Sessions.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}

func (s *AuthService) IsRegisterEmailRateLimited(
	ctx context.Context,
	email string,
) (bool, time.Duration, error) {
	digest := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))

	checks := []rateLimitCheck{
		{
			key:    fmt.Sprintf("rl:register:email:hour:%s", digest),
			window: time.Minute * 15,
			limit:  3,
		},
	}

	return s.checkRateLimits(ctx, checks)
}

func (s *AuthService) IsRegisterIPRateLimited(
	ctx context.Context,
	ip string,
) (bool, time.Duration, error) {
	checks := []rateLimitCheck{
		{
			key:    fmt.Sprintf("rl:register:ip:hour:%s", ip),
			window: time.Hour,
			limit:  30,
		},
		{
			key:    fmt.Sprintf("rl:register:ip:minute:%s", ip),
			window: time.Minute,
			limit:  10,
		},
	}

	return s.checkRateLimits(ctx, checks)
}

func (s *AuthService) checkRateLimits(
	ctx context.Context,
	checks []rateLimitCheck,
) (bool, time.Duration, error) {
	for _, check := range checks {
		attempts, ttl, err := s.repos.Auth.IncrementAttempts(
			ctx,
			check.key,
			check.window,
		)
		if err != nil {
			return false, 0, fmt.Errorf("check rate limit: %w", err)
		}

		if attempts > check.limit && ttl > 0 {
			return true, ttl, nil
		}
	}

	return false, 0, nil
}
