package model

import (
	"time"
)

type Oferta struct {
	Id                 int64      `json:"id_oferta"`
	Nome               string     `json:"nome"`
	DataCriacao        time.Time  `json:"data_criacao"`
	DataInicio         *time.Time `json:"data_inicio"`
	DataFim            *time.Time `json:"data_fim"`
	ValorFixo          *float64   `json:"valor_fixo"`
	PercentualDesconto *int       `json:"percentual_desconto"`
}

type OfertaCreate struct {
	Nome               string     `json:"nome"`
	DataInicio         *time.Time `json:"data_inicio"`
	DataFim            *time.Time `json:"data_fim"`
	ValorFixo          *float64   `json:"valor_fixo"`
	PercentualDesconto *int       `json:"percentual_desconto"`
}

func (oc OfertaCreate) ToOferta() Oferta {
	return Oferta{
		Nome:               oc.Nome,
		DataInicio:         oc.DataInicio,
		DataFim:            oc.DataFim,
		ValorFixo:          oc.ValorFixo,
		PercentualDesconto: oc.PercentualDesconto,
	}
}
