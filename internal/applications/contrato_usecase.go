package applications

import (
	"context"
	"fmt"
	"edna/internal/domain"
)

// DTOs for Contrato operations
type CreateContratoRequest struct {
	Valor           float64 `json:"valor"`
	NomeResponsavel string  `json:"nome_responsavel"`
	GraficaContID   int     `json:"grafica_cont_id"`
}

type UpdateContratoRequest struct {
	Valor           float64 `json:"valor"`
	NomeResponsavel string  `json:"nome_responsavel"`
}

type ContratoUsecase interface {
	List(ctx context.Context) ([]ContratoResponse, error)
	Create(ctx context.Context, req CreateContratoRequest) (*ContratoResponse, error)
	Get(ctx context.Context, id int) (*ContratoResponse, error)
	Update(ctx context.Context, id int, req UpdateContratoRequest) (*ContratoResponse, error)
	Delete(ctx context.Context, id int) error
	GetByGraficaContID(ctx context.Context, graficaContID int) ([]ContratoResponse, error)
	GetByResponsavel(ctx context.Context, nomeResponsavel string) ([]ContratoResponse, error)
	GetByValueRange(ctx context.Context, minValue, maxValue float64) ([]ContratoResponse, error)
	GetContractAnalysis(ctx context.Context) (*ContractAnalysisResponse, error)
}

type ContractAnalysisResponse struct {
	TotalContracts int     `json:"total_contracts"`
	TotalValue     float64 `json:"total_value"`
	AverageValue   float64 `json:"average_value"`
	MinValue       float64 `json:"min_value"`
	MaxValue       float64 `json:"max_value"`
}

type contratoService struct {
	contratoRepo            domain.ContratoRepository
	contratadaRepo          domain.ContratadaRepository
	graficaRepo             domain.GraficaRepository
	printingService         *domain.PrintingService
	contractAnalysisService *domain.ContractAnalysisService
}

func NewContratoUsecase(
	contratoRepo domain.ContratoRepository,
	contratadaRepo domain.ContratadaRepository,
	graficaRepo domain.GraficaRepository,
	particularRepo domain.ParticularRepository,
	imprimeRepo domain.ImprimeRepository,
	printingJobRepo domain.PrintingJobRepository,
) ContratoUsecase {
	printingService := domain.NewPrintingService(
		graficaRepo, particularRepo, contratadaRepo, contratoRepo, imprimeRepo, printingJobRepo,
	)
	contractAnalysisService := domain.NewContractAnalysisService(contratoRepo)
	
	return &contratoService{
		contratoRepo:            contratoRepo,
		contratadaRepo:          contratadaRepo,
		graficaRepo:             graficaRepo,
		printingService:         printingService,
		contractAnalysisService: contractAnalysisService,
	}
}

func (s *contratoService) List(ctx context.Context) ([]ContratoResponse, error) {
	contratos, err := s.contratoRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts: %w", err)
	}

	var response []ContratoResponse
	for _, contrato := range contratos {
		response = append(response, s.toContratoResponse(contrato))
	}

	return response, nil
}

func (s *contratoService) Create(ctx context.Context, req CreateContratoRequest) (*ContratoResponse, error) {
	// Validate request
	if err := s.validateCreateContratoRequest(req); err != nil {
		return nil, err
	}

	// Create domain entity (ID will be set by repository)
	contrato := &domain.Contrato{
		ID:              0, // Will be set by repository
		Valor:           req.Valor,
		NomeResponsavel: req.NomeResponsavel,
		GraficaContID:   req.GraficaContID,
	}

	// Use domain service to create contract with validation
	if err := s.printingService.CreateContract(ctx, contrato); err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	response := s.toContratoResponse(contrato)
	return &response, nil
}

func (s *contratoService) Get(ctx context.Context, id int) (*ContratoResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	contrato, err := s.contratoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewContractNotFoundError(id)
	}

	response := s.toContratoResponse(contrato)
	return &response, nil
}

func (s *contratoService) Update(ctx context.Context, id int, req UpdateContratoRequest) (*ContratoResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	// Validate request
	if err := s.validateUpdateContratoRequest(req); err != nil {
		return nil, err
	}

	// Get existing contract
	existingContrato, err := s.contratoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewContractNotFoundError(id)
	}

	// Update fields
	existingContrato.Valor = req.Valor
	existingContrato.NomeResponsavel = req.NomeResponsavel

	// Validate updated contract
	if !existingContrato.IsValid() {
		return nil, fmt.Errorf("invalid contract data after update")
	}

	// Update in repository
	if err := s.contratoRepo.Update(ctx, existingContrato); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	response := s.toContratoResponse(existingContrato)
	return &response, nil
}

