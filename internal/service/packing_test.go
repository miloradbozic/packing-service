package service

import (
	"testing"

	"github.com/miloradbozic/packing-service/internal/database"
)

// Mock repository for testing
type mockPackSizeRepository struct {
	sizes []int
}

func (m *mockPackSizeRepository) GetAllActive() ([]int, error) {
	return m.sizes, nil
}

func (m *mockPackSizeRepository) GetAll() ([]database.PackSize, error) {
	return nil, nil
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

func (m *mockPackSizeRepository) MigrateFromConfig(sizes []int) error {
	return nil
}

func TestPackingService_CalculatePacks(t *testing.T) {
	packSizes := []int{250, 500, 1000, 2000, 5000}
	mockRepo := &mockPackSizeRepository{sizes: packSizes}
	service := NewPackingService(mockRepo)

	tests := []struct {
		name          string
		itemsOrdered  int
		expectedPacks map[int]int
		expectedItems int
	}{
		{
			name:         "Exactly 1 item",
			itemsOrdered: 1,
			expectedPacks: map[int]int{
				250: 1,
			},
			expectedItems: 250,
		},
		{
			name:         "Exactly 250 items",
			itemsOrdered: 250,
			expectedPacks: map[int]int{
				250: 1,
			},
			expectedItems: 250,
		},
		{
			name:         "251 items",
			itemsOrdered: 251,
			expectedPacks: map[int]int{
				500: 1,
			},
			expectedItems: 500,
		},
		{
			name:         "501 items",
			itemsOrdered: 501,
			expectedPacks: map[int]int{
				500: 1,
				250: 1,
			},
			expectedItems: 750,
		},
		{
			name:         "12001 items",
			itemsOrdered: 12001,
			expectedPacks: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
			expectedItems: 12250,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solution, err := service.CalculatePacks(tt.itemsOrdered)
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
		})
	}
}
