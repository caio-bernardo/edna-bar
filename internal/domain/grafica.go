package domain

import "strings"

// Grafica represents a printing company in the system
type Grafica struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// NewGrafica creates a new Grafica with validated fields
func NewGrafica(id int, nome string) *Grafica {
	return &Grafica{
		ID:   id,
		Nome: strings.TrimSpace(nome),
	}
}

// ValidateID checks if the ID is valid
func (g *Grafica) ValidateID() bool {
	return g.ID > 0
}

// ValidateNome checks if the name is valid
func (g *Grafica) ValidateNome() bool {
	return len(g.Nome) > 0 && len(g.Nome) <= 255
}

// IsValid checks if all grafica fields are valid
func (g *Grafica) IsValid() bool {
	return g.ValidateID() && g.ValidateNome()
}

// GetFullInfo returns formatted grafica information
func (g *Grafica) GetFullInfo() string {
	return g.Nome + " (ID: " + string(rune(g.ID)) + ")"
}

// Particular represents a private printing company
type Particular struct {
	GraficaID int `json:"grafica_id"`
}

// NewParticular creates a new Particular with validated fields
func NewParticular(graficaID int) *Particular {
	return &Particular{
		GraficaID: graficaID,
	}
}

// ValidateGraficaID checks if the printing company ID is valid
func (p *Particular) ValidateGraficaID() bool {
	return p.GraficaID > 0
}

// IsValid checks if the particular printing company is valid
func (p *Particular) IsValid() bool {
	return p.ValidateGraficaID()
}

// Contratada represents a contracted printing company
type Contratada struct {
	GraficaID int    `json:"grafica_id"`
	Endereco  string `json:"endereco"`
}

// NewContratada creates a new Contratada with validated fields
func NewContratada(graficaID int, endereco string) *Contratada {
	return &Contratada{
		GraficaID: graficaID,
		Endereco:  strings.TrimSpace(endereco),
	}
}

// ValidateGraficaID checks if the printing company ID is valid
func (c *Contratada) ValidateGraficaID() bool {
	return c.GraficaID > 0
}

// ValidateEndereco checks if the address is valid
func (c *Contratada) ValidateEndereco() bool {
	return len(c.Endereco) > 0 && len(c.Endereco) <= 255
}

// IsValid checks if the contracted printing company is valid
func (c *Contratada) IsValid() bool {
	return c.ValidateGraficaID() && c.ValidateEndereco()
}

// GetFullInfo returns formatted contracted printing company information
func (c *Contratada) GetFullInfo() string {
	return "Contracted Grafica (ID: " + string(rune(c.GraficaID)) + ") at " + c.Endereco
}