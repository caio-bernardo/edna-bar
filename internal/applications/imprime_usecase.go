package applications

import (
	"context"
	"fmt"
	"time"
	"edna/internal/domain"
)

// DTOs for Imprime operations
type CreateImprimeRequest struct {
	LISBN       string    `json:"lisbn"`
	GraficaID   int       `json:"grafica_id"`
	NtoCopias   int       `json:"nto_copias"`
	DataEntrega time.Time `json:"data_entrega"`
}

type UpdateImprimeRequest struct {
	NtoCopias   int       `json:"nto_copias"`
	DataEntrega time.Time `json:"data_entrega"`
}

type ImprimeResponse struct {
	LISBN        string    `json:"lisbn"`
	BookTitle    string    `json:"book_title,omitempty"`
	GraficaID    int       `json:"grafica_id"`
	GraficaName  string    `json:"grafica_name,omitempty"`
	NtoCopias    int       `json:"nto_copias"`
	DataEntrega  time.Time `json:"data_entrega"`
	Status       string    `json:"status"`
	DaysUntilDelivery int  `json:"days_until_delivery"`
	IsOverdue    bool      `json:"is_overdue"`
}

type PrintingStatisticsResponse struct {
	Period              string `json:"period"`
	TotalJobs           int    `json:"total_jobs"`
	TotalCopies         int    `json:"total_copies"`
	MostActiveGraficaID int    `json:"most_active_grafica_id"`
	OverdueJobs         int    `json:"overdue_jobs"`
	PendingJobs         int    `json:"pending_jobs"`
}

type ImprimeUsecase interface {
	List(ctx context.Context) ([]ImprimeResponse, error)
	Create(ctx context.Context, req CreateImprimeRequest) (*ImprimeResponse, error)
	Get(ctx context.Context, lisbn string, graficaID int) (*ImprimeResponse, error)
	Update(ctx context.Context, lisbn string, graficaID int, req UpdateImprimeRequest) (*ImprimeResponse, error)
	Delete(ctx context.Context, lisbn string, graficaID int) error
	GetByISBN(ctx context.Context, lisbn string) ([]ImprimeResponse, error)
	GetByGrafica(ctx context.Context, graficaID int) ([]ImprimeResponse, error)
	GetByDeliveryDateRange(ctx context.Context, start, end time.Time) ([]ImprimeResponse, error)
	GetOverdueJobs(ctx context.Context) ([]ImprimeResponse, error)
	GetPendingJobs(ctx context.Context) ([]ImprimeResponse, error)
	GetPrintingStatistics(ctx context.Context, start, end time.Time) (*PrintingStatisticsResponse, error)
	MarkAsCompleted(ctx context.Context, lisbn string, graficaID int) error
}

type imprimeService struct {
	imprimeRepo      domain.ImprimeRepository
	livroRepo        domain.LivroRepository
	graficaRepo      domain.GraficaRepository
	printingService  *domain.PrintingService
	reportingService *domain.ReportingService
}

func NewImprimeUsecase(
	imprimeRepo domain.ImprimeRepository,
	livroRepo domain.LivroRepository,
	graficaRepo domain.GraficaRepository,
	particularRepo domain.ParticularRepository,
	contratadaRepo domain.ContratadaRepository,
	contratoRepo domain.ContratoRepository,
	printingJobRepo domain.PrintingJobRepository,
) ImprimeUsecase {
	printingService := domain.NewPrintingService(
		graficaRepo, particularRepo, contratadaRepo, contratoRepo, imprimeRepo, printingJobRepo,
	)
	reportingService := domain.NewReportingService(printingJobRepo, imprimeRepo, contratoRepo)
	
	return &imprimeService{
		imprimeRepo:      imprimeRepo,
		livroRepo:        livroRepo,
		graficaRepo:      graficaRepo,
		printingService:  printingService,
		reportingService: reportingService,
	}
}

func (s *imprimeService) List(ctx context.Context) ([]ImprimeResponse, error) {
	imprimes, err := s.imprimeRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list printing jobs: %w", err)
	}

	var response []ImprimeResponse
	for _, imprime := range imprimes {
		resp, err := s.toImprimeResponse(ctx, imprime)
		if err != nil {
			// Log error but continue with other items
			continue
		}
		response = append(response, *resp)
	}

	return response, nil
}

func (s *imprimeService) Create(ctx context.Context, req CreateImprimeRequest) (*ImprimeResponse, error) {
	// Validate request
	if err := s.validateCreateImprimeRequest(req); err != nil {
		return nil, err
	}

	// Create domain entity
	imprime := domain.NewImprime(req.LISBN, req.GraficaID, req.NtoCopias, req.DataEntrega)

	// Use domain service to schedule printing job with validation
	if err := s.printingService.SchedulePrintingJob(ctx, imprime); err != nil {
		return nil, fmt.Errorf("failed to create printing job: %w", err)
	}

	response, err := s.toImprimeResponse(ctx, imprime)
	if err != nil {
		return nil, fmt.Errorf("failed to convert printing job to response: %w", err)
	}

	return response, nil
}

