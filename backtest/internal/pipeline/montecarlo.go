package pipeline

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/akagr/finance-tools/backtest/internal/engine"
	"github.com/akagr/finance-tools/backtest/internal/report"
	"github.com/akagr/finance-tools/backtest/internal/series"
)

// BuildMonteCarlo estimates how much of a strategy's result could be luck. It
// runs the strategy once, takes its daily returns, and then bootstraps
// (resamples with replacement) those returns `trials` times to build a whole
// distribution of plausible outcomes. The spread answers the question a single
// backtest can't: "if the same daily edge had shown up in a different order and
// mix, how good — or bad — could it have looked?"
//
// A tight distribution comfortably above zero is reassuring; one that is wide or
// straddles zero means the headline number leaned heavily on luck. The bootstrap
// is i.i.d. (it shuffles days independently), so it captures sampling luck but
// not autocorrelation or regime structure — read it alongside walk-forward.
func BuildMonteCarlo(opts Options, trials int, seed int64) (report.MonteCarlo, error) {
	if trials < 100 {
		return report.MonteCarlo{}, fmt.Errorf("pipeline: monte-carlo needs >= 100 trials, got %d", trials)
	}
	if opts.Strategy == "all" {
		return report.MonteCarlo{}, fmt.Errorf("pipeline: monte-carlo needs a single strategy, got %q", opts.Strategy)
	}

	all, err := series.Load(opts.PricesPath)
	if err != nil {
		return report.MonteCarlo{}, err
	}
	if len(all) == 0 {
		return report.MonteCarlo{}, fmt.Errorf("pipeline: no price series in %s", opts.PricesPath)
	}
	s, err := pick(all, opts.Symbol)
	if err != nil {
		return report.MonteCarlo{}, err
	}
	if len(s.Points) < 2 {
		return report.MonteCarlo{}, fmt.Errorf("pipeline: series %q has %d bars, need >=2", s.Label, len(s.Points))
	}

	strat, err := buildStrategy(opts)
	if err != nil {
		return report.MonteCarlo{}, err
	}
	if strat, err = maybeVolTarget(strat, opts); err != nil {
		return report.MonteCarlo{}, err
	}

	capital := opts.InitialCapital
	if capital <= 0 {
		capital = 100000
	}
	costs := opts.Costs
	if costs == (engine.Costs{}) {
		costs = engine.DefaultCosts()
	}
	cfg := engine.Config{InitialCapital: capital, Costs: costs}

	res, err := engine.Run(s, strat, cfg)
	if err != nil {
		return report.MonteCarlo{}, err
	}

	// Daily simple returns of the actual equity curve.
	rets := make([]float64, 0, len(res.Equity)-1)
	for i := 1; i < len(res.Equity); i++ {
		if res.Equity[i-1] > 0 {
			rets = append(rets, res.Equity[i]/res.Equity[i-1]-1)
		}
	}
	if len(rets) < 2 {
		return report.MonteCarlo{}, fmt.Errorf("pipeline: not enough return observations for monte-carlo")
	}
	years := yearsBetween(res.Dates[0], res.Dates[len(res.Dates)-1])

	// Actual (observed) metrics.
	actRet, actCAGR, actDD, actSharpe := pathMetrics(rets, years)

	// Bootstrap.
	rng := rand.New(rand.NewSource(seed))
	n := len(rets)
	sampleRet := make([]float64, trials)
	sampleCAGR := make([]float64, trials)
	sampleDD := make([]float64, trials)
	sampleSharpe := make([]float64, trials)
	profitable := 0
	draw := make([]float64, n)
	for t := 0; t < trials; t++ {
		for i := 0; i < n; i++ {
			draw[i] = rets[rng.Intn(n)]
		}
		r, c, dd, sh := pathMetrics(draw, years)
		sampleRet[t], sampleCAGR[t], sampleDD[t], sampleSharpe[t] = r, c, dd, sh
		if r > 0 {
			profitable++
		}
	}

	mkMetric := func(name string, actual float64, xs []float64, pct bool) report.MCMetric {
		sort.Float64s(xs)
		return report.MCMetric{
			Name: name, Percent: pct, Actual: actual,
			P5: percentile(xs, 0.05), P25: percentile(xs, 0.25), P50: percentile(xs, 0.50),
			P75: percentile(xs, 0.75), P95: percentile(xs, 0.95),
		}
	}

	// Build metrics first (this sorts each sample slice in place), then read the
	// now-sorted return distribution for the summary note.
	mcMetrics := []report.MCMetric{
		mkMetric("Total return", actRet, sampleRet, true),
		mkMetric("CAGR", actCAGR, sampleCAGR, true),
		mkMetric("Max drawdown", actDD, sampleDD, true),
		mkMetric("Sharpe", actSharpe, sampleSharpe, false),
	}

	profShare := float64(profitable) / float64(trials)
	notes := []string{fmt.Sprintf(
		"%.0f%% of %d bootstrap trials were profitable, and the middle 90%% of total returns span %s to %s. A distribution that hugs or crosses zero means the headline result leaned on luck; a tight, clearly-positive spread is what a real edge looks like.",
		profShare*100, trials, pctString(percentile(sampleRet, 0.05)), pctString(percentile(sampleRet, 0.95)))}
	notes = append(notes, "Bootstrap shuffles days independently, so it measures sampling luck, not regime or autocorrelation effects — read it next to walk-forward, not instead of it.")

	return report.MonteCarlo{
		Meta: report.MCMeta{
			Symbol:   s.Label,
			Strategy: res.Strategy,
			Start:    res.Dates[0],
			End:      res.Dates[len(res.Dates)-1],
			Bars:     len(s.Points),
			Trials:   trials,
			Seed:     seed,
			Notes:    notes,
		},
		Metrics: mcMetrics,
	}, nil
}

