# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build -d 2>/dev/null; then \
		echo "Running Docker Compose" ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		echo "Shutdown Docker Compose" ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Generate swagger documentation
swagger:
	@echo "Generating swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal; \
		echo "Swagger documentation generated in docs/"; \
	else \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal; \
		echo "Swagger documentation generated in docs/"; \
	fi

# Serve swagger UI (requires server to be running)
swagger-ui:
	@echo "Swagger UI available at:"
	@echo "  http://localhost:8080/docs/"
	@echo "  http://localhost:8080/swagger/"
	@echo "OpenAPI spec available at:"
	@echo "  http://localhost:8080/docs/swagger.yaml"

# Build and run with swagger
run-with-docs: swagger run

.PHONY: all build run test clean watch docker-run docker-down itest swagger swagger-ui run-with-docs
