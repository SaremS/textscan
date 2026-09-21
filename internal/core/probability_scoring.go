package core

import (
	"fmt"
	"math"
)

type ProbabilityScoring struct {
	items []float64
}

func newProbabilityScoring(capacity int) *ProbabilityScoring {
	items := make([]float64, capacity)
	for i := range items {
		items[i] = -math.MaxFloat64
	}

	return &ProbabilityScoring{items: items}
}

func (s *ProbabilityScoring) InitInRange(start, end int, probability float64) error {
	if start < 0 || end > len(s.items) || start > end {
		return fmt.Errorf("invalid range [%d:%d)", start, end)
	}

	value := math.Log(clampProbability(probability))
	for index := start; index < end; index++ {
		s.items[index] = value
	}

	return nil
}

func (s *ProbabilityScoring) MultiplyInRange(start, end int, probability float64) error {
	if start < 0 || end > len(s.items) || start > end {
		return fmt.Errorf("invalid range [%d:%d)", start, end)
	}

	value := math.Log(clampProbability(probability))
	for index := start; index < end; index++ {
		s.items[index] += value
	}

	return nil
}

func (s *ProbabilityScoring) Items() []float64 {
	output := make([]float64, len(s.items))
	for index, value := range s.items {
		output[index] = math.Exp(value)
	}

	return output
}

func clampProbability(probability float64) float64 {
	if probability <= 0 {
		return 1e-6
	}
	if probability >= 1 {
		return 1 - 1e-6
	}
	return probability
}
