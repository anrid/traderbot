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

	farm, err := trade.NewLPFarm(a, b, coingecko.USD, initialInvestment, startDate, 99.0 /* APR */)
	if err != nil {
		log.Fatal(err)
	}

	farm.SetAPRDailyDecay(0.15) // Lower APR by 0.15 percentage points every day.

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

Using cached data : terra-luna-365-days-usd
Using cached data : osmosis-365-days-usd

[2021-07-02] position  :   8,752.79  (IL:  -0.20 , hodl:   8,735.00 , APR:  98.85 % , a:       5.93 , b:       3.37)
[2021-07-03] position  :   8,669.43  (IL:  -0.53 , hodl:   8,623.95 , APR:  98.70 % , a:       5.74 , b:       3.40)
[2021-07-04] position  :   8,816.55  (IL:  -0.80 , hodl:   8,746.45 , APR:  98.55 % , a:       5.81 , b:       3.45)
[2021-07-05] position  :   9,057.06  (IL:  -1.07 , hodl:   8,961.16 , APR:  98.40 % , a:       5.97 , b:       3.53)
[2021-07-06] position  :   9,429.79  (IL:  -1.35 , hodl:   9,303.98 , APR:  98.25 % , a:       6.15 , b:       3.69)
[2021-07-07] position  :   8,979.46  (IL:  -0.97 , hodl:   8,892.94 , APR:  98.10 % , a:       6.48 , b:       3.16)
[2021-07-08] position  :   9,067.94  (IL:   0.13 , hodl:   9,080.18 , APR:  97.95 % , a:       7.12 , b:       2.92)
[2021-07-09] position  :   8,081.00  (IL:   1.19 , hodl:   8,178.40 , APR:  97.80 % , a:       6.71 , b:       2.44)
[2021-07-10] position  :   8,616.12  (IL:   5.06 , hodl:   9,074.94 , APR:  97.65 % , a:       8.17 , b:       2.27)
[2021-07-11] position  :   8,720.93  (IL:   3.50 , hodl:   9,037.32 , APR:  97.50 % , a:       7.94 , b:       2.38)
[2021-07-12] position  :   8,864.39  (IL:   4.83 , hodl:   9,314.60 , APR:  97.35 % , a:       8.43 , b:       2.31)
[2021-07-13] position  :   7,934.92  (IL:   6.66 , hodl:   8,501.37 , APR:  97.20 % , a:       7.94 , b:       1.95)
[2021-07-14] position  :   7,047.61  (IL:   6.75 , hodl:   7,558.03 , APR:  97.05 % , a:       7.10 , b:       1.71)
[2021-07-15] position  :   6,502.32  (IL:   9.06 , hodl:   7,150.44 , APR:  96.90 % , a:       6.94 , b:       1.48)
[2021-07-16] position  :   6,147.09  (IL:   9.26 , hodl:   6,774.32 , APR:  96.75 % , a:       6.61 , b:       1.38)
[2021-07-17] position  :   6,118.43  (IL:   8.95 , hodl:   6,719.97 , APR:  96.60 % , a:       6.55 , b:       1.38)
[2021-07-18] position  :   6,041.54  (IL:   9.02 , hodl:   6,640.47 , APR:  96.45 % , a:       6.49 , b:       1.35)
[2021-07-19] position  :   6,444.11  (IL:   6.85 , hodl:   6,918.14 , APR:  96.30 % , a:       6.61 , b:       1.50)
[2021-07-20] position  :   5,904.55  (IL:   6.62 , hodl:   6,322.82 , APR:  96.15 % , a:       6.04 , b:       1.37)
[2021-07-21] position  :   5,662.75  (IL:   7.29 , hodl:   6,108.00 , APR:  96.00 % , a:       5.90 , b:       1.28)
[2021-07-22] position  :   6,549.15  (IL:   6.11 , hodl:   6,975.26 , APR:  95.85 % , a:       6.67 , b:       1.51)
[2021-07-23] position  :   7,038.19  (IL:   7.23 , hodl:   7,586.51 , APR:  95.70 % , a:       7.37 , b:       1.57)
[2021-07-24] position  :   7,293.90  (IL:   7.82 , hodl:   7,912.94 , APR:  95.55 % , a:       7.76 , b:       1.59)
[2021-07-25] position  :   7,635.74  (IL:   9.50 , hodl:   8,436.91 , APR:  95.40 % , a:       8.44 , b:       1.60)
[2021-07-26] position  :   7,627.04  (IL:   8.47 , hodl:   8,332.46 , APR:  95.25 % , a:       8.27 , b:       1.62)
[2021-07-27] position  :   7,878.49  (IL:   8.09 , hodl:   8,572.16 , APR:  95.10 % , a:       8.49 , b:       1.67)
[2021-07-28] position  :   8,328.71  (IL:  10.48 , hodl:   9,303.51 , APR:  94.95 % , a:       9.45 , b:       1.67)
[2021-07-29] position  :   9,038.38  (IL:  14.13 , hodl:  10,526.18 , APR:  94.80 % , a:      11.05 , b:       1.67)
[2021-07-30] position  :   9,207.72  (IL:  13.15 , hodl:  10,601.48 , APR:  94.65 % , a:      11.06 , b:       1.72)
[2021-07-31] position  :   9,235.99  (IL:  11.94 , hodl:  10,488.69 , APR:  94.50 % , a:      10.86 , b:       1.76)
[2021-08-01] position  :   9,342.09  (IL:  11.26 , hodl:  10,527.24 , APR:  94.35 % , a:      10.86 , b:       1.79)
[2021-08-02] position  :   9,549.41  (IL:  14.45 , hodl:  11,162.06 , APR:  94.20 % , a:      11.82 , b:       1.71)
[2021-08-03] position  :   9,880.23  (IL:  15.19 , hodl:  11,649.72 , APR:  94.05 % , a:      12.42 , b:       1.73)
[2021-08-04] position  :  10,660.43  (IL:  19.67 , hodl:  13,271.56 , APR:  93.90 % , a:      14.58 , b:       1.71)
[2021-08-05] position  :  10,818.10  (IL:  18.69 , hodl:  13,305.12 , APR:  93.75 % , a:      14.55 , b:       1.75)
[2021-08-06] position  :  10,695.76  (IL:  19.28 , hodl:  13,251.12 , APR:  93.60 % , a:      14.56 , b:       1.71)
[2021-08-07] position  :  10,931.07  (IL:  17.90 , hodl:  13,314.32 , APR:  93.45 % , a:      14.53 , b:       1.78)
[2021-08-08] position  :  11,144.43  (IL:  17.09 , hodl:  13,441.59 , APR:  93.30 % , a:      14.61 , b:       1.83)
[2021-08-09] position  :  10,363.35  (IL:  16.52 , hodl:  12,413.85 , APR:  93.15 % , a:      13.47 , b:       1.70)
[2021-08-10] position  :  10,607.96  (IL:  17.14 , hodl:  12,802.04 , APR:  93.00 % , a:      13.96 , b:       1.71)
[2021-08-11] position  :  11,871.58  (IL:  20.72 , hodl:  14,974.34 , APR:  92.85 % , a:      16.68 , b:       1.79)
[2021-08-12] position  :  12,488.23  (IL:  19.27 , hodl:  15,468.90 , APR:  92.70 % , a:      17.12 , b:       1.92)
[2021-08-13] position  :  12,258.23  (IL:  18.53 , hodl:  15,045.85 , APR:  92.55 % , a:      16.60 , b:       1.90)
[2021-08-14] position  :  13,447.07  (IL:  17.01 , hodl:  16,204.09 , APR:  92.40 % , a:      17.74 , b:       2.12)
[2021-08-15] position  :  13,761.00  (IL:  14.41 , hodl:  16,077.00 , APR:  92.25 % , a:      17.35 , b:       2.26)
[2021-08-16] position  :  14,687.44  (IL:  16.53 , hodl:  17,595.38 , APR:  92.10 % , a:      19.26 , b:       2.31)
[2021-08-17] position  :  15,646.44  (IL:  20.65 , hodl:  19,718.26 , APR:  91.95 % , a:      22.10 , b:       2.27)
[2021-08-18] position  :  16,919.72  (IL:  23.52 , hodl:  22,123.24 , APR:  91.80 % , a:      25.16 , b:       2.32)
[2021-08-19] position  :  20,012.77  (IL:  25.96 , hodl:  27,029.58 , APR:  91.65 % , a:      31.11 , b:       2.62)
[2021-08-20] position  :  20,256.76  (IL:  22.42 , hodl:  26,112.16 , APR:  91.50 % , a:      29.60 , b:       2.80)
[2021-08-21] position  :  20,862.43  (IL:  22.30 , hodl:  26,850.09 , APR:  91.35 % , a:      30.45 , b:       2.88)
[2021-08-22] position  :  21,670.96  (IL:  14.45 , hodl:  25,330.06 , APR:  91.20 % , a:      27.59 , b:       3.41)
[2021-08-23] position  :  21,446.34  (IL:  15.11 , hodl:  25,264.23 , APR:  91.05 % , a:      27.66 , b:       3.31)
[2021-08-24] position  :  22,803.48  (IL:  15.60 , hodl:  27,016.84 , APR:  90.90 % , a:      29.70 , b:       3.47)
[2021-08-25] position  :  21,681.71  (IL:  18.78 , hodl:  26,696.03 , APR:  90.75 % , a:      29.88 , b:       3.10)
[2021-08-26] position  :  21,974.10  (IL:  18.53 , hodl:  26,972.36 , APR:  90.60 % , a:      30.18 , b:       3.14)
[2021-08-27] position  :  19,997.42  (IL:  18.25 , hodl:  24,461.31 , APR:  90.45 % , a:      27.36 , b:       2.85)
[2021-08-28] position  :  23,387.75  (IL:  19.34 , hodl:  28,995.53 , APR:  90.30 % , a:      32.64 , b:       3.26)
[2021-08-29] position  :  24,611.71  (IL:  20.26 , hodl:  30,864.06 , APR:  90.15 % , a:      34.92 , b:       3.35)
[2021-08-30] position  :  24,511.95  (IL:  17.81 , hodl:  29,823.95 , APR:  90.00 % , a:      33.38 , b:       3.46)
[2021-08-31] position  :  25,134.07  (IL:  16.25 , hodl:  30,010.98 , APR:  89.85 % , a:      33.36 , b:       3.63)
[2021-09-01] position  :  24,970.54  (IL:  13.63 , hodl:  28,909.88 , APR:  89.70 % , a:      31.73 , b:       3.74)
[2021-09-02] position  :  25,905.72  (IL:  11.13 , hodl:  29,148.50 , APR:  89.55 % , a:      31.57 , b:       4.03)

