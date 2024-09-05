package util

import (
	"net/http"

	"github.com/bwise1/your_care_api/util/values"
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
