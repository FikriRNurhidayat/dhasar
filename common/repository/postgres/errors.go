package postgres_repository

import (
	"net/http"

	common_errors "github.com/fikrirnurhidayat/dhasar/common/errors"
)

var (
	ErrInvalidSchema = &common_errors.DynamicError{
		Code:     http.StatusInternalServerError,
		Reason:   "INVALID_DATABASE_SCHEMA_ERROR",
		Template: "Mismatch or missing column: %s, Expected: %s, Found: %s",
	}
)
