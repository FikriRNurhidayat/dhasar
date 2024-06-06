package dhasar

type OffsetSpecification struct {
	Offset uint32
}

func WithOffsetSpecs(offset uint32) Specification {
	return OffsetSpecification{
		Offset: offset,
	}
}
