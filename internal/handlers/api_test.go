package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/miloradbozic/packing-service/internal/database"
	"github.com/miloradbozic/packing-service/internal/models"
	"github.com/miloradbozic/packing-service/internal/service"
)

// Mock repository for testing
type mockPackSizeRepository struct {
	packSizes []database.PackSize
	nextID    int
}

func (m *mockPackSizeRepository) GetAll() ([]database.PackSize, error) {
	return m.packSizes, nil
}

func (m *mockPackSizeRepository) GetByID(id int) (*database.PackSize, error) {
	for _, ps := range m.packSizes {
		if ps.ID == id {
			return &ps, nil
		}
	}
	return nil, fmt.Errorf("pack size with id %d not found", id)
}

func (m *mockPackSizeRepository) Create(size int) (*database.PackSize, error) {
	m.nextID++
	newPack := database.PackSize{
		ID:        m.nextID,
		Size:      size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.packSizes = append(m.packSizes, newPack)
	return &newPack, nil
}

func (m *mockPackSizeRepository) Update(id int, size int) (*database.PackSize, error) {
	for i, ps := range m.packSizes {
		if ps.ID == id {
			m.packSizes[i].Size = size
			m.packSizes[i].UpdatedAt = time.Now()
			return &m.packSizes[i], nil
		}
	}
	return nil, fmt.Errorf("pack size with id %d not found", id)
}

func (m *mockPackSizeRepository) Delete(id int) error {
	for i, ps := range m.packSizes {
		if ps.ID == id {
			m.packSizes = append(m.packSizes[:i], m.packSizes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("pack size with id %d not found", id)
}


func setupTestHandler() *APIHandler {
	mockRepo := &mockPackSizeRepository{
		packSizes: []database.PackSize{
			{ID: 1, Size: 250, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, Size: 500, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 3, Size: 1000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		nextID: 3,
	}
	packingService := service.NewPackingService(mockRepo)
	return NewAPIHandler(packingService, mockRepo)
}

func setupTestHandlerWithPackSizes(packSizes []database.PackSize) *APIHandler {
	mockRepo := &mockPackSizeRepository{
		packSizes: packSizes,
		nextID:    len(packSizes),
	}
	packingService := service.NewPackingService(mockRepo)
	return NewAPIHandler(packingService, mockRepo)
}

// Helper function to create a request with mux variables
func createRequestWithVars(method, url string, body *bytes.Buffer, vars map[string]string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, body)
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	
	// Add mux variables to the request context
	req = mux.SetURLVars(req, vars)
	return req
}

func TestAPIHandler_Calculate(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		requestBody    models.CalculateRequest
		requestBodyRaw string
		expectedStatus int
		expectedItems  int
		expectedPacks  int
		checkResponse  bool
	}{
		{
			name:           "Valid calculation - 1 item",
			requestBody:    models.CalculateRequest{Items: 1},
			expectedStatus: http.StatusOK,
			expectedItems:  250,
			expectedPacks:  1,
			checkResponse:  true,
		},
		{
			name:           "Valid calculation - 501 items",
			requestBody:    models.CalculateRequest{Items: 501},
			expectedStatus: http.StatusOK,
			expectedItems:  750,
			expectedPacks:  2,
			checkResponse:  true,
		},
		{
			name:           "Valid calculation - large order",
			requestBody:    models.CalculateRequest{Items: 12001},
			expectedStatus: http.StatusOK,
			expectedItems:  12250,
			expectedPacks:  13,
			checkResponse:  true,
		},
		{
			name:           "Invalid calculation - negative items",
			requestBody:    models.CalculateRequest{Items: -1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid calculation - zero items",
			requestBody:    models.CalculateRequest{Items: 0},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON body",
			requestBodyRaw: `{invalid json}`,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:           "Missing items field",
			requestBodyRaw: `{}`,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:           "Invalid items type",
			requestBodyRaw: `{"items": "not a number"}`,
			expectedStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.requestBodyRaw != "" {
				body = []byte(tt.requestBodyRaw)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Calculate(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse && tt.expectedStatus == http.StatusOK {
				var response models.CalculateResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				// Only check expected values if they're set (> 0)
				if tt.expectedItems > 0 && response.TotalItems != tt.expectedItems {
					t.Errorf("expected total items %d, got %d", tt.expectedItems, response.TotalItems)
				}

				if tt.expectedPacks > 0 && response.TotalPacks != tt.expectedPacks {
					t.Errorf("expected total packs %d, got %d", tt.expectedPacks, response.TotalPacks)
				}

				// Verify response structure
				if response.Items != tt.requestBody.Items {
					t.Errorf("expected items ordered %d, got %d", tt.requestBody.Items, response.Items)
				}

				if response.ExcessItems != response.TotalItems-response.Items {
					t.Errorf("excess items mismatch: expected %d, got %d", response.TotalItems-response.Items, response.ExcessItems)
				}

				if len(response.Packs) == 0 {
					t.Error("expected packs in response, got empty")
				}

				// Verify total items >= items ordered
				if response.TotalItems < response.Items {
					t.Errorf("total items %d should be >= items ordered %d", response.TotalItems, response.Items)
				}
			}
		})
	}

}

func TestAPIHandler_Calculate_EmptyPackSizes(t *testing.T) {
	handler := setupTestHandlerWithPackSizes([]database.PackSize{})
	body, _ := json.Marshal(models.CalculateRequest{Items: 100})
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Calculate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Verify error message
	var response models.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	expectedError := "no pack sizes configured"
	if response.Error != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, response.Error)
	}
}

func TestAPIHandler_GetConfig(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/api/v1/config", nil)
	w := httptest.NewRecorder()

	handler.GetConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.ConfigResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedSizes := []int{250, 500, 1000}
	if len(response.PackSizes) != len(expectedSizes) {
		t.Errorf("expected %d pack sizes, got %d", len(expectedSizes), len(response.PackSizes))
	}

	for i, expectedSize := range expectedSizes {
		if response.PackSizes[i] != expectedSize {
			t.Errorf("expected pack size %d at index %d, got %d", expectedSize, i, response.PackSizes[i])
		}
	}
}

func TestAPIHandler_ListPackSizes(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/api/v1/pack-sizes", nil)
	w := httptest.NewRecorder()

	handler.ListPackSizes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PackSizeListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.PackSizes) != 3 {
		t.Errorf("expected 3 pack sizes, got %d", len(response.PackSizes))
	}

	expectedSizes := []int{250, 500, 1000}
	for i, expectedSize := range expectedSizes {
		if response.PackSizes[i].Size != expectedSize {
			t.Errorf("expected pack size %d at index %d, got %d", expectedSize, i, response.PackSizes[i].Size)
		}
	}
}

func TestAPIHandler_GetPackSize(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		packID         string
		expectedStatus int
		expectedSize   int
	}{
		{
			name:           "Valid pack size ID",
			packID:         "1",
			expectedStatus: http.StatusOK,
			expectedSize:   250,
		},
		{
			name:           "Invalid pack size ID",
			packID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Non-numeric pack size ID",
			packID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := map[string]string{"id": tt.packID}
			req := createRequestWithVars("GET", fmt.Sprintf("/api/v1/pack-sizes/%s", tt.packID), nil, vars)
			w := httptest.NewRecorder()

			handler.GetPackSize(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.PackSizeResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if response.Size != tt.expectedSize {
					t.Errorf("expected pack size %d, got %d", tt.expectedSize, response.Size)
				}
			}
		})
	}
}

func TestAPIHandler_CreatePackSize(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		requestBody    models.CreatePackSizeRequest
		expectedStatus int
		expectedSize   int
	}{
		{
			name:           "Valid pack size creation",
			requestBody:    models.CreatePackSizeRequest{Size: 750},
			expectedStatus: http.StatusCreated,
			expectedSize:   750,
		},
		{
			name:           "Invalid pack size - negative",
			requestBody:    models.CreatePackSizeRequest{Size: -1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid pack size - zero",
			requestBody:    models.CreatePackSizeRequest{Size: 0},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/pack-sizes", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreatePackSize(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response models.PackSizeResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if response.Size != tt.expectedSize {
					t.Errorf("expected pack size %d, got %d", tt.expectedSize, response.Size)
				}
			}
		})
	}
}

