package handlers

import (
	"context"
	"edna/internal/applications"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type GraficaHandler struct {
	graficaUsecase applications.GraficaUsecase
}

func NewGraficaHandler(usecase applications.GraficaUsecase) *GraficaHandler {
	return &GraficaHandler{
		graficaUsecase: usecase,
	}
}

// RegisterRoutes registers all grafica routes
func (h *GraficaHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /graficas", h.handleGetAllGraficas)
	mux.HandleFunc("POST /graficas", h.handleCreateGrafica)
	mux.HandleFunc("GET /graficas/{id}", h.handleGetGraficaByID)
	mux.HandleFunc("GET /graficas/{id}/contracts", h.handleGetGraficaWithContracts)
	mux.HandleFunc("GET /graficas/{id}/jobs", h.handleGetGraficaWithPrintingJobs)
	mux.HandleFunc("PUT /graficas/{id}", h.handleUpdateGrafica)
	mux.HandleFunc("DELETE /graficas/{id}", h.handleDeleteGrafica)
	mux.HandleFunc("GET /graficas/search/name", h.handleGetGraficasByName)
	mux.HandleFunc("GET /graficas/search/type", h.handleGetGraficasByType)
}

// GET /graficas
func (h *GraficaHandler) handleGetAllGraficas(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	graficas, err := h.graficaUsecase.List(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to list printing companies", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, graficas, http.StatusOK)
}

// POST /graficas
func (h *GraficaHandler) handleCreateGrafica(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req applications.CreateGraficaRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	grafica, err := h.graficaUsecase.Create(ctx, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create printing company", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, grafica, http.StatusCreated)
}

// GET /graficas/{id}
func (h *GraficaHandler) handleGetGraficaByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	grafica, err := h.graficaUsecase.Get(ctx, id)
	if err != nil {
		h.sendErrorResponse(w, "Printing company not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, grafica, http.StatusOK)
}

// GET /graficas/{id}/contracts
func (h *GraficaHandler) handleGetGraficaWithContracts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	grafica, err := h.graficaUsecase.GetWithContracts(ctx, id)
	if err != nil {
		h.sendErrorResponse(w, "Printing company not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, grafica, http.StatusOK)
}

// GET /graficas/{id}/jobs
func (h *GraficaHandler) handleGetGraficaWithPrintingJobs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	grafica, err := h.graficaUsecase.GetWithPrintingJobs(ctx, id)
	if err != nil {
		h.sendErrorResponse(w, "Printing company not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, grafica, http.StatusOK)
}

// PUT /graficas/{id}
func (h *GraficaHandler) handleUpdateGrafica(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	var req applications.UpdateGraficaRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	grafica, err := h.graficaUsecase.Update(ctx, id, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to update printing company", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, grafica, http.StatusOK)
}

// DELETE /graficas/{id}
func (h *GraficaHandler) handleDeleteGrafica(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	if err := h.graficaUsecase.Delete(ctx, id); err != nil {
		h.sendErrorResponse(w, "Failed to delete printing company", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Printing company deleted successfully"}, http.StatusOK)
}

// GET /graficas/search/name?name=ABC
func (h *GraficaHandler) handleGetGraficasByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")
	if name == "" {
		h.sendErrorResponse(w, "Name parameter is required", http.StatusBadRequest, nil)
		return
	}

	graficas, err := h.graficaUsecase.GetByName(ctx, name)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get printing companies by name", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, graficas, http.StatusOK)
}

// GET /graficas/search/type?type=particular
func (h *GraficaHandler) handleGetGraficasByType(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	graficaType := r.URL.Query().Get("type")
	if graficaType == "" {
		h.sendErrorResponse(w, "Type parameter is required", http.StatusBadRequest, nil)
		return
	}

	if graficaType != "particular" && graficaType != "contratada" {
		h.sendErrorResponse(w, "Type must be 'particular' or 'contratada'", http.StatusBadRequest, nil)
		return
	}

	graficas, err := h.graficaUsecase.GetByType(ctx, graficaType)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get printing companies by type", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, graficas, http.StatusOK)
}

// Helper methods

func (h *GraficaHandler) decodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func (h *GraficaHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *GraficaHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
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