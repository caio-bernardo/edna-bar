package produto

import (
	"context"
	"edna/internal/model"
	"edna/internal/util"
	"net/http"
)

type Handler struct {
	store ProdutoStore
}

type ProdutoStore interface {
	GetAllComercial(ctx context.Context) ([]model.Comercial, error)
	GetAllEstrutural(ctx context.Context) ([]model.Estrutural, error)
	CreateComercial(ctx context.Context, props *model.Comercial) error
	CreateEstrutural(ctx context.Context, props *model.Estrutural) error
	UpdateComercial(ctx context.Context, id int64, props *model.Comercial) error
	UpdateEstrutural(ctx context.Context, id int64, props *model.Estrutural) error
	GetComercialByID(ctx context.Context, id int64) (*model.Produto, error)
	GetEstruturalByID(ctx context.Context, id int64) (*model.Produto, error)
	Delete(ctx context.Context, id int64) error
}

func NewHandler(store ProdutoStore) Handler {
	return Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /produtos/comercial", h.getAllComercialHandler)
	mux.HandleFunc("GET /produtos/estrutural", h.getAllEstruturalHandler)
	mux.HandleFunc("POST /produtos/comercial", h.createComercialHandler)
	mux.HandleFunc("POST /produtos/estrutural", h.createEstruturalHandler)
	mux.HandleFunc("PUT /produtos/comercial/{id}", h.updateComercialHandler)
	mux.HandleFunc("PUT /produtos/estrutural/{id}", h.updateEstruturalHandler)
	mux.HandleFunc("GET /produtos/comercial/{id}", h.getComercialHandler)
	mux.HandleFunc("GET /produtos/estrutural/{id}", h.getEstruturalHandler)
	mux.HandleFunc("DELETE /produtos/{id}", h.deleleteProdutoHandler)
}

func (h *Handler) getAllComercialHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	produtos, err := h.store.GetAllComercial(ctx)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}

	if err = util.WriteJSON(w, http.StatusOK, produtos); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h *Handler) getAllEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	produtos, err := h.store.GetAllEstrutural(ctx)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}

	if err = util.WriteJSON(w, http.StatusOK, produtos); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h *Handler) createComercialHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	produto := model.Comercial{}
	if err := util.ReadJSON(r, &produto); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := h.store.CreateComercial(ctx, &produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusCreated, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) createEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	produto := model.Estrutural{}
	if err := util.ReadJSON(r, &produto); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := h.store.CreateEstrutural(ctx, &produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusCreated, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) updateComercialHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	produto := model.Comercial{}
	if err := util.ReadJSON(r, &produto); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateComercial(ctx, id, &produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) updateEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	produto := model.Estrutural{}
	if err := util.ReadJSON(r, &produto); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateEstrutural(ctx, id, &produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getComercialHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	produto, err := h.store.GetComercialByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h *Handler) getEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	produto, err := h.store.GetEstruturalByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h *Handler) deleleteProdutoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	if err := h.store.Delete(ctx, id); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
