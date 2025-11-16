package produto

import (
	"context"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"net/http"
)

type Handler struct {
	store ProdutoStore
}

type ProdutoStore interface {
	GetAll(ctx context.Context, filter *util.Filter) ([]model.UnionProduto, error)
	GetAllComercial(ctx context.Context, filter *util.Filter) ([]model.Comercial, error)
	GetAllEstrutural(ctx context.Context, filter *util.Filter) ([]model.Produto, error)
	CreateComercial(ctx context.Context, props *model.Comercial) error
	Create(ctx context.Context, props *model.Produto) error
	UpdateComercial(ctx context.Context, props *model.Comercial) error
	Update(ctx context.Context, props *model.Produto) error
	GetComercialByID(ctx context.Context, id int64) (*model.Comercial, error)
	GetByID(ctx context.Context, id int64) (*model.Produto, error)
	GetQntByID(ctx context.Context, id int64) (uint64, error)
	Delete(ctx context.Context, id int64) error
}

func NewHandler(store ProdutoStore) Handler {
	return Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /produtos", h.getAll)
	mux.HandleFunc("POST /produtos", h.createEstruturalHandler)
	mux.HandleFunc("GET /produtos/{id}", h.getEstruturalHandler)
	mux.HandleFunc("PUT /produtos/{id}", h.updateEstruturalHandler)
	mux.HandleFunc("DELETE /produtos/{id}", h.deleteProdutoHandler)

	mux.HandleFunc("GET /produtos/estrutural", h.getAllEstruturalHandler)
	mux.HandleFunc("GET /produtos/comercial", h.getAllComercialHandler)
	mux.HandleFunc("POST /produtos/comercial", h.createComercialHandler)
	mux.HandleFunc("GET /produtos/comercial/{id}", h.getComercialHandler)
	mux.HandleFunc("PUT /produtos/comercial/{id}", h.updateComercialHandler)

	mux.HandleFunc("GET /produtos/{id}/quantidades", h.getQuantidadeHandler)
}

 // @Summary List Produtos (all types)
 // @Tags Produtos
 // @Produce json
 // @Param filter-nome query string false "Filter by nome. Format: <op>.<value>. Ops: like, ilike, eq, ne"
 // @Param filter-categoria query string false "Filter by categoria. Format: <op>.<value>. Ops: like, ilike, eq, ne"
 // @Param filter-marca query string false "Filter by marca. Format: <op>.<value>. Ops: like, ilike, eq, ne"
 // @Param sort query string false "Sort by attribute. Allowed: nome, categoria, marca. Prefix '-' for desc. Comma separated"
 // @Param offset query int false "Pagination offset (default 0)"
 // @Param limit query int false "Pagination limit (default 0)"
 // @Success 200 {array} model.UnionProduto
 // @Failure 400 {object} types.ErrorResponse
 // @Failure 500 {object} types.ErrorResponse
 // @Router /produtos [get]
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), util.RequestTimeout)
	defer cancel()

	// WARN: Não é possivel acessar atributos do comercial
	filter, err := NewProdutoFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	produtos, err := h.store.GetAll(ctx, &filter)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, http.StatusOK, produtos)
}

// @Summary List Comercial products
// @Tags Produtos
// @Produce json
// @Param filter-nome query string false "Filter by nome. Format: <op>.<value>. Ops: like, ilike, eq, ne"
// @Param filter-categoria query string false "Filter by categoria. Format: <op>.<value>. Ops: like, ilike, eq, ne"
// @Param filter-marca query string false "Filter by marca. Format: <op>.<value>. Ops: like, ilike, eq, ne"
// @Param filter-preco_venda query number false "Filter by preco_venda. Format: <op>.<value>. Ops: eq, ne, lt, gt, le, ge"
// @Param sort query string false "Sort fields: nome, categoria, marca, preco_venda. Prefix '-' for desc. Comma separated"
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 0)"
// @Success 200 {array} model.Comercial
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/comercial [get]
func (h *Handler) getAllComercialHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filter, err := NewComercialFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	produtos, err := h.store.GetAllComercial(ctx, &filter)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, produtos); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

