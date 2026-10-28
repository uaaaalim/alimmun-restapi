package router

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/database"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers/v1/auth"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers/v1/projects"
	"github.com/uaaaalim/alimmun-restapi/internal/repository"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
	"gorm.io/gorm"
)

func CreateRouter(db *gorm.DB, redis *database.RedisConnection) *http.ServeMux {
	mux := http.NewServeMux()

	// GET /
	mux.HandleFunc("GET /{$}", handlers.Home)

	// repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	loginAttemptRepo := repository.NewLoginAttemptRepository(redis)
	registrationLimitRepo := repository.NewRegistrationLimitRepository(redis)

	// services
	userService := services.NewAuthService(userRepo, sessionRepo, loginAttemptRepo, registrationLimitRepo)

	// handlers
	auth.NewHandler(userService).RegisterRoutes(mux) // for /v1/auth
	projects.NewHandler().RegisterRoutes(mux)        // for /v1/projects

	return mux
}
