package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"edna/internal/database"
	"edna/internal/services/aplica_oferta"
	"edna/internal/services/cliente"
	"edna/internal/services/fornecedor"
	"edna/internal/services/funcionario"
	"edna/internal/services/item_oferta"
	"edna/internal/services/item_venda"
	"edna/internal/services/lote"
	"edna/internal/services/oferta"
	"edna/internal/services/produto"
	"edna/internal/services/relatorio"
	"edna/internal/services/venda"
)

type Server struct {
	port int

	db                database.Service
	fornecedorStore   *fornecedor.Store
	produtoStore      *produto.Store
	clienteStore      *cliente.Store
	loteStore         *lote.Store
	ofertaStore       *oferta.Store
	vendaStore        *venda.Store
	funcionarioStore  *funcionario.Store
	relatorioStore    *relatorio.Store
	itemOfertaStore   *item_oferta.Store
	itemVendaStore    *item_venda.Store
	aplicaOfertaStore *aplica_oferta.Store
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	if port == 0 {
		port = 8080
	}
	db := database.New()
	NewServer := &Server{
		port: port,

		db:                db,
		fornecedorStore:   fornecedor.NewStore(db.Conn()),
		produtoStore:      produto.NewStore(db.Conn()),
		clienteStore:      cliente.NewStore(db.Conn()),
		loteStore:         lote.NewStore(db.Conn()),
		ofertaStore:       oferta.NewStore(db.Conn()),
		vendaStore:        venda.NewStore(db.Conn()),
		itemVendaStore:    item_venda.NewStore(db.Conn()),
		itemOfertaStore:   item_oferta.NewStore(db.Conn()),
		aplicaOfertaStore: aplica_oferta.NewStore(db.Conn()),
		funcionarioStore:  funcionario.NewStore(db.Conn()),
		relatorioStore:    relatorio.NewStore(db.Conn()),
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
