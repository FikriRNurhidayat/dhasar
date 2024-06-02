package postgres_repository

import (
	"net/http"

	dhasar_errors "github.com/fikrirnurhidayat/dhasar/errors"
)

var (
	ErrInvalidSchema = &dhasar_errors.DynamicError{
		Code:     http.StatusInternalServerError,
		Reason:   "INVALID_DATABASE_SCHEMA_ERROR",
		Template: "Mismatch or missing column: %s, Expected: %s, Found: %s",
	}
)
