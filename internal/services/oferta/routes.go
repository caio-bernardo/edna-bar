package oferta

import (
	"context"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)

type Handler struct {
	store OfertaStore
}

type OfertaStore interface {
	GetAll(ctx context.Context, filter util.Filter) ([]model.Oferta, error)
	Create(ctx context.Context, props *model.Oferta) error
	GetByID(ctx context.Context, id int64) (*model.Oferta, error)
	Update(ctx context.Context, props *model.Oferta) error
	Delete(ctx context.Context, id int64) (*model.Oferta, error)
}

func NewHandler(store OfertaStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /ofertas", h.getAll)
	mux.HandleFunc("POST /ofertas", h.create)
	mux.HandleFunc("GET /ofertas/{id}", h.fetch)
	mux.HandleFunc("PUT /ofertas/{id}", h.update)
	mux.HandleFunc("DELETE /ofertas/{id}", h.delete)
}

// @Summary List Ofertas
// @Tags Oferta
// @Produce json
// @Param filter-nome query string false "Filter by nome using operators: like, ilike, eq, ne. Format: operator.value (e.g. like.João)"
// @Param filter-cnpj query string false "Filter by cnpj using operators: eq, ne, like, ilike. Format: operator.value (e.g. eq.123456789)"
// @Param sort query string false "Sort fields: nome, cnpj. Prefix with '-' for desc. Comma separated for multiple fields (e.g. -nome,cnpj)"
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 10)"
// @Success 200 {array} model.Oferta
// @Failure 500 {object} types.ErrorResponse
// @Router /ofertas [get]
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filters, err := NewOfertaFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ofertas, err := h.store.GetAll(ctx, filters)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = util.WriteJSON(w, http.StatusOK, ofertas)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Create Oferta
// @Tags Oferta
// @Accept json
// @Produce json
// @Param fornecedor body model.OfertaCreate true "Oferta payload"
// @Success 201 {object} model.Oferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /ofertas [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No body in the request", http.StatusBadRequest)
		return
	}

	var payload model.OfertaCreate
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToOferta()
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

// @Summary Get Oferta by ID
// @Tags Oferta
// @Produce json
// @Param id path int true "Oferta ID"
// @Success 200 {object} model.Oferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /ofertas/{id} [get]
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
			util.ErrorJSON(w, "Oferta not found.", http.StatusNotFound)
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

// @Summary Update Oferta
// @Tags Oferta
// @Accept json
// @Produce json
// @Param id path int true "Oferta ID"
// @Param fornecedor body model.OfertaCreate true "Oferta payload"
// @Success 200 {object} model.Oferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /ofertas/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.OfertaCreate
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToOferta()
	model.Id = id
	err = h.store.Update(ctx, &model)
	if err != nil {
		if err == types.ErrNotFound {
			util.ErrorJSON(w, "Oferta not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}

// @Summary Delete Oferta
// @Tags Oferta
// @Produce json
// @Param id path int true "Oferta ID"
// @Success 200 {object} model.Oferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /ofertas/{id} [delete]
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
			util.ErrorJSON(w, "Oferta not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}
