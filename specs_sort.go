package dhasar

type SortArg struct {
	Column    string
	Direction string
}

type SortSpecification struct {
	Args []SortArg
}

func Sort(args ...SortArg) Specification {
	return SortSpecification{
		Args: args,
	}
}
