package item_oferta

import (
	"edna/internal/util"
	"net/url"
)

func NewItemOfertaFilter(params url.Values) (util.Filter, error) {
	var filter util.Filter

	if err := filter.GetOffset(params); err != nil {
		return filter, err
	}

	if err := filter.GetLimit(params); err != nil {
		return filter, err
	}

	attrs := []string{"quantidade", "id_produto", "id_oferta"}

	if err := filter.GetSorts(params, attrs); err != nil {
		return filter, err
	}

	// Filtro de int
	if err := filter.GetFilterInt(params, "quantidade"); err != nil {
		return filter, err
	}
	// Filtro de int
	if err := filter.GetFilterInt(params, "id_produto"); err != nil {
		return filter, err
	}
	// Filtro de int
	if err := filter.GetFilterInt(params, "id_oferta"); err != nil {
		return filter, err
	}

	return filter, nil
}