func (s *imprimeService) Get(ctx context.Context, lisbn string, graficaID int) (*ImprimeResponse, error) {
	if lisbn == "" {
		return nil, domain.NewFieldRequiredError("lisbn")
	}
	if graficaID <= 0 {
		return nil, domain.NewFieldRequiredError("grafica_id")
	}

	imprime, err := s.imprimeRepo.FindByISBNAndGraficaID(ctx, lisbn, graficaID)
	if err != nil {
		return nil, domain.NewPrintingJobNotFoundError(lisbn, graficaID)
	}

	response, err := s.toImprimeResponse(ctx, imprime)
	if err != nil {
		return nil, fmt.Errorf("failed to convert printing job to response: %w", err)
	}

	return response, nil
}

func (s *imprimeService) Update(ctx context.Context, lisbn string, graficaID int, req UpdateImprimeRequest) (*ImprimeResponse, error) {
	if lisbn == "" {
		return nil, domain.NewFieldRequiredError("lisbn")
	}
	if graficaID <= 0 {
		return nil, domain.NewFieldRequiredError("grafica_id")
	}

	// Validate request
	if err := s.validateUpdateImprimeRequest(req); err != nil {
		return nil, err
	}

	// Get existing printing job
	existingImprime, err := s.imprimeRepo.FindByISBNAndGraficaID(ctx, lisbn, graficaID)
	if err != nil {
		return nil, domain.NewPrintingJobNotFoundError(lisbn, graficaID)
	}

	// Update fields
	existingImprime.NtoCopias = req.NtoCopias
	existingImprime.DataEntrega = req.DataEntrega

	// Validate updated printing job
	if !existingImprime.IsValid() {
		return nil, fmt.Errorf("invalid printing job data after update")
	}

	// Update in repository
	if err := s.imprimeRepo.Update(ctx, existingImprime); err != nil {
		return nil, fmt.Errorf("failed to update printing job: %w", err)
	}

	response, err := s.toImprimeResponse(ctx, existingImprime)
	if err != nil {
		return nil, fmt.Errorf("failed to convert printing job to response: %w", err)
	}

	return response, nil
}

func (s *imprimeService) Delete(ctx context.Context, lisbn string, graficaID int) error {
	if lisbn == "" {
		return domain.NewFieldRequiredError("lisbn")
	}
	if graficaID <= 0 {
		return domain.NewFieldRequiredError("grafica_id")
	}

	// Check if printing job exists
	_, err := s.imprimeRepo.FindByISBNAndGraficaID(ctx, lisbn, graficaID)
	if err != nil {
		return domain.NewPrintingJobNotFoundError(lisbn, graficaID)
	}

	// Delete the printing job
	if err := s.imprimeRepo.Delete(ctx, lisbn, graficaID); err != nil {
		return fmt.Errorf("failed to delete printing job: %w", err)
	}

	return nil
}

func (s *imprimeService) GetByISBN(ctx context.Context, lisbn string) ([]ImprimeResponse, error) {
	if lisbn == "" {
		return nil, domain.NewFieldRequiredError("lisbn")
	}

	imprimes, err := s.imprimeRepo.FindByISBN(ctx, lisbn)
	if err != nil {
		return nil, fmt.Errorf("failed to get printing jobs by ISBN: %w", err)
	}

	var response []ImprimeResponse
	for _, imprime := range imprimes {
		resp, err := s.toImprimeResponse(ctx, imprime)
		if err != nil {
			continue // Skip items with errors
		}
		response = append(response, *resp)
	}

	return response, nil
}

func (s *imprimeService) GetByGrafica(ctx context.Context, graficaID int) ([]ImprimeResponse, error) {
	if graficaID <= 0 {
		return nil, domain.NewFieldRequiredError("grafica_id")
	}

	imprimes, err := s.printingService.GetPrintingJobsByGrafica(ctx, graficaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get printing jobs by grafica: %w", err)
	}

	var response []ImprimeResponse
	for _, imprime := range imprimes {
		resp, err := s.toImprimeResponse(ctx, imprime)
		if err != nil {
			continue // Skip items with errors
		}
		response = append(response, *resp)
	}

	return response, nil
}

func (s *imprimeService) GetByDeliveryDateRange(ctx context.Context, start, end time.Time) ([]ImprimeResponse, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("start and end dates are required")
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	imprimes, err := s.imprimeRepo.FindByDeliveryDateRange(ctx, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get printing jobs by delivery date range: %w", err)
	}

	var response []ImprimeResponse
	for _, imprime := range imprimes {
		resp, err := s.toImprimeResponse(ctx, imprime)
		if err != nil {
			continue // Skip items with errors
		}
		response = append(response, *resp)
	}

	return response, nil
}

