package applications

import (
	"edna/internal/domain"
)

// ApplicationService aggregates all use cases and provides a single entry point
// for the application layer
type ApplicationService struct {
	LivroUsecase    LivroUsecase
	AutorUsecase    AutorUsecase
	EditoraUsecase  EditoraUsecase
	GraficaUsecase  GraficaUsecase
	ContratoUsecase ContratoUsecase
	ImprimeUsecase  ImprimeUsecase
}

// ApplicationServiceConfig holds the configuration for creating ApplicationService
type ApplicationServiceConfig struct {
	// Repository interfaces
	LivroRepo       domain.LivroRepository
	AutorRepo       domain.AutorRepository
	EditoraRepo     domain.EditoraRepository
	GraficaRepo     domain.GraficaRepository
	ParticularRepo  domain.ParticularRepository
	ContratadaRepo  domain.ContratadaRepository
	ContratoRepo    domain.ContratoRepository
	EscreveRepo     domain.EscreveRepository
	ImprimeRepo     domain.ImprimeRepository
	
	// Aggregate repositories
	BookAuthorsRepo domain.BookAuthorsRepository
	PrintingJobRepo domain.PrintingJobRepository
}

// NewApplicationService creates a new ApplicationService with all use cases
func NewApplicationService(config ApplicationServiceConfig) *ApplicationService {
	
	// Create Livro use case
	livroUsecase := NewLivroUsecase(
		config.LivroRepo,
		config.EditoraRepo,
		config.AutorRepo,
		config.EscreveRepo,
		config.BookAuthorsRepo,
	)
	
	// Create Autor use case
	autorUsecase := NewAutorUsecase(
		config.AutorRepo,
		config.EscreveRepo,
		config.BookAuthorsRepo,
	)
	
	// Create Editora use case
	editoraUsecase := NewEditoraUsecase(
		config.EditoraRepo,
		config.LivroRepo,
	)
	
	// Create Grafica use case
	graficaUsecase := NewGraficaUsecase(
		config.GraficaRepo,
		config.ParticularRepo,
		config.ContratadaRepo,
		config.ContratoRepo,
		config.ImprimeRepo,
		config.PrintingJobRepo,
	)
	
	// Create Contrato use case
	contratoUsecase := NewContratoUsecase(
		config.ContratoRepo,
		config.ContratadaRepo,
		config.GraficaRepo,
		config.ParticularRepo,
		config.ImprimeRepo,
		config.PrintingJobRepo,
	)
	
	// Create Imprime use case
	imprimeUsecase := NewImprimeUsecase(
		config.ImprimeRepo,
		config.LivroRepo,
		config.GraficaRepo,
		config.ParticularRepo,
		config.ContratadaRepo,
		config.ContratoRepo,
		config.PrintingJobRepo,
	)
	
	return &ApplicationService{
		LivroUsecase:    livroUsecase,
		AutorUsecase:    autorUsecase,
		EditoraUsecase:  editoraUsecase,
		GraficaUsecase:  graficaUsecase,
		ContratoUsecase: contratoUsecase,
		ImprimeUsecase:  imprimeUsecase,
	}
}

// HealthCheck verifies that all use cases are properly initialized
func (a *ApplicationService) HealthCheck() error {
	if a.LivroUsecase == nil {
		return domain.NewDomainError("MISSING_USECASE", "LivroUsecase not initialized")
	}
	if a.AutorUsecase == nil {
		return domain.NewDomainError("MISSING_USECASE", "AutorUsecase not initialized")
	}
	if a.EditoraUsecase == nil {
		return domain.NewDomainError("MISSING_USECASE", "EditoraUsecase not initialized")
	}
	if a.GraficaUsecase == nil {
		return domain.NewDomainError("MISSING_USECASE", "GraficaUsecase not initialized")
	}
	if a.ContratoUsecase == nil {
		return domain.NewDomainError("MISSING_USECASE", "ContratoUsecase not initialized")
	}
	if a.ImprimeUsecase == nil {
		return domain.NewDomainError("MISSING_USECASE", "ImprimeUsecase not initialized")
	}
	
	return nil
}

// GetAllUsecases returns a map of all use cases for reflection or debugging
func (a *ApplicationService) GetAllUsecases() map[string]interface{} {
	return map[string]interface{}{
		"livro":    a.LivroUsecase,
		"autor":    a.AutorUsecase,
		"editora":  a.EditoraUsecase,
		"grafica":  a.GraficaUsecase,
		"contrato": a.ContratoUsecase,
		"imprime":  a.ImprimeUsecase,
	}
}