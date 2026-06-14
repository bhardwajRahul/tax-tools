// Package ibkr ingests Interactive Brokers data.
//
// v1 parses a downloaded Activity Flex Query in XML form (offline mode, M1).
// The Flex Web Service online pull (SendRequest/GetStatement) lands in M6.
package ibkr

import (
	"errors"
	"io"

	"github.com/akagr/tax-tools/schedule-fa/internal/model"
)

// ErrNotImplemented is returned by stubs not yet built.
var ErrNotImplemented = errors.New("ibkr: not implemented")

// ParseFlexXML reads an IBKR Activity Flex Query (XML output) and returns the
// statement constrained to `year`. Sections consumed: AccountInformation,
// OpenPositions, Trades/Lots, CashTransactions (dividends + withholding),
// SecuritiesInfo (instrument metadata). Implemented in M1.
func ParseFlexXML(r io.Reader, year int) (*model.Statement, error) {
	return nil, ErrNotImplemented
}
