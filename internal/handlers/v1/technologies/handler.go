package technologies

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

type Handler struct {
	AuthService       *services.AuthService
	TechnologyService *services.TechnologyService
}

func NewHandler(authService *services.AuthService, technologyService *services.TechnologyService) *Handler {
	return &Handler{AuthService: authService, TechnologyService: technologyService}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/technologies", h.List)
	mux.HandleFunc("POST /v1/technologies", h.Create)
	mux.HandleFunc("DELETE /v1/technologies/{id}", h.Delete)
}

func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) bool {
	_, err := h.AuthService.ValidateAdminSession(r.Context(), r.Header.Get("Authorization"))
	if errors.Is(err, services.AuthInvalidSession) {
		response.UnauthorizedToken(w)
		return false
	}
	if errors.Is(err, services.AuthAdminRequired) {
		response.Fault(w, http.StatusForbidden, "FORBIDDEN", "Admin role required")
		return false
	}
	if err != nil {
		response.InternalError(w, r, err)
		return false
	}
	return true
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	technologies, err := h.TechnologyService.List(r.Context())
	if err != nil {
		response.InternalError(w, r, err)
		return
	}
	response.Success(w, http.StatusOK, technologies)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.authorize(w, r) {
		return
	}
	var body struct {
		Name  string  `json:"name"`
		Slug  string  `json:"slug"`
		Image *string `json:"image"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		response.Fault(w, http.StatusBadRequest, "INVALID_BODY", "Invalid technology JSON")
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		response.Fault(w, http.StatusBadRequest, "INVALID_BODY", "Body must contain one JSON object")
		return
	}
	technology, err := h.TechnologyService.Create(r.Context(), body.Name, body.Slug, body.Image)
	if errors.Is(err, services.TechnologyInvalid) || errors.Is(err, services.TechnologyInvalidImage) {
		response.Fault(w, http.StatusBadRequest, "INVALID_TECHNOLOGY", "Invalid technology")
		return
	}
	if errors.Is(err, services.TechnologyAlreadyExists) {
		response.Fault(w, http.StatusConflict, "TECHNOLOGY_ALREADY_EXISTS", "Technology already exists")
		return
	}
	if err != nil {
		response.InternalError(w, r, err)
		return
	}
	response.Success(w, http.StatusCreated, technology)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.authorize(w, r) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fault(w, http.StatusBadRequest, "INVALID_ID", "Technology ID must be a positive integer")
		return
	}
	if err := h.TechnologyService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, services.TechnologyNotFound) {
			response.Fault(w, http.StatusNotFound, "TECHNOLOGY_NOT_FOUND", "T	echnology was not found")
			return
		}
		response.InternalError(w, r, err)
		return
	}
	response.Success(w, http.StatusOK, nil)
}
