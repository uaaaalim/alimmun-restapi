package services

import "github.com/uaaaalim/alimmun-restapi/internal/repositories"

type Services struct {
	Auth *AuthService
}

func NewServices(repos *repositories.Repositories) *Services {
	return &Services{
		Auth: NewAuthService(repos),
	}
}
