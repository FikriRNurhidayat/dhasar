package common_errors

import "fmt"

type DynamicError struct {
	Code     int
	Reason   string
	Template string
}

func (e *DynamicError) Format(args ...any) *Error {
	return &Error{
		Code:    e.Code,
		Reason:  e.Reason,
		Message: fmt.Sprintf(e.Template, args...),
	}
}
