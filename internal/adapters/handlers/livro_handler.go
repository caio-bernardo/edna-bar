package handlers

import (
	"context"
	"edna/internal/applications"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type LivroHandler struct {
	livroUsecase applications.LivroUsecase
}

func NewLivroHandler(usecase applications.LivroUsecase) *LivroHandler {
	return &LivroHandler{
		livroUsecase: usecase,
	}
}

// RegisterRoutes registers all livro routes
func (h *LivroHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /livros", h.handleGetAllLivros)
	mux.HandleFunc("POST /livros", h.handleCreateLivro)
	mux.HandleFunc("GET /livros/{isbn}", h.handleGetLivroByISBN)
	mux.HandleFunc("GET /livros/{isbn}/authors", h.handleGetLivroWithAuthors)
	mux.HandleFunc("PUT /livros/{isbn}", h.handleUpdateLivro)
	mux.HandleFunc("DELETE /livros/{isbn}", h.handleDeleteLivro)
	// mux.HandleFunc("GET /livros/editora/{editoraId}", h.handleGetLivrosByEditora)
	mux.HandleFunc("GET /livros/search/date-range", h.handleGetLivrosByDateRange)
	mux.HandleFunc("POST /livros/{isbn}/authors/{authorRG}", h.handleAddAuthorToLivro)
	mux.HandleFunc("DELETE /livros/{isbn}/authors/{authorRG}", h.handleRemoveAuthorFromLivro)
}

// GET /livros
func (h *LivroHandler) handleGetAllLivros(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	livros, err := h.livroUsecase.List(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to list books", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, livros, http.StatusOK)
}

// POST /livros
func (h *LivroHandler) handleCreateLivro(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req applications.CreateLivroRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	livro, err := h.livroUsecase.Create(ctx, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create book", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, livro, http.StatusCreated)
}

// GET /livros/{isbn}
func (h *LivroHandler) handleGetLivroByISBN(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	livro, err := h.livroUsecase.Get(ctx, isbn)
	if err != nil {
		h.sendErrorResponse(w, "Book not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, livro, http.StatusOK)
}

// GET /livros/{isbn}/authors
func (h *LivroHandler) handleGetLivroWithAuthors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	livro, err := h.livroUsecase.GetWithAuthors(ctx, isbn)
	if err != nil {
		h.sendErrorResponse(w, "Book not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, livro, http.StatusOK)
}

// PUT /livros/{isbn}
func (h *LivroHandler) handleUpdateLivro(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	var req applications.UpdateLivroRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	livro, err := h.livroUsecase.Update(ctx, isbn, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to update book", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, livro, http.StatusOK)
}

// DELETE /livros/{isbn}
func (h *LivroHandler) handleDeleteLivro(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	if err := h.livroUsecase.Delete(ctx, isbn); err != nil {
		h.sendErrorResponse(w, "Failed to delete book", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Book deleted successfully"}, http.StatusOK)
}

// GET /livros/editora/{editoraId}
func (h *LivroHandler) handleGetLivrosByEditora(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	editoraIdStr := r.PathValue("editoraId")
	editoraId, err := strconv.Atoi(editoraIdStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid editora ID", http.StatusBadRequest, err)
		return
	}

	livros, err := h.livroUsecase.GetByEditora(ctx, editoraId)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get books by editora", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, livros, http.StatusOK)
}

// GET /livros/search/date-range?start=2020-01-01&end=2023-12-31
func (h *LivroHandler) handleGetLivrosByDateRange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		h.sendErrorResponse(w, "Start and end dates are required", http.StatusBadRequest, nil)
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid start date format (expected YYYY-MM-DD)", http.StatusBadRequest, err)
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid end date format (expected YYYY-MM-DD)", http.StatusBadRequest, err)
		return
	}

	livros, err := h.livroUsecase.GetByPublicationDateRange(ctx, start, end)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get books by date range", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, livros, http.StatusOK)
}

// POST /livros/{isbn}/authors/{authorRG}
func (h *LivroHandler) handleAddAuthorToLivro(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	authorRG := r.PathValue("authorRG")

	if isbn == "" || authorRG == "" {
		h.sendErrorResponse(w, "ISBN and author RG are required", http.StatusBadRequest, nil)
		return
	}

	if err := h.livroUsecase.AddAuthor(ctx, isbn, authorRG); err != nil {
		h.sendErrorResponse(w, "Failed to add author to book", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Author added to book successfully"}, http.StatusOK)
}

// DELETE /livros/{isbn}/authors/{authorRG}
func (h *LivroHandler) handleRemoveAuthorFromLivro(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	authorRG := r.PathValue("authorRG")

	if isbn == "" || authorRG == "" {
		h.sendErrorResponse(w, "ISBN and author RG are required", http.StatusBadRequest, nil)
		return
	}

	if err := h.livroUsecase.RemoveAuthor(ctx, isbn, authorRG); err != nil {
		h.sendErrorResponse(w, "Failed to remove author from book", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Author removed from book successfully"}, http.StatusOK)
}

// Helper methods

func (h *LivroHandler) decodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func (h *LivroHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *LivroHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
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
