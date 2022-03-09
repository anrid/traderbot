//
// Simulates yield farming LUNA/OSMO LP on Osmosis DEX.
// - Average APR: 75%.
// - Starting from Jul 1, 2021.
// - Daily harvesting and compounding for max 365 days.
//
package main

import (
	"log"
	"time"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	startDate := "2021-07-01"     // Date of initial investment. We start farming from this date.
	harvestDays := 365            // Number of days to harvest and compound yields.
	initialInvestment := 10_000.0 // Initial investment in USD.
	apr := 99.0                   // Farm APR.

	cg := coingecko.New(coingecko.USD)

	a, err := cg.MarketChartWithCache("terra-luna", uint(harvestDays), jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	b, err := cg.MarketChartWithCache("osmosis", uint(harvestDays), jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	farm, err := trade.NewLPFarm(a, b, coingecko.USD, initialInvestment, startDate, apr)
	if err != nil {
		log.Fatal(err)
	}

	farm.SetAPRChangeRateAtHarvest(0.15) // Lower APR by 0.15 percentage points every day.

	from := timeseries.ToTime(startDate)
	for i := 1; i <= harvestDays; i++ {
		current := from.Add(time.Duration(i) * 24 * time.Hour)
		if current.After(time.Now()) {
			break
		}

		date := timeseries.ToDate(current)

		yield, err := farm.Harvest(date)
		if err != nil {
			log.Fatal(err)
		}

		farm.AddLP(date, yield) // Compound yield!
	}
}
