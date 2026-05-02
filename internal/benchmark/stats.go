package benchmark

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"bees_knapsack/internal/algorithm"
	"bees_knapsack/internal/problem"
)

func RandomProblem(n int, seed int64) *problem.Problem {
	rng := rand.New(rand.NewSource(seed))
	items := make([]problem.Item, n)
	total := 0
	for i := range items {
		items[i] = problem.Item{
			Name:   fmt.Sprintf("item%d", i),
			Weight: rng.Intn(100) + 1,
			Value:  rng.Intn(100) + 1,
		}
		total += items[i].Weight
	}
	return &problem.Problem{Items: items, Capacity: total / 2}
}

func MeasureRun(strategy algorithm.ParallelStrategy) float64 {
	runtime.GC()
	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	start := time.Now()
	_ = strategy.Run()
	return float64(time.Since(start).Nanoseconds()) / 1e6
}

func MeasureSequentialRun(p *problem.Problem, params algorithm.Params, seed int64) float64 {
	runtime.GC()
	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	ba := algorithm.NewBeesAlgorithm(p, params, seed)
	start := time.Now()
	_ = ba.Run()
	return float64(time.Since(start).Nanoseconds()) / 1e6
}

func ComputeStats(times []float64) (mean float64, std float64) {
	n := len(times)
	if n == 0 {
		return 0, 0
	}
	sorted := make([]float64, n)
	copy(sorted, times)
	sort.Float64s(sorted)
	if n >= 10 {
		sorted = sorted[1 : n-1]
	}
	sum := 0.0
	for _, v := range sorted {
		sum += v
	}
	mean = sum / float64(len(sorted))
	sqSum := 0.0
	for _, v := range sorted {
		diff := v - mean
		sqSum += diff * diff
	}
	std = math.Sqrt(sqSum / float64(len(sorted)))
	return mean, std
}

func MakeStrategy(name string, p *problem.Problem, params algorithm.Params, numWorkers int, seed int64) algorithm.ParallelStrategy {
	switch name {
	case "Sequential":
		return algorithm.NewBeesAlgorithm(p, params, seed)
	case "WorkerPool":
		return algorithm.NewWorkerPoolStrategy(p, params, numWorkers, seed)
	case "GoroutinePerTask":
		return algorithm.NewGoroutinePerTaskStrategy(p, params, numWorkers, seed)
	case "BatchedWorkerPool":
		return algorithm.NewBatchedWorkerPoolStrategy(p, params, numWorkers, seed)
	}
	return nil
}

func DefaultParams() algorithm.Params {
	return algorithm.Params{
		NumScouts:        100,
		NumBestSites:     20,
		NumEliteSites:    5,
		NumEliteForagers: 50,
		NumBestForagers:  20,
		InitPatchSize:    5,
		MaxIterations:    100,
	}
}
