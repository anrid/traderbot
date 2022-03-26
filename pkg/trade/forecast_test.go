package trade

import (
	"testing"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/stretchr/testify/require"
)

func TestForecast(t *testing.T) {
	r := require.New(t)

	forecastDays := 10

	fc := NewForecast(coingecko.USD, 10_000.0, forecastDays)

	aaPrice := 100.0

	aa := fc.CreateMarket("Asset A", "AA", aaPrice, []PriceChange{
		{IncPct: 20, IncDays: 5},
		{DecPct: 20, DecDays: 5},
	})

	r.Equal(11, len(aa.Prices))

	{
		dailyIncrease := aaPrice * (20.0 / 5 / 100.0)

		r.Equal(aaPrice, aa.Prices[0].V)

		r.Equal(aaPrice+(1*dailyIncrease), aa.Prices[1].V)
		r.Equal(aaPrice+(2*dailyIncrease), aa.Prices[2].V)
		r.Equal(aaPrice+(3*dailyIncrease), aa.Prices[3].V)
		r.Equal(aaPrice+(4*dailyIncrease), aa.Prices[4].V)
		r.Equal(aaPrice+(5*dailyIncrease), aa.Prices[5].V)

		aaPrice = aa.Prices[5].V

		dailyDecrease := -aaPrice * (20.0 / 5 / 100.0)

		r.Equal(aaPrice+(1*dailyDecrease), aa.Prices[6].V)
		r.Equal(aaPrice+(2*dailyDecrease), aa.Prices[7].V)
		r.Equal(aaPrice+(3*dailyDecrease), aa.Prices[8].V)
		r.Equal(aaPrice+(4*dailyDecrease), aa.Prices[9].V)
		r.Equal(aaPrice+(5*dailyDecrease), aa.Prices[10].V)
	}

	abPrice := 100.0

	ab := fc.CreateMarket("Asset B", "AB", abPrice, []PriceChange{
		{IncPct: 25, IncDays: 5},
		{DecPct: 20, DecDays: 5},
	})

	r.Equal(11, len(ab.Prices))

	{
		dailyIncrease := abPrice * (25.0 / 5 / 100.0)

		r.Equal(abPrice, ab.Prices[0].V)

		r.Equal(abPrice+(1*dailyIncrease), ab.Prices[1].V)
		r.Equal(abPrice+(2*dailyIncrease), ab.Prices[2].V)
		r.Equal(abPrice+(3*dailyIncrease), ab.Prices[3].V)
		r.Equal(abPrice+(4*dailyIncrease), ab.Prices[4].V)
		r.Equal(abPrice+(5*dailyIncrease), ab.Prices[5].V)

		abPrice = ab.Prices[5].V

		dailyDecrease := -abPrice * (20.0 / 5 / 100.0)

		r.Equal(abPrice+(1*dailyDecrease), ab.Prices[6].V)
		r.Equal(abPrice+(2*dailyDecrease), ab.Prices[7].V)
		r.Equal(abPrice+(3*dailyDecrease), ab.Prices[8].V)
		r.Equal(abPrice+(4*dailyDecrease), ab.Prices[9].V)
		r.Equal(abPrice+(5*dailyDecrease), ab.Prices[10].V)

		r.Equal(100.0, ab.Prices[10].V)
	}
}
