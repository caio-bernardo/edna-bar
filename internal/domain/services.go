package domain

import (
	"context"
	"fmt"
	"time"
)

// BookService handles business logic related to books
type BookService struct {
	livroRepo       LivroRepository
	autorRepo       AutorRepository
	editoraRepo     EditoraRepository
	escreveRepo     EscreveRepository
	bookAuthorsRepo BookAuthorsRepository
}

// NewBookService creates a new BookService
func NewBookService(
	livroRepo LivroRepository, 
	autorRepo AutorRepository, 
	editoraRepo EditoraRepository,
	escreveRepo EscreveRepository,
	bookAuthorsRepo BookAuthorsRepository,
) *BookService {
	return &BookService{
		livroRepo:       livroRepo,
		autorRepo:       autorRepo,
		editoraRepo:     editoraRepo,
		escreveRepo:     escreveRepo,
		bookAuthorsRepo: bookAuthorsRepo,
	}
}

// CreateBook creates a new book with validation
func (bs *BookService) CreateBook(ctx context.Context, livro *Livro) error {
	if !livro.IsValid() {
		return fmt.Errorf("invalid book data")
	}

	// Check if book already exists
	existing, err := bs.livroRepo.FindByISBN(ctx, livro.ISBN)
	if err == nil && existing != nil {
		return fmt.Errorf("book with ISBN %s already exists", livro.ISBN)
	}

	// Validate that the editora exists
	_, err = bs.editoraRepo.FindByID(ctx, livro.EditoraID)
	if err != nil {
		return fmt.Errorf("editora with ID %d not found", livro.EditoraID)
	}

	return bs.livroRepo.Save(ctx, livro)
}

// AddAuthorToBook adds an author to a book
func (bs *BookService) AddAuthorToBook(ctx context.Context, isbn, rg string) error {
	// Validate book exists
	_, err := bs.livroRepo.FindByISBN(ctx, isbn)
	if err != nil {
		return fmt.Errorf("book with ISBN %s not found", isbn)
	}

	// Validate author exists
	_, err = bs.autorRepo.FindByRG(ctx, rg)
	if err != nil {
		return fmt.Errorf("author with RG %s not found", rg)
	}

	// Check if relationship already exists
	existing, err := bs.escreveRepo.FindByISBNAndRG(ctx, isbn, rg)
	if err == nil && existing != nil {
		return fmt.Errorf("author %s is already associated with book %s", rg, isbn)
	}

	escreve := NewEscreve(isbn, rg)
	return bs.escreveRepo.Save(ctx, escreve)
}

// GetBookWithAuthors returns a book with its authors
func (bs *BookService) GetBookWithAuthors(ctx context.Context, isbn string) (*Livro, []*Autor, error) {
	livro, err := bs.livroRepo.FindByISBN(ctx, isbn)
	if err != nil {
		return nil, nil, err
	}

	authors, err := bs.bookAuthorsRepo.FindAuthorsByBook(ctx, isbn)
	if err != nil {
		return nil, nil, err
	}

	return livro, authors, nil
}

// PrintingService handles business logic related to printing operations
type PrintingService struct {
	graficaRepo        GraficaRepository
	particularRepo     ParticularRepository
	contratadaRepo     ContratadaRepository
	contratoRepo       ContratoRepository
	imprimeRepo        ImprimeRepository
	printingJobRepo    PrintingJobRepository
}

// NewPrintingService creates a new PrintingService
func NewPrintingService(
	graficaRepo GraficaRepository,
	particularRepo ParticularRepository,
	contratadaRepo ContratadaRepository,
	contratoRepo ContratoRepository,
	imprimeRepo ImprimeRepository,
	printingJobRepo PrintingJobRepository,
) *PrintingService {
	return &PrintingService{
		graficaRepo:     graficaRepo,
		particularRepo:  particularRepo,
		contratadaRepo:  contratadaRepo,
		contratoRepo:    contratoRepo,
		imprimeRepo:     imprimeRepo,
		printingJobRepo: printingJobRepo,
	}
}

// CreatePrintingCompany creates a new printing company
func (ps *PrintingService) CreatePrintingCompany(ctx context.Context, grafica *Grafica, isParticular bool, endereco string) error {
	if !grafica.IsValid() {
		return fmt.Errorf("invalid grafica data")
	}

	// Save the base grafica
	err := ps.graficaRepo.Save(ctx, grafica)
	if err != nil {
		return err
	}

	// Create specific type
	if isParticular {
		particular := NewParticular(grafica.ID)
		return ps.particularRepo.Save(ctx, particular)
	} else {
		if endereco == "" {
			return fmt.Errorf("endereco is required for contratada grafica")
		}
		contratada := NewContratada(grafica.ID, endereco)
		return ps.contratadaRepo.Save(ctx, contratada)
	}
}

