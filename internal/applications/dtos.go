package applications

import (
	"time"
)

// Common response types used across multiple use cases

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortDir  string `json:"sort_dir" form:"sort_dir"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ListResponse represents a paginated list response
type ListResponse[T any] struct {
	Data       []T                 `json:"data"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

// SuccessResponse represents a successful operation response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SearchRequest represents search parameters
type SearchRequest struct {
	Query      string `json:"query" form:"query"`
	Field      string `json:"field" form:"field"`
	Exact      bool   `json:"exact" form:"exact"`
	CaseSensitive bool `json:"case_sensitive" form:"case_sensitive"`
}

// DateRangeRequest represents date range parameters
type DateRangeRequest struct {
	StartDate time.Time `json:"start_date" form:"start_date"`
	EndDate   time.Time `json:"end_date" form:"end_date"`
}

// FilterRequest represents generic filter parameters
type FilterRequest struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, ne, gt, gte, lt, lte, like, in
	Value    interface{} `json:"value"`
}

// BulkOperationRequest represents bulk operation parameters
type BulkOperationRequest struct {
	IDs       []string `json:"ids"`
	Operation string   `json:"operation"` // delete, update, etc.
	Data      interface{} `json:"data,omitempty"`
}

// BulkOperationResponse represents bulk operation results
type BulkOperationResponse struct {
	TotalRequested int      `json:"total_requested"`
	TotalProcessed int      `json:"total_processed"`
	TotalFailed    int      `json:"total_failed"`
	SuccessIDs     []string `json:"success_ids"`
	FailedIDs      []string `json:"failed_ids"`
	Errors         []ErrorResponse `json:"errors,omitempty"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationResponse represents validation results
type ValidationResponse struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// HealthCheckResponse represents health check status
type HealthCheckResponse struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	Timestamp  time.Time         `json:"timestamp"`
	Services   map[string]string `json:"services"`
	Dependencies map[string]bool `json:"dependencies"`
}

// StatisticsResponse represents general statistics
type StatisticsResponse struct {
	Entity string                 `json:"entity"`
	Period string                 `json:"period,omitempty"`
	Count  int                    `json:"count"`
	Data   map[string]interface{} `json:"data"`
}

// RelationshipRequest represents relationship operations
type RelationshipRequest struct {
	ParentID   string `json:"parent_id"`
	ChildID    string `json:"child_id"`
	Type       string `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// RelationshipResponse represents relationship data
type RelationshipResponse struct {
	ParentID   string                 `json:"parent_id"`
	ChildID    string                 `json:"child_id"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at,omitempty"`
}

// AuditLogResponse represents audit log entries
type AuditLogResponse struct {
	ID        string    `json:"id"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Action    string    `json:"action"`
	Changes   map[string]interface{} `json:"changes,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Constants for common operations
const (
	// Sort directions
	SortAsc  = "asc"
	SortDesc = "desc"
	
	// Filter operators
	OpEqual              = "eq"
	OpNotEqual           = "ne"
	OpGreaterThan        = "gt"
	OpGreaterThanOrEqual = "gte"
	OpLessThan           = "lt"
	OpLessThanOrEqual    = "lte"
	OpLike               = "like"
	OpIn                 = "in"
	OpNotIn              = "nin"
	OpBetween            = "between"
	
	// Bulk operations
	BulkOpDelete = "delete"
	BulkOpUpdate = "update"
	BulkOpCreate = "create"
	
	// Health check statuses
	HealthStatusHealthy   = "healthy"
	HealthStatusUnhealthy = "unhealthy"
	HealthStatusDegraded  = "degraded"
	
	// Default pagination values
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Helper functions for DTOs

// NewListResponse creates a new paginated list response
func NewListResponse[T any](data []T, pagination *PaginationResponse) *ListResponse[T] {
	return &ListResponse[T]{
		Data:       data,
		Pagination: pagination,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewFieldErrorResponse creates a new field-specific error response
func NewFieldErrorResponse(code, message, field string, value interface{}) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Field:   field,
		Value:   value,
	}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Message: message,
		Data:    data,
	}
}

// ValidatePaginationRequest validates and normalizes pagination parameters
func ValidatePaginationRequest(req *PaginationRequest) {
	if req.Page <= 0 {
		req.Page = DefaultPage
	}
	if req.PageSize <= 0 {
		req.PageSize = DefaultPageSize
	}
	if req.PageSize > MaxPageSize {
		req.PageSize = MaxPageSize
	}
	if req.SortDir != SortAsc && req.SortDir != SortDesc {
		req.SortDir = SortAsc
	}
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize, total int) *PaginationResponse {
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}
	
	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// ValidateDateRange validates that start date is before end date
func ValidateDateRange(start, end time.Time) bool {
	return !start.IsZero() && !end.IsZero() && start.Before(end)
}