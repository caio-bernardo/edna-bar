package handlers

import (
	"context"
	"edna/internal/applications"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type ContratoHandler struct {
	contratoUsecase applications.ContratoUsecase
}

func NewContratoHandler(usecase applications.ContratoUsecase) *ContratoHandler {
	return &ContratoHandler{
		contratoUsecase: usecase,
	}
}

// RegisterRoutes registers all contrato routes
func (h *ContratoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /contratos", h.handleGetAllContratos)
	mux.HandleFunc("POST /contratos", h.handleCreateContrato)
	mux.HandleFunc("GET /contratos/{id}", h.handleGetContratoByID)
	mux.HandleFunc("PUT /contratos/{id}", h.handleUpdateContrato)
	mux.HandleFunc("DELETE /contratos/{id}", h.handleDeleteContrato)
	mux.HandleFunc("GET /contratos/grafica/{graficaId}", h.handleGetContratosByGrafica)
	mux.HandleFunc("GET /contratos/search/responsavel", h.handleGetContratosByResponsavel)
	mux.HandleFunc("GET /contratos/search/value-range", h.handleGetContratosByValueRange)
	mux.HandleFunc("GET /contratos/analysis", h.handleGetContractAnalysis)
}

// GET /contratos
func (h *ContratoHandler) handleGetAllContratos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	contratos, err := h.contratoUsecase.List(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to list contracts", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, contratos, http.StatusOK)
}

// POST /contratos
func (h *ContratoHandler) handleCreateContrato(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req applications.CreateContratoRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	contrato, err := h.contratoUsecase.Create(ctx, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create contract", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, contrato, http.StatusCreated)
}

// GET /contratos/{id}
func (h *ContratoHandler) handleGetContratoByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	contrato, err := h.contratoUsecase.Get(ctx, id)
	if err != nil {
		h.sendErrorResponse(w, "Contract not found", http.StatusNotFound, err)
		return
	}

	h.sendJSONResponse(w, contrato, http.StatusOK)
}

// PUT /contratos/{id}
func (h *ContratoHandler) handleUpdateContrato(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	var req applications.UpdateContratoRequest
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	contrato, err := h.contratoUsecase.Update(ctx, id, req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to update contract", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, contrato, http.StatusOK)
}

// DELETE /contratos/{id}
func (h *ContratoHandler) handleDeleteContrato(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, err)
		return
	}

	if err := h.contratoUsecase.Delete(ctx, id); err != nil {
		h.sendErrorResponse(w, "Failed to delete contract", http.StatusBadRequest, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Contract deleted successfully"}, http.StatusOK)
}

// GET /contratos/grafica/{graficaId}
func (h *ContratoHandler) handleGetContratosByGrafica(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	graficaIdStr := r.PathValue("graficaId")
	graficaId, err := strconv.Atoi(graficaIdStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid grafica ID format", http.StatusBadRequest, err)
		return
	}

	contratos, err := h.contratoUsecase.GetByGraficaContID(ctx, graficaId)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get contracts by grafica", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, contratos, http.StatusOK)
}

// GET /contratos/search/responsavel?name=John
func (h *ContratoHandler) handleGetContratosByResponsavel(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	responsavel := r.URL.Query().Get("name")
	if responsavel == "" {
		h.sendErrorResponse(w, "Name parameter is required", http.StatusBadRequest, nil)
		return
	}

	contratos, err := h.contratoUsecase.GetByResponsavel(ctx, responsavel)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get contracts by responsible person", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, contratos, http.StatusOK)
}

// GET /contratos/search/value-range?min=1000&max=5000
func (h *ContratoHandler) handleGetContratosByValueRange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")

	if minStr == "" || maxStr == "" {
		h.sendErrorResponse(w, "Both min and max parameters are required", http.StatusBadRequest, nil)
		return
	}

	minValue, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		h.sendErrorResponse(w, "Invalid min value format", http.StatusBadRequest, err)
		return
	}

	maxValue, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		h.sendErrorResponse(w, "Invalid max value format", http.StatusBadRequest, err)
		return
	}

	contratos, err := h.contratoUsecase.GetByValueRange(ctx, minValue, maxValue)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get contracts by value range", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, contratos, http.StatusOK)
}

// GET /contratos/analysis
func (h *ContratoHandler) handleGetContractAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	analysis, err := h.contratoUsecase.GetContractAnalysis(ctx)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get contract analysis", http.StatusInternalServerError, err)
		return
	}

	h.sendJSONResponse(w, analysis, http.StatusOK)
}

// Helper methods

func (h *ContratoHandler) decodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func (h *ContratoHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ContratoHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
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