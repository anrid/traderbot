package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/trade"
	"github.com/spf13/pflag"
)

func main() {
	listOnly := pflag.BoolP("list", "l", false, "List top 100 markets (coins) on CoinGecko")
	useEMS921 := pflag.BoolP("ems921", "9", true, "Use '9-Day/21-Day EMS CrossOver' strategy (default: true)")
	ids := pflag.StringSliceP("ids", "i",
		[]string{"terra-luna", "solana", "bitcoin", "ethereum"},
		"Coin IDs to trade (default: [\"terra-luna\", \"solana\", \"bitcoin\", \"ethereum\"]",
	)
	periodInDays := pflag.UintP("days", "d", 365, "Start trading X number of days ago (default: 365)")

	pflag.Parse()

	cg := coingecko.New(coingecko.USD)

	if *listOnly {
		markets, err := cg.Markets()
		if err != nil {
			log.Fatal(err)
		}

		for i, m := range markets {
			fmt.Printf("%4d. %s (%s)  --  id: %s\n", i+1, m.Name, strings.ToUpper(m.Symbol), m.ID)
		}
		fmt.Println("")
		return
	}

	for _, id := range *ids {
		m, err := cg.MarketChartWithCache(id, *periodInDays, jsoncache.InvalidateDaily)
		if err != nil {
			log.Fatal(err)
		}

		if *useEMS921 {
			ema9d := trade.NewEMAIndicator(9, m.Prices)
			ema21d := trade.NewEMAIndicator(21, m.Prices)

			s, err := trade.NewEMACrossOverStrategy(ema9d, ema21d, m, m)
			if err != nil {
				log.Fatal(err)
			}

			initialInvestment := 10_000.0 // USD.

			trade.ExecuteTradesAndPrint("9-Day/21-Day EMS CrossOver Strategy", initialInvestment, s.Trades)
		}
	}
}
