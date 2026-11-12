package auth

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

type Handler struct {
	services *services.Services
}

func NewAuthHandler(appServices *services.Services) *Handler {
	return &Handler{services: appServices}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Route /auth
	mux.HandleFunc("POST /auth/login", h.login)
	mux.HandleFunc("POST /auth/register", h.register)
	mux.HandleFunc("GET /auth/me", h.me)
}
