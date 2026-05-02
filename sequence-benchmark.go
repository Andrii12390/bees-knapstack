//go:build benchmark

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	sizes := []int{2500, 5000, 7500, 10000, 12500, 15000, 17500, 20000}
	const runs = 3

	params := Params{
		NumScouts:        100,
		NumBestSites:     20,
		NumEliteSites:    5,
		NumEliteForagers: 50,
		NumBestForagers:  20,
		InitPatchSize:    5,
		MaxIterations:    100,
	}

	fmt.Printf("%-10s %-15s %-15s\n", "Items", "Time (ms)", "Best Value")
	fmt.Println("--------------------------------------------")

	for _, n := range sizes {
		problem := randomProblem(n)

		var totalNs int64
		var totalValue int

		for r := 0; r < runs; r++ {
			ba := NewBeesAlgorithm(problem, params, int64(r+1))
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

func randomProblem(n int) *Problem {
	rng := rand.New(rand.NewSource(42))
	items := make([]Item, n)
	total := 0
	for i := range items {
		items[i] = Item{
			Name:   fmt.Sprintf("item%d", i),
			Weight: rng.Intn(100) + 1,
			Value:  rng.Intn(100) + 1,
		}
		total += items[i].Weight
	}
	return &Problem{Items: items, Capacity: total / 2}
}