func (s *imprimeService) GetOverdueJobs(ctx context.Context) ([]ImprimeResponse, error) {
	imprimes, err := s.printingService.GetOverdueJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue printing jobs: %w", err)
	}

	var response []ImprimeResponse
	for _, imprime := range imprimes {
		resp, err := s.toImprimeResponse(ctx, imprime)
		if err != nil {
			continue // Skip items with errors
		}
		response = append(response, *resp)
	}

	return response, nil
}

func (s *imprimeService) GetPendingJobs(ctx context.Context) ([]ImprimeResponse, error) {
	imprimes, err := s.imprimeRepo.FindPendingDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending printing jobs: %w", err)
	}

	var response []ImprimeResponse
	for _, imprime := range imprimes {
		resp, err := s.toImprimeResponse(ctx, imprime)
		if err != nil {
			continue // Skip items with errors
		}
		response = append(response, *resp)
	}

	return response, nil
}

func (s *imprimeService) GetPrintingStatistics(ctx context.Context, start, end time.Time) (*PrintingStatisticsResponse, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("start and end dates are required")
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	stats, err := s.reportingService.GetPrintingStatistics(ctx, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get printing statistics: %w", err)
	}

	// Get additional statistics
	overdueJobs, err := s.imprimeRepo.FindOverdueDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue jobs count: %w", err)
	}

	pendingJobs, err := s.imprimeRepo.FindPendingDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending jobs count: %w", err)
	}

	return &PrintingStatisticsResponse{
		Period:              stats.Period,
		TotalJobs:           stats.TotalJobs,
		TotalCopies:         stats.TotalCopies,
		MostActiveGraficaID: stats.MostActiveGraficaID,
		OverdueJobs:         len(overdueJobs),
		PendingJobs:         len(pendingJobs),
	}, nil
}

func (s *imprimeService) MarkAsCompleted(ctx context.Context, lisbn string, graficaID int) error {
	if lisbn == "" {
		return domain.NewFieldRequiredError("lisbn")
	}
	if graficaID <= 0 {
		return domain.NewFieldRequiredError("grafica_id")
	}

	// Check if printing job exists
	_, err := s.imprimeRepo.FindByISBNAndGraficaID(ctx, lisbn, graficaID)
	if err != nil {
		return domain.NewPrintingJobNotFoundError(lisbn, graficaID)
	}

	// For now, we can delete the job to mark it as completed
	// In a more complex system, we might have a status field
	if err := s.imprimeRepo.Delete(ctx, lisbn, graficaID); err != nil {
		return fmt.Errorf("failed to mark printing job as completed: %w", err)
	}

	return nil
}

// Helper methods

func (s *imprimeService) validateCreateImprimeRequest(req CreateImprimeRequest) error {
	if req.LISBN == "" {
		return domain.NewFieldRequiredError("lisbn")
	}
	if req.GraficaID <= 0 {
		return domain.NewFieldRequiredError("grafica_id")
	}
	if req.NtoCopias <= 0 {
		return domain.NewZeroValueError("nto_copias")
	}
	if req.DataEntrega.IsZero() {
		return domain.NewFieldRequiredError("data_entrega")
	}

	if len(req.LISBN) > 20 {
		return domain.NewFieldTooLongError("lisbn", 20, len(req.LISBN))
	}

	if req.DataEntrega.Before(time.Now()) {
		return domain.NewInvalidFormatError("data_entrega", "delivery date cannot be in the past")
	}

	return nil
}

func (s *imprimeService) validateUpdateImprimeRequest(req UpdateImprimeRequest) error {
	if req.NtoCopias <= 0 {
		return domain.NewZeroValueError("nto_copias")
	}
	if req.DataEntrega.IsZero() {
		return domain.NewFieldRequiredError("data_entrega")
	}

	return nil
}

func (s *imprimeService) toImprimeResponse(ctx context.Context, imprime *domain.Imprime) (*ImprimeResponse, error) {
	response := &ImprimeResponse{
		LISBN:             imprime.LISBN,
		GraficaID:         imprime.GraficaID,
		NtoCopias:         imprime.NtoCopias,
		DataEntrega:       imprime.DataEntrega,
		DaysUntilDelivery: imprime.DaysUntilDelivery(),
		IsOverdue:         imprime.IsDeliveryOverdue(),
	}

	// Set status
	if response.IsOverdue {
		response.Status = "overdue"
	} else if response.DaysUntilDelivery <= 7 {
		response.Status = "urgent"
	} else {
		response.Status = "pending"
	}

	// Get book title
	if livro, err := s.livroRepo.FindByISBN(ctx, imprime.LISBN); err == nil {
		response.BookTitle = livro.Titulo
	}

	// Get grafica name
	if grafica, err := s.graficaRepo.FindByID(ctx, imprime.GraficaID); err == nil {
		response.GraficaName = grafica.Nome
	}

	return response, nil
}