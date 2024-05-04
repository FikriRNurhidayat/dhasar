package common_specification

type LimitSpecification struct {
	Limit uint32
}

func WithLimit(limit uint32) Specification {
	return LimitSpecification{
		Limit: limit,
	}
}
