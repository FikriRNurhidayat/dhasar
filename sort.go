package dhasar

type SortDirection int

type SortColumn string
type SortColumns []string

const (
	SortAscending SortDirection = iota
	SortDescending
)

type SortArgument struct {
	Column    SortColumn
	Direction SortDirection
}

type SortArguments []SortArgument
