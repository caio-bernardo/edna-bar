package domain

// Escreve represents the relationship between authors and books
// This is a many-to-many relationship table
type Escreve struct {
	ISBN string `json:"isbn"`
	RG   string `json:"rg"`
}

// NewEscreve creates a new Escreve relationship
func NewEscreve(isbn, rg string) *Escreve {
	return &Escreve{
		ISBN: isbn,
		RG:   rg,
	}
}

// ValidateISBN checks if the ISBN is valid (basic validation)
func (e *Escreve) ValidateISBN() bool {
	return len(e.ISBN) > 0 && len(e.ISBN) <= 20
}

// ValidateRG checks if the RG is valid (basic validation)
func (e *Escreve) ValidateRG() bool {
	return len(e.RG) > 0 && len(e.RG) <= 20
}

// IsValid checks if the relationship is valid
func (e *Escreve) IsValid() bool {
	return e.ValidateISBN() && e.ValidateRG()
}