package server

import (
	"edna/internal/services/cliente"
	"edna/internal/services/fornecedor"
	"edna/internal/services/lote"
	"edna/internal/services/oferta"
	"edna/internal/services/produto"
	"encoding/json"
	"log"
	"net/http"

	_ "edna/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {

	v1 := http.NewServeMux()
	mux := http.NewServeMux()

	fornecedorHandler := fornecedor.NewHandler(s.fornecedorStore)
	produtoHandler := produto.NewHandler(s.produtoStore)
	clienteHandler := cliente.NewHandler(s.clienteStore)
	loteHandler := lote.NewHandler(s.loteStore)
	ofertaHandler := oferta.NewHandler(s.ofertaStore)

	mux.HandleFunc("/health", s.healthHandler)
	fornecedorHandler.RegisterRoutes(mux)
	produtoHandler.RegisterRoutes(mux)
	clienteHandler.RegisterRoutes(mux)
	loteHandler.RegisterRoutes(mux)
	ofertaHandler.RegisterRoutes(mux)

	// Register routes
	v1.HandleFunc("/", s.trailingSlashHandler)
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))
	v1.Handle("/swagger/", httpSwagger.Handler())
	// Wrap the mux with CORS middleware
	return s.logMiddleware(s.corsMiddleware(v1))
}

// @Summary Unmatched path handler
// @Description Returns a 404 JSON response for unmatched routes.
// @Tags Server
// @Produce json
// @Success 404 {object} map[string]string
// @Router / [get]
func (s *Server) trailingSlashHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	resp := map[string]string{"message": "Unmatched path, please check your url path and try again."}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// @Summary Check health of the system
// @Description Returns the health status of the application and dependencies.
// @Tags Server
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /health [get]
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
