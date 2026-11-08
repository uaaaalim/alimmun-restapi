package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
)

var (
	ErrBodyTooLarge = errors.New("request body too large")
	ErrInvalidJSON  = errors.New("invalid JSON")
)

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any, ackHttp bool) error {
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return classifyError(w, err, ackHttp)
	}

	if err := decoder.Decode(new(any)); err != io.EOF {
		return classifyError(w, err, ackHttp)
	}

	return nil
}

func classifyError(w http.ResponseWriter, err error, ackHttp bool) error {
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		if ackHttp {
			response.Fail(w, http.StatusRequestEntityTooLarge, "body_too_large",
				"The body request is too large",
				nil,
			)
		}
		return ErrBodyTooLarge
	}

	if ackHttp {
		response.Fail(w, http.StatusBadRequest, "invalid_json",
			"Body contains invalid json",
			nil,
		)
	}
	return ErrInvalidJSON
}
