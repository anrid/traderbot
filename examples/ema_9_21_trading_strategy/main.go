//
// Executes a basic trading strategy based on the Exponential Moving Average indicator (EMA).
// - When 9-day EMA crosses over the 21-day EMA indicator from below: Buy.
// - When 9-day EMA crosses over the 21-day EMA indicator from above: Sell.
//
package main

import (
	"log"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/trade"
	"github.com/spf13/pflag"
)

func main() {
	id := pflag.String("id", "terra-luna", "CoinGecko ID (default: terra-luna)")
	initialInvestment := pflag.Float64("invest", 10_000.00, "Initial investment (default: 10,000.00 USD)")
	startDaysAgo := pflag.UintP("days-ago", "d", 365, "Start our trading strategy 365 days ago counting from today (default: 365)")

	pflag.Parse()

	cg := coingecko.New(coingecko.USD)

	m, err := cg.MarketChartWithCache(*id, *startDaysAgo, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	ema9d := trade.NewEMAIndicator(9, m.Prices)
	ema21d := trade.NewEMAIndicator(21, m.Prices)

	s, err := trade.NewEMACrossOverStrategy(ema9d, ema21d, m)
	if err != nil {
		log.Fatal(err)
	}

	trade.ExecuteTradesAndPrint("9-Day/21-Day EMS CrossOver Strategy", *initialInvestment, s.Trades)
}
