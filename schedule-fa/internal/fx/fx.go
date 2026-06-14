// Package fx converts foreign-currency amounts to INR using the SBI TT Buying
// Rate (TTBR), and records an audit trail for every conversion.
//
// Rule: use the TTBR published for the valuation date; if none was published
// (weekend/holiday), fall back to the nearest PRECEDING working day. The actual
// date used is captured in Conversion.RateDate so the report can show its work.
package fx

import (
	"errors"
	"math/big"
	"time"

	"github.com/akagr/tax-tools/schedule-fa/internal/model"
)

// ErrNotImplemented is returned by stubs not yet built (M2).
var ErrNotImplemented = errors.New("fx: not implemented")

// Rate is a single day's TTBR for one currency, expressed as INR per 1 unit.
type Rate struct {
	Currency   model.Currency
	Date       time.Time // the date this rate was published for
	INRPerUnit *big.Rat
}

// Conversion is the audit record behind one INR figure in the report.
type Conversion struct {
	Source   model.Money // the original foreign-currency amount
	Rate     Rate        // the rate applied
	RateDate time.Time   // date actually used (== Rate.Date; may differ from the valuation date via fallback)
	Result   model.Money // the resulting INR amount
}

// Store provides TTBR lookups with preceding-working-day fallback.
type Store interface {
	// RateOn returns the TTBR for cur to apply to a value dated `date`.
	RateOn(cur model.Currency, date time.Time) (Rate, error)
}

// Convert applies the TTBR for `date` to `amount` and returns both the INR
// result and the audit record. (M2 will implement Store; this helper is final.)
func Convert(s Store, amount model.Money, date time.Time) (Conversion, error) {
	if amount.Currency == model.INR {
		return Conversion{Source: amount, Result: amount, RateDate: date}, nil
	}
	r, err := s.RateOn(amount.Currency, date)
	if err != nil {
		return Conversion{}, err
	}
	src := amount.Amount
	if src == nil {
		src = new(big.Rat)
	}
	inr := new(big.Rat).Mul(src, r.INRPerUnit)
	return Conversion{
		Source:   amount,
		Rate:     r,
		RateDate: r.Date,
		Result:   model.NewMoney(model.INR, inr),
	}, nil
}

// CSVStore loads TTBR rates from CSV files (bundled or user-supplied). M2.
type CSVStore struct {
	// rates[currency] -> sorted-by-date slice, for preceding-day lookup.
	rates map[model.Currency][]Rate
}

// NewCSVStore will load rate CSVs from the given paths. Stub until M2.
func NewCSVStore(paths ...string) (*CSVStore, error) {
	return nil, ErrNotImplemented
}

// RateOn implements Store. Stub until M2.
func (c *CSVStore) RateOn(cur model.Currency, date time.Time) (Rate, error) {
	return Rate{}, ErrNotImplemented
}
