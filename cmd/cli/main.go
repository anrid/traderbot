package main

import (
	"log"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	cg := coingecko.New(coingecko.USD)

	var period uint = 30 * 12

	for _, id := range []string{"terra-luna", "solana", "bitcoin"} {
		c, err := cg.MarketChartWithCache(id, period, jsoncache.InvalidateDaily)
		if err != nil {
			log.Fatal(err)
		}

		ema9d := trade.NewEMAIndicator(9, c.Prices)
		ema21d := trade.NewEMAIndicator(21, c.Prices)
		// fmt.Printf("EMA 9d  : %.04f\n", ema9d.ForDate("2021-11-10"))
		// fmt.Printf("EMA 21d : %.04f\n", ema21d.ForDate("2021-11-10"))

		strat, err := trade.NewEMACrossOverStrategy(ema9d, ema21d, c)
		if err != nil {
			log.Fatal(err)
		}

		trade.ExecuteTrades("EMS 9/21 Cross Over", 10_000, strat.Trades)
	}
}
