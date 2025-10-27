package applications

import (
	"context"
	"fmt"
	"edna/internal/domain"
	"time"
)

// DTOs for Livro operations
type CreateLivroRequest struct {
	ISBN             string    `json:"isbn"`
	Titulo           string    `json:"titulo"`
	DataDePublicacao time.Time `json:"data_de_publicacao"`
	EditoraID        int       `json:"editora_id"`
}

type UpdateLivroRequest struct {
	Titulo           string    `json:"titulo"`
	DataDePublicacao time.Time `json:"data_de_publicacao"`
	EditoraID        int       `json:"editora_id"`
}

type LivroResponse struct {
	ISBN             string    `json:"isbn"`
	Titulo           string    `json:"titulo"`
	DataDePublicacao time.Time `json:"data_de_publicacao"`
	EditoraID        int       `json:"editora_id"`
	Editora          *EditoraResponse `json:"editora,omitempty"`
	Autores          []AutorResponse  `json:"autores,omitempty"`
}

type LivroUsecase interface {
	List(ctx context.Context) ([]LivroResponse, error)
	Create(ctx context.Context, req CreateLivroRequest) (*LivroResponse, error)
	Get(ctx context.Context, isbn string) (*LivroResponse, error)
	GetWithAuthors(ctx context.Context, isbn string) (*LivroResponse, error)
	Update(ctx context.Context, isbn string, req UpdateLivroRequest) (*LivroResponse, error)
	Delete(ctx context.Context, isbn string) error
	GetByEditora(ctx context.Context, editoraID int) ([]LivroResponse, error)
	GetByPublicationDateRange(ctx context.Context, start, end time.Time) ([]LivroResponse, error)
	AddAuthor(ctx context.Context, isbn, authorRG string) error
	RemoveAuthor(ctx context.Context, isbn, authorRG string) error
}

type livroService struct {
	livroRepo       domain.LivroRepository
	editoraRepo     domain.EditoraRepository
	autorRepo       domain.AutorRepository
	escreveRepo     domain.EscreveRepository
	bookAuthorsRepo domain.BookAuthorsRepository
	bookService     *domain.BookService
}

func NewLivroUsecase(
	livroRepo domain.LivroRepository,
	editoraRepo domain.EditoraRepository,
	autorRepo domain.AutorRepository,
	escreveRepo domain.EscreveRepository,
	bookAuthorsRepo domain.BookAuthorsRepository,
) LivroUsecase {
	bookService := domain.NewBookService(livroRepo, autorRepo, editoraRepo, escreveRepo, bookAuthorsRepo)
	
	return &livroService{
		livroRepo:       livroRepo,
		editoraRepo:     editoraRepo,
		autorRepo:       autorRepo,
		escreveRepo:     escreveRepo,
		bookAuthorsRepo: bookAuthorsRepo,
		bookService:     bookService,
	}
}

func (s *livroService) List(ctx context.Context) ([]LivroResponse, error) {
	livros, err := s.livroRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %w", err)
	}

	var response []LivroResponse
	for _, livro := range livros {
		response = append(response, s.toLivroResponse(livro))
	}

	return response, nil
}

func (s *livroService) Create(ctx context.Context, req CreateLivroRequest) (*LivroResponse, error) {
	// Validate request
	if err := s.validateCreateLivroRequest(req); err != nil {
		return nil, err
	}

	// Create domain entity
	livro := domain.NewLivro(req.ISBN, req.Titulo, req.DataDePublicacao, req.EditoraID)

	// Use domain service to create book with validation
	if err := s.bookService.CreateBook(ctx, livro); err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	response := s.toLivroResponse(livro)
	return &response, nil
}

func (s *livroService) Get(ctx context.Context, isbn string) (*LivroResponse, error) {
	if isbn == "" {
		return nil, domain.NewFieldRequiredError("isbn")
	}

	livro, err := s.livroRepo.FindByISBN(ctx, isbn)
	if err != nil {
		return nil, domain.NewBookNotFoundError(isbn)
	}

	response := s.toLivroResponse(livro)
	return &response, nil
}

func (s *livroService) GetWithAuthors(ctx context.Context, isbn string) (*LivroResponse, error) {
	if isbn == "" {
		return nil, domain.NewFieldRequiredError("isbn")
	}

	livro, authors, err := s.bookService.GetBookWithAuthors(ctx, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to get book with authors: %w", err)
	}

	response := s.toLivroResponse(livro)
	
	// Add authors to response
	for _, author := range authors {
		response.Autores = append(response.Autores, AutorResponse{
			RG:       author.RG,
			Nome:     author.Nome,
			Endereco: author.Endereco,
		})
	}

	return &response, nil
}

