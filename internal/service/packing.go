package service

import (
	"fmt"
	"sort"
)

type PackingService struct {
	packSizes []int
}

func NewPackingService(packSizes []int) *PackingService {
	// Create a copy and sort in descending order for optimization
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	return &PackingService{
		packSizes: sizes,
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

	if len(ps.packSizes) == 0 {
		return nil, fmt.Errorf("no pack sizes configured")
	}

	solution := ps.findOptimalSolution(itemsOrdered)

	if solution == nil {
		return nil, fmt.Errorf("unable to fulfill order with current pack sizes")
	}

	return solution, nil
}

func (ps *PackingService) findOptimalSolution(target int) *PackSolution {
	// Dynamic programming approach
	maxItems := target + ps.packSizes[0]
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

		for _, packSize := range ps.packSizes {
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

func (ps *PackingService) GetPackSizes() []int {
	sizes := make([]int, len(ps.packSizes))
	copy(sizes, ps.packSizes)
	sort.Ints(sizes) // Return in ascending order
	return sizes
}
