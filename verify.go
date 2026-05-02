package main

import (
	"fmt"
	"strings"
)

func bruteForce(problem *Problem) ([]int, int) {
	n := len(problem.Items)
	bestValue := 0
	bestBits := make([]int, n)

	for mask := 0; mask < (1 << n); mask++ {
		bits := make([]int, n)
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				bits[i] = 1
			}
		}
		value := problem.Evaluate(bits)
		if value > bestValue {
			bestValue = value
			bestBits = bits
		}
	}

	return bestBits, bestValue
}

type TestCase struct {
	Name        string
	Problem     *Problem
	ExpectedMin int
}

func RunVerification() {
	testCases := []TestCase{
		{
			Name: "All items fit",
			Problem: &Problem{
				Capacity: 100,
				Items: []Item{
					{Name: "A", Weight: 5, Value: 10},
					{Name: "B", Weight: 5, Value: 20},
					{Name: "C", Weight: 5, Value: 30},
				},
			},
		},
		{
			Name: "Only one item fits",
			Problem: &Problem{
				Capacity: 6,
				Items: []Item{
					{Name: "A", Weight: 5, Value: 20},
					{Name: "B", Weight: 5, Value: 35},
					{Name: "C", Weight: 5, Value: 50},
				},
			},
		},
		{
			Name: "Greedy trap (ratio heuristic fails)",
			Problem: &Problem{
				Capacity: 6,
				Items: []Item{
					{Name: "A (ratio=10)", Weight: 1, Value: 10},
					{Name: "B (ratio=3)", Weight: 3, Value: 9},
					{Name: "C (ratio=2.3)", Weight: 3, Value: 7},
				},
			},
		},
		{
			Name: "Nothing fits",
			Problem: &Problem{
				Capacity: 2,
				Items: []Item{
					{Name: "A", Weight: 5, Value: 100},
					{Name: "B", Weight: 8, Value: 200},
				},
			},
		},
		{
			Name: "Medium instance (10 items)",
			Problem: &Problem{
				Capacity: 20,
				Items: []Item{
					{Name: "item1", Weight: 6, Value: 10},
					{Name: "item2", Weight: 4, Value: 40},
					{Name: "item3", Weight: 12, Value: 30},
					{Name: "item4", Weight: 1, Value: 50},
					{Name: "item5", Weight: 5, Value: 35},
					{Name: "item6", Weight: 3, Value: 30},
					{Name: "item7", Weight: 7, Value: 15},
					{Name: "item8", Weight: 2, Value: 25},
					{Name: "item9", Weight: 4, Value: 20},
					{Name: "item10", Weight: 1, Value: 18},
				},
			},
		},
	}

	passes := 0
	separator := strings.Repeat("─", 52)

	for i, tc := range testCases {
		fmt.Printf("\n%s\n", separator)
		fmt.Printf("Test %d: %s\n", i+1, tc.Name)
		fmt.Printf("Capacity: %d | Items: %d\n", tc.Problem.Capacity, len(tc.Problem.Items))
		fmt.Printf("%s\n", separator)

		bfBits, bfValue := bruteForce(tc.Problem)
		bfWeight := tc.Problem.TotalWeight(bfBits)

		baValue := 0
		for run := 0; run < 5; run++ {
			params := Params{
				NumScouts:        20,
				NumBestSites:     6,
				NumEliteSites:    2,
				NumEliteForagers: 8,
				NumBestForagers:  4,
				InitPatchSize:    3,
				MaxIterations:    300,
			}
			ba := NewBeesAlgorithm(tc.Problem, params, int64(run+42))
			result := ba.Run()
			if result.Fitness > baValue {
				baValue = result.Fitness
			}
		}

		fmt.Printf("Brute force > value: %d, weight: %d/%d\n",
			bfValue, bfWeight, tc.Problem.Capacity)
		fmt.Printf("  Items: %s\n", itemNames(tc.Problem.TakenItems(bfBits)))

		fmt.Printf("Bees Algorithm > best value found: %d\n", baValue)

		if baValue == bfValue {
			fmt.Printf("PASS — BA found the optimal solution\n")
			passes++
		} else if baValue > 0 && float64(baValue)/float64(bfValue) >= 0.95 {
			fmt.Printf("NEAR-OPTIMAL — BA found %.1f%% of optimal\n",
				float64(baValue)/float64(bfValue)*100)
		} else {
			fmt.Printf("FAIL — BA found %d, optimal is %d\n", baValue, bfValue)
		}
	}

	fmt.Printf("\n%s\n", separator)
	fmt.Printf("Result: %d / %d tests passed\n", passes, len(testCases))
	fmt.Printf("%s\n\n", separator)
}

func itemNames(items []Item) string {
	if len(items) == 0 {
		return "(none)"
	}
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}
	return strings.Join(names, ", ")
}