func (s *livroService) Update(ctx context.Context, isbn string, req UpdateLivroRequest) (*LivroResponse, error) {
	if isbn == "" {
		return nil, domain.NewFieldRequiredError("isbn")
	}

	// Validate request
	if err := s.validateUpdateLivroRequest(req); err != nil {
		return nil, err
	}

	// Get existing book
	existingLivro, err := s.livroRepo.FindByISBN(ctx, isbn)
	if err != nil {
		return nil, domain.NewBookNotFoundError(isbn)
	}

	// Update fields
	existingLivro.Titulo = req.Titulo
	existingLivro.DataDePublicacao = req.DataDePublicacao
	existingLivro.EditoraID = req.EditoraID

	// Validate updated book
	if !existingLivro.IsValid() {
		return nil, fmt.Errorf("invalid book data after update")
	}

	// Validate that new editora exists
	if req.EditoraID != 0 {
		_, err := s.editoraRepo.FindByID(ctx, req.EditoraID)
		if err != nil {
			return nil, domain.NewPublisherNotFoundError(req.EditoraID)
		}
	}

	// Update in repository
	if err := s.livroRepo.Update(ctx, existingLivro); err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	response := s.toLivroResponse(existingLivro)
	return &response, nil
}

func (s *livroService) Delete(ctx context.Context, isbn string) error {
	if isbn == "" {
		return domain.NewFieldRequiredError("isbn")
	}

	// Check if book exists
	_, err := s.livroRepo.FindByISBN(ctx, isbn)
	if err != nil {
		return domain.NewBookNotFoundError(isbn)
	}

	// Delete author relationships first
	if err := s.escreveRepo.DeleteByISBN(ctx, isbn); err != nil {
		return fmt.Errorf("failed to delete author relationships: %w", err)
	}

	// Delete the book
	if err := s.livroRepo.Delete(ctx, isbn); err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

func (s *livroService) GetByEditora(ctx context.Context, editoraID int) ([]LivroResponse, error) {
	if editoraID <= 0 {
		return nil, domain.NewFieldRequiredError("editora_id")
	}

	livros, err := s.livroRepo.FindByEditora(ctx, editoraID)
	if err != nil {
		return nil, fmt.Errorf("failed to get books by editora: %w", err)
	}

	var response []LivroResponse
	for _, livro := range livros {
		response = append(response, s.toLivroResponse(livro))
	}

	return response, nil
}

func (s *livroService) GetByPublicationDateRange(ctx context.Context, start, end time.Time) ([]LivroResponse, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("start and end dates are required")
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	livros, err := s.livroRepo.FindByPublicationDateRange(ctx, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get books by date range: %w", err)
	}

	var response []LivroResponse
	for _, livro := range livros {
		response = append(response, s.toLivroResponse(livro))
	}

	return response, nil
}

func (s *livroService) AddAuthor(ctx context.Context, isbn, authorRG string) error {
	if isbn == "" {
		return domain.NewFieldRequiredError("isbn")
	}
	if authorRG == "" {
		return domain.NewFieldRequiredError("author_rg")
	}

	return s.bookService.AddAuthorToBook(ctx, isbn, authorRG)
}

func (s *livroService) RemoveAuthor(ctx context.Context, isbn, authorRG string) error {
	if isbn == "" {
		return domain.NewFieldRequiredError("isbn")
	}
	if authorRG == "" {
		return domain.NewFieldRequiredError("author_rg")
	}

	// Check if this is the last author (business rule: book must have at least one author)
	authors, err := s.bookAuthorsRepo.FindAuthorsByBook(ctx, isbn)
	if err != nil {
		return fmt.Errorf("failed to get book authors: %w", err)
	}

	if len(authors) <= 1 {
		return domain.NewDomainError(domain.CodeCannotRemoveLastAuthor, "cannot remove the last author from a book")
	}

	return s.escreveRepo.Delete(ctx, isbn, authorRG)
}

// Helper methods

func (s *livroService) validateCreateLivroRequest(req CreateLivroRequest) error {
	if req.ISBN == "" {
		return domain.NewFieldRequiredError("isbn")
	}
	if req.Titulo == "" {
		return domain.NewFieldRequiredError("titulo")
	}
	if req.DataDePublicacao.IsZero() {
		return domain.NewFieldRequiredError("data_de_publicacao")
	}
	if req.EditoraID <= 0 {
		return domain.NewFieldRequiredError("editora_id")
	}

	if len(req.ISBN) > 20 {
		return domain.NewFieldTooLongError("isbn", 20, len(req.ISBN))
	}
	if len(req.Titulo) > 255 {
		return domain.NewFieldTooLongError("titulo", 255, len(req.Titulo))
	}

	if req.DataDePublicacao.After(time.Now()) {
		return domain.NewInvalidFormatError("data_de_publicacao", "publication date cannot be in the future")
	}

	return nil
}

func (s *livroService) validateUpdateLivroRequest(req UpdateLivroRequest) error {
	if req.Titulo == "" {
		return domain.NewFieldRequiredError("titulo")
	}
	if req.DataDePublicacao.IsZero() {
		return domain.NewFieldRequiredError("data_de_publicacao")
	}
	if req.EditoraID <= 0 {
		return domain.NewFieldRequiredError("editora_id")
	}

	if len(req.Titulo) > 255 {
		return domain.NewFieldTooLongError("titulo", 255, len(req.Titulo))
	}

	if req.DataDePublicacao.After(time.Now()) {
		return domain.NewInvalidFormatError("data_de_publicacao", "publication date cannot be in the future")
	}

	return nil
}

func (s *livroService) toLivroResponse(livro *domain.Livro) LivroResponse {
	return LivroResponse{
		ISBN:             livro.ISBN,
		Titulo:           livro.Titulo,
		DataDePublicacao: livro.DataDePublicacao,
		EditoraID:        livro.EditoraID,
	}
}
