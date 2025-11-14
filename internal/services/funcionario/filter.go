package funcionario

import (
	"edna/internal/util"
	"net/url"
)

func NewFuncionarioFilter(params url.Values) (util.Filter, error) {
	var filter util.Filter

	if err := filter.GetOffset(params); err != nil {
		return filter, err
	}

	if err := filter.GetLimit(params); err != nil {
		return filter, err
	}

	attrs := []string{"data_contratacao", "salario", "expediente", "tipo", "CPF", "nome", "id_funcionario"}

	if err := filter.GetSorts(params, attrs); err != nil {
		return filter, err
	}

	for _, attr := range attrs {
		if err := filter.GetFilterStr(params, attr); err != nil {
			return filter, err
		}
	}
	return filter, nil
}
