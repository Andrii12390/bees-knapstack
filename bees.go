package main

import (
	"math/rand"
	"sort"
)

type Params struct {
	NumScouts        int
	NumBestSites     int
	NumEliteSites    int
	NumEliteForagers int
	NumBestForagers  int
	InitPatchSize    int
	MaxIterations    int
}

type BeesAlgorithm struct {
	problem *Problem
	params  Params
	rng     *rand.Rand
}

func NewBeesAlgorithm(problem *Problem, params Params, seed int64) *BeesAlgorithm {
	return &BeesAlgorithm{
		problem: problem,
		params:  params,
		rng:     rand.New(rand.NewSource(seed)),
	}
}

func (ba *BeesAlgorithm) Run() Solution {
	scouts := ba.initializeScouts()
	best := scouts[0].clone()

	for iter := 0; iter < ba.params.MaxIterations; iter++ {
		sort.Slice(scouts, func(i, j int) bool {
			return scouts[i].Fitness > scouts[j].Fitness
		})

		if scouts[0].Fitness > best.Fitness {
			best = scouts[0].clone()
		}

		nextScouts := make([]Solution, ba.params.NumScouts)

		for i := 0; i < ba.params.NumEliteSites; i++ {
			nextScouts[i] = ba.localSearch(scouts[i], ba.params.NumEliteForagers)
		}

		for i := ba.params.NumEliteSites; i < ba.params.NumBestSites; i++ {
			nextScouts[i] = ba.localSearch(scouts[i], ba.params.NumBestForagers)
		}

		for i := ba.params.NumBestSites; i < ba.params.NumScouts; i++ {
			nextScouts[i] = ba.newRandomSolution()
		}

		scouts = nextScouts
	}

	return best
}

func (ba *BeesAlgorithm) initializeScouts() []Solution {
	numItems := len(ba.problem.Items)
	scouts := make([]Solution, ba.params.NumScouts)

	for i := range scouts {
		scouts[i] = ba.newRandomSolution()
		_ = numItems
	}

	return scouts
}

func (ba *BeesAlgorithm) newRandomSolution() Solution {
	bits := randomSolution(len(ba.problem.Items), ba.rng)
	return Solution{
		Bits:      bits,
		Fitness:   ba.problem.Evaluate(bits),
		PatchSize: ba.params.InitPatchSize,
	}
}

func (ba *BeesAlgorithm) localSearch(site Solution, numForagers int) Solution {
	best := site.clone()

	for f := 0; f < numForagers; f++ {
		candidateBits := neighborSolution(site.Bits, site.PatchSize, ba.rng)
		candidateFitness := ba.problem.Evaluate(candidateBits)

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
