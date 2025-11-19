package aplica_oferta

import (
	"context"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)

type Handler struct {
	store AplicaOfertaStore
}

type AplicaOfertaStore interface {
	GetAll(ctx context.Context, filter util.Filter) ([]model.AplicaOferta, error)
	GetByID(ctx context.Context, id int64) (*model.AplicaOferta, error)
	Create(ctx context.Context, c *model.AplicaOferta) error
	Update(ctx context.Context, c *model.AplicaOferta) error
	Delete(ctx context.Context, id int64) (*model.AplicaOferta, error)
}

func NewHandler(store AplicaOfertaStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /aplica_oferta", h.getAll)
	mux.HandleFunc("POST /aplica_oferta", h.create)
	mux.HandleFunc("GET /aplica_oferta/{id}", h.fetch)
	mux.HandleFunc("PUT /aplica_oferta/{id}", h.update)
	mux.HandleFunc("DELETE /aplica_oferta/{id}", h.delete)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filters, err := NewAplicaOfertaFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aplicaOfertas, err := h.store.GetAll(ctx, filters)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, aplicaOfertas)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No request body", http.StatusBadRequest)
		return
	}

	var payload model.AplicaOfertaResponse
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToAplicaOferta()
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

func (h *Handler) fetch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	aplicaOferta, err := h.store.GetByID(ctx, id)
	if err != nil {
		if err == types.ErrNotFound {
			util.ErrorJSON(w, "Oferta not found.", http.StatusNotFound)
			return
		}
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, aplicaOferta); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.AplicaOfertaResponse
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToAplicaOferta()
	model.IDAplicaOferta = id
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
