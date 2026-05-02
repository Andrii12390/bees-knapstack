package algorithm

import (
	"math/rand"
	"sort"
	"sync"

	"bees_knapsack/internal/problem"
)

type BatchedWorkerPoolStrategy struct {
	problem    *problem.Problem
	params     Params
	numWorkers int
	seed       int64
}

func NewBatchedWorkerPoolStrategy(p *problem.Problem, params Params, numWorkers int, seed int64) *BatchedWorkerPoolStrategy {
	return &BatchedWorkerPoolStrategy{
		problem:    p,
		params:     params,
		numWorkers: numWorkers,
		seed:       seed,
	}
}

func (s *BatchedWorkerPoolStrategy) Name() string { return "BatchedWorkerPool" }

func (s *BatchedWorkerPoolStrategy) Run() problem.Solution {
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

		tasks := s.buildTasks(iter, scouts)
		scouts = s.processTasks(iter, tasks)
	}

	return best
}

func (s *BatchedWorkerPoolStrategy) initScouts(rng *rand.Rand) []problem.Solution {
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

func (s *BatchedWorkerPoolStrategy) buildTasks(iter int, scouts []problem.Solution) []task {
	tasks := make([]task, s.params.NumScouts)
	for i := 0; i < s.params.NumEliteSites; i++ {
		tasks[i] = task{kind: taskLocalSearch, scoutIndex: i, site: scouts[i], numForagers: s.params.NumEliteForagers}
	}
	for i := s.params.NumEliteSites; i < s.params.NumBestSites; i++ {
		tasks[i] = task{kind: taskLocalSearch, scoutIndex: i, site: scouts[i], numForagers: s.params.NumBestForagers}
	}
	for i := s.params.NumBestSites; i < s.params.NumScouts; i++ {
		tasks[i] = task{kind: taskRandomSearch, scoutIndex: i}
	}

	shuffleRng := rand.New(rand.NewSource(s.seed + int64(iter)))
	shuffleRng.Shuffle(len(tasks), func(a, b int) {
		tasks[a], tasks[b] = tasks[b], tasks[a]
	})

	return tasks
}

func (s *BatchedWorkerPoolStrategy) processTasks(iter int, tasks []task) []problem.Solution {
	nextScouts := make([]problem.Solution, s.params.NumScouts)
	var wg sync.WaitGroup

	batchSize := (s.params.NumScouts + s.numWorkers - 1) / s.numWorkers
	for w := 0; w < s.numWorkers; w++ {
		lo := w * batchSize
		if lo >= s.params.NumScouts {
			break
		}
		hi := lo + batchSize
		if hi > s.params.NumScouts {
			hi = s.params.NumScouts
		}
		wg.Add(1)
		go func(lo, hi int) {
			defer wg.Done()
			for idx := lo; idx < hi; idx++ {
				t := tasks[idx]
				rng := rand.New(rand.NewSource(s.seed + int64(iter)*int64(s.params.NumScouts) + int64(t.scoutIndex) + 1))
				var sol problem.Solution
				switch t.kind {
				case taskLocalSearch:
					sol = DoLocalSearch(s.problem, t.site, t.numForagers, rng)
				case taskRandomSearch:
					sol = DoRandomSearch(s.problem, s.params, rng)
				}
				nextScouts[t.scoutIndex] = sol
			}
		}(lo, hi)
	}

	wg.Wait()
	return nextScouts
}
