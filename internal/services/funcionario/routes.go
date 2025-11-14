package funcionario

import (
	"context"
	"edna/internal/model"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)

type Handler struct {
	store FuncionarioStore
}

type FuncionarioStore interface {
	GetAll(ctx context.Context, filter util.Filter) ([]model.Funcionario, error)
	Create(ctx context.Context, props *model.Funcionario) error
	GetByID(ctx context.Context, id int64) (*model.Funcionario, error)
	Update(ctx context.Context, props *model.Funcionario) error
	Delete(ctx context.Context, id int64) (*model.Funcionario, error)
}

func NewHandler(store FuncionarioStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /funcionarios", h.getAll)
	mux.HandleFunc("POST /funcionarios", h.create)
	mux.HandleFunc("GET /funcionarios/{id}", h.fetch)
	mux.HandleFunc("PUT /funcionarios/{id}", h.update)
	mux.HandleFunc("DELETE /funcionarios/{id}", h.delete)
}

// @Summary List Funcionarios
// @Tags Funcionario
// @Produce json
// @Param filter-nome query string false "Filter by nome using operators: like, ilike, eq, ne. Format: operator.value (e.g. like.João)"
// @Param filter-CPF query string false "Filter by CPF using operators: eq, ne, like, ilike. Format: operator.value (e.g. eq.123456789)"
// @Param sort query string false "Sort fields: nome, CPF. Prefix with '-' for desc. Comma separated for multiple fields (e.g. -nome,CPF)"
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 10)"
// @Success 200 {array} model.Funcionario
// @Failure 500 {object} types.ErrorResponse
// @Router /funcionarios [get]
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filters, err := NewFuncionarioFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	funcionarios, err := h.store.GetAll(ctx, filters)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = util.WriteJSON(w, http.StatusOK, funcionarios)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Create Funcionario
// @Tags Funcionario
// @Accept json
// @Produce json
// @Param funcionario body model.FuncionarioCreate true "Funcionario payload"
// @Success 201 {object} model.Funcionario
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /funcionarios [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No body in the request", http.StatusBadRequest)
		return
	}

	var payload model.FuncionarioCreate
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToFuncionario()
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

// @Summary Get Funcionario by ID
// @Tags Funcionario
// @Produce json
// @Param id path int true "Funcionario ID"
// @Success 200 {object} model.Funcionario
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /funcionarios/{id} [get]
func (h *Handler) fetch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	funcionario, err := h.store.GetByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if funcionario == nil {
		util.ErrorJSON(w, "Funcionario not found.", http.StatusNotFound)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, funcionario); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update Funcionario
// @Tags Funcionario
// @Accept json
// @Produce json
// @Param id path int true "Funcionario ID"
// @Param funcionario body model.FuncionarioCreate true "Funcionario payload"
// @Success 200 {object} model.Funcionario
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /funcionarios/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.FuncionarioCreate
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToFuncionario()
	model.Id = id
	err = h.store.Update(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}

// @Summary Delete Funcionario
// @Tags Funcionario
// @Produce json
// @Param id path int true "Funcionario ID"
// @Success 200 {object} model.Funcionario
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /funcionarios/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model, err := h.store.Delete(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}
