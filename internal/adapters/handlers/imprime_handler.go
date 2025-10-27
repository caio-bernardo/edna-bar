package handlers

import (
	"context"
	"edna/internal/applications"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type ImprimeHandler struct {
	imprimeUsecase applications.ImprimeUsecase
}

func NewImprimeHandler(usecase applications.ImprimeUsecase) *ImprimeHandler {
	return &ImprimeHandler{
		imprimeUsecase: usecase,
	}
}

// RegisterRoutes registers all imprime routes
func (h *ImprimeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /printing-jobs", h.handleGetAllImprimes)
	mux.HandleFunc("POST /printing-jobs", h.handleCreateImprime)
	mux.HandleFunc("GET /printing-jobs/{isbn}/{graficaId}", h.handleGetImprimeByKey)
	mux.HandleFunc("PUT /printing-jobs/{isbn}/{graficaId}", h.handleUpdateImprime)
	mux.HandleFunc("DELETE /printing-jobs/{isbn}/{graficaId}", h.handleDeleteImprime)
	mux.HandleFunc("GET /printing-jobs/book/{isbn}", h.handleGetImprimesByISBN)
	mux.HandleFunc("GET /printing-jobs/grafica/{graficaId}", h.handleGetImprimesByGrafica)
	mux.HandleFunc("GET /printing-jobs/search/date-range", h.handleGetImprimesByDateRange)
	mux.HandleFunc("GET /printing-jobs/overdue", h.handleGetOverdueJobs)
	mux.HandleFunc("GET /printing-jobs/pending", h.handleGetPendingJobs)
	mux.HandleFunc("GET /printing-jobs/statistics", h.handleGetPrintingStatistics)
	mux.HandleFunc("POST /printing-jobs/{isbn}/{graficaId}/complete", h.handleMarkAsCompleted)
}

// GET /printing-jobs
func (h *ImprimeHandler) handleGetAllImprimes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	imprimes, err := h.imprimeUsecase.List(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to list printing jobs", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, imprimes, http.StatusOK)
}

// POST /printing-jobs
func (h *ImprimeHandler) handleCreateImprime(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req applications.CreateImprimeRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	imprime, err := h.imprimeUsecase.Create(ctx, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create printing job", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, imprime, http.StatusCreated)
}

// GET /printing-jobs/{isbn}/{graficaId}
func (h *ImprimeHandler) handleGetImprimeByKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	graficaIdStr := r.PathValue("graficaId")

	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	graficaId, err := strconv.Atoi(graficaIdStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid grafica ID format", http.StatusBadRequest, err)
		return
	}

	imprime, err := h.imprimeUsecase.Get(ctx, isbn, graficaId)
	if err != nil {
		h.sendErrorResponse(w, "Printing job not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, imprime, http.StatusOK)
}

// PUT /printing-jobs/{isbn}/{graficaId}
func (h *ImprimeHandler) handleUpdateImprime(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	graficaIdStr := r.PathValue("graficaId")

	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	graficaId, err := strconv.Atoi(graficaIdStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid grafica ID format", http.StatusBadRequest, err)
		return
	}

	var req applications.UpdateImprimeRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	imprime, err := h.imprimeUsecase.Update(ctx, isbn, graficaId, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to update printing job", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, imprime, http.StatusOK)
}

// DELETE /printing-jobs/{isbn}/{graficaId}
func (h *ImprimeHandler) handleDeleteImprime(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	graficaIdStr := r.PathValue("graficaId")

	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	graficaId, err := strconv.Atoi(graficaIdStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid grafica ID format", http.StatusBadRequest, err)
		return
	}

	if err := h.imprimeUsecase.Delete(ctx, isbn, graficaId); err != nil {
		h.sendErrorResponse(w, "Failed to delete printing job", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Printing job deleted successfully"}, http.StatusOK)
}

// GET /printing-jobs/book/{isbn}
func (h *ImprimeHandler) handleGetImprimesByISBN(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	imprimes, err := h.imprimeUsecase.GetByISBN(ctx, isbn)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get printing jobs by ISBN", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, imprimes, http.StatusOK)
}

// GET /printing-jobs/grafica/{graficaId}
func (h *ImprimeHandler) handleGetImprimesByGrafica(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	graficaIdStr := r.PathValue("graficaId")
	graficaId, err := strconv.Atoi(graficaIdStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid grafica ID format", http.StatusBadRequest, err)
		return
	}

	imprimes, err := h.imprimeUsecase.GetByGrafica(ctx, graficaId)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get printing jobs by grafica", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, imprimes, http.StatusOK)
}

// GET /printing-jobs/search/date-range?start=2020-01-01&end=2023-12-31
func (h *ImprimeHandler) handleGetImprimesByDateRange(w http.ResponseWriter, r *http.Request) {
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

	imprimes, err := h.imprimeUsecase.GetByDeliveryDateRange(ctx, start, end)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get printing jobs by date range", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, imprimes, http.StatusOK)
}

// GET /printing-jobs/overdue
func (h *ImprimeHandler) handleGetOverdueJobs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	imprimes, err := h.imprimeUsecase.GetOverdueJobs(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get overdue printing jobs", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, imprimes, http.StatusOK)
}

// GET /printing-jobs/pending
func (h *ImprimeHandler) handleGetPendingJobs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	imprimes, err := h.imprimeUsecase.GetPendingJobs(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get pending printing jobs", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, imprimes, http.StatusOK)
}

// GET /printing-jobs/statistics?start=2020-01-01&end=2023-12-31
func (h *ImprimeHandler) handleGetPrintingStatistics(w http.ResponseWriter, r *http.Request) {
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

	statistics, err := h.imprimeUsecase.GetPrintingStatistics(ctx, start, end)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get printing statistics", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, statistics, http.StatusOK)
}

// POST /printing-jobs/{isbn}/{graficaId}/complete
func (h *ImprimeHandler) handleMarkAsCompleted(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isbn := r.PathValue("isbn")
	graficaIdStr := r.PathValue("graficaId")

	if isbn == "" {
		h.sendErrorResponse(w, "ISBN is required", http.StatusBadRequest, nil)
		return
	}

	graficaId, err := strconv.Atoi(graficaIdStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid grafica ID format", http.StatusBadRequest, err)
		return
	}

	if err := h.imprimeUsecase.MarkAsCompleted(ctx, isbn, graficaId); err != nil {
		h.sendErrorResponse(w, "Failed to mark printing job as completed", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Printing job marked as completed successfully"}, http.StatusOK)
}

// Helper methods

func (h *ImprimeHandler) decodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func (h *ImprimeHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ImprimeHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
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