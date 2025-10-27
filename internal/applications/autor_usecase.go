package applications

import (
	"context"
	"fmt"
	"edna/internal/domain"
)

// DTOs for Autor operations
type CreateAutorRequest struct {
	RG       string `json:"rg"`
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
}

type UpdateAutorRequest struct {
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
}

type AutorResponse struct {
	RG       string `json:"rg"`
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
	Livros   []LivroResponse `json:"livros,omitempty"`
}

type AutorUsecase interface {
	List(ctx context.Context) ([]AutorResponse, error)
	Create(ctx context.Context, req CreateAutorRequest) (*AutorResponse, error)
	Get(ctx context.Context, rg string) (*AutorResponse, error)
	GetWithBooks(ctx context.Context, rg string) (*AutorResponse, error)
	Update(ctx context.Context, rg string, req UpdateAutorRequest) (*AutorResponse, error)
	Delete(ctx context.Context, rg string) error
	GetByName(ctx context.Context, name string) ([]AutorResponse, error)
	AddToBook(ctx context.Context, rg, isbn string) error
	RemoveFromBook(ctx context.Context, rg, isbn string) error
}

type autorService struct {
	autorRepo       domain.AutorRepository
	escreveRepo     domain.EscreveRepository
	bookAuthorsRepo domain.BookAuthorsRepository
	authorService   *domain.AuthorService
}

func NewAutorUsecase(
	autorRepo domain.AutorRepository,
	escreveRepo domain.EscreveRepository,
	bookAuthorsRepo domain.BookAuthorsRepository,
) AutorUsecase {
	authorService := domain.NewAuthorService(autorRepo, escreveRepo, bookAuthorsRepo)
	
	return &autorService{
		autorRepo:       autorRepo,
		escreveRepo:     escreveRepo,
		bookAuthorsRepo: bookAuthorsRepo,
		authorService:   authorService,
	}
}

func (s *autorService) List(ctx context.Context) ([]AutorResponse, error) {
	autores, err := s.autorRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	var response []AutorResponse
	for _, autor := range autores {
		response = append(response, s.toAutorResponse(autor))
	}

	return response, nil
}

func (s *autorService) Create(ctx context.Context, req CreateAutorRequest) (*AutorResponse, error) {
	// Validate request
	if err := s.validateCreateAutorRequest(req); err != nil {
		return nil, err
	}

	// Create domain entity
	autor := domain.NewAutor(req.RG, req.Nome, req.Endereco)

	// Use domain service to create author with validation
	if err := s.authorService.CreateAuthor(ctx, autor); err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}

	response := s.toAutorResponse(autor)
	return &response, nil
}

func (s *autorService) Get(ctx context.Context, rg string) (*AutorResponse, error) {
	if rg == "" {
		return nil, domain.NewFieldRequiredError("rg")
	}

	autor, err := s.autorRepo.FindByRG(ctx, rg)
	if err != nil {
		return nil, domain.NewAuthorNotFoundError(rg)
	}

	response := s.toAutorResponse(autor)
	return &response, nil
}

