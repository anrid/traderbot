package main

import (
	"fmt"
	"log"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	cg := coingecko.New()

	c, err := cg.MarketChartWithCache("bitcoin", 180, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	c.Prices.Print()

	ema9d := trade.NewEMAIndicator(9, c.Prices)
	ema21d := trade.NewEMAIndicator(21, c.Prices)

	fmt.Printf("EMA 9d  : %.04f\n", ema9d.ForDate("2021-11-10"))
	fmt.Printf("EMA 21d : %.04f\n", ema21d.ForDate("2021-11-10"))

	trade.NewEMACrossOverStrategy(ema9d, ema21d, c.Prices)
}
