package handlers

import (
	"edna/internal/applications"
	"net/http"
	"time"
)

// HandlerRegistry coordinates all HTTP handlers and provides a unified API
type HandlerRegistry struct {
	livroHandler    *LivroHandler
	autorHandler    *AutorHandler
	editoraHandler  *EditoraHandler
	graficaHandler  *GraficaHandler
	contratoHandler *ContratoHandler
	imprimeHandler  *ImprimeHandler
	swaggerHandler  *SwaggerHandler
}

// HandlerConfig holds the configuration for creating HandlerRegistry
type HandlerConfig struct {
	ApplicationService *applications.ApplicationService
	RequestTimeout     time.Duration
}

// NewHandlerRegistry creates a new HandlerRegistry with all handlers
func NewHandlerRegistry(config HandlerConfig) *HandlerRegistry {
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 30 * time.Second
	}

	return &HandlerRegistry{
		livroHandler:    NewLivroHandler(config.ApplicationService.LivroUsecase),
		autorHandler:    NewAutorHandler(config.ApplicationService.AutorUsecase),
		editoraHandler:  NewEditoraHandler(config.ApplicationService.EditoraUsecase),
		graficaHandler:  NewGraficaHandler(config.ApplicationService.GraficaUsecase),
		contratoHandler: NewContratoHandler(config.ApplicationService.ContratoUsecase),
		imprimeHandler:  NewImprimeHandler(config.ApplicationService.ImprimeUsecase),
		swaggerHandler:  NewSwaggerHandler(),
	}
}

// RegisterRoutes registers all handler routes with the provided ServeMux
func (hr *HandlerRegistry) RegisterRoutes(mux *http.ServeMux) {
	// Apply middleware chain to all routes
	middleware := Chain(
		Recovery(),
		Logger(),
		CORS(),
		ContentType(),
		Timeout(30*time.Second),
	)

	// Create a new ServeMux for API routes
	apiMux := http.NewServeMux()

	// Register all handler routes
	hr.livroHandler.RegisterRoutes(apiMux)
	hr.autorHandler.RegisterRoutes(apiMux)
	hr.editoraHandler.RegisterRoutes(apiMux)
	hr.graficaHandler.RegisterRoutes(apiMux)
	hr.contratoHandler.RegisterRoutes(apiMux)
	hr.imprimeHandler.RegisterRoutes(apiMux)

	// Add health check endpoint
	apiMux.HandleFunc("GET /health", hr.handleHealthCheck)

	// Add API info endpoint
	apiMux.HandleFunc("GET /", hr.handleAPIInfo)

	// Register swagger documentation routes (outside API prefix)
	hr.swaggerHandler.RegisterRoutes(mux)

	// Mount API routes with middleware under /api prefix
	mux.Handle("/api/", http.StripPrefix("/api", middleware(apiMux)))
}

// handleHealthCheck provides a health check endpoint
func (hr *HandlerRegistry) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `","version":"1.0.0"}`))
}

// handleAPIInfo provides API information and available endpoints
func (hr *HandlerRegistry) handleAPIInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	apiInfoJSON := `{
		"name": "Edna Bar Book Printing API",
		"version": "1.0.0",
		"description": "RESTful API for managing book printing operations",
		"endpoints": {
			"books": {
				"base": "/api/livros",
				"operations": ["GET", "POST", "PUT", "DELETE"],
				"description": "Manage books and their authors"
			},
			"authors": {
				"base": "/api/autores", 
				"operations": ["GET", "POST", "PUT", "DELETE"],
				"description": "Manage authors and their books"
			},
			"publishers": {
				"base": "/api/editoras",
				"operations": ["GET", "POST", "PUT", "DELETE"], 
				"description": "Manage publishers and their catalogs"
			},
			"printing-companies": {
				"base": "/api/graficas",
				"operations": ["GET", "POST", "PUT", "DELETE"],
				"description": "Manage printing companies (particular/contracted)"
			},
			"contracts": {
				"base": "/api/contratos",
				"operations": ["GET", "POST", "PUT", "DELETE"],
				"description": "Manage contracts with printing companies"
			},
			"printing-jobs": {
				"base": "/api/printing-jobs",
				"operations": ["GET", "POST", "PUT", "DELETE"],
				"description": "Manage printing jobs and track delivery status"
			}
		},
		"documentation": {
			"swagger_ui": "/docs/",
			"openapi_spec": "/docs/swagger.yaml",
			"alternative_ui": "/swagger/"
		}
	}`
	
	w.Write([]byte(apiInfoJSON))
}

// GetLivroHandler returns the livro handler
func (hr *HandlerRegistry) GetLivroHandler() *LivroHandler {
	return hr.livroHandler
}

// GetAutorHandler returns the autor handler
func (hr *HandlerRegistry) GetAutorHandler() *AutorHandler {
	return hr.autorHandler
}

// GetEditoraHandler returns the editora handler
func (hr *HandlerRegistry) GetEditoraHandler() *EditoraHandler {
	return hr.editoraHandler
}

// GetGraficaHandler returns the grafica handler
func (hr *HandlerRegistry) GetGraficaHandler() *GraficaHandler {
	return hr.graficaHandler
}

// GetContratoHandler returns the contrato handler
func (hr *HandlerRegistry) GetContratoHandler() *ContratoHandler {
	return hr.contratoHandler
}

// GetImprimeHandler returns the imprime handler
func (hr *HandlerRegistry) GetImprimeHandler() *ImprimeHandler {
	return hr.imprimeHandler
}

// GetSwaggerHandler returns the swagger handler
func (hr *HandlerRegistry) GetSwaggerHandler() *SwaggerHandler {
	return hr.swaggerHandler
}