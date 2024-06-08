package dhasar

type SortSpecification struct {
	Arguments SortArguments
}

func Sort(args ...SortArgument) Specification {
	return SortSpecification{
		Arguments: args,
	}
}
