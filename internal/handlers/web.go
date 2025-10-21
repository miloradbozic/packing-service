package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"github.com/miloradbozic/packing-service/internal/database"
	"github.com/miloradbozic/packing-service/internal/models"
	"github.com/miloradbozic/packing-service/internal/service"
)

type WebHandler struct {
	service      *service.PackingService
	packSizeRepo *database.PackSizeRepository
	templates    *template.Template
}

func NewWebHandler(packingService *service.PackingService, packSizeRepo *database.PackSizeRepository) (*WebHandler, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &WebHandler{
		service:      packingService,
		packSizeRepo: packSizeRepo,
		templates:    tmpl,
	}, nil
}

func (h *WebHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.handleCalculate(w, r)
		return
	}

	packSizes, err := h.service.GetPackSizes()
	if err != nil {
		http.Error(w, "Failed to get pack sizes", http.StatusInternalServerError)
		return
	}

	data := struct {
		PackSizes []int
		Results   *models.CalculateResponse
		Error     string
		Items     string
	}{
		PackSizes: packSizes,
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func (h *WebHandler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	itemsStr := r.FormValue("items")
	items, err := strconv.Atoi(itemsStr)

	packSizes, packErr := h.service.GetPackSizes()
	if packErr != nil {
		http.Error(w, "Failed to get pack sizes", http.StatusInternalServerError)
		return
	}

	data := struct {
		PackSizes []int
		Results   *models.CalculateResponse
		Error     string
		Items     string
	}{
		PackSizes: packSizes,
		Items:     itemsStr,
	}

	if err != nil || items <= 0 {
		data.Error = "Please enter a valid positive number"
		h.templates.ExecuteTemplate(w, "index.html", data)
		return
	}

	solution, err := h.service.CalculatePacks(items)
	if err != nil {
		data.Error = err.Error()
		h.templates.ExecuteTemplate(w, "index.html", data)
		return
	}

	// Convert solution to display format
	packs := make([]models.Pack, 0)
	for size, qty := range solution.Packs {
		if qty > 0 {
			packs = append(packs, models.Pack{
				Size:     size,
				Quantity: qty,
			})
		}
	}

	// Sort packs by size (largest first)
	sort.Slice(packs, func(i, j int) bool {
		return packs[i].Size > packs[j].Size
	})

	data.Results = &models.CalculateResponse{
		Items:       items,
		TotalItems:  solution.TotalItems,
		TotalPacks:  solution.TotalPacks,
		Packs:       packs,
		ExcessItems: solution.TotalItems - items,
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
