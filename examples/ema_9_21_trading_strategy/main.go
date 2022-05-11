//
// Executes a basic trading strategy based on the Exponential Moving Average indicator (EMA).
// - When 9-day EMA crosses over the 21-day EMA indicator from below: Buy.
// - When 9-day EMA crosses over the 21-day EMA indicator from above: Sell.
//
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/trade"
	"github.com/spf13/pflag"
)

func main() {
	trackID := pflag.String("track", "terra-luna", "CoinGecko ID of market to track EMS 9/21-day crossover indicator (default: terra-luna)")
	tradeID := pflag.String("trade", "terra-luna", "CoinGecko ID of market to trade (default: terra-luna)")
	initialInvestment := pflag.Float64("invest", 10_000.00, "Initial investment (default: 10,000.00 USD)")
	startDaysAgo := pflag.UintP("start", "d", 365, "Start our trading strategy 365 days ago counting from today (default: 365)")
	startDate := pflag.String("date", "", "Start our trading strategy from this date, format YYYY-MM-DD")

	pflag.Parse()

	cg := coingecko.New(coingecko.USD)

	if *startDate != "" {
		t, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			log.Fatal(err)
		}

		*startDaysAgo = uint(time.Since(t).Hours() / 24)
	}

	fmt.Printf("Simulating trading %d days ago\n", *startDaysAgo)

	tracking, err := cg.MarketChartWithCache(*trackID, *startDaysAgo, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	trading, err := cg.MarketChartWithCache(*tradeID, *startDaysAgo, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	ema9d := trade.NewEMAIndicator(9, tracking.Prices)
	ema21d := trade.NewEMAIndicator(21, tracking.Prices)

	s, err := trade.NewEMACrossOverStrategy(ema9d, ema21d, tracking, trading)
	if err != nil {
		log.Fatal(err)
	}

	trade.ExecuteTradesAndPrint("9-Day/21-Day EMS CrossOver Strategy", *initialInvestment, s.Trades)
}
