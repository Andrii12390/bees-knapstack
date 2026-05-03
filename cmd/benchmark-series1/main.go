package main

import (
	"fmt"
	"runtime"

	"bees_knapsack/internal/benchmark"
)

func main() {
	const n = 50000
	workerCounts := []int{1, 2, 4, 6, 8, 10, 12, 16}
	const runs = 20
	const strategyName = "BatchedWorkerPool"

	params := benchmark.DefaultParams()

	originalProcs := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(originalProcs)

	problem := benchmark.RandomProblem(n, 42)

	fmt.Println("=== Series 2: Scaling by worker count ===")
	fmt.Printf("Problem size: %d, Strategy: %s, Runs per point: %d\n", n, strategyName, runs)
	fmt.Printf("CPU cores available: %d (original GOMAXPROCS=%d)\n", runtime.NumCPU(), originalProcs)

	seqTimes := make([]float64, 0, runs)
	for r := 0; r < runs; r++ {
		seqTimes = append(seqTimes, benchmark.MeasureSequentialRun(problem, params, int64(r+1)))
	}
	seqMean, seqStd := benchmark.ComputeStats(seqTimes)
	fmt.Printf("Sequential baseline: %.2f ± %.2f ms\n\n", seqMean, seqStd)

	fmt.Printf("%-10s %-25s %-12s %s\n", "Workers", "Time mean±std (ms)", "Speedup", "Efficiency")
	fmt.Println("---------------------------------------------------------")

	for _, w := range workerCounts {
		runtime.GOMAXPROCS(w)

		parTimes := make([]float64, 0, runs)
		for r := 0; r < runs; r++ {
			strategy := benchmark.MakeStrategy(strategyName, problem, params, w, int64(r+1))
			parTimes = append(parTimes, benchmark.MeasureRun(strategy))
		}
		parMean, parStd := benchmark.ComputeStats(parTimes)
		speedup := seqMean / parMean
		efficiency := speedup / float64(w) * 100

		fmt.Printf("%-10d %-25s %-12.2f %.0f%%\n",
			w,
			fmt.Sprintf("%.2f ± %.2f", parMean, parStd),
			speedup,
			efficiency,
		)

		runtime.GOMAXPROCS(originalProcs)
	}
}
