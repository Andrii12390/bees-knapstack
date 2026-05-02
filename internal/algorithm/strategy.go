package algorithm

import (
	"math/rand"

	"bees_knapsack/internal/problem"
)

type ParallelStrategy interface {
	Name() string
	Run() problem.Solution
}

type taskKind int

const (
	taskLocalSearch taskKind = iota
	taskRandomSearch
)

type task struct {
	kind        taskKind
	scoutIndex  int
	site        problem.Solution
	numForagers int
}

func DoLocalSearch(p *problem.Problem, site problem.Solution, numForagers int, rng *rand.Rand) problem.Solution {
	best := site.Clone()

	for f := 0; f < numForagers; f++ {
		candidateBits := problem.NeighborSolution(site.Bits, site.PatchSize, rng)
		candidateFitness := p.Evaluate(candidateBits)

		if candidateFitness > best.Fitness {
			best.Bits = candidateBits
			best.Fitness = candidateFitness
		}
	}

	if best.Fitness <= site.Fitness && best.PatchSize > 1 {
		best.PatchSize--
	}

	return best
}

func DoRandomSearch(p *problem.Problem, params Params, rng *rand.Rand) problem.Solution {
	bits := problem.RandomSolution(len(p.Items), rng)
	return problem.Solution{
		Bits:      bits,
		Fitness:   p.Evaluate(bits),
		PatchSize: params.InitPatchSize,
	}
}
