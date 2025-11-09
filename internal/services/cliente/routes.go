package cliente

import (
	"context"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)

type Handler struct {
	store ClienteStore
}

type ClienteStore interface {
	GetAll(ctx context.Context, filter util.Filter) ([]model.Cliente, error)
	Create(ctx context.Context, props *model.Cliente) error
	GetByID(ctx context.Context, id int64) (*model.Cliente, error)
	Update(ctx context.Context, props *model.Cliente) error
	Delete(ctx context.Context, id int64) (*model.Cliente, error)
}

func NewHandler(store ClienteStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /clientes", h.getAll)
	mux.HandleFunc("POST /clientes", h.create)
	mux.HandleFunc("GET /clientes/{id}", h.fetch)
	mux.HandleFunc("PUT /clientes/{id}", h.update)
	mux.HandleFunc("DELETE /clientes/{id}", h.delete)
}

// @Summary List Clients
// @Tags Cliente
// @Produce json
// @Param filter-nome query string false "Filter by nome using operators: like, ilike, eq, ne. Format: operator.value (e.g. like.João)"
// @Param filter-cnpj query string false "Filter by cnpj using operators: eq, ne, like, ilike. Format: operator.value (e.g. eq.123456789)"
// @Param sort query string false "Sort fields: nome, cnpj. Prefix with '-' for desc. Comma separated for multiple fields (e.g. -nome,cnpj)"
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 10)"
// @Success 200 {array} model.Cliente
// @Failure 500 {object} types.ErrorResponse
// @Router /clientes [get]
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filters, err := NewClienteFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clientes, err := h.store.GetAll(ctx, filters)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = util.WriteJSON(w, http.StatusOK, clientes)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Create Cliente
// @Tags Cliente
// @Accept json
// @Produce json
// @Param fornecedor body model.ClienteCreate true "Cliente payload"
// @Success 201 {object} model.Cliente
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /clientes [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No body in the request", http.StatusBadRequest)
		return
	}

	var payload model.ClienteCreate
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToCliente()
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

// @Summary Get Cliente by ID
// @Tags Cliente
// @Produce json
// @Param id path int true "Cliente ID"
// @Success 200 {object} model.Cliente
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /clientes/{id} [get]
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
			util.ErrorJSON(w, "Cliente not found.", http.StatusNotFound)
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

// @Summary Update Cliente
// @Tags Cliente
// @Accept json
// @Produce json
// @Param id path int true "Cliente ID"
// @Param fornecedor body model.ClienteCreate true "Cliente payload"
// @Success 200 {object} model.Cliente
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /clientes/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.ClienteCreate
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToCliente()
	model.Id = id
	err = h.store.Update(ctx, &model)
	if err != nil {
		if err == types.ErrNotFound {
			util.ErrorJSON(w, "Cliente not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}

// @Summary Delete Cliente
// @Tags Cliente
// @Produce json
// @Param id path int true "Cliente ID"
// @Success 200 {object} model.Cliente
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /clientes/{id} [delete]
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
			util.ErrorJSON(w, "Cliente not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}
