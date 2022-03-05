package trade

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/stretchr/testify/require"
)

func TestLPFarm(t *testing.T) {
	r := require.New(t)

	now := time.Now().UTC()
	_24h := 24 * time.Hour

	day1 := now.UnixMilli()
	day2 := now.Add(_24h * 1).UnixMilli()
	day3 := now.Add(_24h * 2).UnixMilli()
	day4 := now.Add(_24h * 3).UnixMilli()
	day5 := now.Add(_24h * 4).UnixMilli()

	day1Date := timeseries.FromTSToDate(day1)

	pricesA := timeseries.FromTuples([][]interface{}{
		{float64(day1), 10.0},
		{float64(day2), 20.0},
		{float64(day3), 30.0},
		{float64(day4), 40.0},
		{float64(day5), 50.0},
	})

	pricesB1 := timeseries.FromTuples([][]interface{}{
		{float64(day1), 1.0},
		{float64(day2), 2.0},
		{float64(day3), 3.0},
		{float64(day4), 4.0},
		{float64(day5), 5.0},
	})

	assetA := &coingecko.Market{
		Currency: coingecko.USD,
		Name:     "Asset A",
		Symbol:   "AAA",
		Prices:   pricesA,
	}

	assetB := &coingecko.Market{
		Currency: coingecko.USD,
		Name:     "Asset B",
		Symbol:   "BBB",
		Prices:   pricesB1,
	}

	// Base case: Asset prices change proportionally over time. No IL. 0% APR.
	{
		f, err := NewLPFarm(assetA, assetB, coingecko.USD, 10_000, day1Date, 0.0)
		r.NoError(err)

		fmt.Printf("a = %f  b = %f  total = %f\n", f.UnitsA, f.UnitsB, f.TotalValue)

		f.Harvest(timeseries.FromTSToDate(day2))

		fmt.Printf("a = %f  b = %f  total = %f\n", f.UnitsA, f.UnitsB, f.TotalValue)

		r.Equal(f.InitialUnitsA, f.UnitsA)
		r.Equal(f.InitialUnitsB, f.UnitsB)
		r.Equal(20_000.0, f.TotalValue)
	}

	// 365% APR case: Asset prices change proportionally over time. No IL. 365% APR.
	{
		f, err := NewLPFarm(assetA, assetB, coingecko.USD, 10_000, day1Date, 365.0)
		r.NoError(err)

		fmt.Printf("a = %f  b = %f  total = %f\n", f.UnitsA, f.UnitsB, f.TotalValue)

		f.Harvest(timeseries.FromTSToDate(day2))

		fmt.Printf("a = %f  b = %f  total = %f\n", f.UnitsA, f.UnitsB, f.TotalValue)

		r.Equal(f.InitialUnitsA+(f.InitialUnitsA*0.01), f.UnitsA)
		r.Equal(f.InitialUnitsB+(f.InitialUnitsB*0.01), f.UnitsB)
		r.Equal(20_000.0+200.0, f.TotalValue)
	}

	// Farming X/stable pair case:
	// - Asset A prices increase 500% over 5 days
	// - Asset B is a stable coin.
	// - Lots of IL. 0% APR.
	{
		pricesB2 := timeseries.FromTuples([][]interface{}{
			{float64(day1), 1.0},
			{float64(day2), 1.0},
			{float64(day3), 1.0},
			{float64(day4), 1.0},
			{float64(day5), 1.0},
		})

		assetB.Prices = pricesB2

		f, err := NewLPFarm(assetA, assetB, coingecko.USD, 10_000, day1Date, 0.0)
		r.NoError(err)

		fmt.Printf("a = %f  b = %f  total = %f\n", f.UnitsA, f.UnitsB, f.TotalValue)

		f.Harvest(timeseries.FromTSToDate(day2))

		fmt.Printf("a = %f  b = %f  total = %f\n", f.UnitsA, f.UnitsB, f.TotalValue)

		r.Equal(354.0, math.Round(f.UnitsA))
		r.Equal(7_071.0, math.Round(f.UnitsB))
		r.Equal(14_142.0, math.Round(f.TotalValue))
		r.Equal(15_000.0, f.TotalValueHODL)
	}

}
