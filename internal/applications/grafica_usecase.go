package applications

import (
	"context"
	"fmt"
	"edna/internal/domain"
)

// DTOs for Grafica operations
type CreateGraficaRequest struct {
	Nome        string `json:"nome"`
	Type        string `json:"type"` // "particular" or "contratada"
	Endereco    string `json:"endereco,omitempty"` // Required only for contratada
}

type UpdateGraficaRequest struct {
	Nome     string `json:"nome"`
	Endereco string `json:"endereco,omitempty"` // Only for contratada
}

type GraficaResponse struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Type     string `json:"type"`
	Endereco string `json:"endereco,omitempty"`
	Contratos []ContratoResponse `json:"contratos,omitempty"`
	PrintingJobs []PrintingJobResponse `json:"printing_jobs,omitempty"`
}

type ContratoResponse struct {
	ID              int     `json:"id"`
	Valor           float64 `json:"valor"`
	NomeResponsavel string  `json:"nome_responsavel"`
	GraficaContID   int     `json:"grafica_cont_id"`
}

type PrintingJobResponse struct {
	ISBN        string `json:"isbn"`
	BookTitle   string `json:"book_title,omitempty"`
	GraficaID   int    `json:"grafica_id"`
	NtoCopias   int    `json:"nto_copias"`
	DataEntrega string `json:"data_entrega"`
	Status      string `json:"status,omitempty"`
}

type GraficaUsecase interface {
	List(ctx context.Context) ([]GraficaResponse, error)
	Create(ctx context.Context, req CreateGraficaRequest) (*GraficaResponse, error)
	Get(ctx context.Context, id int) (*GraficaResponse, error)
	GetWithContracts(ctx context.Context, id int) (*GraficaResponse, error)
	GetWithPrintingJobs(ctx context.Context, id int) (*GraficaResponse, error)
	Update(ctx context.Context, id int, req UpdateGraficaRequest) (*GraficaResponse, error)
	Delete(ctx context.Context, id int) error
	GetByName(ctx context.Context, name string) ([]GraficaResponse, error)
	GetByType(ctx context.Context, graficaType string) ([]GraficaResponse, error)
}

type graficaService struct {
	graficaRepo        domain.GraficaRepository
	particularRepo     domain.ParticularRepository
	contratadaRepo     domain.ContratadaRepository
	contratoRepo       domain.ContratoRepository
	imprimeRepo        domain.ImprimeRepository
	printingService    *domain.PrintingService
}

func NewGraficaUsecase(
	graficaRepo domain.GraficaRepository,
	particularRepo domain.ParticularRepository,
	contratadaRepo domain.ContratadaRepository,
	contratoRepo domain.ContratoRepository,
	imprimeRepo domain.ImprimeRepository,
	printingJobRepo domain.PrintingJobRepository,
) GraficaUsecase {
	printingService := domain.NewPrintingService(
		graficaRepo, particularRepo, contratadaRepo, contratoRepo, imprimeRepo, printingJobRepo,
	)
	
	return &graficaService{
		graficaRepo:     graficaRepo,
		particularRepo:  particularRepo,
		contratadaRepo:  contratadaRepo,
		contratoRepo:    contratoRepo,
		imprimeRepo:     imprimeRepo,
		printingService: printingService,
	}
}

func (s *graficaService) List(ctx context.Context) ([]GraficaResponse, error) {
	graficas, err := s.graficaRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list printing companies: %w", err)
	}

	var response []GraficaResponse
	for _, grafica := range graficas {
		graficaResp, err := s.toGraficaResponse(ctx, grafica)
		if err != nil {
			return nil, fmt.Errorf("failed to convert grafica to response: %w", err)
		}
		response = append(response, *graficaResp)
	}

	return response, nil
}

func (s *graficaService) Create(ctx context.Context, req CreateGraficaRequest) (*GraficaResponse, error) {
	// Validate request
	if err := s.validateCreateGraficaRequest(req); err != nil {
		return nil, err
	}

	// Create domain entity (ID will be set by repository)
	grafica := domain.NewGrafica(0, req.Nome)

	// Determine if it's particular or contratada
	isParticular := req.Type == "particular"
	
	// Use domain service to create printing company with validation
	if err := s.printingService.CreatePrintingCompany(ctx, grafica, isParticular, req.Endereco); err != nil {
		return nil, fmt.Errorf("failed to create printing company: %w", err)
	}

	response, err := s.toGraficaResponse(ctx, grafica)
	if err != nil {
		return nil, fmt.Errorf("failed to convert grafica to response: %w", err)
	}

	return response, nil
}

