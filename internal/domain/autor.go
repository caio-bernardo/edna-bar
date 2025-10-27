package domain

import "strings"

type Autor struct {
	RG       string `json:"rg"`
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
}

// NewAutor creates a new Autor with validated fields
func NewAutor(rg, nome, endereco string) *Autor {
	return &Autor{
		RG:       strings.TrimSpace(rg),
		Nome:     strings.TrimSpace(nome),
		Endereco: strings.TrimSpace(endereco),
	}
}

// ValidateRG checks if the RG is valid
func (a *Autor) ValidateRG() bool {
	return len(a.RG) > 0 && len(a.RG) <= 20
}

// ValidateNome checks if the name is valid
func (a *Autor) ValidateNome() bool {
	return len(a.Nome) > 0 && len(a.Nome) <= 255
}

// ValidateEndereco checks if the address is valid
func (a *Autor) ValidateEndereco() bool {
	return len(a.Endereco) > 0 && len(a.Endereco) <= 255
}

// IsValid checks if all author fields are valid
func (a *Autor) IsValid() bool {
	return a.ValidateRG() && a.ValidateNome() && a.ValidateEndereco()
}

// GetFullInfo returns formatted author information
func (a *Autor) GetFullInfo() string {
	return a.Nome + " (RG: " + a.RG + ")"
}
