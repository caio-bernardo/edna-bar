package domain

import "time"

// Imprime represents the relationship between books and printing companies
// This tracks which books are printed by which printing companies
type Imprime struct {
	LISBN       string    `json:"lisbn"`
	GraficaID   int       `json:"grafica_id"`
	NtoCopias   int       `json:"nto_copias"`
	DataEntrega time.Time `json:"data_entrega"`
}

// NewImprime creates a new Imprime relationship
func NewImprime(lisbn string, graficaID int, ntoCopias int, dataEntrega time.Time) *Imprime {
	return &Imprime{
		LISBN:       lisbn,
		GraficaID:   graficaID,
		NtoCopias:   ntoCopias,
		DataEntrega: dataEntrega,
	}
}

// ValidateISBN checks if the ISBN is valid (basic validation)
func (i *Imprime) ValidateISBN() bool {
	return len(i.LISBN) > 0 && len(i.LISBN) <= 20
}

// ValidateNtoCopias checks if the number of copies is valid
func (i *Imprime) ValidateNtoCopias() bool {
	return i.NtoCopias > 0
}

// ValidateGraficaID checks if the printing company ID is valid
func (i *Imprime) ValidateGraficaID() bool {
	return i.GraficaID > 0
}

// ValidateDataEntrega checks if the delivery date is valid
func (i *Imprime) ValidateDataEntrega() bool {
	return !i.DataEntrega.IsZero()
}

// IsValid checks if all fields are valid
func (i *Imprime) IsValid() bool {
	return i.ValidateISBN() && 
		   i.ValidateGraficaID() && 
		   i.ValidateNtoCopias() && 
		   i.ValidateDataEntrega()
}

// IsDeliveryOverdue checks if the delivery date has passed
func (i *Imprime) IsDeliveryOverdue() bool {
	return time.Now().After(i.DataEntrega)
}

// DaysUntilDelivery returns the number of days until delivery
func (i *Imprime) DaysUntilDelivery() int {
	duration := time.Until(i.DataEntrega)
	return int(duration.Hours() / 24)
}