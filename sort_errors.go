package dhasar

import "net/http"

var (
	ErrInvalidSortParams = &Error{
		Code:    http.StatusBadRequest,
		Reason:  "INVALID_SORT_PARAMS",
		Message: "Sort parameter is not valid. Please pass valid sort parameters.",
	}
)
