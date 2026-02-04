package models

// ListOptions contains common pagination and filtering options
type ListOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Order    string // "asc" or "desc"
	Filters  map[string]any
}
