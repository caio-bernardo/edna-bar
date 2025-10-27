package domain

import (
	"strings"
	"time"
)

type Livro struct {
	ISBN               string    `json:"isbn"`
	Titulo            string    `json:"titulo"`
	DataDePublicacao  time.Time `json:"data_de_publicacao"`
	EditoraID         int       `json:"editora_id"`
}

// NewLivro creates a new Livro with validated fields
func NewLivro(isbn, titulo string, dataDePublicacao time.Time, editoraID int) *Livro {
	return &Livro{
		ISBN:             strings.TrimSpace(isbn),
		Titulo:           strings.TrimSpace(titulo),
		DataDePublicacao: dataDePublicacao,
		EditoraID:        editoraID,
	}
}

// ValidateISBN checks if the ISBN is valid
func (l *Livro) ValidateISBN() bool {
	return len(l.ISBN) > 0 && len(l.ISBN) <= 20
}

// ValidateTitulo checks if the title is valid
func (l *Livro) ValidateTitulo() bool {
	return len(l.Titulo) > 0 && len(l.Titulo) <= 255
}

// ValidateDataDePublicacao checks if the publication date is valid
func (l *Livro) ValidateDataDePublicacao() bool {
	return !l.DataDePublicacao.IsZero() && l.DataDePublicacao.Before(time.Now().AddDate(0, 0, 1))
}

// ValidateEditoraID checks if the publisher ID is valid
func (l *Livro) ValidateEditoraID() bool {
	return l.EditoraID > 0
}

// IsValid checks if all book fields are valid
func (l *Livro) IsValid() bool {
	return l.ValidateISBN() && 
		   l.ValidateTitulo() && 
		   l.ValidateDataDePublicacao() && 
		   l.ValidateEditoraID()
}

// GetAge returns the age of the book in years
func (l *Livro) GetAge() int {
	return int(time.Since(l.DataDePublicacao).Hours() / 24 / 365)
}

// IsNewRelease checks if the book was published within the last year
func (l *Livro) IsNewRelease() bool {
	return l.GetAge() < 1
}

// GetFullInfo returns formatted book information
func (l *Livro) GetFullInfo() string {
	return l.Titulo + " (ISBN: " + l.ISBN + ")"
}
