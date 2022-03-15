package main

import (
	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	fc := trade.NewForecast(coingecko.USD, 10_000.0, 100)

	luna := fc.CreateMarket("Terra", "LUNA", 100.0, []trade.PriceChange{
		{IncPct: 150.0, IncDays: 60},
		{DecPct: 50.0, DecDays: 90},
	})

	// luna.Prices.PrintSample(30)

	ust := fc.CreateMarket("Terra USD", "UST", 1.0, nil)

	// ust.Prices.PrintSample(30)

	fc.AddLPFarm(luna, ust, 100.0, 40.0, 1000.0 /* Invest an additional $1,000 USD into the farm at the beginning of every month */)
}
