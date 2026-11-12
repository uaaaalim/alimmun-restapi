package router

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/handlers/auth"
	"github.com/uaaaalim/alimmun-restapi/internal/middlewares"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

func NewGeneralRouter(appServices *services.Services) http.Handler {
	mux := http.NewServeMux()

	// routes
	auth.NewAuthHandler(appServices).RegisterRoutes(mux)

	// middlewares
	handler := middlewares.RealIP(mux)
	return middlewares.Logger(handler)
}
