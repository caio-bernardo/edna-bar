package aplica_oferta

import (
	"edna/internal/util"
	"net/url"
)

func NewAplicaOfertaFilter(params url.Values) (util.Filter, error) {
	var filter util.Filter
	if err := filter.GetOffset(params); err != nil {
		return filter, err
	}

	if err := filter.GetLimit(params); err != nil {
		return filter, err
	}

	attrs := []string{"id_oferta", "id_venda", "id_item_venda"}

	if err := filter.GetSorts(params, attrs); err != nil {
		return filter, err
	}

	if err := filter.GetFilterInt(params, "id_oferta"); err != nil {
		return filter, err
	}

	if err := filter.GetFilterInt(params, "id_venda"); err != nil {
		return filter, err
	}

	if err := filter.GetFilterInt(params, "id_item_venda"); err != nil {
		return filter, err
	}

	return filter, nil
}
