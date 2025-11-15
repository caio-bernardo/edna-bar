package cliente

import (
	"edna/internal/util"
	"net/url"
)

func NewClienteFilter(params url.Values) (util.Filter, error) {
	var filter util.Filter

	if err := filter.GetOffset(params); err != nil {
		return filter, err
	}

	if err := filter.GetLimit(params); err != nil {
		return filter, err
	}

	attrs := []string{"nome", "cpf", "data_nascimento"}
	if err := filter.GetSorts(params, attrs); err != nil {
		return filter, err
	}

	// Filtro de string (like, ilike, eq, ne)
	if err := filter.GetFilterStr(params, "nome"); err != nil {
		return filter, err
	}
	// Filtro de string (like, ilike, eq, ne)
	if err := filter.GetFilterStr(params, "cpf"); err != nil {
		return filter, err
	}

	// Filtro tempo (eq, neq, gt, lt, ...)
	if err := filter.GetFilterTime(params, "data_nascimento"); err != nil {
		return filter, err
	}

	return filter, nil
}
