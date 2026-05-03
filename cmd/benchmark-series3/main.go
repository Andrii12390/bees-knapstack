package main

import (
	"fmt"
	"runtime"

	"bees_knapsack/internal/benchmark"
)

func main() {
	const n = 50000
	const numWorkers = 12
	const runs = 20

	strategies := []string{"WorkerPool", "GoroutinePerTask", "BatchedWorkerPool"}
	params := benchmark.DefaultParams()

	problem := benchmark.RandomProblem(n, 42)

	fmt.Println("=== Series 3: Comparison of parallelization strategies ===")
	fmt.Printf("Problem size: %d, Workers: %d, Runs per strategy: %d\n", n, numWorkers, runs)
	fmt.Printf("CPU cores available: %d\n", runtime.NumCPU())

	seqTimes := make([]float64, 0, runs)
	for r := 0; r < runs; r++ {
		seqTimes = append(seqTimes, benchmark.MeasureSequentialRun(problem, params, int64(r+1)))
	}
	seqMean, seqStd := benchmark.ComputeStats(seqTimes)
	fmt.Printf("Sequential baseline: %.2f ± %.2f ms\n\n", seqMean, seqStd)

	fmt.Printf("%-22s %-12s %-12s %s\n", "Strategy", "Mean (ms)", "Std (ms)", "Speedup")
	fmt.Println("-----------------------------------------------------------")

	for _, name := range strategies {
		times := make([]float64, 0, runs)
		for r := 0; r < runs; r++ {
			strategy := benchmark.MakeStrategy(name, problem, params, numWorkers, int64(r+1))
			times = append(times, benchmark.MeasureRun(strategy))
		}
		mean, std := benchmark.ComputeStats(times)
		speedup := seqMean / mean

		fmt.Printf("%-22s %-12.2f %-12.2f %.2f\n", name, mean, std, speedup)
	}
}
