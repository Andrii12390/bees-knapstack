package problem

type Item struct {
	Name   string
	Weight int
	Value  int
}

type Problem struct {
	Items    []Item
	Capacity int
}

func (p *Problem) Evaluate(solution []uint8) int {
	totalWeight := 0
	totalValue := 0

	for i, taken := range solution {
		if taken == 1 {
			totalWeight += p.Items[i].Weight
			totalValue += p.Items[i].Value
		}
	}

	if totalWeight > p.Capacity {
		return 0
	}

	return totalValue
}

func (p *Problem) TakenItems(solution []uint8) []Item {
	taken := make([]Item, 0)
	for i, bit := range solution {
		if bit == 1 {
			taken = append(taken, p.Items[i])
		}
	}
	return taken
}

func (p *Problem) TotalWeight(solution []uint8) int {
	total := 0
	for i, taken := range solution {
		if taken == 1 {
			total += p.Items[i].Weight
		}
	}
	return total
}
