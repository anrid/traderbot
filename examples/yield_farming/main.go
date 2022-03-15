// Simulates yield farming an LP pair. See flags for details.
//
package main

import (
	"log"
	"os"
	"time"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/anrid/traderbot/pkg/trade"
	"github.com/spf13/pflag"
)

func main() {
	path := pflag.StringP("path", "p", "", "path to output dir (required, e.g. /mnt/c/Users/whatever/)")
	assetAID := pflag.StringP("asset-a", "a", "terra-luna", "CoinGecko ID of asset A (default: terra-luna)")
	assetBID := pflag.StringP("asset-b", "b", "osmosis", "CoinGecko ID of asset B (default: osmosis)")
	startDate := pflag.StringP("start-date", "d", "2021-07-01", "start farming on date (default: 2021-07-01)")
	harvestDays := pflag.Int("harvest-days", 365, "number of days to harvest and compound yields (default: 365)")
	apr := pflag.Float64("apr", 100.0, "APR to use for farm (default: 100.0)")
	finalAPR := pflag.Float64("final-apr", 0.0, "APR will gradually change to reach this final value at the last harvest date (ignored if <= 0)")

	pflag.Parse()

	if *path == "" {
		pflag.PrintDefaults()
		os.Exit(-1)
	}

	cg := coingecko.New(coingecko.USD)

	a, err := cg.MarketChartWithCache(*assetAID, uint(*harvestDays), jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	b, err := cg.MarketChartWithCache(*assetBID, uint(*harvestDays), jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	initialInvestment := 10_000.0

	farm, err := trade.NewLPFarm(a, b, coingecko.USD, initialInvestment, *startDate, *apr)
	if err != nil {
		log.Fatal(err)
	}

	var harvestDates []string
	from := timeseries.ToTime(*startDate)
	for i := 1; i <= *harvestDays; /* number of days to harvest and compound yields */ i++ {
		cur := from.Add(time.Duration(i) * 24 * time.Hour)
		if cur.After(time.Now()) {
			break
		}
		harvestDates = append(harvestDates, timeseries.ToDate(cur))
	}

	// Calculate APR change based on the number of times we harvest.
	// The goal is to have the APR reach a certain target value by the final
	// harvest date.
	{
		if *finalAPR > 0.0 {
			change := (*apr - *finalAPR) / float64(len(harvestDates))
			farm.SetAPRChangeRateAtHarvest(change)
		}
	}

	for _, d := range harvestDates {
		yield, err := farm.Harvest(d)
		if err != nil {
			log.Fatal(err)
		}
		farm.AddLP(d, yield, false) // Compound yield!
	}

	err = trade.RenderYieldFarmingPerformanceChart(*path, farm)
	if err != nil {
		log.Fatal(err)
	}
}
