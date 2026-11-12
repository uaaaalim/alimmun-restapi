package auth

import (
	"errors"
	"net/http"
	"time"
	"uuid"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	session, err := h.services.Auth.AuthorizeWithWrite(
		r.Context(), w,
		r.UserAgent(),
		r.Context().Value("real-ip").(string),
		r.Header.Get("Authorization"),
	)
	if err != nil {
		return
	}

	user, err := h.services.Auth.GetUserBySession(r.Context(), session)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			response.Fail(w, http.StatusNotFound, "USER_NOT_FOUND", "User was not found", nil)
			return
		}
		if errors.Is(err, services.ErrBannedUser) {
			var fields map[string]string
			if user.BanReason != nil {
				fields = map[string]string{"reason": *user.BanReason}
			}
			response.Fail(
				w, http.StatusForbidden, "BANNED_USER",
				"User is banned", fields,
			)
			return
		}
		if errors.Is(err, services.ErrUserIsNotActivated) {
			response.Fail(
				w, http.StatusNotFound, "ACCOUNT_NOT_ACTIVE",
				"Account is not activated", nil,
			)
			return
		}
		response.InternalError(w, err)
		return
	}

	res := struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		IsActivated bool      `json:"is_activated"`
		IsBanned    bool      `json:"is_banned"`
		BanReason   *string   `json:"ban_reason"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		IsActivated: user.IsActivated,
		IsBanned:    user.IsBanned,
		BanReason:   user.BanReason,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	response.OK(w, http.StatusOK, res)
}
