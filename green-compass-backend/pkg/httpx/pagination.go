package httpx

import (
	"math"
)

// PaginationParams holds request pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// ValidatePagination ensures pagination parameters are within bounds
func ValidatePagination(page, limit int) PaginationParams {
	const (
		minLimit = 1
		maxLimit = 100
		minPage  = 1
	)

	if page < minPage {
		page = minPage
	}
	if limit < minLimit {
		limit = 20 // default
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// CalculatePaginationMeta computes pagination metadata
func CalculatePaginationMeta(page, limit, total int) *PaginationMeta {
	return &PaginationMeta{
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasNext: page < int(math.Ceil(float64(total)/float64(limit))),
	}
}

// CalculateOffset computes database OFFSET for a page
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}
