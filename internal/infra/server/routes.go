package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// RegisterRoutes registers all routes and returns the configured handler
func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register API routes through handler registry
	s.handlerRegistry.RegisterRoutes(mux)

	// Add root handler for basic info
	// mux.HandleFunc("GET /", s.handleRoot)

	// Add database health check endpoint
	mux.HandleFunc("GET /health/db", s.handleDatabaseHealth)

	// Add application health check endpoint
	mux.HandleFunc("GET /health/app", s.handleApplicationHealth)

	// Add comprehensive health check endpoint
	mux.HandleFunc("GET /health", s.handleComprehensiveHealth)

	// Add graceful shutdown info endpoint
	mux.HandleFunc("GET /status", s.handleServerStatus)

	// Legacy endpoints for backwards compatibility
	mux.HandleFunc("/legacy/hello", s.HelloWorldHandler)
	mux.HandleFunc("/legacy/health", s.healthHandler)

	log.Printf("Routes registered successfully")
	return mux
}

// Legacy handlers for backwards compatibility
func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"message": "Hello World - Edna Bar Book Printing API",
		"version": "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"api_endpoint": "/api/",
		"note": "This is a legacy endpoint. Use /api/ for the main API.",
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	dbHealth := s.db.Health()

	// Enhanced health response with application status
	healthData := map[string]interface{}{
		"database": dbHealth,
		"application": map[string]interface{}{
			"status": "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		"server": map[string]interface{}{
			"status": "running",
			"port": s.port,
		},
		"note": "This is a legacy endpoint. Use /health for comprehensive health check.",
	}

	// Check application health
	if err := s.applicationService.HealthCheck(); err != nil {
		healthData["application"] = map[string]interface{}{
			"status": "unhealthy",
			"error": err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	resp, err := json.Marshal(healthData)
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
