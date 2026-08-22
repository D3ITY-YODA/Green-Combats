package httpx

// PaginationParams holds request pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
	Limit    int // Alias for PageSize — callers may use either
}

// MaxPageSize is the upper bound for page sizes accepted by the API.
const MaxPageSize = 100

// ValidatePagination ensures pagination parameters are within bounds
func ValidatePagination(page, pageSize int) PaginationParams {
	const (
		minPageSize = 1
		maxPageSize = 100
		minPage     = 1
	)

	if page < minPage {
		page = minPage
	}
	if pageSize < minPageSize {
		pageSize = 20 // default
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Limit:    pageSize,
	}
}

// CalculatePaginationMeta computes pagination metadata
func CalculatePaginationMeta(page, pageSize, total int) *PaginationMeta {
	return &PaginationMeta{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		NextCursor: "",
	}
}

// CalculateOffset computes database OFFSET for a page
func CalculateOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
