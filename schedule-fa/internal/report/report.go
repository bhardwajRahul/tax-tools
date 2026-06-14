// Package report renders a schedulefa.Report to various formats. Each format
// includes a companion audit view (USD figure, TTBR, rate date) behind every
// INR value, plus a reconciliation summary and manual-review flags.
package report

import (
	"errors"
	"io"

	"github.com/akagr/tax-tools/schedule-fa/internal/schedulefa"
)

// ErrNotImplemented is returned by stubs not yet built.
var ErrNotImplemented = errors.New("report: not implemented")

// Format is an output format.
type Format string

const (
	Markdown Format = "md"
	CSV      Format = "csv"
	JSON     Format = "json"
)

// Renderer writes a report in one format.
type Renderer interface {
	Render(w io.Writer, r *schedulefa.Report) error
}

// For returns the Renderer for a format. Implemented in M7 (M3 ships md+csv+json).
func For(f Format) (Renderer, error) {
	return nil, ErrNotImplemented
}
