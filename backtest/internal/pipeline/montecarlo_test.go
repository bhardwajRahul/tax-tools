package pipeline

import (
	"testing"
)

func mcOpts() Options {
	o := baseOpts()
	o.Strategy = "sma-cross"
	o.Fast = 3
	o.Slow = 8
	return o
}

func TestMonteCarloValidation(t *testing.T) {
	if _, err := BuildMonteCarlo(mcOpts(), 50, 1); err == nil {
		t.Error("expected error for < 100 trials")
	}
	all := baseOpts()
	all.Strategy = "all"
	if _, err := BuildMonteCarlo(all, 200, 1); err == nil {
		t.Error("expected error for --strategy all")
	}
}

func TestMonteCarloDeterministicWithSeed(t *testing.T) {
	a, err := BuildMonteCarlo(mcOpts(), 200, 42)
	if err != nil {
		t.Fatal(err)
	}
	b, err := BuildMonteCarlo(mcOpts(), 200, 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Metrics) != len(b.Metrics) {
		t.Fatal("metric count differs")
	}
	for i := range a.Metrics {
		if a.Metrics[i].P50 != b.Metrics[i].P50 || a.Metrics[i].P5 != b.Metrics[i].P5 {
			t.Errorf("metric %q not reproducible with same seed", a.Metrics[i].Name)
		}
	}
}

func TestMonteCarloDiffersBySeed(t *testing.T) {
	a, _ := BuildMonteCarlo(mcOpts(), 200, 1)
	b, _ := BuildMonteCarlo(mcOpts(), 200, 2)
	// Different seeds should give different distributions (P5 at least).
	same := true
	for i := range a.Metrics {
		if a.Metrics[i].P5 != b.Metrics[i].P5 {
			same = false
			break
		}
	}
	if same {
		t.Error("different seeds produced identical distributions")
	}
}

func TestMonteCarloPercentilesOrdered(t *testing.T) {
	mc, err := BuildMonteCarlo(mcOpts(), 500, 7)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range mc.Metrics {
		if !(m.P5 <= m.P25 && m.P25 <= m.P50 && m.P50 <= m.P75 && m.P75 <= m.P95) {
			t.Errorf("metric %q percentiles not monotonic: %v", m.Name,
				[]float64{m.P5, m.P25, m.P50, m.P75, m.P95})
		}
	}
}

func TestPercentileHelper(t *testing.T) {
	xs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if got := percentile(xs, 0.0); got != 1 {
		t.Errorf("p0 = %v, want 1", got)
	}
	if got := percentile(xs, 1.0); got != 10 {
		t.Errorf("p100 = %v, want 10", got)
	}
	if got := percentile(xs, 0.5); got != 6 { // nearest-rank of 0.5*(9)=4.5 -> round 5 -> xs[5]=6
		t.Errorf("p50 = %v, want 6", got)
	}
}