// CreateContract creates a new contract for a contracted printing company
func (ps *PrintingService) CreateContract(ctx context.Context, contrato *Contrato) error {
	if !contrato.IsValid() {
		return fmt.Errorf("invalid contract data")
	}

	// Validate that the grafica is contratada
	_, err := ps.contratadaRepo.FindByGraficaID(ctx, contrato.GraficaContID)
	if err != nil {
		return fmt.Errorf("grafica with ID %d is not a contracted company", contrato.GraficaContID)
	}

	return ps.contratoRepo.Save(ctx, contrato)
}

// SchedulePrintingJob schedules a printing job
func (ps *PrintingService) SchedulePrintingJob(ctx context.Context, imprime *Imprime) error {
	if !imprime.IsValid() {
		return fmt.Errorf("invalid printing job data")
	}

	// Validate grafica exists
	_, err := ps.graficaRepo.FindByID(ctx, imprime.GraficaID)
	if err != nil {
		return fmt.Errorf("grafica with ID %d not found", imprime.GraficaID)
	}

	// Check if there's already a printing job for this book at this grafica
	existing, err := ps.imprimeRepo.FindByISBNAndGraficaID(ctx, imprime.LISBN, imprime.GraficaID)
	if err == nil && existing != nil {
		return fmt.Errorf("printing job already exists for book %s at grafica %d", imprime.LISBN, imprime.GraficaID)
	}

	return ps.imprimeRepo.Save(ctx, imprime)
}

// GetOverdueJobs returns all overdue printing jobs
func (ps *PrintingService) GetOverdueJobs(ctx context.Context) ([]*Imprime, error) {
	return ps.imprimeRepo.FindOverdueDeliveries(ctx)
}

// GetPrintingJobsByGrafica returns all printing jobs for a specific grafica
func (ps *PrintingService) GetPrintingJobsByGrafica(ctx context.Context, graficaID int) ([]*Imprime, error) {
	return ps.imprimeRepo.FindByGraficaID(ctx, graficaID)
}

// PublisherService handles business logic related to publishers
type PublisherService struct {
	editoraRepo EditoraRepository
	livroRepo   LivroRepository
}

// NewPublisherService creates a new PublisherService
func NewPublisherService(editoraRepo EditoraRepository, livroRepo LivroRepository) *PublisherService {
	return &PublisherService{
		editoraRepo: editoraRepo,
		livroRepo:   livroRepo,
	}
}

// CreatePublisher creates a new publisher
func (ps *PublisherService) CreatePublisher(ctx context.Context, editora *Editora) error {
	if !editora.IsValid() {
		return fmt.Errorf("invalid publisher data")
	}

	// Check if publisher with same name already exists
	existing, err := ps.editoraRepo.FindByName(ctx, editora.Nome)
	if err == nil && len(existing) > 0 {
		return fmt.Errorf("publisher with name %s already exists", editora.Nome)
	}

	return ps.editoraRepo.Save(ctx, editora)
}

// GetPublisherBooks returns all books published by a specific publisher
func (ps *PublisherService) GetPublisherBooks(ctx context.Context, editoraID int) ([]*Livro, error) {
	// Validate publisher exists
	_, err := ps.editoraRepo.FindByID(ctx, editoraID)
	if err != nil {
		return nil, fmt.Errorf("publisher with ID %d not found", editoraID)
	}

	return ps.livroRepo.FindByEditora(ctx, editoraID)
}

// AuthorService handles business logic related to authors
type AuthorService struct {
	autorRepo       AutorRepository
	escreveRepo     EscreveRepository
	bookAuthorsRepo BookAuthorsRepository
}

// NewAuthorService creates a new AuthorService
func NewAuthorService(
	autorRepo AutorRepository, 
	escreveRepo EscreveRepository,
	bookAuthorsRepo BookAuthorsRepository,
) *AuthorService {
	return &AuthorService{
		autorRepo:       autorRepo,
		escreveRepo:     escreveRepo,
		bookAuthorsRepo: bookAuthorsRepo,
	}
}

