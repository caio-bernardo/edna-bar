package venda

import (
	"context"
	"edna/internal/model"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)

type Handler struct {
	store VendaStore
}

type VendaStore interface {
	GetAll(ctx context.Context, filter util.Filter) ([]model.Venda, error)
	Create(ctx context.Context, props *model.Venda) error
	GetByID(ctx context.Context, id int64) (*model.Venda, error)
	Update(ctx context.Context, props *model.Venda) error
	Delete(ctx context.Context, id int64) (*model.Venda, error)
}

func NewHandler(store VendaStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /vendas", h.getAll)
	mux.HandleFunc("POST /vendas", h.create)
	mux.HandleFunc("GET /vendas/{id}", h.fetch)
	mux.HandleFunc("PUT /vendas/{id}", h.update)
	mux.HandleFunc("DELETE /vendas/{id}", h.delete)
}

// @Summary List Vendas
// @Tags Venda
// @Produce json
// @Param filter-tipoPagamento query string false "Filter by tipoPagamento using operators: eq, ne, like, ilike"
// @Param filter-dataHoraVenda query string false "Filter by dataHoraVenda using operators: eq, ne, gt, lt"
// @Param sort query string false "Sort fields: dataHoraVenda, dataHoraPagamento, tipoPagamento. Prefix with '-' for desc."
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 10)"
// @Success 200 {array} model.Venda
// @Failure 500 {object} types.ErrorResponse
// @Router /vendas [get]
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filters, err := NewVendaFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vendas, err := h.store.GetAll(ctx, filters)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = util.WriteJSON(w, http.StatusOK, vendas)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Create Venda
// @Tags Venda
// @Accept json
// @Produce json
// @Param venda body model.VendaCreate true "Venda payload"
// @Success 201 {object} model.Venda
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /vendas [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No body in the request", http.StatusBadRequest)
		return
	}

	var payload model.VendaCreate
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToVenda()
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

// @Summary Get Venda by ID
// @Tags Venda
// @Produce json
// @Param id path int true "Venda ID"
// @Success 200 {object} model.Venda
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /vendas/{id} [get]
func (h *Handler) fetch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	venda, err := h.store.GetByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if venda == nil {
		util.ErrorJSON(w, "Venda not found.", http.StatusNotFound)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, venda); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update Venda
// @Tags Venda
// @Accept json
// @Produce json
// @Param id path int true "Venda ID"
// @Param venda body model.VendaCreate true "Venda payload"
// @Success 200 {object} model.Venda
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /vendas/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.VendaCreate
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToVenda()
	model.Id = id
	err = h.store.Update(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}

// @Summary Delete Venda
// @Tags Venda
// @Produce json
// @Param id path int true "Venda ID"
// @Success 200 {object} model.Venda
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /vendas/{id} [delete]
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
