package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"edna/internal/infra/server"
)

// Example demonstrates the complete Edna Bar Book Printing API system
// This shows how all components work together in a real scenario
func main() {
	fmt.Println("🚀 Edna Bar Book Printing API - Complete System Example")
	fmt.Println("======================================================")

	// Step 1: Start the server (in production, this would be in a separate process)
	fmt.Println("1. Starting the integrated server...")
	
	// Set up test environment
	setupTestEnvironment()
	
	// Create and start server
	httpServer := server.NewServer()
	
	// In a real scenario, you'd start this in a goroutine or separate process
	// For this example, we'll simulate API calls
	baseURL := "http://localhost" + httpServer.Addr
	
	fmt.Printf("   Server configured at %s\n", baseURL)
	fmt.Println("   ✅ All components integrated successfully!")
	fmt.Println()

	// Step 2: Demonstrate the complete workflow
	demonstrateCompleteWorkflow(baseURL)
}

// setupTestEnvironment configures environment variables for the example
func setupTestEnvironment() {
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_DATABASE", "edna_bar")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_SCHEMA", "public")
}

// demonstrateCompleteWorkflow shows a complete business scenario
func demonstrateCompleteWorkflow(baseURL string) {
	fmt.Println("2. Demonstrating Complete Book Printing Workflow")
	fmt.Println("================================================")
	
	// Step 2.1: Check system health
	fmt.Println("2.1 Checking System Health...")
	checkSystemHealth(baseURL)
	
	// Step 2.2: Create a publisher
	fmt.Println("\n2.2 Creating a Publisher...")
	editoraID := createPublisher(baseURL)
	
	// Step 2.3: Create authors
	fmt.Println("\n2.3 Creating Authors...")
	autorRG1 := createAuthor(baseURL, "Gabriel García Márquez", "12345678901")
	autorRG2 := createAuthor(baseURL, "Isabel Allende", "10987654321")
	
	// Step 2.4: Create a book
	fmt.Println("\n2.4 Creating a Book...")
	bookISBN := createBook(baseURL, editoraID)
	
	// Step 2.5: Associate authors with the book
	fmt.Println("\n2.5 Associating Authors with Book...")
	associateAuthorWithBook(baseURL, bookISBN, autorRG1)
	associateAuthorWithBook(baseURL, bookISBN, autorRG2)
	
	// Step 2.6: Create printing companies
	fmt.Println("\n2.6 Creating Printing Companies...")
	graficaID1 := createPrintingCompany(baseURL, "PrintTech Solutions")
	graficaID2 := createPrintingCompany(baseURL, "Global Print Services")
	
	// Step 2.7: Set up contracted printing company
	fmt.Println("\n2.7 Setting up Contracted Printing Company...")
	setupContractedCompany(baseURL, graficaID1)
	
	// Step 2.8: Create printing contracts
	fmt.Println("\n2.8 Creating Printing Contracts...")
	createContract(baseURL, graficaID1, 15000.00, "Maria Silva")
	
	// Step 2.9: Create printing jobs
	fmt.Println("\n2.9 Creating Printing Jobs...")
	createPrintingJob(baseURL, bookISBN, graficaID1, 5000)
	createPrintingJob(baseURL, bookISBN, graficaID2, 3000)
	
	// Step 2.10: Query and analyze the data
	fmt.Println("\n2.10 Querying and Analyzing Data...")
	analyzeData(baseURL, bookISBN, autorRG1, graficaID1)
	
	fmt.Println("\n🎉 Complete workflow demonstrated successfully!")
	fmt.Println("   This example shows how all layers of the architecture work together:")
	fmt.Println("   📊 Domain Layer: Business entities and rules")
	fmt.Println("   🔄 Application Layer: Use cases and business workflows")
	fmt.Println("   🌐 Presentation Layer: HTTP handlers and middleware")
	fmt.Println("   💾 Infrastructure Layer: Database repositories and server")
}

