package common_errors

import "net/http"

var (
	ErrInternalServer = &Error{
		Code:    http.StatusInternalServerError,
		Reason:  "INTERNAL_SERVER_ERROR",
		Message: "Internal server error.",
	}

	ErrBadRequest = &Error{
		Code:    http.StatusBadRequest,
		Reason:  "BAD_REQUEST",
		Message: "Bad request. Please pass valid request parameters.",
	}

	ErrNotFound = &DynamicError{
		Code:     http.StatusNotFound,
		Reason:   "ROUTE_NOT_FOUND",
		Template: "Route '%s %s' not found.",
	}

	ErrInvalidUUID = &Error{
		Code:    http.StatusBadRequest,
		Reason:  "INVALID_UUID_ERROR",
		Message: "UUID is not valid. Please pass valid UUID.",
	}
)
