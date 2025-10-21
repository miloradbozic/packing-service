package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/miloradbozic/packing-service/internal/database"
	"github.com/miloradbozic/packing-service/internal/models"
	"github.com/miloradbozic/packing-service/internal/service"
)

type APIHandler struct {
	service      *service.PackingService
	packSizeRepo database.PackSizeRepositoryInterface
}

func NewAPIHandler(packingService *service.PackingService, packSizeRepo database.PackSizeRepositoryInterface) *APIHandler {
	return &APIHandler{
		service:      packingService,
		packSizeRepo: packSizeRepo,
	}
}

func (h *APIHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	var req models.CalculateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	solution, err := h.service.CalculatePacks(req.Items)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert solution to response format
	packs := make([]models.Pack, 0)
	for size, qty := range solution.Packs {
		if qty > 0 {
			packs = append(packs, models.Pack{
				Size:     size,
				Quantity: qty,
			})
		}
	}

	response := models.CalculateResponse{
		Items:       req.Items,
		TotalItems:  solution.TotalItems,
		TotalPacks:  solution.TotalPacks,
		Packs:       packs,
		ExcessItems: solution.TotalItems - req.Items,
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *APIHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	packSizes, err := h.service.GetPackSizes()
	if err != nil {
		h.sendError(w, "Failed to get pack sizes", http.StatusInternalServerError)
		return
	}

	response := models.ConfigResponse{
		PackSizes: packSizes,
	}
	h.sendJSON(w, response, http.StatusOK)
}

func (h *APIHandler) sendError(w http.ResponseWriter, message string, status int) {
	response := models.ErrorResponse{
		Error: message,
	}
	h.sendJSON(w, response, status)
}

// Pack size management endpoints

func (h *APIHandler) ListPackSizes(w http.ResponseWriter, r *http.Request) {
	packSizes, err := h.packSizeRepo.GetAll()
	if err != nil {
		h.sendError(w, "Failed to get pack sizes", http.StatusInternalServerError)
		return
	}

	response := models.PackSizeListResponse{
		PackSizes: make([]models.PackSizeResponse, len(packSizes)),
	}

	for i, ps := range packSizes {
		response.PackSizes[i] = models.PackSizeResponse{
			ID:        ps.ID,
			Size:      ps.Size,
			CreatedAt: ps.CreatedAt.Format(time.RFC3339),
			UpdatedAt: ps.UpdatedAt.Format(time.RFC3339),
		}
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *APIHandler) GetPackSize(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, "Invalid pack size ID", http.StatusBadRequest)
		return
	}

	packSize, err := h.packSizeRepo.GetByID(id)
	if err != nil {
		h.sendError(w, "Pack size not found", http.StatusNotFound)
		return
	}

	response := models.PackSizeResponse{
		ID:        packSize.ID,
		Size:      packSize.Size,
		CreatedAt: packSize.CreatedAt.Format(time.RFC3339),
		UpdatedAt: packSize.UpdatedAt.Format(time.RFC3339),
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *APIHandler) CreatePackSize(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePackSizeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Size <= 0 {
		h.sendError(w, "Pack size must be positive", http.StatusBadRequest)
		return
	}

	packSize, err := h.packSizeRepo.Create(req.Size)
	if err != nil {
		h.sendError(w, fmt.Sprintf("Failed to create pack size: %v", err), http.StatusBadRequest)
		return
	}

	response := models.PackSizeResponse{
		ID:        packSize.ID,
		Size:      packSize.Size,
		CreatedAt: packSize.CreatedAt.Format(time.RFC3339),
		UpdatedAt: packSize.UpdatedAt.Format(time.RFC3339),
	}

	h.sendJSON(w, response, http.StatusCreated)
}

func (h *APIHandler) UpdatePackSize(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, "Invalid pack size ID", http.StatusBadRequest)
		return
	}

	var req models.UpdatePackSizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Size <= 0 {
		h.sendError(w, "Pack size must be positive", http.StatusBadRequest)
		return
	}

	packSize, err := h.packSizeRepo.Update(id, req.Size)
	if err != nil {
		h.sendError(w, fmt.Sprintf("Failed to update pack size: %v", err), http.StatusBadRequest)
		return
	}

	response := models.PackSizeResponse{
		ID:        packSize.ID,
		Size:      packSize.Size,
		CreatedAt: packSize.CreatedAt.Format(time.RFC3339),
		UpdatedAt: packSize.UpdatedAt.Format(time.RFC3339),
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *APIHandler) DeletePackSize(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, "Invalid pack size ID", http.StatusBadRequest)
		return
	}

	if err := h.packSizeRepo.Delete(id); err != nil {
		h.sendError(w, fmt.Sprintf("Failed to delete pack size: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
