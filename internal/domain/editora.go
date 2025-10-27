package domain

import "strings"

type Editora struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
}

// NewEditora creates a new Editora with validated fields
func NewEditora(id int, nome, endereco string) *Editora {
	return &Editora{
		ID:       id,
		Nome:     strings.TrimSpace(nome),
		Endereco: strings.TrimSpace(endereco),
	}
}

// ValidateID checks if the ID is valid
func (e *Editora) ValidateID() bool {
	return e.ID > 0
}

// ValidateNome checks if the name is valid
func (e *Editora) ValidateNome() bool {
	return len(e.Nome) > 0 && len(e.Nome) <= 255
}

// ValidateEndereco checks if the address is valid
func (e *Editora) ValidateEndereco() bool {
	return len(e.Endereco) > 0 && len(e.Endereco) <= 255
}

// IsValid checks if all editora fields are valid
func (e *Editora) IsValid() bool {
	return e.ValidateID() && e.ValidateNome() && e.ValidateEndereco()
}

// GetFullInfo returns formatted editora information
func (e *Editora) GetFullInfo() string {
	return e.Nome + " (ID: " + string(rune(e.ID)) + ")"
}