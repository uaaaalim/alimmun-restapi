package projects

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
)

func (h *Handler) ListAllProjects(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "HELLO YOPTA")
}