func (s *graficaService) Get(ctx context.Context, id int) (*GraficaResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	grafica, err := s.graficaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewGraficaNotFoundError(id)
	}

	response, err := s.toGraficaResponse(ctx, grafica)
	if err != nil {
		return nil, fmt.Errorf("failed to convert grafica to response: %w", err)
	}

	return response, nil
}

func (s *graficaService) GetWithContracts(ctx context.Context, id int) (*GraficaResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	// Get printing company
	grafica, err := s.graficaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewGraficaNotFoundError(id)
	}

	response, err := s.toGraficaResponse(ctx, grafica)
	if err != nil {
		return nil, fmt.Errorf("failed to convert grafica to response: %w", err)
	}

	// Check if it's contratada and get contracts
	contratada, err := s.contratadaRepo.FindByGraficaID(ctx, id)
	if err == nil && contratada != nil {
		contracts, err := s.contratoRepo.FindByGraficaContID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get contracts: %w", err)
		}

		for _, contract := range contracts {
			response.Contratos = append(response.Contratos, ContratoResponse{
				ID:              contract.ID,
				Valor:           contract.Valor,
				NomeResponsavel: contract.NomeResponsavel,
				GraficaContID:   contract.GraficaContID,
			})
		}
	}

	return response, nil
}

func (s *graficaService) GetWithPrintingJobs(ctx context.Context, id int) (*GraficaResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	// Get printing company
	grafica, err := s.graficaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewGraficaNotFoundError(id)
	}

	response, err := s.toGraficaResponse(ctx, grafica)
	if err != nil {
		return nil, fmt.Errorf("failed to convert grafica to response: %w", err)
	}

	// Get printing jobs
	jobs, err := s.printingService.GetPrintingJobsByGrafica(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get printing jobs: %w", err)
	}

	for _, job := range jobs {
		status := "pending"
		if job.IsDeliveryOverdue() {
			status = "overdue"
		}

		response.PrintingJobs = append(response.PrintingJobs, PrintingJobResponse{
			ISBN:        job.LISBN,
			GraficaID:   job.GraficaID,
			NtoCopias:   job.NtoCopias,
			DataEntrega: job.DataEntrega.Format("2006-01-02"),
			Status:      status,
		})
	}

	return response, nil
}

func (s *graficaService) Update(ctx context.Context, id int, req UpdateGraficaRequest) (*GraficaResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	// Validate request
	if err := s.validateUpdateGraficaRequest(req); err != nil {
		return nil, err
	}

	// Get existing printing company
	existingGrafica, err := s.graficaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewGraficaNotFoundError(id)
	}

	// Update fields
	existingGrafica.Nome = req.Nome

	// Validate updated printing company
	if !existingGrafica.IsValid() {
		return nil, fmt.Errorf("invalid printing company data after update")
	}

	// Update in repository
	if err := s.graficaRepo.Update(ctx, existingGrafica); err != nil {
		return nil, fmt.Errorf("failed to update printing company: %w", err)
	}

	// Update address if it's contratada
	if req.Endereco != "" {
		contratada, err := s.contratadaRepo.FindByGraficaID(ctx, id)
		if err == nil && contratada != nil {
			contratada.Endereco = req.Endereco
			if err := s.contratadaRepo.Update(ctx, contratada); err != nil {
				return nil, fmt.Errorf("failed to update contracted company address: %w", err)
			}
		}
	}

	response, err := s.toGraficaResponse(ctx, existingGrafica)
	if err != nil {
		return nil, fmt.Errorf("failed to convert grafica to response: %w", err)
	}

	return response, nil
}

func (s *graficaService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.NewFieldRequiredError("id")
	}

	// Check if printing company exists
	_, err := s.graficaRepo.FindByID(ctx, id)
	if err != nil {
		return domain.NewGraficaNotFoundError(id)
	}

	// Check if printing company has printing jobs (business rule: cannot delete with active jobs)
	jobs, err := s.imprimeRepo.FindByGraficaID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check printing jobs: %w", err)
	}

	if len(jobs) > 0 {
		return domain.NewDomainError("GRAFICA_HAS_JOBS", "cannot delete printing company with active printing jobs")
	}

	// Delete contracts if it's contratada
	contracts, err := s.contratoRepo.FindByGraficaContID(ctx, id)
	if err == nil && len(contracts) > 0 {
		for _, contract := range contracts {
			if err := s.contratoRepo.Delete(ctx, contract.ID); err != nil {
				return fmt.Errorf("failed to delete contract: %w", err)
			}
		}
	}

	// Delete particular or contratada record
	_ = s.particularRepo.Delete(ctx, id)
	_ = s.contratadaRepo.Delete(ctx, id)

	// Delete the printing company
	if err := s.graficaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete printing company: %w", err)
	}

	return nil
}

