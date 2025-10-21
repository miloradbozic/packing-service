package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/miloradbozic/packing-service/internal/models"
	"github.com/miloradbozic/packing-service/internal/service"
)

type APIHandler struct {
	service *service.PackingService
}

func NewAPIHandler(packingService *service.PackingService) *APIHandler {
	return &APIHandler{
		service: packingService,
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
	response := models.ConfigResponse{
		PackSizes: h.service.GetPackSizes(),
	}
	h.sendJSON(w, response, http.StatusOK)
}

func (h *APIHandler) sendError(w http.ResponseWriter, message string, status int) {
	response := models.ErrorResponse{
		Error: message,
	}
	h.sendJSON(w, response, status)
}

func (h *APIHandler) sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
