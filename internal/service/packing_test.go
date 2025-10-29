package service

import (
	"testing"
	"time"

	"github.com/miloradbozic/packing-service/internal/database"
)

// Mock repository for testing
type mockPackSizeRepository struct {
	sizes []int
}

func (m *mockPackSizeRepository) GetAll() ([]database.PackSize, error) {
	// Convert sizes to PackSize objects for testing
	packSizes := make([]database.PackSize, len(m.sizes))
	for i, size := range m.sizes {
		packSizes[i] = database.PackSize{
			ID:        i + 1,
			Size:      size,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return packSizes, nil
}

func (m *mockPackSizeRepository) GetByID(id int) (*database.PackSize, error) {
	return nil, nil
}

func (m *mockPackSizeRepository) Create(size int) (*database.PackSize, error) {
	return nil, nil
}

func (m *mockPackSizeRepository) Update(id int, size int) (*database.PackSize, error) {
	return nil, nil
}

func (m *mockPackSizeRepository) Delete(id int) error {
	return nil
}


func TestPackingService_CalculatePacks(t *testing.T) {
	defaultPackSizes := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name          string
		packSizes     []int
		itemsOrdered  int
		expectedPacks map[int]int
		expectedItems int
		expectError   bool
	}{
		{
			name:         "Exactly 1 item",
			packSizes:    defaultPackSizes,
			itemsOrdered: 1,
			expectedPacks: map[int]int{
				250: 1,
			},
			expectedItems: 250,
		},
		{
			name:         "Exactly 250 items",
			packSizes:    defaultPackSizes,
			itemsOrdered: 250,
			expectedPacks: map[int]int{
				250: 1,
			},
			expectedItems: 250,
		},
		{
			name:         "251 items",
			packSizes:    defaultPackSizes,
			itemsOrdered: 251,
			expectedPacks: map[int]int{
				500: 1,
			},
			expectedItems: 500,
		},
		{
			name:         "501 items",
			packSizes:    defaultPackSizes,
			itemsOrdered: 501,
			expectedPacks: map[int]int{
				500: 1,
				250: 1,
			},
			expectedItems: 750,
		},
		{
			name:         "12001 items",
			packSizes:    defaultPackSizes,
			itemsOrdered: 12001,
			expectedPacks: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
			expectedItems: 12250,
		},
		{
			name:         "Empty pack sizes",
			packSizes:    []int{},
			itemsOrdered: 100,
			expectError:  true,
		},
		{
			name:         "Zero items ordered",
			packSizes:    []int{250, 500},
			itemsOrdered: 0,
			expectError:  true,
		},
		{
			name:         "Negative items ordered",
			packSizes:    []int{250, 500},
			itemsOrdered: -1,
			expectError:  true,
		},
		{
			name:         "Negative pack size",
			packSizes:    []int{250, -100, 500},
			itemsOrdered: 100,
			expectError:  true,
		},
		{
			name:         "Zero pack size",
			packSizes:    []int{250, 0, 500},
			itemsOrdered: 100,
			expectError:  true,
		},
		{
			name:         "Single pack size - exact match",
			packSizes:    []int{500},
			itemsOrdered: 500,
			expectedPacks: map[int]int{
				500: 1,
			},
			expectedItems: 500,
		},
		{
			name:         "Single pack size - needs multiple",
			packSizes:    []int{250},
			itemsOrdered: 501,
			expectedPacks: map[int]int{
				250: 3,
			},
			expectedItems: 750,
		},
		{
			name:         "Very large pack sizes",
			packSizes:    []int{10000, 20000},
			itemsOrdered: 1000,
			expectedPacks: map[int]int{
				10000: 1,
			},
			expectedItems: 10000,
		},
		{
			name:         "Non-standard pack sizes",
			packSizes:    []int{19, 47, 101},
			itemsOrdered: 200,
			expectedPacks: map[int]int{
				101: 2,
			},
			expectedItems: 202,
		},
		{
			name:         "Small pack sizes",
			packSizes:    []int{1, 2, 5},
			itemsOrdered: 7,
			expectedPacks: map[int]int{
				5: 1,
				2: 1,
			},
			expectedItems: 7,
		},
		{
			name:         "Large gaps between pack sizes",
			packSizes:    []int{1, 1000, 10000},
			itemsOrdered: 500,
			expectedPacks: map[int]int{
				1: 500,
			},
			expectedItems: 500,
		},
		{
			name:         "Duplicate pack sizes",
			packSizes:    []int{250, 250, 500},
			itemsOrdered: 300,
			expectedPacks: map[int]int{
				500: 1,
			},
			expectedItems: 500,
		},
		{
			name:         "Boundary - just under pack size",
			packSizes:    []int{250, 500},
			itemsOrdered: 249,
			expectedPacks: map[int]int{
				250: 1,
			},
			expectedItems: 250,
		},
		{
			name:         "Boundary - just over pack size",
			packSizes:    []int{250, 500},
			itemsOrdered: 251,
			expectedPacks: map[int]int{
				500: 1,
			},
			expectedItems: 500,
		},
		{
			name:         "Very small order with large packs",
			packSizes:    []int{1000, 2000, 5000},
			itemsOrdered: 1,
			expectedPacks: map[int]int{
				1000: 1,
			},
			expectedItems: 1000,
		},
		{
			name:         "Large order with prime pack sizes",
			packSizes:    []int{23, 31, 53},
			itemsOrdered: 500000,
			expectedPacks: map[int]int{
				23: 2,
				31: 7,
				53: 9429,
			},
			expectedItems: 500000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPackSizeRepository{sizes: tt.packSizes}
			service := NewPackingService(mockRepo)

			solution, err := service.CalculatePacks(tt.itemsOrdered)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if solution.TotalItems != tt.expectedItems {
				t.Errorf("expected total items %d, got %d", tt.expectedItems, solution.TotalItems)
			}

			for packSize, expectedQty := range tt.expectedPacks {
				if solution.Packs[packSize] != expectedQty {
					t.Errorf("expected %d packs of size %d, got %d",
						expectedQty, packSize, solution.Packs[packSize])
				}
			}

			// Verify no unexpected pack sizes
			for packSize, qty := range solution.Packs {
				if expectedQty, exists := tt.expectedPacks[packSize]; !exists || expectedQty == 0 {
					t.Errorf("unexpected pack size %d with quantity %d", packSize, qty)
				}
			}
		})
	}
}

func TestPackingService_GetPackSizes(t *testing.T) {
	packSizes := []int{250, 500, 1000, 2000, 5000}
	mockRepo := &mockPackSizeRepository{sizes: packSizes}
	service := NewPackingService(mockRepo)

	sizes, err := service.GetPackSizes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sizes) != len(packSizes) {
		t.Errorf("expected %d pack sizes, got %d", len(packSizes), len(sizes))
	}

	for i, expectedSize := range packSizes {
		if sizes[i] != expectedSize {
			t.Errorf("expected pack size %d at index %d, got %d", expectedSize, i, sizes[i])
		}
	}
}
