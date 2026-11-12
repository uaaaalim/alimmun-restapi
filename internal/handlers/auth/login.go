package auth

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
	"github.com/uaaaalim/alimmun-restapi/internal/security"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
	"github.com/uaaaalim/alimmun-restapi/internal/utils"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var input loginRequest
	if err := utils.ReadJSON(w, r, &input, true); err != nil {
		return
	}

	fields := make(map[string]string)

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	address, err := mail.ParseAddress(input.Email)
	switch {
	case input.Email == "":
		fields["email"] = "required"
	case len(input.Email) > 254:
		fields["email"] = "maximum length is 254 bytes"
	case err != nil:
		fields["email"] = "invalid email"
	case address.Address != input.Email:
		fields["email"] = "invalid email"
	}

	if input.Password == "" {
		fields["password"] = "required"
	} else {
		switch {
		case utf8.RuneCountInString(input.Password) < 8:
			fields["password"] = "minimum length is 8 characters"
		case len(input.Password) > 1024:
			fields["password"] = "maximum length is 1024 bytes"
		}
	}

	if len(fields) > 0 {
		response.Fail(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Check the provided fields",
			fields,
		)
		return
	}

	token, err := security.GenerateToken()
	if err != nil {
		response.InternalError(w, err)
		return
	}

	session, err := h.services.Auth.Login(
		r.Context(),
		input.Email,
		input.Password,
		r.UserAgent(),
		r.Header.Get("X-Real-IP"),
		security.HashToken(token),
	)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			response.Fail(
				w, http.StatusUnauthorized, "INVALID_CREDENTIALS",
				"Password or email are incorrect",
				nil,
			)
			return
		}
		response.InternalError(w, err)
		return
	}

	var res = struct {
		ID           uuid.UUID `json:"id"`
		Token        string    `json:"token"`
		ExpiresAfter int       `json:"expires_after"` // seconds
	}{
		ID:           session.ID,
		Token:        token,
		ExpiresAfter: max(0, int(time.Until(session.ExpiresAt).Seconds())),
	}

	response.OK(w, http.StatusOK, res)
}
