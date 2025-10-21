package service

import (
	"fmt"
	"sort"

	"github.com/miloradbozic/packing-service/internal/database"
)

type PackingService struct {
	packSizeRepo *database.PackSizeRepository
}

func NewPackingService(packSizeRepo *database.PackSizeRepository) *PackingService {
	return &PackingService{
		packSizeRepo: packSizeRepo,
	}
}

type PackSolution struct {
	Packs      map[int]int
	TotalItems int
	TotalPacks int
}

func (ps *PackingService) CalculatePacks(itemsOrdered int) (*PackSolution, error) {
	if itemsOrdered <= 0 {
		return nil, fmt.Errorf("items ordered must be positive")
	}

	// Get active pack sizes from database
	packSizes, err := ps.packSizeRepo.GetAllActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}

	if len(packSizes) == 0 {
		return nil, fmt.Errorf("no pack sizes configured")
	}

	solution := ps.findOptimalSolution(itemsOrdered, packSizes)

	if solution == nil {
		return nil, fmt.Errorf("unable to fulfill order with current pack sizes")
	}

	return solution, nil
}

func (ps *PackingService) findOptimalSolution(target int, packSizes []int) *PackSolution {
	// Sort pack sizes in descending order for optimization
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	// Dynamic programming approach
	maxItems := target + sizes[0]
	dp := make([]int, maxItems+1)
	parent := make([]int, maxItems+1)

	// Initialize with impossible values
	for i := range dp {
		dp[i] = -1
		parent[i] = -1
	}
	dp[0] = 0

	// Build the DP table
	for i := 0; i <= maxItems; i++ {
		if dp[i] == -1 {
			continue
		}

		for _, packSize := range sizes {
			next := i + packSize
			if next <= maxItems {
				if dp[next] == -1 || dp[next] > dp[i]+packSize {
					dp[next] = dp[i] + packSize
					parent[next] = packSize
				}
			}
		}
	}

	// Find the minimum items >= target
	minItems := -1
	for i := target; i <= maxItems; i++ {
		if dp[i] != -1 {
			minItems = i
			break
		}
	}

	if minItems == -1 {
		return nil
	}

	// Reconstruct the solution
	packs := make(map[int]int)
	current := minItems
	totalPacks := 0

	for current > 0 && parent[current] != -1 {
		packSize := parent[current]
		packs[packSize]++
		totalPacks++
		current -= packSize
	}

	return &PackSolution{
		Packs:      packs,
		TotalItems: minItems,
		TotalPacks: totalPacks,
	}
}

func (ps *PackingService) GetPackSizes() ([]int, error) {
	return ps.packSizeRepo.GetAllActive()
}