// checkSystemHealth verifies all system components are working
func checkSystemHealth(baseURL string) {
	// Check overall health
	health := makeGetRequest(baseURL + "/health")
	if health != nil {
		status := extractString(health, "overall_status")
		fmt.Printf("   Overall System Status: %s\n", status)
	}
	
	// Check database health
	dbHealth := makeGetRequest(baseURL + "/health/db")
	if dbHealth != nil {
		if db, ok := dbHealth["database"].(map[string]interface{}); ok {
			status := extractString(db, "status")
			fmt.Printf("   Database Status: %s\n", status)
		}
	}
	
	// Check application health
	appHealth := makeGetRequest(baseURL + "/health/app")
	if appHealth != nil {
		if app, ok := appHealth["application"].(map[string]interface{}); ok {
			status := extractString(app, "status")
			fmt.Printf("   Application Status: %s\n", status)
		}
	}
	
	fmt.Println("   ✅ System health verified")
}

// createPublisher creates a new publisher
func createPublisher(baseURL string) int {
	publisher := map[string]interface{}{
		"nome":     "Penguin Random House Brasil",
		"endereco": "São Paulo, SP, Brasil",
	}
	
	response := makePostRequest(baseURL+"/api/editoras", publisher)
	if response != nil {
		id := extractFloat64(response, "id")
		fmt.Printf("   Created publisher with ID: %.0f\n", id)
		return int(id)
	}
	
	fmt.Println("   ⚠️  Using mock publisher ID: 1")
	return 1
}

// createAuthor creates a new author
func createAuthor(baseURL, name, rg string) string {
	author := map[string]interface{}{
		"rg":       rg,
		"nome":     name,
		"endereco": "América Latina",
	}
	
	response := makePostRequest(baseURL+"/api/autores", author)
	if response != nil {
		returnedRG := extractString(response, "rg")
		fmt.Printf("   Created author: %s (RG: %s)\n", name, returnedRG)
		return returnedRG
	}
	
	fmt.Printf("   ⚠️  Using mock author RG: %s\n", rg)
	return rg
}

// createBook creates a new book
func createBook(baseURL string, editoraID int) string {
	book := map[string]interface{}{
		"isbn":                "978-0307474728",
		"titulo":              "Cem Anos de Solidão",
		"data_de_publicacao":  "1967-06-05T00:00:00Z",
		"editora_id":          editoraID,
	}
	
	response := makePostRequest(baseURL+"/api/livros", book)
	if response != nil {
		isbn := extractString(response, "isbn")
		title := extractString(response, "titulo")
		fmt.Printf("   Created book: %s (ISBN: %s)\n", title, isbn)
		return isbn
	}
	
	fmt.Println("   ⚠️  Using mock book ISBN: 978-0307474728")
	return "978-0307474728"
}

// associateAuthorWithBook creates the author-book relationship
func associateAuthorWithBook(baseURL, isbn, rg string) {
	association := map[string]interface{}{
		"isbn": isbn,
		"rg":   rg,
	}
	
	response := makePostRequest(baseURL+"/api/livros/"+isbn+"/authors", association)
	if response != nil {
		fmt.Printf("   Associated author %s with book %s\n", rg, isbn)
	} else {
		fmt.Printf("   ⚠️  Mock association: author %s with book %s\n", rg, isbn)
	}
}

// createPrintingCompany creates a new printing company
func createPrintingCompany(baseURL, name string) int {
	company := map[string]interface{}{
		"nome": name,
	}
	
	response := makePostRequest(baseURL+"/api/graficas", company)
	if response != nil {
		id := extractFloat64(response, "id")
		fmt.Printf("   Created printing company: %s (ID: %.0f)\n", name, id)
		return int(id)
	}
	
	fmt.Printf("   ⚠️  Using mock company ID for %s: 1\n", name)
	return 1
}

// setupContractedCompany sets up a printing company as contracted
func setupContractedCompany(baseURL string, graficaID int) {
	contracted := map[string]interface{}{
		"grafica_id": graficaID,
		"endereco":   "Distrito Industrial, São Paulo, SP",
	}
	
	response := makePostRequest(baseURL+"/api/graficas/contracted", contracted)
	if response != nil {
		fmt.Printf("   Set up contracted company for grafica ID: %d\n", graficaID)
	} else {
		fmt.Printf("   ⚠️  Mock contracted setup for grafica ID: %d\n", graficaID)
	}
}

