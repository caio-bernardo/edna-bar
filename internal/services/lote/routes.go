package lote

import (
	"context"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)

type Handler struct {
	store LoteStore
}

type LoteStore interface {
	GetAll(ctx context.Context, filter util.Filter) ([]model.Lote, error)
	Create(ctx context.Context, props *model.Lote) error
	GetByID(ctx context.Context, id int64) (*model.Lote, error)
	Update(ctx context.Context, props *model.Lote) error
	Delete(ctx context.Context, id int64) (*model.Lote, error)
}

func NewHandler(store LoteStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /lotes", h.getAll)
	mux.HandleFunc("POST /lotes", h.create)
	mux.HandleFunc("GET /lotes/{id}", h.fetch)
	mux.HandleFunc("PUT /lotes/{id}", h.update)
	mux.HandleFunc("DELETE /lotes/{id}", h.delete)
}

// @Summary List Lotes
// @Tags Lote
// @Produce json
// @Param filter-nome query string false "Filter by nome using operators: like, ilike, eq, ne. Format: operator.value (e.g. like.João)"
// @Param filter-cnpj query string false "Filter by cnpj using operators: eq, ne, like, ilike. Format: operator.value (e.g. eq.123456789)"
// @Param sort query string false "Sort fields: nome, cnpj. Prefix with '-' for desc. Comma separated for multiple fields (e.g. -nome,cnpj)"
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 10)"
// @Success 200 {array} model.Lote
// @Failure 500 {object} types.ErrorResponse
// @Router /lotes [get]
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filters, err := NewLoteFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lotes, err := h.store.GetAll(ctx, filters)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = util.WriteJSON(w, http.StatusOK, lotes)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Create Lote
// @Tags Lote
// @Accept json
// @Produce json
// @Param fornecedor body model.LoteCreate true "Lote payload"
// @Success 201 {object} model.Lote
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /lotes [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No body in the request", http.StatusBadRequest)
		return
	}

	var payload model.LoteCreate
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToLote()
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

// @Summary Get Lote by ID
// @Tags Lote
// @Produce json
// @Param id path int true "Lote ID"
// @Success 200 {object} model.Lote
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /lotes/{id} [get]
func (h *Handler) fetch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	cliente, err := h.store.GetByID(ctx, id)
	if err != nil {
		if err == types.ErrNotFound {
			util.ErrorJSON(w, "Lote not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, cliente); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update Lote
// @Tags Lote
// @Accept json
// @Produce json
// @Param id path int true "Lote ID"
// @Param fornecedor body model.LoteCreate true "Lote payload"
// @Success 200 {object} model.Lote
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /lotes/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.LoteCreate
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToLote()
	model.Id = id
	err = h.store.Update(ctx, &model)
	if err != nil {
		if err == types.ErrNotFound {
			util.ErrorJSON(w, "Lote not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}

// @Summary Delete Lote
// @Tags Lote
// @Produce json
// @Param id path int true "Lote ID"
// @Success 200 {object} model.Lote
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /lotes/{id} [delete]
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
		if err == types.ErrNotFound {
			util.ErrorJSON(w, "Lote not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}
