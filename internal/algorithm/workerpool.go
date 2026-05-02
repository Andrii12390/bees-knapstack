package algorithm

import (
	"math/rand"
	"sort"
	"sync"

	"bees_knapsack/internal/problem"
)

type WorkerPoolStrategy struct {
	problem    *problem.Problem
	params     Params
	numWorkers int
	seed       int64
}

func NewWorkerPoolStrategy(p *problem.Problem, params Params, numWorkers int, seed int64) *WorkerPoolStrategy {
	return &WorkerPoolStrategy{
		problem:    p,
		params:     params,
		numWorkers: numWorkers,
		seed:       seed,
	}
}

func (s *WorkerPoolStrategy) Name() string { return "WorkerPool" }

type result struct {
	scoutIndex int
	solution   problem.Solution
}

func (s *WorkerPoolStrategy) Run() problem.Solution {
	tasks := make(chan task, s.params.NumScouts)
	results := make(chan result, s.params.NumScouts)

	wg := s.startWorkers(tasks, results)

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

		s.dispatchTasks(tasks, scouts)
		scouts = s.collectResults(results)
	}

	close(tasks)
	wg.Wait()

	return best
}

func (s *WorkerPoolStrategy) startWorkers(tasks chan task, results chan result) *sync.WaitGroup {
	var wg sync.WaitGroup
	for w := 0; w < s.numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(s.seed + int64(workerID) + 1))
			s.worker(tasks, results, rng)
		}(w)
	}
	return &wg
}

func (s *WorkerPoolStrategy) initScouts(rng *rand.Rand) []problem.Solution {
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

func (s *WorkerPoolStrategy) dispatchTasks(tasks chan<- task, scouts []problem.Solution) {
	for i := 0; i < s.params.NumEliteSites; i++ {
		tasks <- task{
			kind:        taskLocalSearch,
			scoutIndex:  i,
			site:        scouts[i],
			numForagers: s.params.NumEliteForagers,
		}
	}
	for i := s.params.NumEliteSites; i < s.params.NumBestSites; i++ {
		tasks <- task{
			kind:        taskLocalSearch,
			scoutIndex:  i,
			site:        scouts[i],
			numForagers: s.params.NumBestForagers,
		}
	}
	for i := s.params.NumBestSites; i < s.params.NumScouts; i++ {
		tasks <- task{
			kind:       taskRandomSearch,
			scoutIndex: i,
		}
	}
}

func (s *WorkerPoolStrategy) collectResults(results <-chan result) []problem.Solution {
	nextScouts := make([]problem.Solution, s.params.NumScouts)
	for collected := 0; collected < s.params.NumScouts; collected++ {
		r := <-results
		nextScouts[r.scoutIndex] = r.solution
	}
	return nextScouts
}

func (s *WorkerPoolStrategy) worker(tasks <-chan task, results chan<- result, rng *rand.Rand) {
	for t := range tasks {
		var solution problem.Solution
		switch t.kind {
		case taskLocalSearch:
			solution = DoLocalSearch(s.problem, t.site, t.numForagers, rng)
		case taskRandomSearch:
			solution = DoRandomSearch(s.problem, s.params, rng)
		}
		results <- result{scoutIndex: t.scoutIndex, solution: solution}
	}
}
