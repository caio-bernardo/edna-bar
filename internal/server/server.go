package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"edna/internal/database"
	"edna/internal/services/cliente"
	"edna/internal/services/fornecedor"
	"edna/internal/services/lote"
	"edna/internal/services/oferta"
	"edna/internal/services/produto"
)

type Server struct {
	port int

	db              database.Service
	fornecedorStore *fornecedor.Store
	produtoStore    *produto.Store
	clienteStore    *cliente.Store
	loteStore       *lote.Store
	ofertaStore     *oferta.Store
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	NewServer := &Server{
		port: port,

		db:              db,
		fornecedorStore: fornecedor.NewStore(db.Conn()),
		produtoStore:    produto.NewStore(db.Conn()),
		clienteStore:    cliente.NewStore(db.Conn()),
		loteStore:       lote.NewStore(db.Conn()),
		ofertaStore:     oferta.NewStore(db.Conn()),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
