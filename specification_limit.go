package dhasar

type LimitSpecification struct {
	Limit uint32
}

func WithLimitSpecs(limit uint32) Specification {
	return LimitSpecification{
		Limit: limit,
	}
}
