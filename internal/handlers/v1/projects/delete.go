package projects

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	_, err := h.AuthService.ValidateAdminSession(r.Context(), r.Header.Get("Authorization"))
	if errors.Is(err, services.AuthInvalidSession) {
		response.UnauthorizedToken(w)
		return
	}
	if errors.Is(err, services.AuthAdminRequired) {
		response.Fault(w, http.StatusForbidden, "FORBIDDEN", "Admin role required")
		return
	}
	if err != nil {
		response.InternalError(w, r, err)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fault(w, http.StatusBadRequest, "INVALID_ID", "Project ID must be a positive integer")
		return
	}
	err = h.ProjectService.Delete(r.Context(), id)
	if errors.Is(err, services.ProjectNotFound) {
		response.Fault(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}
	if err != nil {
		response.InternalError(w, r, err)
		return
	}
	response.Success(w, http.StatusOK, nil)
}
