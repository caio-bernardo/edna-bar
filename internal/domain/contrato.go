package domain

import "math"

// Contrato represents a contract with a contracted printing company
type Contrato struct {
	ID               int     `json:"id"`
	Valor            float64 `json:"valor"`
	NomeResponsavel  string  `json:"nome_responsavel"`
	GraficaContID    int     `json:"grafica_cont_id"`
}

// ValidateValor checks if the contract value is valid
func (c *Contrato) ValidateValor() bool {
	return c.Valor > 0 && !math.IsInf(c.Valor, 0) && !math.IsNaN(c.Valor)
}

// ValidateNomeResponsavel checks if the responsible person's name is valid
func (c *Contrato) ValidateNomeResponsavel() bool {
	return len(c.NomeResponsavel) > 0 && len(c.NomeResponsavel) <= 255
}

// IsValid checks if all contract fields are valid
func (c *Contrato) IsValid() bool {
	return c.ValidateValor() && c.ValidateNomeResponsavel() && c.GraficaContID > 0
}