// @Summary List Estrutural products
// @Tags Produtos
// @Produce json
// @Param filter-nome query string false "Filter by nome. Format: <op>.<value>. Ops: like, ilike, eq, ne"
// @Param filter-categoria query string false "Filter by categoria. Format: <op>.<value>. Ops: like, ilike, eq, ne"
// @Param filter-marca query string false "Filter by marca. Format: <op>.<value>. Ops: like, ilike, eq, ne"
// @Param sort query string false "Sort fields: nome, categoria, marca. Prefix '-' for desc. Comma separated"
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 0)"
// @Success 200 {array} model.Produto
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/estrutural [get]
func (h *Handler) getAllEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filter, err := NewProdutoFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	produtos, err := h.store.GetAllEstrutural(ctx, &filter)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, produtos); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

// @Summary Create Comercial Produto
// @Tags Produtos
// @Accept json
// @Produce json
// @Param produto body model.ComercialCreate true "Comercial product payload"
// @Success 201 {object} model.Comercial
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/comercial [post]
func (h *Handler) createComercialHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	payload := model.ComercialCreate{}
	if err := util.ReadJSON(r, &payload); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	produto := payload.ToComercial()
	if err := h.store.CreateComercial(ctx, &produto); err != nil {
		status := http.StatusInternalServerError
		if err == types.ErrNotFound {
			status = http.StatusNotFound
		}
		util.ErrorJSON(w, err.Error(), status)
		return
	}

	if err := util.WriteJSON(w, http.StatusCreated, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Create Produto
// @Tags Produtos
// @Accept json
// @Produce json
// @Param produto body model.ProdutoCreate true "Product payload"
// @Success 201 {object} model.Produto
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos [post]
func (h *Handler) createEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	payload := model.ProdutoCreate{}
	if err := util.ReadJSON(r, &payload); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	produto := payload.ToProduto()
	if err := h.store.Create(ctx, &produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusCreated, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Update Comercial Produto
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "Produto ID"
// @Param produto body model.ComercialCreate true "Comercial product payload"
// @Success 200 {object} model.Comercial
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/comercial/{id} [put]
func (h *Handler) updateComercialHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	payload := model.ComercialCreate{}
	if err := util.ReadJSON(r, &payload); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	produto := payload.ToComercial()
	produto.Id = id
	if err := h.store.UpdateComercial(ctx, &produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Update Produto
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "Produto ID"
// @Param produto body model.ProdutoCreate true "Product payload"
// @Success 200 {object} model.Produto
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/{id} [put]
func (h *Handler) updateEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	payload := model.ProdutoCreate{}
	if err := util.ReadJSON(r, &payload); err != nil {
		util.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	produto := payload.ToProduto()
	produto.Id = id
	if err := h.store.Update(ctx, &produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, produto); err != nil {
		util.ErrorJSON(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get Comercial Produto by ID
// @Tags Produtos
// @Produce json
// @Param id path int true "Produto ID"
// @Success 200 {object} model.Comercial
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/comercial/{id} [get]
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

// @Summary Get Produto by ID
// @Tags Produtos
// @Produce json
// @Param id path int true "Produto ID"
// @Success 200 {object} model.Produto
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/{id} [get]
func (h *Handler) getEstruturalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	produto, err := h.store.GetByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, produto); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

// @Summary Delete Produto
// @Tags Produtos
// @Param id path int true "Produto ID"
// @Success 204 {string} string
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/{id} [delete]
func (h *Handler) deleteProdutoHandler(w http.ResponseWriter, r *http.Request) {
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

// @Summary Get Produto Quantidade
// @Tags Produtos
// @Produce json
// @Param id path int true "Produto ID"
// @Success 200 {object} uint64
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /produtos/{id}/quantidades [get]
func (h *Handler) getQuantidadeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	qnt, err := h.store.GetQntByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, qnt); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
