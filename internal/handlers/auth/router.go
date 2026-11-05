package auth

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/cache"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
)

type Handler struct {
	DB    *database.Service
	Cache *cache.Service
}

func NewHandler(db *database.Service, cache *cache.Service) *Handler {
	return &Handler{DB: db, Cache: cache}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/login", h.login)
}
