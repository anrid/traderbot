package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	days := 10

	fc := trade.NewForecast(coingecko.USD, 10_000.0, days)

	luna := fc.CreateMarket("Terra", "LUNA", 100.0, []trade.PriceChange{
		{IncPct: 150.0, IncDays: 60},
		{DecPct: 50.0, DecDays: 90},
	})

	osmo := fc.CreateMarket("Osmosis", "OSMO", 8.00, []trade.PriceChange{
		{IncPct: 125.0, IncDays: 30},
		{DecPct: 40.0, DecDays: 60},
	})

	ust := fc.CreateMarket("Terra USD", "UST", 1.0, nil)

	fc.AddLPFarm(luna, ust, 100.0 /* APR */, 40.0 /* Final APR */, 1000.0 /* Invest an additional $1,000 USD into the farm at the beginning of every month */)
	fc.AddLPFarm(luna, osmo, 125.0 /* APR */, 60.0 /* Final APR */, 0.0)

	j, err := fc.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	js := "const data = " + j + "\n\nexport default data\n"

	_, filename, _, _ := runtime.Caller(0)
	jsFile := filepath.Join(filepath.Dir(filename), "..", "..", "web", "src", "data", "forecast1.js")

	err = ioutil.WriteFile(jsFile, []byte(js), 0777)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Wrote forecast data to: %s\n", jsFile)
}
