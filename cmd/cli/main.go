package main

import (
	"log"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
)

func main() {
	cg := coingecko.New()

	c, err := cg.MarketChartWithCache("bitcoin", 90, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	c.Prices.Print()

	coingecko.NewEMAIndicator(3, c.Prices)
}
