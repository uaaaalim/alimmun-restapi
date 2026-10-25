package router

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/handlers"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers/v1/auth"
	"github.com/uaaaalim/alimmun-restapi/internal/repository"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
	"gorm.io/gorm"
)

func CreateRouter(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// GET /
	mux.HandleFunc("GET /{$}", handlers.Home)

	// repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// services
	userService := services.NewAuthService(userRepo, sessionRepo)

	// handlers
	auth.NewHandler(userService).RegisterRoutes(mux) // for /v1/auth

	return mux
}
