package exists

import (
	"time"
)

func String(value string) bool {
	return value != ""
}

func Date(value time.Time) bool {
	return !value.IsZero()
}

func Number(value uint32) bool {
	return value != 0
}
