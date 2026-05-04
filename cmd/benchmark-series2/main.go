package main

import (
	"fmt"
	"runtime"

	"bees_knapsack/internal/benchmark"
)

func main() {
	sizes := []int{1000, 5000, 10000, 25000, 50000, 100000, 150000, 200000, 250000, 300000}
	const numWorkers = 12
	const runs = 20
	const strategyName = "WorkerPool"

	params := benchmark.DefaultParams()

	fmt.Println("=== Series 2: Scaling by problem size ===")
	fmt.Printf("Workers: %d, Strategy: %s, Runs per point: %d\n", numWorkers, strategyName, runs)
	fmt.Printf("CPU cores available: %d\n\n", runtime.NumCPU())

	fmt.Printf("%-10s %-25s %-25s %s\n", "Items", "Seq mean±std (ms)", "Par mean±std (ms)", "Speedup")
	fmt.Println("-------------------------------------------------------------------")

	for _, n := range sizes {
		problem := benchmark.RandomProblem(n, 42)
		seqTimes := make([]float64, 0, runs)
		parTimes := make([]float64, 0, runs)

		for r := 0; r < runs; r++ {
			seed := int64(r + 1)
			if r%2 == 0 {
				seqTimes = append(seqTimes, benchmark.MeasureSequentialRun(problem, params, seed))
				strategy := benchmark.MakeStrategy(strategyName, problem, params, numWorkers, seed)
				parTimes = append(parTimes, benchmark.MeasureRun(strategy))
			} else {
				strategy := benchmark.MakeStrategy(strategyName, problem, params, numWorkers, seed)
				parTimes = append(parTimes, benchmark.MeasureRun(strategy))
				seqTimes = append(seqTimes, benchmark.MeasureSequentialRun(problem, params, seed))
			}
		}

		seqMean, seqStd := benchmark.ComputeStats(seqTimes)
		parMean, parStd := benchmark.ComputeStats(parTimes)
		speedup := seqMean / parMean

		fmt.Printf("%-10d %-25s %-25s %.2f\n",
			n,
			fmt.Sprintf("%.2f ± %.2f", seqMean, seqStd),
			fmt.Sprintf("%.2f ± %.2f", parMean, parStd),
			speedup,
		)
	}
}