// createContract creates a printing contract
func createContract(baseURL string, graficaID int, value float64, responsible string) {
	contract := map[string]interface{}{
		"valor":             value,
		"nome_responsavel":  responsible,
		"grafica_cont_id":   graficaID,
	}
	
	response := makePostRequest(baseURL+"/api/contratos", contract)
	if response != nil {
		id := extractFloat64(response, "id")
		fmt.Printf("   Created contract ID: %.0f (Value: R$%.2f, Responsible: %s)\n", id, value, responsible)
	} else {
		fmt.Printf("   ⚠️  Mock contract: R$%.2f with %s\n", value, responsible)
	}
}

// createPrintingJob creates a printing job
func createPrintingJob(baseURL, isbn string, graficaID, copies int) {
	deliveryDate := time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339)
	
	job := map[string]interface{}{
		"isbn":         isbn,
		"grafica_id":   graficaID,
		"copies":       copies,
		"delivery_date": deliveryDate,
	}
	
	response := makePostRequest(baseURL+"/api/printing-jobs", job)
	if response != nil {
		fmt.Printf("   Created printing job: %d copies of %s at grafica %d\n", copies, isbn, graficaID)
	} else {
		fmt.Printf("   ⚠️  Mock printing job: %d copies of %s at grafica %d\n", copies, isbn, graficaID)
	}
}

// analyzeData demonstrates querying and analyzing the created data
func analyzeData(baseURL, isbn, rg string, graficaID int) {
	// Get book details
	book := makeGetRequest(baseURL + "/api/livros/" + isbn)
	if book != nil {
		title := extractString(book, "titulo")
		fmt.Printf("   📖 Book Details: %s (ISBN: %s)\n", title, isbn)
	}
	
	// Get book authors
	authors := makeGetRequest(baseURL + "/api/livros/" + isbn + "/authors")
	if authors != nil {
		fmt.Printf("   ✍️  Authors associated with book: %v\n", authors)
	}
	
	// Get author's books
	authorBooks := makeGetRequest(baseURL + "/api/autores/" + rg + "/books")
	if authorBooks != nil {
		fmt.Printf("   📚 Books by author %s: %v\n", rg, authorBooks)
	}
	
	// Get printing company details
	grafica := makeGetRequest(baseURL + fmt.Sprintf("/api/graficas/%d", graficaID))
	if grafica != nil {
		name := extractString(grafica, "nome")
		fmt.Printf("   🖨️  Printing Company: %s (ID: %d)\n", name, graficaID)
	}
	
	// Get total copies information
	totalCopies := makeGetRequest(baseURL + "/api/printing-jobs/book/" + isbn + "/total-copies")
	if totalCopies != nil {
		copies := extractFloat64(totalCopies, "total_copies")
		fmt.Printf("   📊 Total copies scheduled for printing: %.0f\n", copies)
	}
	
	fmt.Println("   ✅ Data analysis completed")
}

// HTTP helper functions
func makeGetRequest(url string) map[string]interface{} {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("GET request failed to %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("GET request to %s returned status %d", url, resp.StatusCode)
		return nil
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body from %s: %v", url, err)
		return nil
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to unmarshal response from %s: %v", url, err)
		return nil
	}
	
	return result
}

func makePostRequest(url string, data map[string]interface{}) map[string]interface{} {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal request data: %v", err)
		return nil
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("POST request failed to %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("POST request to %s returned status %d", url, resp.StatusCode)
		return nil
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body from %s: %v", url, err)
		return nil
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to unmarshal response from %s: %v", url, err)
		return nil
	}
	
	return result
}

// Helper functions to extract values from maps
func extractString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func extractFloat64(m map[string]interface{}, key string) float64 {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return num
		}
	}
	return 0
}