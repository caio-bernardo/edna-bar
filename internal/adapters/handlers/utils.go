package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// HTTPUtils provides common utilities for HTTP handlers
type HTTPUtils struct{}

// NewHTTPUtils creates a new HTTPUtils instance
func NewHTTPUtils() *HTTPUtils {
	return &HTTPUtils{}
}

// DecodeJSONBody decodes JSON request body into the provided struct
func (h *HTTPUtils) DecodeJSONBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	
	return nil
}

// SendJSONResponse sends a JSON response with the specified status code
func (h *HTTPUtils) SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// SendErrorResponse sends a standardized error response
func (h *HTTPUtils) SendErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResponse := map[string]interface{}{
		"error":     message,
		"code":      statusCode,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	if err != nil {
		errorResponse["details"] = err.Error()
	}
	
	json.NewEncoder(w).Encode(errorResponse)
}

// SendSuccessResponse sends a standardized success response
func (h *HTTPUtils) SendSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := map[string]interface{}{
		"success":   true,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	if data != nil {
		response["data"] = data
	}
	
	h.SendJSONResponse(w, response, http.StatusOK)
}

// ParseIntParam parses an integer parameter from URL path or query
func (h *HTTPUtils) ParseIntParam(value string, paramName string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("%s is required", paramName)
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: must be an integer", paramName)
	}
	
	if intValue <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", paramName)
	}
	
	return intValue, nil
}

// ParseFloatParam parses a float parameter from URL query
func (h *HTTPUtils) ParseFloatParam(value string, paramName string) (float64, error) {
	if value == "" {
		return 0, fmt.Errorf("%s is required", paramName)
	}
	
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: must be a number", paramName)
	}
	
	return floatValue, nil
}

// ParseDateParam parses a date parameter in YYYY-MM-DD format
func (h *HTTPUtils) ParseDateParam(value string, paramName string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%s is required", paramName)
	}
	
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s format: expected YYYY-MM-DD", paramName)
	}
	
	return date, nil
}

// ParseDateTimeParam parses a datetime parameter in RFC3339 format
func (h *HTTPUtils) ParseDateTimeParam(value string, paramName string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%s is required", paramName)
	}
	
	datetime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s format: expected RFC3339 format", paramName)
	}
	
	return datetime, nil
}

// ValidateRequiredParam checks if a required parameter is present
func (h *HTTPUtils) ValidateRequiredParam(value string, paramName string) error {
	if value == "" {
		return fmt.Errorf("%s is required", paramName)
	}
	return nil
}

// ValidateStringLength validates string parameter length
func (h *HTTPUtils) ValidateStringLength(value string, paramName string, minLength, maxLength int) error {
	if len(value) < minLength {
		return fmt.Errorf("%s must be at least %d characters long", paramName, minLength)
	}
	
	if maxLength > 0 && len(value) > maxLength {
		return fmt.Errorf("%s must be at most %d characters long", paramName, maxLength)
	}
	
	return nil
}

// ValidateDateRange validates that start date is before end date
func (h *HTTPUtils) ValidateDateRange(start, end time.Time) error {
	if start.IsZero() || end.IsZero() {
		return fmt.Errorf("both start and end dates are required")
	}
	
	if start.After(end) {
		return fmt.Errorf("start date cannot be after end date")
	}
	
	return nil
}

// ExtractPaginationParams extracts pagination parameters from query string
func (h *HTTPUtils) ExtractPaginationParams(r *http.Request) (page, pageSize int, err error) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	
	// Default values
	page = 1
	pageSize = 20
	
	if pageStr != "" {
		if page, err = strconv.Atoi(pageStr); err != nil {
			return 0, 0, fmt.Errorf("invalid page parameter")
		}
		if page < 1 {
			page = 1
		}
	}
	
	if pageSizeStr != "" {
		if pageSize, err = strconv.Atoi(pageSizeStr); err != nil {
			return 0, 0, fmt.Errorf("invalid page_size parameter")
		}
		if pageSize < 1 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100 // Max page size limit
		}
	}
	
	return page, pageSize, nil
}

// ExtractSortParams extracts sorting parameters from query string
func (h *HTTPUtils) ExtractSortParams(r *http.Request) (sortBy, sortDir string) {
	sortBy = r.URL.Query().Get("sort_by")
	sortDir = r.URL.Query().Get("sort_dir")
	
	// Default values
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc"
	}
	
	return sortBy, sortDir
}

// SetCacheHeaders sets appropriate cache headers for the response
func (h *HTTPUtils) SetCacheHeaders(w http.ResponseWriter, maxAge int) {
	if maxAge > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}
}

// SetSecurityHeaders sets common security headers
func (h *HTTPUtils) SetSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

// IsValidEnum checks if a value is in the allowed enum values
func (h *HTTPUtils) IsValidEnum(value string, allowedValues []string) bool {
	for _, allowed := range allowedValues {
		if value == allowed {
			return true
		}
	}
	return false
}

// BuildLocationHeader builds a Location header for created resources
func (h *HTTPUtils) BuildLocationHeader(r *http.Request, resourcePath string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, resourcePath)
}

// Global utility instance
var Utils = NewHTTPUtils()