// CreateAuthor creates a new author
func (as *AuthorService) CreateAuthor(ctx context.Context, autor *Autor) error {
	if !autor.IsValid() {
		return fmt.Errorf("invalid author data")
	}

	// Check if author already exists
	existing, err := as.autorRepo.FindByRG(ctx, autor.RG)
	if err == nil && existing != nil {
		return fmt.Errorf("author with RG %s already exists", autor.RG)
	}

	return as.autorRepo.Save(ctx, autor)
}

// GetAuthorBooks returns all books written by a specific author
func (as *AuthorService) GetAuthorBooks(ctx context.Context, rg string) ([]*Livro, error) {
	// Validate author exists
	_, err := as.autorRepo.FindByRG(ctx, rg)
	if err != nil {
		return nil, fmt.Errorf("author with RG %s not found", rg)
	}

	return as.bookAuthorsRepo.FindBooksByAuthor(ctx, rg)
}

// RemoveAuthorFromBook removes an author from a book
func (as *AuthorService) RemoveAuthorFromBook(ctx context.Context, isbn, rg string) error {
	// Validate relationship exists
	existing, err := as.escreveRepo.FindByISBNAndRG(ctx, isbn, rg)
	if err != nil || existing == nil {
		return fmt.Errorf("author %s is not associated with book %s", rg, isbn)
	}

	return as.escreveRepo.Delete(ctx, isbn, rg)
}

// ReportingService handles business logic for reports and analytics
type ReportingService struct {
	printingJobRepo PrintingJobRepository
	imprimeRepo     ImprimeRepository
	contratoRepo    ContratoRepository
}

// NewReportingService creates a new ReportingService
func NewReportingService(
	printingJobRepo PrintingJobRepository,
	imprimeRepo ImprimeRepository,
	contratoRepo ContratoRepository,
) *ReportingService {
	return &ReportingService{
		printingJobRepo: printingJobRepo,
		imprimeRepo:     imprimeRepo,
		contratoRepo:    contratoRepo,
	}
}

// GetPrintingStatistics returns printing statistics for a given period
func (rs *ReportingService) GetPrintingStatistics(ctx context.Context, startDate, endDate time.Time) (*PrintingStatistics, error) {
	jobs, err := rs.imprimeRepo.FindByDeliveryDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	stats := &PrintingStatistics{
		Period:     fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalJobs:  len(jobs),
		TotalCopies: 0,
	}

	graficaJobCount := make(map[int]int)
	for _, job := range jobs {
		stats.TotalCopies += job.NtoCopias
		graficaJobCount[job.GraficaID]++
	}

	// Find most active grafica
	maxJobs := 0
	for graficaID, jobCount := range graficaJobCount {
		if jobCount > maxJobs {
			maxJobs = jobCount
			stats.MostActiveGraficaID = graficaID
		}
	}

	return stats, nil
}

// PrintingStatistics represents printing statistics for a period
type PrintingStatistics struct {
	Period              string
	TotalJobs           int
	TotalCopies         int
	MostActiveGraficaID int
}

// ContractAnalysisService handles contract analysis and optimization
type ContractAnalysisService struct {
	contratoRepo ContratoRepository
}

// NewContractAnalysisService creates a new ContractAnalysisService
func NewContractAnalysisService(contratoRepo ContratoRepository) *ContractAnalysisService {
	return &ContractAnalysisService{
		contratoRepo: contratoRepo,
	}
}

// GetContractValueAnalysis returns analysis of contract values
func (cas *ContractAnalysisService) GetContractValueAnalysis(ctx context.Context) (*ContractValueAnalysis, error) {
	contracts, err := cas.contratoRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(contracts) == 0 {
		return &ContractValueAnalysis{}, nil
	}

	analysis := &ContractValueAnalysis{
		TotalContracts: len(contracts),
		MinValue:       contracts[0].Valor,
		MaxValue:       contracts[0].Valor,
		TotalValue:     0,
	}

	for _, contract := range contracts {
		analysis.TotalValue += contract.Valor
		if contract.Valor < analysis.MinValue {
			analysis.MinValue = contract.Valor
		}
		if contract.Valor > analysis.MaxValue {
			analysis.MaxValue = contract.Valor
		}
	}

	analysis.AverageValue = analysis.TotalValue / float64(len(contracts))

	return analysis, nil
}

// ContractValueAnalysis represents contract value analysis
type ContractValueAnalysis struct {
	TotalContracts int
	TotalValue     float64
	AverageValue   float64
	MinValue       float64
	MaxValue       float64
}