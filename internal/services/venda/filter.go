package venda

import (
	"edna/internal/util"
	"net/url"
)

func NewVendaFilter(params url.Values) (util.Filter, error) {
	var filter util.Filter
	if err := filter.GetOffset(params); err != nil {
		return filter, err
	}

	if err := filter.GetLimit(params); err != nil {
		return filter, err
	}

	attrs := []string{"data_hora_venda", "data_hora_pagamento", "tipo_pagamento", "id_cliente", "id_funcionario"}

	if err := filter.GetSorts(params, attrs); err != nil {
		return filter, err
	}

	if err := filter.GetFilterStr(params, "tipo_pagamento"); err != nil {
		return filter, err
	}

	for _, attr := range []string{"id_cliente", "id_funcionario"} {
		if err := filter.GetFilterInt(params, attr); err != nil {
			return filter, err
		}
	}

	for _, attr := range []string{"data_hora_venda", "data_hora_pagamento"} {
		if err := filter.GetFilterTime(params, attr); err != nil {
			return filter, err
		}
	}
	return filter, nil
}