func (s *contratoService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.NewFieldRequiredError("id")
	}

	// Check if contract exists
	_, err := s.contratoRepo.FindByID(ctx, id)
	if err != nil {
		return domain.NewContractNotFoundError(id)
	}

	// Delete the contract
	if err := s.contratoRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	return nil
}

func (s *contratoService) GetByGraficaContID(ctx context.Context, graficaContID int) ([]ContratoResponse, error) {
	if graficaContID <= 0 {
		return nil, domain.NewFieldRequiredError("grafica_cont_id")
	}

	contratos, err := s.contratoRepo.FindByGraficaContID(ctx, graficaContID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts by grafica: %w", err)
	}

	var response []ContratoResponse
	for _, contrato := range contratos {
		response = append(response, s.toContratoResponse(contrato))
	}

	return response, nil
}

func (s *contratoService) GetByResponsavel(ctx context.Context, nomeResponsavel string) ([]ContratoResponse, error) {
	if nomeResponsavel == "" {
		return nil, domain.NewFieldRequiredError("nome_responsavel")
	}

	contratos, err := s.contratoRepo.FindByResponsavel(ctx, nomeResponsavel)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts by responsible person: %w", err)
	}

	var response []ContratoResponse
	for _, contrato := range contratos {
		response = append(response, s.toContratoResponse(contrato))
	}

	return response, nil
}

func (s *contratoService) GetByValueRange(ctx context.Context, minValue, maxValue float64) ([]ContratoResponse, error) {
	if minValue < 0 {
		return nil, domain.NewNegativeValueError("min_value", minValue)
	}
	if maxValue < 0 {
		return nil, domain.NewNegativeValueError("max_value", maxValue)
	}
	if minValue > maxValue {
		return nil, fmt.Errorf("min_value cannot be greater than max_value")
	}

	contratos, err := s.contratoRepo.FindByValueRange(ctx, minValue, maxValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts by value range: %w", err)
	}

	var response []ContratoResponse
	for _, contrato := range contratos {
		response = append(response, s.toContratoResponse(contrato))
	}

	return response, nil
}

func (s *contratoService) GetContractAnalysis(ctx context.Context) (*ContractAnalysisResponse, error) {
	analysis, err := s.contractAnalysisService.GetContractValueAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract analysis: %w", err)
	}

	return &ContractAnalysisResponse{
		TotalContracts: analysis.TotalContracts,
		TotalValue:     analysis.TotalValue,
		AverageValue:   analysis.AverageValue,
		MinValue:       analysis.MinValue,
		MaxValue:       analysis.MaxValue,
	}, nil
}

// Helper methods

func (s *contratoService) validateCreateContratoRequest(req CreateContratoRequest) error {
	if req.Valor <= 0 {
		return domain.NewZeroValueError("valor")
	}
	if req.NomeResponsavel == "" {
		return domain.NewFieldRequiredError("nome_responsavel")
	}
	if req.GraficaContID <= 0 {
		return domain.NewFieldRequiredError("grafica_cont_id")
	}

	if len(req.NomeResponsavel) > 255 {
		return domain.NewFieldTooLongError("nome_responsavel", 255, len(req.NomeResponsavel))
	}

	return nil
}

func (s *contratoService) validateUpdateContratoRequest(req UpdateContratoRequest) error {
	if req.Valor <= 0 {
		return domain.NewZeroValueError("valor")
	}
	if req.NomeResponsavel == "" {
		return domain.NewFieldRequiredError("nome_responsavel")
	}

	if len(req.NomeResponsavel) > 255 {
		return domain.NewFieldTooLongError("nome_responsavel", 255, len(req.NomeResponsavel))
	}

	return nil
}

func (s *contratoService) toContratoResponse(contrato *domain.Contrato) ContratoResponse {
	return ContratoResponse{
		ID:              contrato.ID,
		Valor:           contrato.Valor,
		NomeResponsavel: contrato.NomeResponsavel,
		GraficaContID:   contrato.GraficaContID,
	}
}