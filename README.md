# Create Yield Farming Charts

### Farming LUNA/OSMO LP

```bash
# Simulates yield farming LUNA/OSMO LP on Osmosis DEX.
# - Starting APR  : 100%
# - Final APR     : 60%
# - Starting from : Jul 1, 2021
# - Duration      : 365 (number of days we want to harvest and compound yields)
#
$ go run examples/yield_farming/main.go --path /tmp --asset-a terra-luna --asset-b osmosis --apr 100.0 --final-apr 60 --start-date 2021-07-01

[2021-07-01] position  :  10,000.00  (IL:   0.00 , hodl:  10,000.00 , APR: 100.00 % , a:       6.54 , b:       4.01 , units: 764.00 / 1,247.38)
[2021-07-02] position  :   8,753.03  (IL:  -0.21 , hodl:   8,735.00 , APR:  99.84 % , a:       5.93 , b:       3.37 , units: 738.56 / 1,297.42)

... lots of rows removed ...

[2022-03-04] position  : 106,430.21  (IL: -28.04 , hodl:  83,124.51 , APR:  60.16 % , a:      90.56 , b:      11.18 , units: 587.65 / 4,761.86)
[2022-03-05] position  :  97,558.70  (IL: -27.28 , hodl:  76,651.23 , APR:  60.00 % , a:      83.82 , b:      10.11 , units: 581.95 / 4,824.31)

writing chart /tmp/yield-farming-luna-osmo-2021-07-01-2022-03-05.html
```

#### Renders HTML chart:

![screenshot of chart](examples/yield_farming/screens/yield-farming-luna-osmo-2021-07-01-2022-03-05-v2.jpg)

### Farming AVAX/USDC LP

```bash
# Simulates yield farming AVAX/USDC.e LP, but with 0% APR to see the full impact of impermanent loss.
# - Starting APR  : 0%
# - Final APR     : 0%
# - Starting from : Jul 1, 2021
# - Duration      : 365 (number of days we want to harvest and compound yields)
#
$ go run examples/yield_farming/main.go --path /tmp --asset-a avalanche-2 --asset-b usd-coin --apr 0.0 --final-apr 0.0 --start-date 2021-07-01

[2021-07-01] position  :  10,000.00  (IL:   0.00 , hodl:  10,000.00 , APR:   0.00 % , a:      11.98 , b:       1.00 , units: 417.38 / 4,983.67)
[2021-07-02] position  :   9,713.69  (IL:   0.05 , hodl:   9,718.26 , APR:   0.00 % , a:      11.28 , b:       1.00 , units: 430.39 / 4,833.02)

... lots of rows removed ...

[2022-03-04] position  :  25,599.92  (IL:  32.42 , hodl:  37,878.16 , APR:   0.00 % , a:      78.82 , b:       1.00 , units: 162.39 / 12,808.83)
[2022-03-05] position  :  25,726.79  (IL:  32.59 , hodl:  38,165.82 , APR:   0.00 % , a:      79.49 , b:       1.00 , units: 161.82 / 12,854.43)

writing chart /tmp/yield-farming-avax-usdc-2021-07-01-2022-03-05.html
```

#### Renders HTML chart:

![screenshot of chart](examples/yield_farming/screens/yield-farming-avax-usdc-2021-07-01-2022-03-05-v2.jpg)

# EMA 9/21-Day Trading Simulation

```golang
// examples/ema_9_21_trading_strategy/main.go
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

# Yield Farming Simulations

```golang
// examples/yield_farming_simple/main.go
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
	startDate := "2021-07-01"     // Date of initial investment. We start farming from this date.
	harvestDays := 365            // Number of days to harvest and compound yields.
	initialInvestment := 10_000.0 // Initial investment in USD.
	apr := 99.0                   // Farm APR.

	cg := coingecko.New(coingecko.USD)

	a, err := cg.MarketChartWithCache("terra-luna", uint(harvestDays), jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	b, err := cg.MarketChartWithCache("osmosis", uint(harvestDays), jsoncache.InvalidateDaily)
	if err != nil {
		log.Fatal(err)
	}

	farm, err := trade.NewLPFarm(a, b, coingecko.USD, initialInvestment, startDate, apr)
	if err != nil {
		log.Fatal(err)
	}

	farm.SetAPRChangeRateAtHarvest(0.15) // Lower APR by 0.15 percentage points every day.

	from := timeseries.ToTime(startDate)
	for i := 1; i <= harvestDays; i++ {
		current := from.Add(time.Duration(i) * 24 * time.Hour)
		if current.After(time.Now()) {
			break
		}

		date := timeseries.ToDate(current)

		yield, err := farm.Harvest(date)
		if err != nil {
			log.Fatal(err)
		}

		farm.AddLP(date, yield) // Compound yield!
	}
}
```

#### Output:

```bash
$ go run cmd/yield_farming_example/main.go

