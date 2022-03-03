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
	cg := coingecko.New(coingecko.USD)

	var periodInDays uint = 365

	a, err := cg.MarketChartWithCache("terra-luna", periodInDays, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	b, err := cg.MarketChartWithCache("osmosis", periodInDays, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	startDate := "2021-07-01"
	harvestDays := 365
	initialInvestment := 10_000.0

	farm, err := trade.NewLPFarm(a, b, coingecko.USD, initialInvestment, startDate, 99.0 /* APR */)
	if err != nil {
		log.Fatal(err)
	}

	farm.SetAPRChangeRateAtHarvest(0.15) // Lower APR by 0.15 percentage points every day.

	from := timeseries.ToTime(startDate)
	for i := 1; i <= harvestDays; /* number of days to harvest and compound yields */ i++ {
		cur := from.Add(time.Duration(i) * 24 * time.Hour)
		if cur.After(time.Now()) {
			break
		}
		farm.Harvest(timeseries.ToDate(cur))
	}
}
