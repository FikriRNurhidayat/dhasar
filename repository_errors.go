package dhasar

import (
	"net/http"
)

var (
	ErrInvalidSchema = &DynamicError{
		Code:     http.StatusInternalServerError,
		Reason:   "INVALID_DATABASE_SCHEMA_ERROR",
		Template: "Mismatch or missing column: %s, Expected: %s, Found: %s",
	}
)
