package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// MiddlewareFunc represents a middleware function
type MiddlewareFunc func(http.Handler) http.Handler

// Chain chains multiple middleware functions
func Chain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// Logger middleware logs HTTP requests
func Logger() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create a response writer wrapper to capture status code
			wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(wrapper, r)
			
			duration := time.Since(start)
			log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapper.statusCode, duration)
		})
	}
}

// Recovery middleware recovers from panics
func Recovery() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
					
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					
					errorResponse := map[string]interface{}{
						"error": "Internal server error",
						"code":  http.StatusInternalServerError,
					}
					
					json.NewEncoder(w).Encode(errorResponse)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// ContentType middleware sets content type to JSON for API responses
func ContentType() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set JSON content type for API responses
			if strings.HasPrefix(r.URL.Path, "/api/") || 
			   strings.Contains(r.Header.Get("Accept"), "application/json") {
				w.Header().Set("Content-Type", "application/json")
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout middleware adds request timeout
func Timeout(duration time.Duration) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}