package common_specification

type OffsetSpecification struct {
	Offset uint32
}

func WithOffset(offset uint32) Specification {
	return OffsetSpecification{
		Offset: offset,
	}
}
