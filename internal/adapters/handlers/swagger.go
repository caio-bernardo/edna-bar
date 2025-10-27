package handlers

import (
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// SwaggerHandler provides Swagger UI documentation endpoints
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// RegisterRoutes registers swagger documentation routes
func (h *SwaggerHandler) RegisterRoutes(mux *http.ServeMux) {
	// Serve Swagger UI at /docs/
	mux.Handle("GET /docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.yaml"), // Path to swagger spec
	))
	
	// Serve the swagger.yaml file
	mux.HandleFunc("GET /docs/swagger.yaml", h.serveSwaggerSpec)
	
	// Alternative endpoint for swagger UI
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.yaml"),
	))
	
	// Redirect from /docs to /docs/
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	
	// Redirect from /swagger to /swagger/
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
}

// serveSwaggerSpec serves the swagger specification file
func (h *SwaggerHandler) serveSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	// Set appropriate headers
	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// Try to serve the swagger.yaml file from the docs directory first
	if _, err := os.Stat("./docs/swagger.yaml"); err == nil {
		http.ServeFile(w, r, "./docs/swagger.yaml")
		return
	}
	
	// Fallback to embedded specification
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EmbeddedSwaggerSpec))
}

// GetSwaggerUIURL returns the URL where Swagger UI is served
func (h *SwaggerHandler) GetSwaggerUIURL() string {
	return "/docs/"
}

// GetSwaggerSpecURL returns the URL where the OpenAPI spec is served
func (h *SwaggerHandler) GetSwaggerSpecURL() string {
	return "/docs/swagger.yaml"
}