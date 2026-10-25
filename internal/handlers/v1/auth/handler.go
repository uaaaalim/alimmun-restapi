package auth

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

type Handler struct {
	AuthService *services.AuthService
}

func NewHandler(authService *services.AuthService) *Handler {
	return &Handler{AuthService: authService}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/auth/login", h.Login)
	mux.HandleFunc("POST /v1/auth/register", h.Register)
}
