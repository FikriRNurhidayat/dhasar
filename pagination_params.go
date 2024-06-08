package dhasar

import (
	"math"

	"github.com/fikrirnurhidayat/x/exists"
)

type PaginationParams struct {
	Page     uint32
	PageSize uint32
}

func (params PaginationParams) Normalize() PaginationParams {
	if !exists.Number(params.Page) {
		params.Page = 1
	}

	if !exists.Number(params.PageSize) {
		params.PageSize = 10
	}

	return params
}

func (params PaginationParams) Limit() uint32 {
	return params.PageSize
}

func (params PaginationParams) Offset() uint32 {
	return params.PageSize * (params.Page - 1)
}

func NewPaginationParams(page uint32, pageSize uint32) PaginationParams {
	params := PaginationParams{}

	if !exists.Number(page) {
		params.Page = 1
	}

	if !exists.Number(pageSize) {
		params.Page = 10
	}

	return params
}

type PaginationResult struct {
	Size      uint32
	Page      uint32
	PageSize  uint32
	PageCount uint32
}

func NewPaginationResult(params PaginationParams, size uint32) PaginationResult {
	return PaginationResult{
		Size:      size,
		Page:      params.Page,
		PageSize:  params.PageSize,
		PageCount: uint32(math.Ceil(float64(size) / float64(params.PageSize))),
	}
}
