package main

import (
	"fmt"
	"time"

	"bees_knapsack/internal/algorithm"
	"bees_knapsack/internal/benchmark"
)

func main() {
	sizes := []int{2500, 5000, 7500, 10000, 12500, 15000, 17500, 20000}
	const runs = 3

	params := benchmark.DefaultParams()

	fmt.Printf("%-10s %-15s %-15s\n", "Items", "Time (ms)", "Best Value")
	fmt.Println("--------------------------------------------")

	for _, n := range sizes {
		problem := benchmark.RandomProblem(n, 42)

		var totalNs int64
		var totalValue int

		for r := 0; r < runs; r++ {
			ba := algorithm.NewBeesAlgorithm(problem, params, int64(r+1))
			start := time.Now()
			best := ba.Run()
			totalNs += time.Since(start).Nanoseconds()
			totalValue += best.Fitness
		}

		avgMs := float64(totalNs) / float64(runs) / 1e6
		avgValue := totalValue / runs

		fmt.Printf("%-10d %-15.2f %-15d\n", n, avgMs, avgValue)
	}
}