// pathMetrics computes total return, CAGR, max drawdown and Sharpe for a series
// of daily returns compounded from 1.0.
func pathMetrics(rets []float64, years float64) (totalReturn, cagr, maxDD, sharpe float64) {
	equity := 1.0
	peak := 1.0
	var mean float64
	for _, r := range rets {
		equity *= 1 + r
		if equity > peak {
			peak = equity
		}
		if peak > 0 {
			if dd := (peak - equity) / peak; dd > maxDD {
				maxDD = dd
			}
		}
		mean += r
	}
	totalReturn = equity - 1
	mean /= float64(len(rets))
	var ss float64
	for _, r := range rets {
		d := r - mean
		ss += d * d
	}
	if len(rets) > 1 {
		sd := math.Sqrt(ss / float64(len(rets)-1))
		if sd > 0 {
			sharpe = mean / sd * math.Sqrt(tradingDaysPerYear)
		}
	}
	if years > 0 && equity > 0 {
		cagr = math.Pow(equity, 1.0/years) - 1
	}
	return totalReturn, cagr, maxDD, sharpe
}

// percentile returns the p-quantile (0..1) of an already-sorted slice using
// nearest-rank; xs must be non-empty.
func percentile(xs []float64, p float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	idx := int(math.Round(p * float64(len(xs)-1)))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(xs) {
		idx = len(xs) - 1
	}
	return xs[idx]
}

const tradingDaysPerYear = 252.0

// yearsBetween returns the calendar span between two YYYY-MM-DD dates in years.
func yearsBetween(start, end string) float64 {
	const layout = "2006-01-02"
	s, err1 := time.Parse(layout, start)
	e, err2 := time.Parse(layout, end)
	if err1 != nil || err2 != nil {
		return 0
	}
	return e.Sub(s).Hours() / 24.0 / 365.25
}

func pctString(v float64) string {
	return fmt.Sprintf("%.1f%%", v*100)
}