... lots of rows removed ...

[2022-02-12] position  :  66,506.85  (IL: -35.91 , hodl:  48,935.68 , APR:  65.10 % , a:      50.42 , b:       8.35)
[2022-02-13] position  :  67,679.34  (IL: -34.93 , hodl:  50,157.46 , APR:  64.95 % , a:      52.02 , b:       8.35)
[2022-02-14] position  :  66,875.96  (IL: -34.35 , hodl:  49,776.20 , APR:  64.80 % , a:      51.84 , b:       8.15)
[2022-02-15] position  :  68,765.89  (IL: -33.59 , hodl:  51,476.07 , APR:  64.65 % , a:      53.89 , b:       8.26)
[2022-02-16] position  :  73,935.13  (IL: -35.24 , hodl:  54,669.67 , APR:  64.50 % , a:      56.82 , b:       9.03)
[2022-02-17] position  :  73,034.06  (IL: -35.32 , hodl:  53,972.72 , APR:  64.35 % , a:      56.14 , b:       8.88)
[2022-02-18] position  :  68,792.84  (IL: -39.30 , hodl:  49,383.69 , APR:  64.20 % , a:      50.34 , b:       8.76)
[2022-02-19] position  :  69,414.27  (IL: -39.36 , hodl:  49,807.65 , APR:  64.05 % , a:      50.82 , b:       8.80)
[2022-02-20] position  :  69,114.27  (IL: -39.63 , hodl:  49,499.36 , APR:  63.90 % , a:      50.50 , b:       8.75)
[2022-02-21] position  :  67,784.17  (IL: -39.61 , hodl:  48,551.16 , APR:  63.75 % , a:      49.61 , b:       8.54)
[2022-02-22] position  :  66,245.14  (IL: -36.87 , hodl:  48,400.05 , APR:  63.60 % , a:      50.25 , b:       8.02)
[2022-02-23] position  :  70,505.21  (IL: -35.19 , hodl:  52,151.53 , APR:  63.45 % , a:      54.67 , b:       8.33)
[2022-02-24] position  :  75,027.24  (IL: -33.04 , hodl:  56,395.63 , APR:  63.30 % , a:      59.79 , b:       8.59)
[2022-02-25] position  :  78,980.92  (IL: -29.88 , hodl:  60,809.89 , APR:  63.15 % , a:      65.44 , b:       8.67)
[2022-02-26] position  :  85,883.20  (IL: -27.61 , hodl:  67,302.13 , APR:  63.00 % , a:      73.18 , b:       9.13)
[2022-02-27] position  :  92,120.85  (IL: -28.55 , hodl:  71,658.79 , APR:  62.85 % , a:      77.69 , b:       9.86)
[2022-02-28] position  :  85,387.56  (IL: -27.89 , hodl:  66,765.55 , APR:  62.70 % , a:      72.64 , b:       9.03)
[2022-03-01] position  : 100,863.67  (IL: -22.72 , hodl:  82,193.25 , APR:  62.55 % , a:      91.26 , b:      10.00)
[2022-03-02] position  : 100,412.99  (IL: -23.59 , hodl:  81,247.15 , APR:  62.40 % , a:      90.00 , b:      10.01)
[2022-03-03] position  : 105,066.12  (IL: -25.68 , hodl:  83,600.55 , APR:  62.25 % , a:      91.97 , b:      10.69)

```