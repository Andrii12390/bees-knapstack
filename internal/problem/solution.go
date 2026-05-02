package problem

import "math/rand"

type Solution struct {
	Bits      []uint8
	Fitness   int
	PatchSize int
}

func (s Solution) Clone() Solution {
	bits := make([]uint8, len(s.Bits))
	copy(bits, s.Bits)
	return Solution{
		Bits:      bits,
		Fitness:   s.Fitness,
		PatchSize: s.PatchSize,
	}
}

func RandomSolution(numItems int, rng *rand.Rand) []uint8 {
	bits := make([]uint8, numItems)
	for i := range bits {
		bits[i] = uint8(rng.Intn(2))
	}
	return bits
}

func NeighborSolution(bits []uint8, patchSize int, rng *rand.Rand) []uint8 {
	neighbor := make([]uint8, len(bits))
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
