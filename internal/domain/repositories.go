package domain

import (
	"context"
	"time"
)

// EditoraRepository defines the interface for editora persistence operations
type EditoraRepository interface {
	Save(ctx context.Context, editora *Editora) error
	FindByID(ctx context.Context, id int) (*Editora, error)
	FindByName(ctx context.Context, name string) ([]*Editora, error)
	FindAll(ctx context.Context) ([]*Editora, error)
	Update(ctx context.Context, editora *Editora) error
	Delete(ctx context.Context, id int) error
}

// AutorRepository defines the interface for autor persistence operations
type AutorRepository interface {
	Save(ctx context.Context, autor *Autor) error
	FindByRG(ctx context.Context, rg string) (*Autor, error)
	FindByName(ctx context.Context, name string) ([]*Autor, error)
	FindAll(ctx context.Context) ([]*Autor, error)
	Update(ctx context.Context, autor *Autor) error
	Delete(ctx context.Context, rg string) error
}

// LivroRepository defines the interface for livro persistence operations
type LivroRepository interface {
	Save(ctx context.Context, livro *Livro) error
	FindByISBN(ctx context.Context, isbn string) (*Livro, error)
	FindByTitle(ctx context.Context, title string) ([]*Livro, error)
	FindByEditora(ctx context.Context, editoraID int) ([]*Livro, error)
	FindByPublicationDateRange(ctx context.Context, start, end time.Time) ([]*Livro, error)
	FindAll(ctx context.Context) ([]*Livro, error)
	Update(ctx context.Context, livro *Livro) error
	Delete(ctx context.Context, isbn string) error
}

// GraficaRepository defines the interface for grafica persistence operations
type GraficaRepository interface {
	Save(ctx context.Context, grafica *Grafica) error
	FindByID(ctx context.Context, id int) (*Grafica, error)
	FindByName(ctx context.Context, name string) ([]*Grafica, error)
	FindAll(ctx context.Context) ([]*Grafica, error)
	Update(ctx context.Context, grafica *Grafica) error
	Delete(ctx context.Context, id int) error
}

// ParticularRepository defines the interface for particular grafica operations
type ParticularRepository interface {
	Save(ctx context.Context, particular *Particular) error
	FindByGraficaID(ctx context.Context, graficaID int) (*Particular, error)
	FindAll(ctx context.Context) ([]*Particular, error)
	Delete(ctx context.Context, graficaID int) error
}

// ContratadaRepository defines the interface for contratada grafica operations
type ContratadaRepository interface {
	Save(ctx context.Context, contratada *Contratada) error
	FindByGraficaID(ctx context.Context, graficaID int) (*Contratada, error)
	FindByAddress(ctx context.Context, endereco string) ([]*Contratada, error)
	FindAll(ctx context.Context) ([]*Contratada, error)
	Update(ctx context.Context, contratada *Contratada) error
	Delete(ctx context.Context, graficaID int) error
}

// ContratoRepository defines the interface for contrato persistence operations
type ContratoRepository interface {
	Save(ctx context.Context, contrato *Contrato) error
	FindByID(ctx context.Context, id int) (*Contrato, error)
	FindByGraficaContID(ctx context.Context, graficaContID int) ([]*Contrato, error)
	FindByResponsavel(ctx context.Context, nomeResponsavel string) ([]*Contrato, error)
	FindByValueRange(ctx context.Context, minValue, maxValue float64) ([]*Contrato, error)
	FindAll(ctx context.Context) ([]*Contrato, error)
	Update(ctx context.Context, contrato *Contrato) error
	Delete(ctx context.Context, id int) error
}

// EscreveRepository defines the interface for escreve relationship operations
type EscreveRepository interface {
	Save(ctx context.Context, escreve *Escreve) error
	FindByISBN(ctx context.Context, isbn string) ([]*Escreve, error)
	FindByRG(ctx context.Context, rg string) ([]*Escreve, error)
	FindByISBNAndRG(ctx context.Context, isbn, rg string) (*Escreve, error)
	FindAll(ctx context.Context) ([]*Escreve, error)
	Delete(ctx context.Context, isbn, rg string) error
	DeleteByISBN(ctx context.Context, isbn string) error
	DeleteByRG(ctx context.Context, rg string) error
}

// ImprimeRepository defines the interface for imprime relationship operations
type ImprimeRepository interface {
	Save(ctx context.Context, imprime *Imprime) error
	FindByISBN(ctx context.Context, lisbn string) ([]*Imprime, error)
	FindByGraficaID(ctx context.Context, graficaID int) ([]*Imprime, error)
	FindByISBNAndGraficaID(ctx context.Context, lisbn string, graficaID int) (*Imprime, error)
	FindByDeliveryDateRange(ctx context.Context, start, end time.Time) ([]*Imprime, error)
	FindOverdueDeliveries(ctx context.Context) ([]*Imprime, error)
	FindPendingDeliveries(ctx context.Context) ([]*Imprime, error)
	FindAll(ctx context.Context) ([]*Imprime, error)
	Update(ctx context.Context, imprime *Imprime) error
	Delete(ctx context.Context, lisbn string, graficaID int) error
	DeleteByISBN(ctx context.Context, lisbn string) error
	DeleteByGraficaID(ctx context.Context, graficaID int) error
}

// BookAuthorsRepository provides aggregate operations for books and their authors
type BookAuthorsRepository interface {
	FindAuthorsByBook(ctx context.Context, isbn string) ([]*Autor, error)
	FindBooksByAuthor(ctx context.Context, rg string) ([]*Livro, error)
	AddAuthorToBook(ctx context.Context, isbn, rg string) error
	RemoveAuthorFromBook(ctx context.Context, isbn, rg string) error
}

// PrintingJobRepository provides aggregate operations for printing jobs
type PrintingJobRepository interface {
	FindBooksByGrafica(ctx context.Context, graficaID int) ([]*Livro, error)
	FindGraficasByBook(ctx context.Context, isbn string) ([]*Grafica, error)
	GetTotalCopiesByBook(ctx context.Context, isbn string) (int, error)
	GetTotalCopiesByGrafica(ctx context.Context, graficaID int) (int, error)
	CreatePrintingJob(ctx context.Context, isbn string, graficaID int, copies int, deliveryDate time.Time) error
}

// UnitOfWork defines the interface for managing transactions across multiple repositories
type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	EditoraRepository() EditoraRepository
	AutorRepository() AutorRepository
	LivroRepository() LivroRepository
	GraficaRepository() GraficaRepository
	ParticularRepository() ParticularRepository
	ContratadaRepository() ContratadaRepository
	ContratoRepository() ContratoRepository
	EscreveRepository() EscreveRepository
	ImprimeRepository() ImprimeRepository
	BookAuthorsRepository() BookAuthorsRepository
	PrintingJobRepository() PrintingJobRepository
}