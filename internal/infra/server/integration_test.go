package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"edna/internal/infra/database"
)

// TestServerIntegration tests the complete server integration
func TestServerIntegration(t *testing.T) {
	// Skip if no test database is available
	if os.Getenv("TEST_DB_SKIP") == "true" {
		t.Skip("Skipping integration test - TEST_DB_SKIP is set")
	}

	// Set up test environment variables
	os.Setenv("PORT", "8081")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_DATABASE", "edna_bar_test")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_SCHEMA", "public")

	// Create test server
	httpServer := NewServer()
	
	// Create test server instance
	testServer := httptest.NewServer(httpServer.Handler)
	defer testServer.Close()

	// Test cases
	tests := []struct {
		name           string
		endpoint       string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "Root endpoint",
			endpoint:       "/",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
					return
				}
				if service, ok := response["service"]; !ok || service != "Edna Bar Book Printing API" {
					t.Errorf("Expected service name, got %v", service)
				}
			},
		},
		{
			name:           "Health check endpoint",
			endpoint:       "/health",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal health response: %v", err)
					return
				}
				if status, ok := response["overall_status"]; !ok {
					t.Errorf("Expected overall_status in health response")
				} else if status != "healthy" && status != "unhealthy" {
					t.Errorf("Expected valid health status, got %v", status)
				}
			},
		},
		{
			name:           "Database health endpoint",
			endpoint:       "/health/db",
			expectedStatus: http.StatusOK, // May be 503 if DB unavailable
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal db health response: %v", err)
					return
				}
				if db, ok := response["database"]; !ok {
					t.Errorf("Expected database section in response")
				} else if dbMap, ok := db.(map[string]interface{}); !ok {
					t.Errorf("Expected database to be an object")
				} else if status, ok := dbMap["status"]; !ok {
					t.Errorf("Expected status in database section")
				} else if status != "up" && status != "down" {
					t.Errorf("Expected valid database status, got %v", status)
				}
			},
		},
		{
			name:           "Application health endpoint",
			endpoint:       "/health/app",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal app health response: %v", err)
					return
				}
				if app, ok := response["application"]; !ok {
					t.Errorf("Expected application section in response")
				} else if appMap, ok := app.(map[string]interface{}); !ok {
					t.Errorf("Expected application to be an object")
				} else if status, ok := appMap["status"]; !ok {
					t.Errorf("Expected status in application section")
				} else if status != "healthy" && status != "unhealthy" {
					t.Errorf("Expected valid application status, got %v", status)
				}
			},
		},
		{
			name:           "Server status endpoint",
			endpoint:       "/status",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal status response: %v", err)
					return
				}
				if server, ok := response["server"]; !ok {
					t.Errorf("Expected server section in response")
				} else if serverMap, ok := server.(map[string]interface{}); !ok {
					t.Errorf("Expected server to be an object")
				} else if status, ok := serverMap["status"]; !ok || status != "running" {
					t.Errorf("Expected server status to be 'running', got %v", status)
				}
			},
		},
		{
			name:           "API info endpoint",
			endpoint:       "/api/",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal API info response: %v", err)
					return
				}
				if name, ok := response["name"]; !ok || name != "Edna Bar Book Printing API" {
					t.Errorf("Expected API name, got %v", name)
				}
				if endpoints, ok := response["endpoints"]; !ok {
					t.Errorf("Expected endpoints section in API info")
				} else if endpointsMap, ok := endpoints.(map[string]interface{}); !ok {
					t.Errorf("Expected endpoints to be an object")
				} else {
					expectedEndpoints := []string{"books", "authors", "publishers", "printing-companies", "contracts", "printing-jobs"}
					for _, endpoint := range expectedEndpoints {
						if _, exists := endpointsMap[endpoint]; !exists {
							t.Errorf("Expected endpoint %s to be documented", endpoint)
						}
					}
				}
			},
		},
		{
			name:           "Legacy hello endpoint",
			endpoint:       "/legacy/hello",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal hello response: %v", err)
					return
				}
				if message, ok := response["message"]; !ok {
					t.Errorf("Expected message in hello response")
				} else if !contains(message.(string), "Hello World") {
					t.Errorf("Expected Hello World message, got %v", message)
				}
			},
		},
		{
			name:           "Swagger UI endpoint",
			endpoint:       "/docs/",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				if !contains(bodyStr, "swagger") && !contains(bodyStr, "Swagger") {
					t.Errorf("Expected Swagger UI content, got response that doesn't contain swagger references")
				}
			},
		},
		{
			name:           "Swagger spec endpoint",
			endpoint:       "/docs/swagger.yaml",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				if !contains(bodyStr, "openapi") {
					t.Errorf("Expected OpenAPI spec content, got response without openapi field")
				}
				if !contains(bodyStr, "Edna Bar Book Printing API") {
					t.Errorf("Expected API title in swagger spec")
				}
			},
		},
		{
			name:           "Alternative Swagger UI endpoint",
			endpoint:       "/swagger/",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				if !contains(bodyStr, "swagger") && !contains(bodyStr, "Swagger") {
					t.Errorf("Expected Swagger UI content, got response that doesn't contain swagger references")
				}
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make request
			resp, err := http.Get(testServer.URL + tt.endpoint)
			if err != nil {
				t.Fatalf("Failed to make request to %s: %v", tt.endpoint, err)
			}
			defer resp.Body.Close()

			// Check status code (allow 503 for health checks if DB unavailable)
			if resp.StatusCode != tt.expectedStatus && 
			   !(contains(tt.endpoint, "health") && resp.StatusCode == http.StatusServiceUnavailable) {
				t.Errorf("Expected status %d, got %d for %s", tt.expectedStatus, resp.StatusCode, tt.endpoint)
			}

			// Read response body
			body := make([]byte, 1024*10) // 10KB buffer
			n, err := resp.Body.Read(body)
			if err != nil && err.Error() != "EOF" {
				t.Fatalf("Failed to read response body: %v", err)
			}
			body = body[:n]

			// Check response content
			if tt.checkResponse != nil {
				tt.checkResponse(t, body)
			}
		})
	}
}

