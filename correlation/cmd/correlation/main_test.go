package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/akagr/finance-tools/correlation/internal/yahoo"
)

// fakeFetcher returns canned bars per Yahoo symbol (or an error).
type fakeFetcher struct {
	bars map[string][]yahoo.Bar
	errs map[string]error
}

func (f fakeFetcher) Chart(_ context.Context, symbol string, _, _ time.Time) ([]yahoo.Bar, error) {
	if err, ok := f.errs[symbol]; ok {
		return nil, err
	}
	return f.bars[symbol], nil
}

func TestParseFetchTickers(t *testing.T) {
	in := `# comment
VWRA     VWRA.L   USD

Nifty50  ^NSEI    INR
GOLDBEES GOLDBEES.NS
BAD
`
	got, err := parseFetchTickers(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	want := []fetchTicker{
		{Label: "VWRA", Yahoo: "VWRA.L", Currency: "USD"},
		{Label: "Nifty50", Yahoo: "^NSEI", Currency: "INR"},
		{Label: "GOLDBEES", Yahoo: "GOLDBEES.NS", Currency: "USD"}, // default currency
	}
	if len(got) != len(want) {
		t.Fatalf("got %d tickers, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ticker %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestFetchPricesTo(t *testing.T) {
	tickers := []fetchTicker{
		{Label: "VWRA", Yahoo: "VWRA.L", Currency: "USD"},
		{Label: "Nifty50", Yahoo: "^NSEI", Currency: "INR"},
	}
	f := fakeFetcher{bars: map[string][]yahoo.Bar{
		"VWRA.L": {{Date: "2024-01-02", Close: 100}, {Date: "2024-01-03", Close: 101.5}},
		"^NSEI":  {{Date: "2024-01-02", Close: 21500.75}},
	}}

	var buf bytes.Buffer
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	if err := fetchPricesTo(context.Background(), f, &buf, tickers, start, end); err != nil {
		t.Fatal(err)
	}
	want := "date,symbol,close,currency\n" +
		"2024-01-02,VWRA,100.0000,USD\n" +
		"2024-01-03,VWRA,101.5000,USD\n" +
		"2024-01-02,Nifty50,21500.7500,INR\n"
	if buf.String() != want {
		t.Errorf("CSV mismatch:\n got: %q\nwant: %q", buf.String(), want)
	}
}

func TestFetchPricesToErrorsOnFetchFailure(t *testing.T) {
	tickers := []fetchTicker{{Label: "X", Yahoo: "X", Currency: "USD"}}
	f := fakeFetcher{errs: map[string]error{"X": context.DeadlineExceeded}}
	var buf bytes.Buffer
	if err := fetchPricesTo(context.Background(), f, &buf, tickers, time.Now(), time.Now()); err == nil {
		t.Fatal("expected an error when a fetch fails")
	}
}

func TestFetchFXTo(t *testing.T) {
	f := fakeFetcher{bars: map[string][]yahoo.Bar{
		"INR=X": {{Date: "2024-01-02", Close: 83.25}, {Date: "2024-01-03", Close: 83.5}},
	}}
	var buf bytes.Buffer
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	if err := fetchFXTo(context.Background(), f, &buf, []string{"USD:INR=X"}, start, end); err != nil {
		t.Fatal(err)
	}
	want := "date,currency,rate\n" +
		"2024-01-02,USD,83.2500\n" +
		"2024-01-03,USD,83.5000\n"
	if buf.String() != want {
		t.Errorf("CSV mismatch:\n got: %q\nwant: %q", buf.String(), want)
	}
}

func TestFetchFXToBadSpec(t *testing.T) {
	f := fakeFetcher{}
	var buf bytes.Buffer
	if err := fetchFXTo(context.Background(), f, &buf, []string{"nocolon"}, time.Now(), time.Now()); err == nil {
		t.Fatal("expected an error for a spec without a colon")
	}
}

func TestParseRange(t *testing.T) {
	if _, _, code := parseRange("", "2024-12-31"); code != 2 {
		t.Errorf("missing start: code = %d, want 2", code)
	}
	if _, _, code := parseRange("2024-01-01", "bad"); code != 2 {
		t.Errorf("bad end: code = %d, want 2", code)
	}
	if _, _, code := parseRange("2024-12-31", "2024-01-01"); code != 2 {
		t.Errorf("end before start: code = %d, want 2", code)
	}
	if _, _, code := parseRange("2024-01-01", "2024-12-31"); code != 0 {
		t.Errorf("valid range: code = %d, want 0", code)
	}
}

func TestCmdFetchDispatch(t *testing.T) {
	if got := cmdFetch(nil); got != 2 {
		t.Errorf("no subcommand: got %d, want 2", got)
	}
	if got := cmdFetch([]string{"bogus"}); got != 2 {
		t.Errorf("unknown subcommand: got %d, want 2", got)
	}
}
