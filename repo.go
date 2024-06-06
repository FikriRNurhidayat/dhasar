package dhasar

import (
	"context"
)

type ListArgs[T any] struct {
	Filters []T
	Sort    Specification
	Limit   Specification
	Offset  Specification
}

type Repository[Entity any, Specification any] interface {
	Save(context.Context, Entity) error
	Get(context.Context, ...Specification) (Entity, error)
	Exist(context.Context, ...Specification) (bool, error)
	Delete(context.Context, ...Specification) error
	List(context.Context, ListArgs[Specification]) ([]Entity, error)
	Each(context.Context, ListArgs[Specification]) (Iterator[Entity], error)
	Size(context.Context, ...Specification) (uint32, error)
}

type Iterator[Entity any] interface {
	Next() bool
	Current() (Entity, error)
}
