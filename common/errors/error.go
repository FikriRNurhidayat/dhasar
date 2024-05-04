package common_errors

type Error struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
