package applications

import (
	"context"
	"fmt"
	"edna/internal/domain"
)

// DTOs for Editora operations
type CreateEditoraRequest struct {
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
}

type UpdateEditoraRequest struct {
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
}

type EditoraResponse struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
	Livros   []LivroResponse `json:"livros,omitempty"`
}

type EditoraUsecase interface {
	List(ctx context.Context) ([]EditoraResponse, error)
	Create(ctx context.Context, req CreateEditoraRequest) (*EditoraResponse, error)
	Get(ctx context.Context, id int) (*EditoraResponse, error)
	GetWithBooks(ctx context.Context, id int) (*EditoraResponse, error)
	Update(ctx context.Context, id int, req UpdateEditoraRequest) (*EditoraResponse, error)
	Delete(ctx context.Context, id int) error
	GetByName(ctx context.Context, name string) ([]EditoraResponse, error)
}

type editoraService struct {
	editoraRepo     domain.EditoraRepository
	livroRepo       domain.LivroRepository
	publisherService *domain.PublisherService
}

func NewEditoraUsecase(
	editoraRepo domain.EditoraRepository,
	livroRepo domain.LivroRepository,
) EditoraUsecase {
	publisherService := domain.NewPublisherService(editoraRepo, livroRepo)
	
	return &editoraService{
		editoraRepo:      editoraRepo,
		livroRepo:        livroRepo,
		publisherService: publisherService,
	}
}

func (s *editoraService) List(ctx context.Context) ([]EditoraResponse, error) {
	editoras, err := s.editoraRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list publishers: %w", err)
	}

	var response []EditoraResponse
	for _, editora := range editoras {
		response = append(response, s.toEditoraResponse(editora))
	}

	return response, nil
}

func (s *editoraService) Create(ctx context.Context, req CreateEditoraRequest) (*EditoraResponse, error) {
	// Validate request
	if err := s.validateCreateEditoraRequest(req); err != nil {
		return nil, err
	}

	// Create domain entity (ID will be set by repository)
	editora := domain.NewEditora(0, req.Nome, req.Endereco)

	// Use domain service to create publisher with validation
	if err := s.publisherService.CreatePublisher(ctx, editora); err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	response := s.toEditoraResponse(editora)
	return &response, nil
}

func (s *editoraService) Get(ctx context.Context, id int) (*EditoraResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	editora, err := s.editoraRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewPublisherNotFoundError(id)
	}

	response := s.toEditoraResponse(editora)
	return &response, nil
}

func (s *editoraService) GetWithBooks(ctx context.Context, id int) (*EditoraResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	// Get publisher
	editora, err := s.editoraRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewPublisherNotFoundError(id)
	}

	// Get publisher's books
	books, err := s.publisherService.GetPublisherBooks(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get publisher books: %w", err)
	}

	response := s.toEditoraResponse(editora)
	
	// Add books to response
	for _, book := range books {
		response.Livros = append(response.Livros, LivroResponse{
			ISBN:             book.ISBN,
			Titulo:           book.Titulo,
			DataDePublicacao: book.DataDePublicacao,
			EditoraID:        book.EditoraID,
		})
	}

	return &response, nil
}

func (s *editoraService) Update(ctx context.Context, id int, req UpdateEditoraRequest) (*EditoraResponse, error) {
	if id <= 0 {
		return nil, domain.NewFieldRequiredError("id")
	}

	// Validate request
	if err := s.validateUpdateEditoraRequest(req); err != nil {
		return nil, err
	}

	// Get existing publisher
	existingEditora, err := s.editoraRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewPublisherNotFoundError(id)
	}

	// Update fields
	existingEditora.Nome = req.Nome
	existingEditora.Endereco = req.Endereco

	// Validate updated publisher
	if !existingEditora.IsValid() {
		return nil, fmt.Errorf("invalid publisher data after update")
	}

	// Update in repository
	if err := s.editoraRepo.Update(ctx, existingEditora); err != nil {
		return nil, fmt.Errorf("failed to update publisher: %w", err)
	}

	response := s.toEditoraResponse(existingEditora)
	return &response, nil
}

func (s *editoraService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.NewFieldRequiredError("id")
	}

	// Check if publisher exists
	_, err := s.editoraRepo.FindByID(ctx, id)
	if err != nil {
		return domain.NewPublisherNotFoundError(id)
	}

	// Check if publisher has books (business rule: cannot delete publisher with books)
	books, err := s.livroRepo.FindByEditora(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check publisher books: %w", err)
	}

	if len(books) > 0 {
		return domain.NewDomainError("PUBLISHER_HAS_BOOKS", "cannot delete publisher with associated books")
	}

	// Delete the publisher
	if err := s.editoraRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete publisher: %w", err)
	}

	return nil
}

func (s *editoraService) GetByName(ctx context.Context, name string) ([]EditoraResponse, error) {
	if name == "" {
		return nil, domain.NewFieldRequiredError("name")
	}

	editoras, err := s.editoraRepo.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get publishers by name: %w", err)
	}

	var response []EditoraResponse
	for _, editora := range editoras {
		response = append(response, s.toEditoraResponse(editora))
	}

	return response, nil
}

// Helper methods

func (s *editoraService) validateCreateEditoraRequest(req CreateEditoraRequest) error {
	if req.Nome == "" {
		return domain.NewFieldRequiredError("nome")
	}
	if req.Endereco == "" {
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

func (s *editoraService) validateUpdateEditoraRequest(req UpdateEditoraRequest) error {
	if req.Nome == "" {
		return domain.NewFieldRequiredError("nome")
	}
	if req.Endereco == "" {
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

func (s *editoraService) toEditoraResponse(editora *domain.Editora) EditoraResponse {
	return EditoraResponse{
		ID:       editora.ID,
		Nome:     editora.Nome,
		Endereco: editora.Endereco,
	}
}