func (s *autorService) GetWithBooks(ctx context.Context, rg string) (*AutorResponse, error) {
	if rg == "" {
		return nil, domain.NewFieldRequiredError("rg")
	}

	// Get author
	autor, err := s.autorRepo.FindByRG(ctx, rg)
	if err != nil {
		return nil, domain.NewAuthorNotFoundError(rg)
	}

	// Get author's books
	books, err := s.authorService.GetAuthorBooks(ctx, rg)
	if err != nil {
		return nil, fmt.Errorf("failed to get author books: %w", err)
	}

	response := s.toAutorResponse(autor)
	
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

func (s *autorService) Update(ctx context.Context, rg string, req UpdateAutorRequest) (*AutorResponse, error) {
	if rg == "" {
		return nil, domain.NewFieldRequiredError("rg")
	}

	// Validate request
	if err := s.validateUpdateAutorRequest(req); err != nil {
		return nil, err
	}

	// Get existing author
	existingAutor, err := s.autorRepo.FindByRG(ctx, rg)
	if err != nil {
		return nil, domain.NewAuthorNotFoundError(rg)
	}

	// Update fields
	existingAutor.Nome = req.Nome
	existingAutor.Endereco = req.Endereco

	// Validate updated author
	if !existingAutor.IsValid() {
		return nil, fmt.Errorf("invalid author data after update")
	}

	// Update in repository
	if err := s.autorRepo.Update(ctx, existingAutor); err != nil {
		return nil, fmt.Errorf("failed to update author: %w", err)
	}

	response := s.toAutorResponse(existingAutor)
	return &response, nil
}

func (s *autorService) Delete(ctx context.Context, rg string) error {
	if rg == "" {
		return domain.NewFieldRequiredError("rg")
	}

	// Check if author exists
	_, err := s.autorRepo.FindByRG(ctx, rg)
	if err != nil {
		return domain.NewAuthorNotFoundError(rg)
	}

	// Check if author has books (business rule: cannot delete author with books)
	books, err := s.bookAuthorsRepo.FindBooksByAuthor(ctx, rg)
	if err != nil {
		return fmt.Errorf("failed to check author books: %w", err)
	}

	if len(books) > 0 {
		return domain.NewDomainError("AUTHOR_HAS_BOOKS", "cannot delete author with associated books")
	}

	// Delete the author
	if err := s.autorRepo.Delete(ctx, rg); err != nil {
		return fmt.Errorf("failed to delete author: %w", err)
	}

	return nil
}

func (s *autorService) GetByName(ctx context.Context, name string) ([]AutorResponse, error) {
	if name == "" {
		return nil, domain.NewFieldRequiredError("name")
	}

	autores, err := s.autorRepo.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get authors by name: %w", err)
	}

	var response []AutorResponse
	for _, autor := range autores {
		response = append(response, s.toAutorResponse(autor))
	}

	return response, nil
}

func (s *autorService) AddToBook(ctx context.Context, rg, isbn string) error {
	if rg == "" {
		return domain.NewFieldRequiredError("rg")
	}
	if isbn == "" {
		return domain.NewFieldRequiredError("isbn")
	}

	// Check if relationship already exists
	existing, err := s.escreveRepo.FindByISBNAndRG(ctx, isbn, rg)
	if err == nil && existing != nil {
		return domain.NewRelationshipAlreadyExistsError(isbn, rg)
	}

	escreve := domain.NewEscreve(isbn, rg)
	if err := s.escreveRepo.Save(ctx, escreve); err != nil {
		return fmt.Errorf("failed to add author to book: %w", err)
	}

	return nil
}

func (s *autorService) RemoveFromBook(ctx context.Context, rg, isbn string) error {
	if rg == "" {
		return domain.NewFieldRequiredError("rg")
	}
	if isbn == "" {
		return domain.NewFieldRequiredError("isbn")
	}

	return s.authorService.RemoveAuthorFromBook(ctx, isbn, rg)
}

// Helper methods

func (s *autorService) validateCreateAutorRequest(req CreateAutorRequest) error {
	if req.RG == "" {
		return domain.NewFieldRequiredError("rg")
	}
	if req.Nome == "" {
		return domain.NewFieldRequiredError("nome")
	}
	if req.Endereco == "" {
		return domain.NewFieldRequiredError("endereco")
	}

	if len(req.RG) > 20 {
		return domain.NewFieldTooLongError("rg", 20, len(req.RG))
	}
	if len(req.Nome) > 255 {
		return domain.NewFieldTooLongError("nome", 255, len(req.Nome))
	}
	if len(req.Endereco) > 255 {
		return domain.NewFieldTooLongError("endereco", 255, len(req.Endereco))
	}

	return nil
}

func (s *autorService) validateUpdateAutorRequest(req UpdateAutorRequest) error {
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

func (s *autorService) toAutorResponse(autor *domain.Autor) AutorResponse {
	return AutorResponse{
		RG:       autor.RG,
		Nome:     autor.Nome,
		Endereco: autor.Endereco,
	}
}