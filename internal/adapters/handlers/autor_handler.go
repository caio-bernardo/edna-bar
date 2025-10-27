package handlers

import (
	"context"
	"edna/internal/applications"
	"encoding/json"
	"net/http"
	"time"
)

type AutorHandler struct {
	autorUsecase applications.AutorUsecase
}

func NewAutorHandler(usecase applications.AutorUsecase) *AutorHandler {
	return &AutorHandler{
		autorUsecase: usecase,
	}
}

// RegisterRoutes registers all autor routes
func (h *AutorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /autores", h.handleGetAllAutores)
	mux.HandleFunc("POST /autores", h.handleCreateAutor)
	mux.HandleFunc("GET /autores/{rg}", h.handleGetAutorByRG)
	mux.HandleFunc("GET /autores/{rg}/books", h.handleGetAutorWithBooks)
	mux.HandleFunc("PUT /autores/{rg}", h.handleUpdateAutor)
	mux.HandleFunc("DELETE /autores/{rg}", h.handleDeleteAutor)
	mux.HandleFunc("GET /autores/search/name", h.handleGetAutoresByName)
	mux.HandleFunc("POST /autores/{rg}/books/{isbn}", h.handleAddAutorToBook)
	mux.HandleFunc("DELETE /autores/{rg}/books/{isbn}", h.handleRemoveAutorFromBook)
}

// GET /autores
func (h *AutorHandler) handleGetAllAutores(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	autores, err := h.autorUsecase.List(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to list authors", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, autores, http.StatusOK)
}

// POST /autores
func (h *AutorHandler) handleCreateAutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req applications.CreateAutorRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	autor, err := h.autorUsecase.Create(ctx, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create author", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, autor, http.StatusCreated)
}

// GET /autores/{rg}
func (h *AutorHandler) handleGetAutorByRG(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rg := r.PathValue("rg")
	if rg == "" {
		h.sendErrorResponse(w, "RG is required", http.StatusBadRequest, nil)
		return
	}

	autor, err := h.autorUsecase.Get(ctx, rg)
	if err != nil {
		h.sendErrorResponse(w, "Author not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, autor, http.StatusOK)
}

// GET /autores/{rg}/books
func (h *AutorHandler) handleGetAutorWithBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rg := r.PathValue("rg")
	if rg == "" {
		h.sendErrorResponse(w, "RG is required", http.StatusBadRequest, nil)
		return
	}

	autor, err := h.autorUsecase.GetWithBooks(ctx, rg)
	if err != nil {
		h.sendErrorResponse(w, "Author not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, autor, http.StatusOK)
}

// PUT /autores/{rg}
func (h *AutorHandler) handleUpdateAutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rg := r.PathValue("rg")
	if rg == "" {
		h.sendErrorResponse(w, "RG is required", http.StatusBadRequest, nil)
		return
	}

	var req applications.UpdateAutorRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	autor, err := h.autorUsecase.Update(ctx, rg, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to update author", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, autor, http.StatusOK)
}

// DELETE /autores/{rg}
func (h *AutorHandler) handleDeleteAutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rg := r.PathValue("rg")
	if rg == "" {
		h.sendErrorResponse(w, "RG is required", http.StatusBadRequest, nil)
		return
	}

	if err := h.autorUsecase.Delete(ctx, rg); err != nil {
		h.sendErrorResponse(w, "Failed to delete author", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Author deleted successfully"}, http.StatusOK)
}

// GET /autores/search/name?name=John
func (h *AutorHandler) handleGetAutoresByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")
	if name == "" {
		h.sendErrorResponse(w, "Name parameter is required", http.StatusBadRequest, nil)
		return
	}

	autores, err := h.autorUsecase.GetByName(ctx, name)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get authors by name", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, autores, http.StatusOK)
}

// POST /autores/{rg}/books/{isbn}
func (h *AutorHandler) handleAddAutorToBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rg := r.PathValue("rg")
	isbn := r.PathValue("isbn")

	if rg == "" || isbn == "" {
		h.sendErrorResponse(w, "RG and ISBN are required", http.StatusBadRequest, nil)
		return
	}

	if err := h.autorUsecase.AddToBook(ctx, rg, isbn); err != nil {
		h.sendErrorResponse(w, "Failed to add author to book", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Author added to book successfully"}, http.StatusOK)
}

// DELETE /autores/{rg}/books/{isbn}
func (h *AutorHandler) handleRemoveAutorFromBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rg := r.PathValue("rg")
	isbn := r.PathValue("isbn")

	if rg == "" || isbn == "" {
		h.sendErrorResponse(w, "RG and ISBN are required", http.StatusBadRequest, nil)
		return
	}

	if err := h.autorUsecase.RemoveFromBook(ctx, rg, isbn); err != nil {
		h.sendErrorResponse(w, "Failed to remove author from book", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Author removed from book successfully"}, http.StatusOK)
}

// Helper methods

func (h *AutorHandler) decodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func (h *AutorHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *AutorHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResponse := map[string]interface{}{
		"error": message,
		"code":  statusCode,
	}
	
	if err != nil {
		errorResponse["details"] = err.Error()
	}
	
	json.NewEncoder(w).Encode(errorResponse)
}