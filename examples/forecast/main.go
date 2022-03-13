package main

import (
	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	fc := trade.NewForecast(coingecko.USD, 1_000.0, 360)

	luna := fc.CreateMarket("Terra", "LUNA", 100.0, []trade.PriceChange{
		{IncPct: 80.0, IncDays: 30},
		{DecPct: 40.0, DecDays: 60},
		{IncPct: 80.0, IncDays: 30},
		{DecPct: 40.0, DecDays: 60},
		{IncPct: 80.0, IncDays: 30},
		{DecPct: 40.0, DecDays: 60},
		{IncPct: 80.0, IncDays: 30},
		{DecPct: 40.0, DecDays: 60},
	})

	coingecko.Dump(luna.Prices)
}
