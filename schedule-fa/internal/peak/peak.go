// Package peak computes the "peak value of investment during the period" for
// each security — the hardest figure in Schedule FA.
//
//   - Mode C (v1, M3): approximate. Without a daily price series, the position is
//     valued at every point where a real market price is known — each trade
//     (using the trade price, both just before and just after the fill) and the
//     31-Dec close (using the mark price) — and the maximum INR value is taken.
//     This is exact for buy-and-hold and a documented approximation otherwise.
//   - Mode B (M4): exact. Reconstruct a daily share-count series and value it
//     against a daily price series and the day's TTBR.
//
// Peak is maximised in INR (units × price × TTBR), not USD, because the rupee can
// move independently of the share price.
package peak

import (
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/akagr/tax-tools/schedule-fa/internal/fx"
	"github.com/akagr/tax-tools/schedule-fa/internal/model"
)

// Mode selects the computation strategy.
type Mode int

const (
	ModeApprox Mode = iota // C — v1 default
	ModeExact              // B — M4, needs a PriceProvider
)

// Result is the peak valuation for one security.
type Result struct {
	Instrument   model.Instrument
	Peak         fx.Conversion // the winning candidate, with its audit trail
	HasValue     bool          // false if no candidate could be valued (no price/rate)
	Approximated bool          // true in Mode C → manual-review flag in the report
}

// PriceProvider yields the per-unit close for an instrument on a date (Mode B).
type PriceProvider interface {
	PriceOn(inst model.Instrument, date time.Time) (model.Money, error)
}

// Compute returns the peak value per security for the statement's year.
func Compute(s *model.Statement, store fx.Store, mode Mode, prices PriceProvider) ([]Result, error) {
	if mode == ModeExact {
		return nil, fmt.Errorf("peak: exact mode (B) not implemented yet (M4)")
	}

	yearEnd := time.Date(s.Year, time.December, 31, 0, 0, 0, 0, time.UTC)

	type group struct {
		inst       model.Instrument
		trades     []model.Trade
		closingPos *big.Rat
		markPrice  model.Money
		hasClosing bool
	}
	groups := map[string]*group{}
	var order []string
	get := func(inst model.Instrument) *group {
		k := instKey(inst)
		g := groups[k]
		if g == nil {
			g = &group{inst: inst, closingPos: new(big.Rat)}
			groups[k] = g
			order = append(order, k)
		}
		return g
	}
	for _, p := range s.OpenPositions {
		g := get(p.Instrument)
		g.closingPos = p.Quantity
		g.markPrice = p.MarkPrice
		g.hasClosing = true
	}
	for _, t := range s.Trades {
		get(t.Instrument).trades = append(get(t.Instrument).trades, t)
	}
	sort.Strings(order)

	type candidate struct {
		date  time.Time
		units *big.Rat
		price model.Money
	}

	var out []Result
	for _, k := range order {
		g := groups[k]
		sort.SliceStable(g.trades, func(i, j int) bool { return g.trades[i].Date.Before(g.trades[j].Date) })

		// Reconstruct the Jan-1 opening position by unwinding the year's trades
		// from the authoritative year-end position.
		start := new(big.Rat).Set(g.closingPos)
		for _, t := range g.trades {
			start.Sub(start, signedQty(t))
		}

		var cands []candidate
		pos := new(big.Rat).Set(start)
		for _, t := range g.trades {
			pre := new(big.Rat).Set(pos)         // held just before the fill
			pos = new(big.Rat).Add(pos, signedQty(t))
			post := new(big.Rat).Set(pos)        // held just after the fill
			cands = append(cands, candidate{t.Date, pre, t.Price}, candidate{t.Date, post, t.Price})
		}
		if g.hasClosing {
			cands = append(cands, candidate{yearEnd, g.closingPos, g.markPrice})
		}

		var best fx.Conversion
		have := false
		for _, c := range cands {
			if c.units.Sign() <= 0 || c.price.Amount == nil || c.price.Amount.Sign() <= 0 {
				continue
			}
			val := new(big.Rat).Mul(c.units, c.price.Amount)
			conv, err := fx.Convert(store, model.NewMoney(c.price.Currency, val), c.date)
			if err != nil {
				continue // no TTBR for this date; skip (Mode C is approximate anyway)
			}
			if !have || conv.Result.Amount.Cmp(best.Result.Amount) > 0 {
				best, have = conv, true
			}
		}
		out = append(out, Result{Instrument: g.inst, Peak: best, HasValue: have, Approximated: true})
	}
	return out, nil
}

func signedQty(t model.Trade) *big.Rat {
	q := new(big.Rat).Abs(t.Quantity)
	if t.Side == model.Sell {
		return q.Neg(q)
	}
	return q
}

func instKey(in model.Instrument) string {
	if in.ISIN != "" {
		return "isin:" + in.ISIN
	}
	return "sym:" + in.Symbol
}
