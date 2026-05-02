package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"bees_knapsack/internal/algorithm"
	"bees_knapsack/internal/benchmark"
	"bees_knapsack/internal/problem"
	"bees_knapsack/internal/verify"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "verify" {
		strategyName := ""
		if len(os.Args) > 2 {
			strategyName = os.Args[2]
		}
		verify.RunVerification(strategyName, runtime.NumCPU())
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "smoke" {
		runStrategySmokeTest()
		return
	}

	p := &problem.Problem{
		Capacity: 50,
		Items: []problem.Item{
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

	params := algorithm.Params{
		NumScouts:        30,
		NumBestSites:     10,
		NumEliteSites:    4,
		NumEliteForagers: 8,
		NumBestForagers:  4,
		InitPatchSize:    3,
		MaxIterations:    200,
	}

	seed := time.Now().UnixNano()
	ba := algorithm.NewBeesAlgorithm(p, params, seed)

	start := time.Now()
	best := ba.Run()
	elapsed := time.Since(start)

	printResult(p, best, elapsed)
}

func printResult(p *problem.Problem, best problem.Solution, elapsed time.Duration) {
	takenItems := p.TakenItems(best.Bits)
	totalWeight := p.TotalWeight(best.Bits)

	fmt.Println("=== Bees Algorithm — Knapsack Problem ===")
	fmt.Printf("Capacity : %d\n\n", p.Capacity)

	fmt.Println("Selected items:")
	fmt.Printf("  %-16s %8s %8s\n", "Name", "Weight", "Value")
	fmt.Printf("  %-16s %8s %8s\n", "----", "------", "-----")
	for _, item := range takenItems {
		fmt.Printf("  %-16s %8d %8d\n", item.Name, item.Weight, item.Value)
	}

	fmt.Println()
	fmt.Printf("Total weight : %d / %d\n", totalWeight, p.Capacity)
	fmt.Printf("Total value  : %d\n", best.Fitness)
	fmt.Printf("Time elapsed : %v\n", elapsed)
}

func runStrategySmokeTest() {
	const n = 100
	const seed = int64(123)
	const numWorkers = 4
	const tolerance = 0.05

	p := benchmark.RandomProblem(n, 42)
	params := benchmark.DefaultParams()

	fmt.Println("=== Strategy smoke test ===")
	fmt.Printf("Items: %d, seed: %d, workers: %d\n\n", n, seed, numWorkers)

	seqBA := algorithm.NewBeesAlgorithm(p, params, seed)
	seqBest := seqBA.Run()
	fmt.Printf("%-18s fitness=%d\n", "Sequential", seqBest.Fitness)

	strategies := []string{"WorkerPool", "GoroutinePerTask", "BatchedWorkerPool""}
	fitnesses := []int{seqBest.Fitness}

	for _, name := range strategies {
		s := benchmark.MakeStrategy(name, p, params, numWorkers, seed)
		out := s.Run()
		fitnesses = append(fitnesses, out.Fitness)
		fmt.Printf("%-18s fitness=%d\n", name, out.Fitness)
	}

	minF, maxF := fitnesses[0], fitnesses[0]
	for _, f := range fitnesses {
		if f < minF {
			minF = f
		}
		if f > maxF {
			maxF = f
		}
	}

	spread := 0.0
	if maxF > 0 {
		spread = float64(maxF-minF) / float64(maxF)
	}
	fmt.Printf("\nSpread: %.2f%% (min=%d, max=%d)\n", spread*100, minF, maxF)
	if spread <= tolerance {
		fmt.Println("PASS — all strategies within 5%")
	} else {
		fmt.Println("FAIL — strategies diverge more than 5%")
	}
}
