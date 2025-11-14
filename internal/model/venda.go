package model

import (
	"time"
)

type Venda struct {
	Id                int64     `json:"id"`
	DataHoraVenda     time.Time `json:"dataHoraVenda"`
	DataHoraPagamento time.Time `json:"dataHoraPagamento"`
	TipoPagamento     string    `json:"tipoPagamento"`
}

type VendaCreate struct {
	DataHoraVenda     time.Time `json:"dataHoraVenda"`
	DataHoraPagamento time.Time `json:"dataHoraPagamento"`
	TipoPagamento     string    `json:"tipoPagamento"`
}

func (vc *VendaCreate) ToVenda() Venda {
	return Venda{
		DataHoraVenda:     vc.DataHoraVenda,
		DataHoraPagamento: vc.DataHoraPagamento,
		TipoPagamento:     vc.TipoPagamento,
	}
}
