//go:build !benchmark

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "verify" {
		RunVerification()
		return
	}

	problem := &Problem{
		Capacity: 50,
		Items: []Item{
			{Name: "laptop", Weight: 10, Value: 60},
			{Name: "phone", Weight: 4, Value: 40},
			{Name: "headphones", Weight: 3, Value: 25},
			{Name: "camera", Weight: 7, Value: 50},
			{Name: "charger", Weight: 2, Value: 15},
			{Name: "tablet", Weight: 8, Value: 45},
			{Name: "watch", Weight: 1, Value: 10},
			{Name: "keyboard", Weight: 6, Value: 30},
			{Name: "hard drive", Weight: 5, Value: 35},
			{Name: "power bank", Weight: 3, Value: 20},
			{Name: "mouse", Weight: 2, Value: 12},
			{Name: "USB hub", Weight: 1, Value: 8},
			{Name: "monitor", Weight: 15, Value: 70},
			{Name: "speakers", Weight: 9, Value: 40},
			{Name: "webcam", Weight: 4, Value: 22},
		},
	}

	params := Params{
		NumScouts:        30,
		NumBestSites:     10,
		NumEliteSites:    4,
		NumEliteForagers: 8,
		NumBestForagers:  4,
		InitPatchSize:    3,
		MaxIterations:    200,
	}

	seed := time.Now().UnixNano()
	ba := NewBeesAlgorithm(problem, params, seed)

	start := time.Now()
	best := ba.Run()
	elapsed := time.Since(start)

	printResult(problem, best, elapsed)
}

func printResult(problem *Problem, best Solution, elapsed time.Duration) {
	takenItems := problem.TakenItems(best.Bits)
	totalWeight := problem.TotalWeight(best.Bits)

	fmt.Println("=== Bees Algorithm — Knapsack Problem ===")
	fmt.Printf("Capacity : %d\n\n", problem.Capacity)

	fmt.Println("Selected items:")
	fmt.Printf("  %-16s %8s %8s\n", "Name", "Weight", "Value")
	fmt.Printf("  %-16s %8s %8s\n", "----", "------", "-----")
	for _, item := range takenItems {
		fmt.Printf("  %-16s %8d %8d\n", item.Name, item.Weight, item.Value)
	}

	fmt.Println()
	fmt.Printf("Total weight : %d / %d\n", totalWeight, problem.Capacity)
	fmt.Printf("Total value  : %d\n", best.Fitness)
	fmt.Printf("Time elapsed : %v\n", elapsed)
}
