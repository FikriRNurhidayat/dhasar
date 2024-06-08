package dhasar

import (
	"strings"

	"github.com/fikrirnurhidayat/x/exists"
	"github.com/labstack/echo/v4"
)

type SortParams struct {
	Columns   SortColumns
	Arguments SortArguments
}

func NewSortParams(columns SortColumns) *SortParams {
	return &SortParams{
		Columns: columns,
	}
}

func (s *SortParams) columnAllowed(column SortColumn) bool {
	for _, col := range s.Columns {
		if column == SortColumn(col) {
			return true
		}
	}

	return false
}

func (s *SortParams) ParseFromContext(c echo.Context) error {
	sortStr := c.QueryParam("sort")
	if !exists.String(sortStr) {
		return nil
	}

	candidates := strings.Split(sortStr, ",")
	for _, candidate := range candidates {
		direction := SortAscending
		if strings.HasPrefix(candidate, "-") {
			direction = SortDescending
		}

		columnStr := strings.ReplaceAll(candidate, "-", "")
		columnStr = strings.ReplaceAll(columnStr, "+", "")
		column := SortColumn(columnStr)

		if !s.columnAllowed(column) {
			return ErrInvalidSortParams
		}

		s.Arguments = append(s.Arguments, SortArgument{
			Column:    column,
			Direction: direction,
		})
	}

	return nil
}
