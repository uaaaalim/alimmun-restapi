package projects

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

type Handler struct {
	ProjectService *services.ProjectService
	AuthService    *services.AuthService
}

func NewHandler(authService *services.AuthService, projectService *services.ProjectService) *Handler {
	return &Handler{AuthService: authService, ProjectService: projectService}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects", h.ListAllProjects)
	mux.HandleFunc("POST /v1/projects", h.AddProject)
	mux.HandleFunc("DELETE /v1/projects/{id}", h.DeleteProject)
}
