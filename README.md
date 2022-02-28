# Trader Bot!

## Example Usage

#### EMA 9/21-Day Trading Strategy

```golang
// cmd/ema_trading_strategy_example/main.go
//
// Executes a basic trading strategy based on the Exponential Moving Average indicator (EMA).
// - When 9-day EMA crosses over the 21-day EMA indicator from below: Buy.
// - When 9-day EMA crosses over the 21-day EMA indicator from above: Sell.
//
import (
	"log"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	cg := coingecko.New(coingecko.USD)

	id := "terra-luna"          // CoinGecko ID.
	var periodInDays uint = 365 // Start our trading strategy 365 ago counting from today.

	m, err := cg.MarketChartWithCache(id, periodInDays, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	ema9d := trade.NewEMAIndicator(9, m.Prices)
	ema21d := trade.NewEMAIndicator(21, m.Prices)

	s, err := trade.NewEMACrossOverStrategy(ema9d, ema21d, m)
	if err != nil {
		log.Fatal(err)
	}

	initialInvestment := 10_000.0 // USD.

	trade.ExecuteTradesAndPrint("9-Day/21-Day EMS CrossOver Strategy", initialInvestment, s.Trades)
}
```

#### Output:

```bash
$ go run cmd/ema_trading_strategy_example/main.go 


Trading 'Terra' (LUNA) : 9-Day/21-Day EMS CrossOver Strategy
=============================================================

  1. [2021-04-27] buy  terra-luna @        17.6307  --  amount:    10,000.0000 , units:       567.1920
  2. [2021-05-13] sell terra-luna @        14.8174  --  amount:     8,404.2978 , units:       567.1920  [portfolio:     8,404.2978]
  3. [2021-07-07] buy  terra-luna @         6.4806  --  amount:     8,404.2978 , units:     1,296.8387
  4. [2021-07-21] sell terra-luna @         5.9036  --  amount:     7,656.0801 , units:     1,296.8387  [portfolio:     7,656.0801]
  5. [2021-07-23] buy  terra-luna @         7.3690  --  amount:     7,656.0801 , units:     1,038.9540
  6. [2021-09-22] sell terra-luna @        24.9468  --  amount:    25,918.5351 , units:     1,038.9540  [portfolio:    25,918.5351]
  7. [2021-09-24] buy  terra-luna @        36.4343  --  amount:    25,918.5351 , units:       711.3775
  8. [2021-10-17] sell terra-luna @        36.7548  --  amount:    26,146.5548 , units:       711.3775  [portfolio:    26,146.5548]
  9. [2021-10-21] buy  terra-luna @        42.8015  --  amount:    26,146.5548 , units:       610.8788
 10. [2021-11-19] sell terra-luna @        40.2905  --  amount:    24,612.6431 , units:       610.8788  [portfolio:    24,612.6431]
 11. [2021-11-30] buy  terra-luna @        51.5801  --  amount:    24,612.6431 , units:       477.1735
 12. [2022-01-09] sell terra-luna @        66.9792  --  amount:    31,960.7004 , units:       477.1735  [portfolio:    31,960.7004]
 13. [2022-01-16] buy  terra-luna @        87.6574  --  amount:    31,960.7004 , units:       364.6091
 14. [2022-01-22] sell terra-luna @        64.3927  --  amount:    23,478.1657 , units:       364.6091  [portfolio:    23,478.1657]
 15. [2022-02-25] buy  terra-luna @        65.4442  --  amount:    23,478.1657 , units:       358.7510


- Number of txns     : 15
- First buy          : 2021-04-27
- Last sell          : 2022-01-22  (270 days after first buy)
- Initial investment : 10,000.00
- Portfolio value    : 26,999.00
- P/L                : 169.99 %

```

#### Yield Farming

```golang
// cmd/yield_farming_example/main.go
//
// Simulates yield farming LUNA/OSMO LP on Osmosis DEX.
// - Average APR: 75%.
// - Starting from Jul 1, 2021.
// - Daily harvesting and compounding for max 365 days.
//
import (
	"log"
	"time"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/anrid/traderbot/pkg/trade"
)

func main() {
	cg := coingecko.New(coingecko.USD)

	var periodInDays uint = 365

	a, err := cg.MarketChartWithCache("terra-luna", periodInDays, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	b, err := cg.MarketChartWithCache("osmosis", periodInDays, jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	startDate := "2021-07-01"
	harvestDays := 365
	initialInvestment := 10_000.0

	farm, err := trade.NewLPFarm(a, b, coingecko.USD, initialInvestment, startDate, 75.0 /* APR */)
	if err != nil {
		log.Fatal(err)
	}

	from := timeseries.ToTime(startDate)
	for i := 1; i <= harvestDays; /* number of days to harvest and compound yields */ i++ {
		cur := from.Add(time.Duration(i) * 24 * time.Hour)
		if cur.After(time.Now()) {
			break
		}
		farm.Harvest(timeseries.ToDate(cur))
	}
}
```

#### Output:

```bash
$ go run cmd/yield_farming_example/main.go


[2021-07-02] position  :   8,747.09  (IL:  -0.14 , hodl:   8,735.00 , a:       5.93 , b:       3.37)
[2021-07-03] position  :   8,658.17  (IL:  -0.40 , hodl:   8,623.95 , a:       5.74 , b:       3.40)
[2021-07-04] position  :   8,799.44  (IL:  -0.61 , hodl:   8,746.45 , a:       5.81 , b:       3.45)
[2021-07-05] position  :   9,033.70  (IL:  -0.81 , hodl:   8,961.16 , a:       5.97 , b:       3.53)
[2021-07-06] position  :   9,399.49  (IL:  -1.03 , hodl:   9,303.98 , a:       6.15 , b:       3.69)
[2021-07-07] position  :   8,944.96  (IL:  -0.58 , hodl:   8,892.94 , a:       6.48 , b:       3.16)
[2021-07-08] position  :   9,027.44  (IL:   0.58 , hodl:   9,080.18 , a:       7.12 , b:       2.92)
[2021-07-09] position  :   8,039.89  (IL:   1.69 , hodl:   8,178.40 , a:       6.71 , b:       2.44)
[2021-07-10] position  :   8,566.98  (IL:   5.60 , hodl:   9,074.94 , a:       8.17 , b:       2.27)

.. <many rows removed for brevity> ...

[2022-02-21] position  :  65,095.98  (IL: -34.08 , hodl:  48,551.16 , a:      49.61 , b:       8.54)
[2022-02-22] position  :  63,637.82  (IL: -31.48 , hodl:  48,400.05 , a:      50.25 , b:       8.02)
[2022-02-23] position  :  67,751.62  (IL: -29.91 , hodl:  52,151.53 , a:      54.67 , b:       8.33)
[2022-02-24] position  :  72,120.11  (IL: -27.88 , hodl:  56,395.63 , a:      59.79 , b:       8.59)
[2022-02-25] position  :  75,945.20  (IL: -24.89 , hodl:  60,809.89 , a:      65.44 , b:       8.67)
[2022-02-26] position  :  82,609.28  (IL: -22.74 , hodl:  67,302.13 , a:      73.18 , b:       9.13)
[2022-02-27] position  :  88,638.59  (IL: -23.70 , hodl:  71,658.79 , a:      77.69 , b:       9.86)
[2022-02-28] position  :  82,187.47  (IL: -23.10 , hodl:  66,765.55 , a:      72.64 , b:       9.03)

```