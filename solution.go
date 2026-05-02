package main

import "math/rand"

type Solution struct {
	Bits      []int
	Fitness   int
	PatchSize int
}

func (s Solution) clone() Solution {
	bits := make([]int, len(s.Bits))
	copy(bits, s.Bits)
	return Solution{
		Bits:      bits,
		Fitness:   s.Fitness,
		PatchSize: s.PatchSize,
	}
}

func randomSolution(numItems int, rng *rand.Rand) []int {
	bits := make([]int, numItems)
	for i := range bits {
		bits[i] = rng.Intn(2)
	}
	return bits
}

func neighborSolution(bits []int, patchSize int, rng *rand.Rand) []int {
	neighbor := make([]int, len(bits))
	copy(neighbor, bits)

	effectivePatch := patchSize
	if effectivePatch > len(bits) {
		effectivePatch = len(bits)
	}

	numFlips := rng.Intn(effectivePatch) + 1
	positions := rng.Perm(len(bits))[:numFlips]

	for _, pos := range positions {
		neighbor[pos] ^= 1
	}

	return neighbor
}
