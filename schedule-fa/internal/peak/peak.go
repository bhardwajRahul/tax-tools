// Package peak computes the "peak value of investment during the period" for
// each security — the hardest figure in Schedule FA.
//
//   - Mode C (v1, M3): approximate. peak ≈ max(closing value, max trade-day value).
//     Cheap and defensible for buy-and-hold; always flagged "approximate".
//   - Mode B (M4): exact. Reconstruct a daily share-count series from trades and
//     multiply by a daily price series and the day's TTBR, taking the max in INR.
//
// Peak must be maximised in INR, not USD, because the rupee can move independently
// of the share price.
package peak

import (
	"errors"
	"math/big"
	"time"

	"github.com/akagr/tax-tools/schedule-fa/internal/model"
)

// ErrNotImplemented is returned by stubs not yet built.
var ErrNotImplemented = errors.New("peak: not implemented")

// Mode selects the computation strategy.
type Mode int

const (
	ModeApprox Mode = iota // C — v1 default
	ModeExact              // B — M4, needs a PriceProvider
)

// Result is the peak valuation plus whether it was approximated.
type Result struct {
	Instrument   model.Instrument
	ValueINR     *big.Rat
	On           time.Time // date of the peak
	Approximated bool      // true in Mode C → manual-review flag in the report
}

// PriceProvider yields the per-unit close for an instrument on a date (Mode B).
type PriceProvider interface {
	PriceOn(inst model.Instrument, date time.Time) (model.Money, error)
}

// Compute returns the peak value per security for the statement's year.
// Implemented in M3 (approx) / M4 (exact).
func Compute(s *model.Statement, mode Mode, prices PriceProvider) ([]Result, error) {
	return nil, ErrNotImplemented
}
