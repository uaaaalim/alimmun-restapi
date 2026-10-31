package projects

import (
	"net/http"
	"strconv"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"github.com/uaaaalim/alimmun-restapi/internal/response"
)

func (h *Handler) ListAllProjects(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if !query.Has("page") || !query.Has("page_size") {
		response.Fault(w, http.StatusBadRequest, "INVALID_QUERY",
			"\"page\" and \"page_size\" are required queries")
		return
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		response.Fault(w, http.StatusBadRequest, "INVALID_QUERY",
			"\"page\" must be a positive integer")
		return
	}

	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize <= 0 || pageSize > 100 {
		response.Fault(w, http.StatusBadRequest, "INVALID_QUERY",
			"\"page_size\" must be between 1 and 100")
		return
	}

	result, err := h.ProjectService.GetByPage(r.Context(), page, pageSize)
	if err != nil {
		response.InternalError(w, r, err)
		return
	}

	res := struct {
		Page     int              `json:"page"`
		PageSize int              `json:"page_size"`
		Results  []models.Project `json:"results"`
	}{
		Page:     page,
		PageSize: pageSize,
		Results:  result,
	}

	response.Success(w, http.StatusOK, res)
}
