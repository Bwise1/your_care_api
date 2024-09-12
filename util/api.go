package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"

	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/pkg/errors"
)

// StatusCode returns the status code represented
// by the specified status. Note that this function
// returns a status code of 200 by default
func StatusCode(status string) int {
	switch status {
	case values.Error:
		return http.StatusInternalServerError
	case values.Created:
		return http.StatusCreated
	case values.BadRequestBody:
		return http.StatusBadRequest
	case values.Unprocessable:
		return http.StatusUnprocessableEntity
	case values.NotAllowed:
		return http.StatusForbidden
	case values.Conflict:
		return http.StatusConflict
	case values.NotFound:
		return http.StatusNotFound
	case values.NotAuthorised:
		return http.StatusUnauthorized
	case values.ActiveLogin:
		return http.StatusForbidden
	default:
		return http.StatusOK
	}
}

const UserAuth = "user-auth"
const AdminAuth = "admin-auth"

// DecodeJSONBody ...
func DecodeJSONBody(tc *tracing.Context, body io.ReadCloser, target interface{}) error {
	defer func() {
		_ = body.Close()
	}()

	if body == nil {
		return fmt.Errorf("missing request body for request: %v", tc)
	}

	if err := json.NewDecoder(body).Decode(&target); err != nil {
		return errors.Wrapf(err, "Error parsing json body for request: %v", tc)
	}

	return nil
}

func ValidEmail(email string) error {
	if email == "" {
		return errors.New("invalid email address")
	}
	_, err := mail.ParseAddress(email)
	return err
}
