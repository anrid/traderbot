package main

import (
	"log"

	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/messari"
	"github.com/spf13/pflag"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	token := pflag.StringP("token", "t", "", "Messari API token (required)")

	pflag.Parse()

	if *token == "" {
		pflag.PrintDefaults()
		log.Fatal("--token or -t flag required but missing")
	}

	m := messari.New(*token)

	assets, err := m.AssetsWithCache(jsoncache.InvalidateMonthly)
	if err != nil {
		log.Fatal(err)
	}

	pr := message.NewPrinter(language.English)

	var count int
	for _, a := range assets {
		mc := a.Metrics.Supply.Circulating * a.Metrics.MarketData.PriceUSD
		if mc == 0 {
			continue
		}

		isKnownFutureSupply := a.Metrics.Supply.AnnualInflationPct > 0 || a.Metrics.Supply.YPlus10IssuedPct > 0

		if isKnownFutureSupply {
			count++
			pr.Printf("%4d. Asset: %-8s %-30s   --  Price: $%-10.03f  Inflation: %8.02f%%  Y+10 Supply: %8.02f%%  MC: %12.f\n",
				count, a.Symbol, a.Name, a.Metrics.MarketData.PriceUSD,
				a.Metrics.Supply.AnnualInflationPct, a.Metrics.Supply.YPlus10IssuedPct, mc,
			)
		}
	}
}
