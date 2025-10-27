package handlers

import (
	"context"
	"edna/internal/applications"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type EditoraHandler struct {
	editoraUsecase applications.EditoraUsecase
}

func NewEditoraHandler(usecase applications.EditoraUsecase) *EditoraHandler {
	return &EditoraHandler{
		editoraUsecase: usecase,
	}
}

// RegisterRoutes registers all editora routes
func (h *EditoraHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /editoras", h.handleGetAllEditoras)
	mux.HandleFunc("POST /editoras", h.handleCreateEditora)
	mux.HandleFunc("GET /editoras/{id}", h.handleGetEditoraByID)
	mux.HandleFunc("GET /editoras/{id}/books", h.handleGetEditoraWithBooks)
	mux.HandleFunc("PUT /editoras/{id}", h.handleUpdateEditora)
	mux.HandleFunc("DELETE /editoras/{id}", h.handleDeleteEditora)
	mux.HandleFunc("GET /editoras/search/name", h.handleGetEditorasByName)
}

// GET /editoras
func (h *EditoraHandler) handleGetAllEditoras(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	editoras, err := h.editoraUsecase.List(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to list publishers", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, editoras, http.StatusOK)
}

// POST /editoras
func (h *EditoraHandler) handleCreateEditora(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req applications.CreateEditoraRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	editora, err := h.editoraUsecase.Create(ctx, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create publisher", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, editora, http.StatusCreated)
}

// GET /editoras/{id}
func (h *EditoraHandler) handleGetEditoraByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	editora, err := h.editoraUsecase.Get(ctx, id)
	if err != nil {
		h.sendErrorResponse(w, "Publisher not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, editora, http.StatusOK)
}

// GET /editoras/{id}/books
func (h *EditoraHandler) handleGetEditoraWithBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	editora, err := h.editoraUsecase.GetWithBooks(ctx, id)
	if err != nil {
		h.sendErrorResponse(w, "Publisher not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, editora, http.StatusOK)
}

// PUT /editoras/{id}
func (h *EditoraHandler) handleUpdateEditora(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	var req applications.UpdateEditoraRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	editora, err := h.editoraUsecase.Update(ctx, id, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to update publisher", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, editora, http.StatusOK)
}

// DELETE /editoras/{id}
func (h *EditoraHandler) handleDeleteEditora(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	if err := h.editoraUsecase.Delete(ctx, id); err != nil {
		h.sendErrorResponse(w, "Failed to delete publisher", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Publisher deleted successfully"}, http.StatusOK)
}

// GET /editoras/search/name?name=Penguin
func (h *EditoraHandler) handleGetEditorasByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")
	if name == "" {
		h.sendErrorResponse(w, "Name parameter is required", http.StatusBadRequest, nil)
		return
	}

	editoras, err := h.editoraUsecase.GetByName(ctx, name)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get publishers by name", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, editoras, http.StatusOK)
}

// Helper methods

func (h *EditoraHandler) decodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func (h *EditoraHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *EditoraHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
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