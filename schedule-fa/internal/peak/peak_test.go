package peak

import (
	"math/big"
	"testing"
	"time"

	"github.com/akagr/tax-tools/schedule-fa/internal/fx"
	"github.com/akagr/tax-tools/schedule-fa/internal/model"
)

// fakeStore returns a constant TTBR for every date, dated to the query day.
type fakeStore struct{ rate *big.Rat }

func (f fakeStore) RateOn(cur model.Currency, d time.Time) (fx.Rate, error) {
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	return fx.Rate{Currency: cur, Date: day, INRPerUnit: f.rate}, nil
}

func usd(v int64) model.Money { return model.NewMoney(model.USD, big.NewRat(v, 1)) }
func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

var inst = model.Instrument{Symbol: "AAA", ISIN: "X1", AssetClass: "STK", ListingCtry: "US", Currency: model.USD}

func TestComputeBuyAndHold(t *testing.T) {
	yearEnd := day(2024, 12, 31)
	st := &model.Statement{
		Year:          2024,
		OpenPositions: []model.Position{{Instrument: inst, Date: yearEnd, Quantity: big.NewRat(10, 1), MarkPrice: usd(150)}},
		Trades:        []model.Trade{{Instrument: inst, Date: day(2024, 3, 15), Side: model.Buy, Quantity: big.NewRat(10, 1), Price: usd(100)}},
	}
	res, err := Compute(st, fakeStore{big.NewRat(80, 1)}, ModeApprox, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || !res[0].HasValue {
		t.Fatalf("want 1 valued result, got %+v", res)
	}
	// max(post-buy 10*100, close 10*150) = 1500; *80 = 120000, on year-end.
	if got := res[0].Peak.Result.Amount; got.Cmp(big.NewRat(120000, 1)) != 0 {
		t.Errorf("peak = %s, want 120000", got.RatString())
	}
	if !res[0].Peak.RateDate.Equal(yearEnd) {
		t.Errorf("peak date = %v, want year-end", res[0].Peak.RateDate)
	}
	if !res[0].Approximated {
		t.Error("mode C result must be marked approximated")
	}
}

func TestComputePeakBeforeSale(t *testing.T) {
	// Bought 10 @ 100, sold all 10 @ 200 intra-year (fully exited, no closing).
	st := &model.Statement{
		Year: 2024,
		Trades: []model.Trade{
			{Instrument: inst, Date: day(2024, 2, 1), Side: model.Buy, Quantity: big.NewRat(10, 1), Price: usd(100)},
			{Instrument: inst, Date: day(2024, 9, 1), Side: model.Sell, Quantity: big.NewRat(10, 1), Price: usd(200)},
		},
	}
	res, err := Compute(st, fakeStore{big.NewRat(80, 1)}, ModeApprox, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Peak is the pre-sale holding: 10*200 = 2000; *80 = 160000.
	if got := res[0].Peak.Result.Amount; got.Cmp(big.NewRat(160000, 1)) != 0 {
		t.Errorf("peak = %s, want 160000 (value just before the sale)", got.RatString())
	}
}
