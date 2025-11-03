package fornecedor

import (
	"context"
	"edna/internal/model"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)


type Handler struct {
	store FornecedorStore
}

type FornecedorStore interface {
	GetAll(ctx context.Context) ([]model.Fornecedor, error)
	Create(ctx context.Context, props *model.Fornecedor) error
	GetByID(ctx context.Context, id int64) (*model.Fornecedor, error)
	Update(ctx context.Context, props *model.Fornecedor) error
	Delete(ctx context.Context, id int64) error
}


func NewHandler(store FornecedorStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /fornecedores", h.getAllFornecedoresHandler)
	mux.HandleFunc("POST /fornecedores", h.createFornecedoresHandler)
	mux.HandleFunc("GET /fornecedores/{id}", h.fetchFornecedorHandler)
	mux.HandleFunc("PUT /fornecedores/{id}", h.updateFornecedorHandler)
	mux.HandleFunc("DELETE /fornecedores/{id}", h.deleteFornecedorHandler)
}


func (h *Handler) getAllFornecedoresHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	fornecedores, err := h.store.GetAll(ctx)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = util.WriteJSON(w, http.StatusOK, fornecedores)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) createFornecedoresHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No body in the request", http.StatusBadRequest)
		return
	}

	var payload model.FornecedorPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := model.FromPayload(payload)
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

func (h *Handler) fetchFornecedorHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	fornecedor, err := h.store.GetByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if fornecedor == nil {
		util.ErrorJSON(w, "Fornecedor not found.", http.StatusNotFound)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, fornecedor); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateFornecedorHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.FornecedorPayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := model.FromPayload(payload)
	model.Id = id
	err = h.store.Update(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}

func (h *Handler) deleteFornecedorHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.store.Delete(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, nil)
}
