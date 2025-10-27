package repository

import (
	"edna/internal/domain"
	"edna/internal/infra/database"
)

// RepositoryRegistry provides access to all repository implementations
type RepositoryRegistry struct {
	dbService database.Service
	
	// Repository instances (lazy-loaded)
	editoraRepo     domain.EditoraRepository
	autorRepo       domain.AutorRepository
	livroRepo       domain.LivroRepository
	graficaRepo     domain.GraficaRepository
	particularRepo  domain.ParticularRepository
	contratadaRepo  domain.ContratadaRepository
	contratoRepo    domain.ContratoRepository
	escreveRepo     domain.EscreveRepository
	imprimeRepo     domain.ImprimeRepository
	bookAuthorsRepo domain.BookAuthorsRepository
	printingJobRepo domain.PrintingJobRepository
}

// NewRepositoryRegistry creates a new repository registry
func NewRepositoryRegistry(dbService database.Service) *RepositoryRegistry {
	return &RepositoryRegistry{
		dbService: dbService,
	}
}

// Editora returns the Editora repository instance
func (r *RepositoryRegistry) Editora() domain.EditoraRepository {
	if r.editoraRepo == nil {
		r.editoraRepo = NewPgEditoraRepository(r.dbService)
	}
	return r.editoraRepo
}

// Autor returns the Autor repository instance
func (r *RepositoryRegistry) Autor() domain.AutorRepository {
	if r.autorRepo == nil {
		r.autorRepo = NewPgAutorRepository(r.dbService)
	}
	return r.autorRepo
}

// Livro returns the Livro repository instance
func (r *RepositoryRegistry) Livro() domain.LivroRepository {
	if r.livroRepo == nil {
		r.livroRepo = NewPgLivroRepository(r.dbService)
	}
	return r.livroRepo
}

// Grafica returns the Grafica repository instance
func (r *RepositoryRegistry) Grafica() domain.GraficaRepository {
	if r.graficaRepo == nil {
		r.graficaRepo = NewPgGraficaRepository(r.dbService)
	}
	return r.graficaRepo
}

// Particular returns the Particular repository instance
func (r *RepositoryRegistry) Particular() domain.ParticularRepository {
	if r.particularRepo == nil {
		r.particularRepo = NewPgParticularRepository(r.dbService)
	}
	return r.particularRepo
}

// Contratada returns the Contratada repository instance
func (r *RepositoryRegistry) Contratada() domain.ContratadaRepository {
	if r.contratadaRepo == nil {
		r.contratadaRepo = NewPgContratadaRepository(r.dbService)
	}
	return r.contratadaRepo
}

// Contrato returns the Contrato repository instance
func (r *RepositoryRegistry) Contrato() domain.ContratoRepository {
	if r.contratoRepo == nil {
		r.contratoRepo = NewPgContratoRepository(r.dbService)
	}
	return r.contratoRepo
}

// Escreve returns the Escreve repository instance
func (r *RepositoryRegistry) Escreve() domain.EscreveRepository {
	if r.escreveRepo == nil {
		r.escreveRepo = NewPgEscreveRepository(r.dbService)
	}
	return r.escreveRepo
}

// Imprime returns the Imprime repository instance
func (r *RepositoryRegistry) Imprime() domain.ImprimeRepository {
	if r.imprimeRepo == nil {
		r.imprimeRepo = NewPgImprimeRepository(r.dbService)
	}
	return r.imprimeRepo
}

// BookAuthors returns the BookAuthors repository instance
func (r *RepositoryRegistry) BookAuthors() domain.BookAuthorsRepository {
	if r.bookAuthorsRepo == nil {
		r.bookAuthorsRepo = NewPgBookAuthorsRepository(r.dbService)
	}
	return r.bookAuthorsRepo
}

// PrintingJob returns the PrintingJob repository instance
func (r *RepositoryRegistry) PrintingJob() domain.PrintingJobRepository {
	if r.printingJobRepo == nil {
		r.printingJobRepo = NewPgPrintingJobRepository(r.dbService)
	}
	return r.printingJobRepo
}

// Close closes the database connection
func (r *RepositoryRegistry) Close() error {
	return r.dbService.Close()
}

// Health returns the database health status
func (r *RepositoryRegistry) Health() map[string]string {
	return r.dbService.Health()
}