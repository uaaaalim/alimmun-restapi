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

	// /v1/auth
	userRepo := repository.NewUserRepository(db)
	userService := services.NewAuthService(userRepo)

	auth.NewHandler(userService).RegisterRoutes(mux)

	return mux
}
