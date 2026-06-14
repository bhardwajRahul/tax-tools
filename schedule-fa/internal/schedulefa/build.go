// Package schedulefa assembles the final Schedule FA rows (Table A2 + A3) from a
// parsed statement, the FX store, and computed peaks. Every INR figure carries
// its fx.Conversion so the report can print a full audit trail.
package schedulefa

import (
	"errors"

	"github.com/akagr/tax-tools/schedule-fa/internal/fx"
	"github.com/akagr/tax-tools/schedule-fa/internal/model"
	"github.com/akagr/tax-tools/schedule-fa/internal/peak"
)

// ErrNotImplemented is returned by stubs not yet built.
var ErrNotImplemented = errors.New("schedulefa: not implemented")

// A2Row is the custodial-account row (the IBKR account itself).
type A2Row struct {
	Institution    string
	Address        string
	ZIP            string
	CountryCode    string
	AccountNumber  string
	Status         string
	OpenDate       string
	PeakBalance    fx.Conversion
	ClosingBalance fx.Conversion
	GrossCredited  fx.Conversion // interest + dividends credited during the year
}

// A3Row is one security held at any time during the calendar year.
type A3Row struct {
	CountryName   string
	CountryCode   string
	EntityName    string
	Address       string
	ZIP           string
	NatureEntity  string
	AcquiredOn    string
	InitialValue  fx.Conversion
	PeakValue     fx.Conversion
	ClosingValue  fx.Conversion
	GrossDividend fx.Conversion // gross amount paid/credited
	SaleProceeds  fx.Conversion // gross proceeds from sale/redemption
	NeedsReview   bool          // e.g. approximate peak, corporate action, missing metadata
	ReviewNote    string
}

// Report is the complete Schedule FA output for one calendar year.
type Report struct {
	Year int
	A2   []A2Row
	A3   []A3Row
}

// Build assembles the report. Implemented in M3 (A3) / M5 (A2 + edge cases).
func Build(s *model.Statement, store fx.Store, peaks []peak.Result) (*Report, error) {
	return nil, ErrNotImplemented
}
