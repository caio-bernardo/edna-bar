package item_oferta

import (
	"context"
	"edna/internal/model"
	"edna/internal/util"
	"encoding/json"
	"net/http"
)

type Handler struct {
	store ItemOfertaStore
}

type ItemOfertaStore interface {
	GetAll(ctx context.Context, filter util.Filter) ([]model.ItemOferta, error)
	GetItemByID(ctx context.Context, id int64) (*model.ItemOferta, error)
	GetOfertaByID(ctx context.Context, id int64) (*model.ItemOferta, error)
	GetByComposedID(ctx context.Context, id_produto int64, id_oferta int64) (*model.ItemOferta, error)
	Create(ctx context.Context, props *model.ItemOferta) error
	Update(ctx context.Context, props *model.ItemOferta) error
	Delete(ctx context.Context, id_produto int64, id_oferta int64) (*model.ItemOferta, error)
}

func NewHandler(store ItemOfertaStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /item_ofertas", h.getAll)
	mux.HandleFunc("POST /item_ofertas", h.create)
	mux.HandleFunc("GET /item_ofertas/{id_produto}/{id_oferta}", h.fetch)
	mux.HandleFunc("PUT /item_ofertas/{id_produto}/{id_oferta}", h.update)
	mux.HandleFunc("DELETE /item_ofertas/{id_produto}/{id_oferta}", h.delete)
	mux.HandleFunc("GET /item_ofertas/item/{id}", h.getItemByID)
	mux.HandleFunc("GET /item_ofertas/oferta/{id}", h.getOfertaByID)
}

// @Summary List Item Ofertas
// @Tags Item Oferta
// @Produce json
// @Param filter-nome query string false "Filter by nome using operators: like, ilike, eq, ne. Format: operator.value (e.g. like.João)"
// @Param filter-cnpj query string false "Filter by cnpj using operators: eq, ne, like, ilike. Format: operator.value (e.g. eq.123456789)"
// @Param sort query string false "Sort fields: nome, cnpj. Prefix with '-' for desc. Comma separated for multiple fields (e.g. -nome,cnpj)"
// @Param offset query int false "Pagination offset (default 0)"
// @Param limit query int false "Pagination limit (default 10)"
// @Success 200 {array} model.ItemOferta
// @Failure 500 {object} types.ErrorResponse
// @Router /item_ofertas [get]
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	filters, err := NewItemOfertaFilter(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	itemOfertas, err := h.store.GetAll(ctx, filters)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = util.WriteJSON(w, http.StatusOK, itemOfertas)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get ItemOferta by Item ID
// @Tags Item Oferta
// @Produce json
// @Param id path int true "Item (Produto) ID"
// @Success 200 {object} model.ItemOferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /item_ofertas/item/{id} [get]
func (h *Handler) getItemByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	itens, err := h.store.GetItemByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if itens == nil {
		util.ErrorJSON(w, "ItemOferta not found for this item id.", http.StatusNotFound)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, itens)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get ItemOferta by Oferta ID
// @Tags Item Oferta
// @Produce json
// @Param id path int true "Oferta ID"
// @Success 200 {object} model.ItemOferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /item_ofertas/oferta/{id} [get]
func (h *Handler) getOfertaByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id, err := util.GetIDParam(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	itens, err := h.store.GetOfertaByID(ctx, id)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if itens == nil {
		util.ErrorJSON(w, "ItemOferta not found for this oferta id.", http.StatusNotFound)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, itens)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Create ItemOferta
// @Tags ItemOferta
// @Accept json
// @Produce json
// @Param itemOferta body model.ItemOfertaCreate true "ItemOferta payload"
// @Success 201 {object} model.ItemOferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /itemOfertas [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	if r.Body == nil {
		util.ErrorJSON(w, "No body in the request", http.StatusBadRequest)
		return
	}

	var payload model.ItemOfertaCreate
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToItemOferta()
	err = h.store.Create(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusCreated, model)
}

// @Summary Get ItemOferta by composed ID
// @Tags ItemOferta
// @Produce json
// @Param id_produto path int true "Produto ID"
// @Param id_oferta path int true "Oferta ID"
// @Success 200 {object} model.ItemOferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /item_ofertas/{id_produto}/{id_oferta} [get]
func (h *Handler) fetch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	// Use a mesma função do seu handler de update
	id_produto, id_oferta, err := util.GetComposedID(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Chame o novo método do store
	itemOferta, err := h.store.GetByComposedID(ctx, id_produto, id_oferta)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if itemOferta == nil {
		util.ErrorJSON(w, "ItemOferta not found.", http.StatusNotFound)
		return
	}

	if err = util.WriteJSON(w, http.StatusOK, itemOferta); err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update ItemOferta
// @Tags ItemOferta
// @Accept json
// @Produce json
// @Param id path int true "ItemOferta ID"
// @Param item body model.ItemOfertaCreate true "ItemOferta payload"
// @Success 200 {object} model.ItemOferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /item_ofertas/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	id_produto, id_oferta, err := util.GetComposedID(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload model.ItemOfertaCreate
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := payload.ToItemOferta()
	model.IDProduto = id_produto
	model.IDOferta = id_oferta
	err = h.store.Update(ctx, &model)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}

// @Summary Delete ItemOferta
// @Tags Item Oferta
// @Produce json
// @Param id_produto path int true "Produto ID"
// @Param id_oferta path int true "Oferta ID"
// @Success 200 {object} model.ItemOferta
// @Failure 400 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
// @Router /item_ofertas/{id_produto}/{id_oferta} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.RequestTimeout)
	defer cancel()

	// Use o helper para obter os dois IDs
	id_produto, id_oferta, err := util.GetComposedID(r)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Chame o método Delete com os dois IDs
	model, err := h.store.Delete(ctx, id_produto, id_oferta)
	if err != nil {
		util.ErrorJSON(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	util.WriteJSON(w, http.StatusOK, model)
}
