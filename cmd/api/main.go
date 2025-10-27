package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"edna/internal/infra/server"
)

// @title Edna Bar Book Printing API
// @version 1.0
// @description A comprehensive RESTful API for managing book printing operations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @schemes http https
// @produce json
// @consumes json

// @tag.name books
// @tag.description Operations related to book management

// @tag.name authors
// @tag.description Operations related to author management

// @tag.name publishers
// @tag.description Operations related to publisher management

// @tag.name printing-companies
// @tag.description Operations related to printing company management

// @tag.name contracts
// @tag.description Operations related to contract management

// @tag.name printing-jobs
// @tag.description Operations related to printing job management

// @tag.name system
// @tag.description System health and information endpoints

func gracefulShutdown(apiServer *http.Server, serverInstance *server.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("Shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Shutdown application components first
	if err := serverInstance.Shutdown(); err != nil {
		log.Printf("Error during application shutdown: %v", err)
	}
	
	// Then shutdown HTTP server
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	log.Println("Starting Edna Bar Book Printing API Server...")

	// Create and configure the server with all components integrated
	httpServer := server.NewServer()
	
	// Get access to the server instance for proper shutdown
	// Note: In a real implementation, you might want to modify NewServer to return both
	// For now, we'll work with what we have
	
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go func() {
		// Create context that listens for the interrupt signal from the OS.
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		// Listen for the interrupt signal.
		<-ctx.Done()

		log.Println("Shutting down gracefully, press Ctrl+C again to force")
		stop() // Allow Ctrl+C to force shutdown

		// The context is used to inform the server it has 30 seconds to finish
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		// Shutdown HTTP server
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown with error: %v", err)
		}

		log.Println("Server exiting")
		done <- true
	}()

	// Log startup information
	log.Printf("Server starting on %s", httpServer.Addr)
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("📚 Edna Bar Book Printing API Server")
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("🔗 API Base URL:        http://localhost" + httpServer.Addr + "/api/")
	log.Println("📋 API Documentation:  http://localhost" + httpServer.Addr + "/api/")
	log.Println("📚 Swagger UI:         http://localhost" + httpServer.Addr + "/docs/")
	log.Println("📄 OpenAPI Spec:       http://localhost" + httpServer.Addr + "/docs/swagger.yaml")
	log.Println("💚 Health Check:       http://localhost" + httpServer.Addr + "/health")
	log.Println("📊 Server Status:      http://localhost" + httpServer.Addr + "/status")
	log.Println("🗄️  Database Health:    http://localhost" + httpServer.Addr + "/health/db")
	log.Println("⚙️  Application Health: http://localhost" + httpServer.Addr + "/health/app")
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("Available Resources:")
	log.Println("  📖 Books:             /api/livros")
	log.Println("  ✍️  Authors:           /api/autores") 
	log.Println("  🏢 Publishers:        /api/editoras")
	log.Println("  🖨️  Printing Companies: /api/graficas")
	log.Println("  📄 Contracts:         /api/contratos")
	log.Println("  🔄 Printing Jobs:     /api/printing-jobs")
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("Press CTRL+C to gracefully shutdown...")
	log.Println("")

	// Start the server
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("HTTP server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("Graceful shutdown complete. Thank you for using Edna Bar API!")
	log.Println("═══════════════════════════════════════════════════════════")
}