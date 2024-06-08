package dhasar

import "github.com/google/uuid"

type BinderFunc func(values []string) []error

func UUIDBinder(v *uuid.UUID) BinderFunc {
	return func(values []string) []error {
		*v = uuid.MustParse(values[0])
		return nil
	}
}

func MustUUIDBinder(v *uuid.UUID) BinderFunc {
	return func(values []string) []error {
		id, err := uuid.Parse(values[0])
		if err != nil {
			return []error{err}
		}

		*v = id

		return nil
	}
}

func UUIDSliceBinder(v []uuid.UUID) BinderFunc {
	return func(values []string) []error {
		for i, idStr := range values {
			v[i] = uuid.MustParse(idStr)
		}

		return nil
	}
}

func MustUUIDSliceBinder(v []uuid.UUID) BinderFunc {
	return func(values []string) []error {
		errors := []error{}
		for i, idStr := range values {
			id, err := uuid.Parse(idStr)
			if err != nil {
				errors = append(errors, err)
			}

			v[i] = id
		}

		if len(errors) > 0 {
			return errors
		}

		return nil
	}
}
