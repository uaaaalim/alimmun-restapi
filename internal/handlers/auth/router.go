package auth

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

type Handler struct {
	services *services.Services
}

func NewHandler(appServices *services.Services) *Handler {
	return &Handler{services: appServices}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.login)
	mux.HandleFunc("POST /auth/register", h.register)
}