func (s *graficaService) GetByName(ctx context.Context, name string) ([]GraficaResponse, error) {
	if name == "" {
		return nil, domain.NewFieldRequiredError("name")
	}

	graficas, err := s.graficaRepo.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get printing companies by name: %w", err)
	}

	var response []GraficaResponse
	for _, grafica := range graficas {
		graficaResp, err := s.toGraficaResponse(ctx, grafica)
		if err != nil {
			return nil, fmt.Errorf("failed to convert grafica to response: %w", err)
		}
		response = append(response, *graficaResp)
	}

	return response, nil
}

func (s *graficaService) GetByType(ctx context.Context, graficaType string) ([]GraficaResponse, error) {
	if graficaType == "" {
		return nil, domain.NewFieldRequiredError("type")
	}

	if graficaType != "particular" && graficaType != "contratada" {
		return nil, domain.NewInvalidFormatError("type", "must be 'particular' or 'contratada'")
	}

	var response []GraficaResponse

	if graficaType == "particular" {
		particulares, err := s.particularRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get particular printing companies: %w", err)
		}

		for _, particular := range particulares {
			grafica, err := s.graficaRepo.FindByID(ctx, particular.GraficaID)
			if err != nil {
				continue // Skip if grafica not found
			}

			graficaResp, err := s.toGraficaResponse(ctx, grafica)
			if err != nil {
				continue // Skip if conversion fails
			}
			response = append(response, *graficaResp)
		}
	} else {
		contratadas, err := s.contratadaRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get contracted printing companies: %w", err)
		}

		for _, contratada := range contratadas {
			grafica, err := s.graficaRepo.FindByID(ctx, contratada.GraficaID)
			if err != nil {
				continue // Skip if grafica not found
			}

			graficaResp, err := s.toGraficaResponse(ctx, grafica)
			if err != nil {
				continue // Skip if conversion fails
			}
			graficaResp.Endereco = contratada.Endereco
			response = append(response, *graficaResp)
		}
	}

	return response, nil
}

// Helper methods

func (s *graficaService) validateCreateGraficaRequest(req CreateGraficaRequest) error {
	if req.Nome == "" {
		return domain.NewFieldRequiredError("nome")
	}
	if req.Type == "" {
		return domain.NewFieldRequiredError("type")
	}

	if req.Type != "particular" && req.Type != "contratada" {
		return domain.NewInvalidFormatError("type", "must be 'particular' or 'contratada'")
	}

	if req.Type == "contratada" && req.Endereco == "" {
		return domain.NewFieldRequiredError("endereco")
	}

	if len(req.Nome) > 255 {
		return domain.NewFieldTooLongError("nome", 255, len(req.Nome))
	}
	if len(req.Endereco) > 255 {
		return domain.NewFieldTooLongError("endereco", 255, len(req.Endereco))
	}

	return nil
}

func (s *graficaService) validateUpdateGraficaRequest(req UpdateGraficaRequest) error {
	if req.Nome == "" {
		return domain.NewFieldRequiredError("nome")
	}

	if len(req.Nome) > 255 {
		return domain.NewFieldTooLongError("nome", 255, len(req.Nome))
	}
	if len(req.Endereco) > 255 {
		return domain.NewFieldTooLongError("endereco", 255, len(req.Endereco))
	}

	return nil
}

func (s *graficaService) toGraficaResponse(ctx context.Context, grafica *domain.Grafica) (*GraficaResponse, error) {
	response := &GraficaResponse{
		ID:   grafica.ID,
		Nome: grafica.Nome,
	}

	// Check if it's particular
	particular, err := s.particularRepo.FindByGraficaID(ctx, grafica.ID)
	if err == nil && particular != nil {
		response.Type = "particular"
		return response, nil
	}

	// Check if it's contratada
	contratada, err := s.contratadaRepo.FindByGraficaID(ctx, grafica.ID)
	if err == nil && contratada != nil {
		response.Type = "contratada"
		response.Endereco = contratada.Endereco
		return response, nil
	}

	// Default to unknown type
	response.Type = "unknown"
	return response, nil
}