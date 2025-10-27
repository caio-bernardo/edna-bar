package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"edna/internal/adapters/handlers"
	"edna/internal/adapters/repository"
	"edna/internal/applications"
	"edna/internal/infra/database"
)

type Server struct {
	port               int
	db                 database.Service
	repositoryRegistry *repository.RepositoryRegistry
	applicationService *applications.ApplicationService
	handlerRegistry    *handlers.HandlerRegistry
}

// NewServer creates a new server instance with all components properly wired
func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080 // Default port
		log.Printf("PORT not set or invalid, defaulting to %d", port)
	}

	// Initialize database service
	db := database.New()
	log.Printf("Database service initialized")

	// Initialize repository registry
	repositoryRegistry := repository.NewRepositoryRegistry(db)
	log.Printf("Repository registry initialized")

	// Create application service configuration
	appConfig := applications.ApplicationServiceConfig{
		LivroRepo:       repositoryRegistry.Livro(),
		AutorRepo:       repositoryRegistry.Autor(),
		EditoraRepo:     repositoryRegistry.Editora(),
		GraficaRepo:     repositoryRegistry.Grafica(),
		ParticularRepo:  repositoryRegistry.Particular(),
		ContratadaRepo:  repositoryRegistry.Contratada(),
		ContratoRepo:    repositoryRegistry.Contrato(),
		EscreveRepo:     repositoryRegistry.Escreve(),
		ImprimeRepo:     repositoryRegistry.Imprime(),
		BookAuthorsRepo: repositoryRegistry.BookAuthors(),
		PrintingJobRepo: repositoryRegistry.PrintingJob(),
	}

	// Initialize application service
	applicationService := applications.NewApplicationService(appConfig)
	log.Printf("Application service initialized")

	// Verify application service health
	if err := applicationService.HealthCheck(); err != nil {
		log.Fatalf("Application service health check failed: %v", err)
	}

	// Create handler registry configuration
	handlerConfig := handlers.HandlerConfig{
		ApplicationService: applicationService,
		RequestTimeout:     30 * time.Second,
	}

	// Initialize handler registry
	handlerRegistry := handlers.NewHandlerRegistry(handlerConfig)
	log.Printf("Handler registry initialized")

	// Create server instance
	serverInstance := &Server{
		port:               port,
		db:                 db,
		repositoryRegistry: repositoryRegistry,
		applicationService: applicationService,
		handlerRegistry:    handlerRegistry,
	}

	// Configure HTTP server
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        serverInstance.RegisterRoutes(),
		IdleTimeout:    time.Minute,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Printf("HTTP server configured on port %d", port)
	return httpServer
}



// handleRoot provides basic server information
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := `{
		"service": "Edna Bar Book Printing API",
		"version": "1.0.0",
		"status": "running",
		"timestamp": "` + time.Now().Format(time.RFC3339) + `",
		"endpoints": {
			"api": "/api/",
			"health": "/health",
			"database_health": "/health/db",
			"application_health": "/health/app",
			"server_status": "/status"
		},
		"documentation": {
			"swagger_ui": "/docs/",
			"openapi_spec": "/docs/swagger.yaml",
			"alternative_ui": "/swagger/",
			"api_info": "/api/",
			"available_resources": [
				"/api/livros",
				"/api/autores", 
				"/api/editoras",
				"/api/graficas",
				"/api/contratos",
				"/api/printing-jobs"
			]
		}
	}`
	
	w.Write([]byte(response))
}

// handleDatabaseHealth checks database connection health
func (s *Server) handleDatabaseHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	health := s.db.Health()
	
	if health["status"] == "up" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	// Convert health map to JSON manually for better control
	response := `{
		"database": {
			"status": "` + health["status"] + `",
			"message": "` + health["message"] + `",
			"timestamp": "` + time.Now().Format(time.RFC3339) + `"`
	
	if health["open_connections"] != "" {
		response += `,
			"connections": {
				"open": ` + health["open_connections"] + `,
				"in_use": ` + health["in_use"] + `,
				"idle": ` + health["idle"] + `
			}`
	}
	
	response += `
		}
	}`
	
	w.Write([]byte(response))
}

// handleApplicationHealth checks application service health
func (s *Server) handleApplicationHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	err := s.applicationService.HealthCheck()
	
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		response := `{
			"application": {
				"status": "unhealthy",
				"error": "` + err.Error() + `",
				"timestamp": "` + time.Now().Format(time.RFC3339) + `"
			}
		}`
		w.Write([]byte(response))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	response := `{
		"application": {
			"status": "healthy",
			"message": "All use cases initialized successfully",
			"timestamp": "` + time.Now().Format(time.RFC3339) + `",
			"use_cases": {
				"livro": "initialized",
				"autor": "initialized",
				"editora": "initialized", 
				"grafica": "initialized",
				"contrato": "initialized",
				"imprime": "initialized"
			}
		}
	}`
	
	w.Write([]byte(response))
}

// handleComprehensiveHealth provides a comprehensive health check
func (s *Server) handleComprehensiveHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	dbHealth := s.db.Health()
	appErr := s.applicationService.HealthCheck()
	
	isHealthy := dbHealth["status"] == "up" && appErr == nil
	
	if isHealthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	appStatus := "healthy"
	appMessage := "All systems operational"
	if appErr != nil {
		appStatus = "unhealthy"
		appMessage = appErr.Error()
	}
	
	response := `{
		"service": "Edna Bar Book Printing API",
		"overall_status": "`
	
	if isHealthy {
		response += `healthy`
	} else {
		response += `unhealthy`
	}
	
	response += `",
		"timestamp": "` + time.Now().Format(time.RFC3339) + `",
		"components": {
			"database": {
				"status": "` + dbHealth["status"] + `",
				"message": "` + dbHealth["message"] + `"
			},
			"application": {
				"status": "` + appStatus + `",
				"message": "` + appMessage + `"
			},
			"server": {
				"status": "running",
				"port": ` + strconv.Itoa(s.port) + `
			}
		}
	}`
	
	w.Write([]byte(response))
}

// handleServerStatus provides server runtime information
func (s *Server) handleServerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := `{
		"server": {
			"status": "running",
			"port": ` + strconv.Itoa(s.port) + `,
			"timestamp": "` + time.Now().Format(time.RFC3339) + `"
		},
		"environment": {
			"go_version": "` + fmt.Sprintf("%s", "go1.21+") + `",
			"database_type": "PostgreSQL"
		},
		"features": {
			"cors_enabled": true,
			"request_logging": true,
			"panic_recovery": true,
			"request_timeout": "30s"
		}
	}`
	
	w.Write([]byte(response))
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")
	
	// Close database connection
	if err := s.repositoryRegistry.Close(); err != nil {
		log.Printf("Error closing repository registry: %v", err)
		return err
	}
	
	log.Println("Server shutdown completed")
	return nil
}

// GetApplicationService returns the application service instance
func (s *Server) GetApplicationService() *applications.ApplicationService {
	return s.applicationService
}

// GetRepositoryRegistry returns the repository registry instance
func (s *Server) GetRepositoryRegistry() *repository.RepositoryRegistry {
	return s.repositoryRegistry
}

// GetDatabaseService returns the database service instance
func (s *Server) GetDatabaseService() database.Service {
	return s.db
}