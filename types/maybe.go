package dhasar_types

import (
	"encoding/json"
	"time"

	"github.com/fikrirnurhidayat/x/exists"
)

type Maybe[T any] struct {
	Valid bool
	Value T
}

func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(m.Value)
}

func MaybeTime(value time.Time) Maybe[time.Time] {
  return Maybe[time.Time]{
    Valid: exists.Date(value),
    Value: value,
  }
}

func MaybeString(value string) Maybe[string] {
  return Maybe[string]{
    Valid: exists.String(value),
    Value: value,
  }
}

func MaybeNumber(value uint32) Maybe[uint32] {
  return Maybe[uint32]{
    Valid: exists.Number(value),
    Value: value,
  }
}