func TestAPIHandler_UpdatePackSize(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		packID         string
		requestBody    models.UpdatePackSizeRequest
		expectedStatus int
		expectedSize   int
	}{
		{
			name:           "Valid pack size update",
			packID:         "1",
			requestBody:    models.UpdatePackSizeRequest{Size: 300},
			expectedStatus: http.StatusOK,
			expectedSize:   300,
		},
		{
			name:           "Invalid pack size ID",
			packID:         "999",
			requestBody:    models.UpdatePackSizeRequest{Size: 300},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid pack size - negative",
			packID:         "1",
			requestBody:    models.UpdatePackSizeRequest{Size: -1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-numeric pack size ID",
			packID:         "invalid",
			requestBody:    models.UpdatePackSizeRequest{Size: 300},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			vars := map[string]string{"id": tt.packID}
			req := createRequestWithVars("PUT", fmt.Sprintf("/api/v1/pack-sizes/%s", tt.packID), bytes.NewBuffer(body), vars)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdatePackSize(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.PackSizeResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if response.Size != tt.expectedSize {
					t.Errorf("expected pack size %d, got %d", tt.expectedSize, response.Size)
				}
			}
		})
	}
}

func TestAPIHandler_DeletePackSize(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		packID         string
		expectedStatus int
	}{
		{
			name:           "Valid pack size deletion",
			packID:         "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Invalid pack size ID",
			packID:         "999",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-numeric pack size ID",
			packID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := map[string]string{"id": tt.packID}
			req := createRequestWithVars("DELETE", fmt.Sprintf("/api/v1/pack-sizes/%s", tt.packID), nil, vars)
			w := httptest.NewRecorder()

			handler.DeletePackSize(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAPIHandler_sendError(t *testing.T) {
	handler := setupTestHandler()

	w := httptest.NewRecorder()

	handler.sendError(w, "Test error", http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response models.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Error != "Test error" {
		t.Errorf("expected error message 'Test error', got '%s'", response.Error)
	}
}

func TestAPIHandler_sendJSON(t *testing.T) {
	handler := setupTestHandler()

	w := httptest.NewRecorder()

	testData := map[string]string{"message": "test"}
	handler.sendJSON(w, testData, http.StatusOK)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["message"] != "test" {
		t.Errorf("expected message 'test', got '%s'", response["message"])
	}
}
