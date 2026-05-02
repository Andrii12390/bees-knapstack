package main

import (
	"fmt"
	"runtime"

	"bees_knapsack/internal/benchmark"
)

func main() {
	const n = 50000
	const runs = 20

	numWorkers := runtime.NumCPU()
	params := benchmark.DefaultParams()
	strategies := []string{"WorkerPool", "GoroutinePerTask", "BatchedWorkerPool"}

	problem := benchmark.RandomProblem(n, 42)

	fmt.Printf("CPU cores: %d, workers: %d\n", runtime.NumCPU(), numWorkers)
	fmt.Printf("Problem size: %d, Runs per strategy: %d\n\n", n, runs)

	seqTimes := make([]float64, 0, runs)
	parTimes := make(map[string][]float64, len(strategies))
	for _, name := range strategies {
		parTimes[name] = make([]float64, 0, runs)
	}

	order := append([]string{"Sequential"}, strategies...)
	reversed := make([]string, len(order))
	for i, name := range order {
		reversed[len(order)-1-i] = name
	}

	for r := 0; r < runs; r++ {
		seed := int64(r + 1)
		exec := order
		if r%2 != 0 {
			exec = reversed
		}
		for _, name := range exec {
			if name == "Sequential" {
				seqTimes = append(seqTimes, benchmark.MeasureSequentialRun(problem, params, seed))
			} else {
				strategy := benchmark.MakeStrategy(name, problem, params, numWorkers, seed)
				parTimes[name] = append(parTimes[name], benchmark.MeasureRun(strategy))
			}
		}
	}

	seqMean, seqStd := benchmark.ComputeStats(seqTimes)

	fmt.Printf("%-22s %-13s %-12s %s\n", "Strategy", "Mean (ms)", "Std (ms)", "Speedup vs Seq")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("%-22s %-13.2f %-12.2f %.2f\n", "Sequential", seqMean, seqStd, 1.0)
	for _, name := range strategies {
		mean, std := benchmark.ComputeStats(parTimes[name])
		speedup := seqMean / mean
		fmt.Printf("%-22s %-13.2f %-12.2f %.2f\n", name, mean, std, speedup)
	}
}
