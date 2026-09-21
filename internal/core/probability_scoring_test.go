package core

import (
	"math"
	"testing"
)

func TestProbabilityScoringMultipliesInLogSpace(t *testing.T) {
	scoring := newProbabilityScoring(3)
	if err := scoring.InitInRange(0, 3, 0.8); err != nil {
		t.Fatal(err)
	}
	if err := scoring.MultiplyInRange(1, 3, 0.5); err != nil {
		t.Fatal(err)
	}

	got := scoring.Items()
	want := []float64{0.8, 0.4, 0.4}
	for index := range want {
		if math.Abs(got[index]-want[index]) > 1e-9 {
			t.Fatalf("score %d = %f, want %f", index, got[index], want[index])
		}
	}
}
