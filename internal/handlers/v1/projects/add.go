package projects

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"github.com/uaaaalim/alimmun-restapi/internal/response"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

func (h *Handler) AddProject(w http.ResponseWriter, r *http.Request) {
	_, err := h.AuthService.ValidateAdminSession(r.Context(), r.Header.Get("Authorization"))
	if errors.Is(err, services.AuthInvalidSession) {
		response.UnauthorizedToken(w)
		return
	} else if errors.Is(err, services.AuthAdminRequired) {
		response.Fault(w, http.StatusForbidden, "FORBIDDEN", "Admin role required")
		return
	} else if err != nil {
		response.InternalError(w, r, err)
		return
	}

	var body struct {
		Name             string               `json:"name"`
		ShortDescription string               `json:"short_description"`
		FullDescription  string               `json:"full_description"`
		Status           models.ProjectStatus `json:"status"`
		StartYear        int                  `json:"start_year"`
		CloseYear        *int                 `json:"close_year"`
		TechnologyIDs    []int64              `json:"technology_ids"`
		SourceCodeLink   *string              `json:"source_code_link"`
		ProductionLink   *string              `json:"production_link"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		response.Fault(w, http.StatusBadRequest, "INVALID_BODY", "Invalid JSON")
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		response.Fault(w, http.StatusBadRequest, "INVALID_BODY", "Body must contain one JSON object")
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	body.ShortDescription = strings.TrimSpace(body.ShortDescription)
	body.FullDescription = strings.TrimSpace(body.FullDescription)
	if body.Name == "" || utf8.RuneCountInString(body.Name) > 100 ||
		body.ShortDescription == "" || utf8.RuneCountInString(body.ShortDescription) > 150 || body.FullDescription == "" {
		response.Fault(w, http.StatusBadRequest, "INVALID_PROJECT",
			"name (1-100), short_description (1-150) and full_description are required")
		return
	}
	switch body.Status {
	case models.Active, models.InDevelopment, models.Suspended, models.Abandoned:
	default:
		response.Fault(w, http.StatusBadRequest, "INVALID_PROJECT", "Unknown project status")
		return
	}
	if body.StartYear < 1 || body.StartYear > 9999 ||
		(body.CloseYear != nil && (*body.CloseYear < body.StartYear || *body.CloseYear > 9999)) {
		response.Fault(w, http.StatusBadRequest, "INVALID_PROJECT",
			"Years must be between 1 and 9999; close_year must not precede start_year")
		return
	}
	for _, link := range []*string{body.SourceCodeLink, body.ProductionLink} {
		if link == nil {
			continue
		}
		parsed, err := url.Parse(*link)
		if err != nil || len(*link) > 255 || parsed.Hostname() == "" || parsed.User != nil ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") {
			response.Fault(w, http.StatusBadRequest, "INVALID_PROJECT",
				"Links must be HTTP or HTTPS URLs of at most 255 characters")
			return
		}
	}

	project := models.Project{
		Name: body.Name, ShortDescription: body.ShortDescription, FullDescription: body.FullDescription,
		Status: body.Status, StartYear: body.StartYear, CloseYear: body.CloseYear,
		SourceCodeLink: body.SourceCodeLink, ProductionLink: body.ProductionLink,
	}
	err = h.ProjectService.Create(r.Context(), &project, body.TechnologyIDs)
	if err != nil && errors.Is(err, services.ProjectTechnologiesInvalid) {
		response.Fault(w, http.StatusBadRequest, "INVALID_PROJECT", err.Error())
		return
	}
	if err != nil {
		response.InternalError(w, r, err)
		return
	}
	response.Success(w, http.StatusCreated, project)
}
