package lote

import (
	"edna/internal/util"
	"net/url"
)

func NewLoteFilter(params url.Values) (util.Filter, error) {
	var filter util.Filter
	if err := filter.GetOffset(params); err != nil {
		return filter, err
	}
	if err := filter.GetLimit(params); err != nil {
		return filter, err
	}

	attrs := []string{"id_fornecedor", "id_produto", "preco_unitario", "estragados", "quantidade_inicial"}
	if err := filter.GetSorts(params, attrs); err != nil {
		return filter, err
	}

	if err := filter.GetFilterInt(params, "id_fornecedor"); err != nil {
		return filter, err
	}
	if err := filter.GetFilterInt(params, "id_produto"); err != nil {
		return filter, err
	}
	if err := filter.GetFilterFloat(params, "preco_unitario"); err != nil {
		return filter, err
	}
	if err := filter.GetFilterInt(params, "estragados"); err != nil {
		return filter, err
	}
	if err := filter.GetFilterInt(params, "quantidade_inicial"); err != nil {
		return filter, err
	}

	return filter, nil
}