CoinGecko says: (V3) To the Moon!
Downloaded data   : terra-luna-365-days-usd
Downloaded data   : osmosis-365-days-usd
[2021-07-01] position  :  10,000.00  (IL:   0.00 , hodl:  10,000.00 , APR:  99.00 % , a:       6.54 , b:       4.01 , units: 764.00 / 1,247.38)
[2021-07-02] position  :   8,729.15  (IL:   0.07 , hodl:   8,735.00 , APR:  98.85 % , a:       5.93 , b:       3.37 , units: 736.54 / 1,293.88)
[2021-07-03] position  :   8,646.05  (IL:  -0.26 , hodl:   8,623.95 , APR:  98.70 % , a:       5.74 , b:       3.40 , units: 753.11 / 1,272.28)
[2021-07-04] position  :   8,792.81  (IL:  -0.53 , hodl:   8,746.45 , APR:  98.55 % , a:       5.81 , b:       3.45 , units: 756.42 / 1,273.57)
[2021-07-05] position  :   9,032.71  (IL:  -0.80 , hodl:   8,961.16 , APR:  98.40 % , a:       5.97 , b:       3.53 , units: 757.03 / 1,279.43)
[2021-07-06] position  :   9,404.48  (IL:  -1.08 , hodl:   9,303.98 , APR:  98.25 % , a:       6.15 , b:       3.69 , units: 764.61 / 1,273.58)
[2021-07-07] position  :   8,955.39  (IL:  -0.70 , hodl:   8,892.94 , APR:  98.10 % , a:       6.48 , b:       3.16 , units: 690.94 / 1,416.98)
[2021-07-08] position  :   9,043.67  (IL:   0.40 , hodl:   9,080.18 , APR:  97.95 % , a:       7.12 , b:       2.92 , units: 634.66 / 1,550.94)
[2021-07-09] position  :   8,059.41  (IL:   1.46 , hodl:   8,178.40 , APR:  97.80 % , a:       6.71 , b:       2.44 , units: 600.11 / 1,649.05)
[2021-07-10] position  :   8,593.13  (IL:   5.31 , hodl:   9,074.94 , APR:  97.65 % , a:       8.17 , b:       2.27 , units: 525.87 / 1,891.94)
[2021-07-11] position  :   8,697.70  (IL:   3.76 , hodl:   9,037.32 , APR:  97.50 % , a:       7.94 , b:       2.38 , units: 547.58 / 1,826.65)
[2021-07-12] position  :   8,840.81  (IL:   5.09 , hodl:   9,314.60 , APR:  97.35 % , a:       8.43 , b:       2.31 , units: 524.53 / 1,917.13)
[2021-07-13] position  :   7,913.85  (IL:   6.91 , hodl:   8,501.37 , APR:  97.20 % , a:       7.94 , b:       1.95 , units: 498.06 / 2,029.81)
[2021-07-14] position  :   7,028.92  (IL:   7.00 , hodl:   7,558.03 , APR:  97.05 % , a:       7.10 , b:       1.71 , units: 495.21 / 2,052.38)

... lots of rows removed ...

[2022-02-21] position  :  67,665.98  (IL: -39.37 , hodl:  48,551.16 , APR:  63.75 % , a:      49.61 , b:       8.54 , units: 682.03 / 3,962.00)
[2022-02-22] position  :  66,129.91  (IL: -36.63 , hodl:  48,400.05 , APR:  63.60 % , a:      50.25 , b:       8.02 , units: 657.99 / 4,121.07)
[2022-02-23] position  :  70,382.86  (IL: -34.96 , hodl:  52,151.53 , APR:  63.45 % , a:      54.67 , b:       8.33 , units: 643.72 / 4,227.15)
[2022-02-24] position  :  74,897.35  (IL: -32.81 , hodl:  56,395.63 , APR:  63.30 % , a:      59.79 , b:       8.59 , units: 626.31 / 4,359.73)
[2022-02-25] position  :  78,844.51  (IL: -29.66 , hodl:  60,809.89 , APR:  63.15 % , a:      65.44 , b:       8.67 , units: 602.38 / 4,548.70)
[2022-02-26] position  :  85,735.22  (IL: -27.39 , hodl:  67,302.13 , APR:  63.00 % , a:      73.18 , b:       9.13 , units: 585.77 / 4,693.84)
[2022-02-27] position  :  91,962.50  (IL: -28.33 , hodl:  71,658.79 , APR:  62.85 % , a:      77.69 , b:       9.86 , units: 591.86 / 4,661.67)
[2022-02-28] position  :  85,241.13  (IL: -27.67 , hodl:  66,765.55 , APR:  62.70 % , a:      72.64 , b:       9.03 , units: 586.71 / 4,718.72)
[2022-03-01] position  : 100,691.11  (IL: -22.51 , hodl:  82,193.25 , APR:  62.55 % , a:      91.26 , b:      10.00 , units: 551.67 / 5,035.71)
[2022-03-02] position  : 100,241.62  (IL: -23.38 , hodl:  81,247.15 , APR:  62.40 % , a:      90.00 , b:      10.01 , units: 556.92 / 5,005.33)
[2022-03-03] position  : 104,887.23  (IL: -25.46 , hodl:  83,600.55 , APR:  62.25 % , a:      91.97 , b:      10.69 , units: 570.23 / 4,905.28)
[2022-03-04] position  : 106,589.57  (IL: -28.23 , hodl:  83,124.51 , APR:  62.10 % , a:      90.56 , b:      11.18 , units: 588.53 / 4,768.99)
[2022-03-05] position  :  97,710.38  (IL: -27.47 , hodl:  76,651.23 , APR:  61.95 % , a:      83.82 , b:      10.11 , units: 582.85 / 4,831.81)
[2022-03-06] position  : 101,427.19  (IL: -27.26 , hodl:  79,699.80 , APR:  61.80 % , a:      87.30 , b:      10.43 , units: 580.93 / 4,864.29)
[2022-03-07] position  :  93,011.33  (IL: -28.26 , hodl:  72,518.26 , APR:  61.65 % , a:      79.19 , b:       9.63 , units: 587.24 / 4,828.31)
[2022-03-08] position  :  92,935.68  (IL: -29.87 , hodl:  71,561.87 , APR:  61.50 % , a:      77.72 , b:       9.76 , units: 597.85 / 4,758.64)
[2022-03-09] position  : 106,407.93  (IL: -22.29 , hodl:  87,014.17 , APR:  61.35 % , a:      97.24 , b:      10.20 , units: 547.12 / 5,217.46)

```
