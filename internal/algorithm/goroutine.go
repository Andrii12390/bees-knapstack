package algorithm

import (
	"math/rand"
	"sort"
	"sync"

	"bees_knapsack/internal/problem"
)

type GoroutinePerTaskStrategy struct {
	problem    *problem.Problem
	params     Params
	numWorkers int
	seed       int64
}

func NewGoroutinePerTaskStrategy(p *problem.Problem, params Params, numWorkers int, seed int64) *GoroutinePerTaskStrategy {
	return &GoroutinePerTaskStrategy{
		problem:    p,
		params:     params,
		numWorkers: numWorkers,
		seed:       seed,
	}
}

func (s *GoroutinePerTaskStrategy) Name() string { return "GoroutinePerTask" }

func (s *GoroutinePerTaskStrategy) Run() problem.Solution {
	mainRng := rand.New(rand.NewSource(s.seed))
	scouts := s.initScouts(mainRng)
	best := scouts[0].Clone()

	for iter := 0; iter < s.params.MaxIterations; iter++ {
		sort.Slice(scouts, func(i, j int) bool {
			return scouts[i].Fitness > scouts[j].Fitness
		})

		if scouts[0].Fitness > best.Fitness {
			best = scouts[0].Clone()
		}

		scouts = s.runIteration(iter, scouts)
	}

	return best
}

func (s *GoroutinePerTaskStrategy) initScouts(rng *rand.Rand) []problem.Solution {
	scouts := make([]problem.Solution, s.params.NumScouts)
	for i := range scouts {
		bits := problem.RandomSolution(len(s.problem.Items), rng)
		scouts[i] = problem.Solution{
			Bits:      bits,
			Fitness:   s.problem.Evaluate(bits),
			PatchSize: s.params.InitPatchSize,
		}
	}
	return scouts
}

func (s *GoroutinePerTaskStrategy) runIteration(iter int, scouts []problem.Solution) []problem.Solution {
	nextScouts := make([]problem.Solution, s.params.NumScouts)
	var wg sync.WaitGroup

	for i := 0; i < s.params.NumEliteSites; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(s.seed + int64(iter*1000+idx) + 1))
			nextScouts[idx] = DoLocalSearch(s.problem, scouts[idx], s.params.NumEliteForagers, rng)
		}(i)
	}
	for i := s.params.NumEliteSites; i < s.params.NumBestSites; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(s.seed + int64(iter*1000+idx) + 1))
			nextScouts[idx] = DoLocalSearch(s.problem, scouts[idx], s.params.NumBestForagers, rng)
		}(i)
	}
	for i := s.params.NumBestSites; i < s.params.NumScouts; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(s.seed + int64(iter*1000+idx) + 1))
			nextScouts[idx] = DoRandomSearch(s.problem, s.params, rng)
		}(i)
	}

	wg.Wait()
	return nextScouts
}
