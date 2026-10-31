package router

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/database"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers/v1/auth"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers/v1/projects"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers/v1/technologies"
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
	projectRepo := repository.NewProjectRepository(db)
	technologyRepo := repository.NewTechnologyRepository(db)

	// services
	authService := services.NewAuthService(userRepo, sessionRepo, loginAttemptRepo, registrationLimitRepo)
	projectService := services.NewProjectService(projectRepo)
	technologyService := services.NewTechnologyService(technologyRepo)

	// handlers
	auth.NewHandler(authService).RegisterRoutes(mux)                     // for /v1/auth
	projects.NewHandler(authService, projectService).RegisterRoutes(mux) // for /v1/projects
	technologies.NewHandler(authService, technologyService).RegisterRoutes(mux)

	return mux
}