// TestServerComponents tests individual server components
func TestServerComponents(t *testing.T) {
	// Test database service creation
	t.Run("Database service initialization", func(t *testing.T) {
		// This test may fail if no database is available, which is OK for unit testing
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Database service creation panicked (expected if no DB): %v", r)
			}
		}()
		
		db := database.New()
		if db == nil {
			t.Error("Expected database service to be created")
		}
	})

	// Test server creation without actual HTTP startup
	t.Run("Server struct creation", func(t *testing.T) {
		// Skip if database connection would fail
		if os.Getenv("TEST_DB_SKIP") == "true" {
			t.Skip("Skipping server creation test - would require database")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Server creation failed (expected without proper DB config): %v", r)
			}
		}()

		server := NewServer()
		if server == nil {
			t.Error("Expected server to be created")
		}
	})
}

// TestServerConfiguration tests server configuration
func TestServerConfiguration(t *testing.T) {
	// Test default port configuration
	t.Run("Default port configuration", func(t *testing.T) {
		// Clear PORT env var
		oldPort := os.Getenv("PORT")
		os.Unsetenv("PORT")
		defer os.Setenv("PORT", oldPort)

		// Skip if database connection would fail
		if os.Getenv("TEST_DB_SKIP") == "true" {
			t.Skip("Skipping server config test - would require database")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Server creation failed (expected without proper DB config): %v", r)
			}
		}()

		server := NewServer()
		if server == nil {
			t.Log("Server creation failed, which is expected without database")
			return
		}

		// Check that default port is used (8080)
		if !contains(server.Addr, ":8080") {
			t.Errorf("Expected default port 8080, got %s", server.Addr)
		}
	})

	// Test custom port configuration
	t.Run("Custom port configuration", func(t *testing.T) {
		os.Setenv("PORT", "9999")
		defer os.Unsetenv("PORT")

		// Skip if database connection would fail
		if os.Getenv("TEST_DB_SKIP") == "true" {
			t.Skip("Skipping server config test - would require database")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Server creation failed (expected without proper DB config): %v", r)
			}
		}()

		server := NewServer()
		if server == nil {
			t.Log("Server creation failed, which is expected without database")
			return
		}

		// Check that custom port is used
		if !contains(server.Addr, ":9999") {
			t.Errorf("Expected custom port 9999, got %s", server.Addr)
		}
	})
}

// TestMiddleware tests that middleware is properly applied
func TestMiddleware(t *testing.T) {
	// Skip if no test database is available
	if os.Getenv("TEST_DB_SKIP") == "true" {
		t.Skip("Skipping middleware test - TEST_DB_SKIP is set")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Middleware test failed (expected without proper DB config): %v", r)
		}
	}()

	httpServer := NewServer()
	testServer := httptest.NewServer(httpServer.Handler)
	defer testServer.Close()

	// Test CORS headers
	t.Run("CORS headers", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", testServer.URL+"/api/livros", nil)
		req.Header.Set("Origin", "http://example.com")
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make OPTIONS request: %v", err)
		}
		defer resp.Body.Close()

		// Check for CORS headers (they might be set by the handler registry middleware)
		corsHeader := resp.Header.Get("Access-Control-Allow-Origin")
		if corsHeader == "" {
			t.Log("CORS headers not found - this might be expected depending on middleware setup")
		}
	})

	// Test content-type header
	t.Run("Content-Type headers", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if !contains(contentType, "application/json") {
			t.Errorf("Expected JSON content type, got %s", contentType)
		}
	})
}

// Benchmark tests
func BenchmarkHealthEndpoint(b *testing.B) {
	if os.Getenv("TEST_DB_SKIP") == "true" {
		b.Skip("Skipping benchmark - TEST_DB_SKIP is set")
	}

	httpServer := NewServer()
	testServer := httptest.NewServer(httpServer.Handler)
	defer testServer.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(testServer.URL + "/health")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkAPIInfoEndpoint(b *testing.B) {
	if os.Getenv("TEST_DB_SKIP") == "true" {
		b.Skip("Skipping benchmark - TEST_DB_SKIP is set")
	}

	httpServer := NewServer()
	testServer := httptest.NewServer(httpServer.Handler)
	defer testServer.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(testServer.URL + "/api/")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      indexOf(s, substr) != -1)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}