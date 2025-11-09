package oferta

import (
	"edna/internal/util"
	"net/url"
)

func NewOfertaFilter(params url.Values) (util.Filter, error) {
	var filter util.Filter
	if err := filter.GetOffset(params); err != nil {
		return filter, err
	}
	if err := filter.GetLimit(params); err != nil {
		return filter, err
	}

	attrs := []string{"nome", "valor_fixo", "percentual_desconto"}
	if err := filter.GetSorts(params, attrs); err != nil {
		return filter, err
	}

	if err := filter.GetFilterStr(params, "nome"); err != nil {
		return filter, err
	}
	if err := filter.GetFilterFloat(params, "valor_fixo"); err != nil {
		return filter, err
	}
	if err := filter.GetFilterInt(params, "percentual_desconto"); err != nil {
		return filter, err
	}
	// Filtros para datas (data_inicio, data_fim) podem ser adicionados
	// como GetFilterStr, mas o ideal seria ter um GetFilterDate no util.Filter
	return filter, nil
}
