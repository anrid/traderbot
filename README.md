# Trader Bot!

### Example Usage

```golang
// cmd/cli/main.go
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

		strategy, err := trade.NewEMACrossOverStrategy(ema9d, ema21d, c)
		if err != nil {
			log.Fatal(err)
		}

		initialInvestment := 10_000.0 // USD.
		trade.ExecuteTradesAndPrint("EMS 9/21 CrossOver Strategy", initialInvestment, strategy.Trades)
	}
}
```

### Output:

```bash
$ go run cmd/cli/main.go 


Trading 'terra-luna' : EMS 9/21 CrossOver Strategy
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
- Portfolio value    : 26,405.18
- P/L                : 164.05 %




Trading 'solana' : EMS 9/21 CrossOver Strategy
=============================================================

  1. [2021-03-28] buy  solana     @        16.5664  --  amount:    10,000.0000 , units:       603.6301
  2. [2021-05-23] sell solana     @        31.6146  --  amount:    19,083.5250 , units:       603.6301  [portfolio:    19,083.5250]
  3. [2021-06-07] buy  solana     @        42.2599  --  amount:    19,083.5250 , units:       451.5754
  4. [2021-06-21] sell solana     @        35.2843  --  amount:    15,933.5102 , units:       451.5754  [portfolio:    15,933.5102]
  5. [2021-07-08] buy  solana     @        36.6573  --  amount:    15,933.5102 , units:       434.6616
  6. [2021-07-10] sell solana     @        33.5319  --  amount:    14,575.0093 , units:       434.6616  [portfolio:    14,575.0093]
  7. [2021-08-01] buy  solana     @        36.2488  --  amount:    14,575.0093 , units:       402.0830
  8. [2021-09-27] sell solana     @       135.7959  --  amount:    54,601.2322 , units:       402.0830  [portfolio:    54,601.2322]
  9. [2021-10-02] buy  solana     @       161.3740  --  amount:    54,601.2322 , units:       338.3520
 10. [2021-11-25] sell solana     @       205.5487  --  amount:    69,547.8061 , units:       338.3520  [portfolio:    69,547.8061]
 11. [2021-12-03] buy  solana     @       232.5722  --  amount:    69,547.8061 , units:       299.0375
 12. [2021-12-04] sell solana     @       212.1339  --  amount:    63,435.9763 , units:       299.0375  [portfolio:    63,435.9763]
 13. [2021-12-27] buy  solana     @       198.1974  --  amount:    63,435.9763 , units:       320.0646
 14. [2021-12-30] sell solana     @       171.2723  --  amount:    54,818.2086 , units:       320.0646  [portfolio:    54,818.2086]


- Number of txns     : 14
- First buy          : 2021-03-28
- Last sell          : 2021-12-30  (277 days after first buy)
- Initial investment : 10,000.00
- Portfolio value    : 54,818.21
- P/L                : 448.18 %




Trading 'bitcoin' : EMS 9/21 CrossOver Strategy
=============================================================

  1. [2021-05-07] buy  bitcoin    @    56,507.7594  --  amount:    10,000.0000 , units:         0.1770
  2. [2021-05-13] sell bitcoin    @    50,004.7622  --  amount:     8,849.1851 , units:         0.1770  [portfolio:     8,849.1851]
  3. [2021-07-26] buy  bitcoin    @    35,456.1247  --  amount:     8,849.1851 , units:         0.2496
  4. [2021-09-11] sell bitcoin    @    44,802.6064  --  amount:    11,181.8920 , units:         0.2496  [portfolio:    11,181.8920]
  5. [2021-09-19] buy  bitcoin    @    48,266.6271  --  amount:    11,181.8920 , units:         0.2317
  6. [2021-09-21] sell bitcoin    @    42,932.9466  --  amount:     9,946.2424 , units:         0.2317  [portfolio:     9,946.2424]
  7. [2021-10-04] buy  bitcoin    @    48,282.9711  --  amount:     9,946.2424 , units:         0.2060
  8. [2021-11-19] sell bitcoin    @    56,987.3223  --  amount:    11,739.3298 , units:         0.2060  [portfolio:    11,739.3298]
  9. [2022-02-07] buy  bitcoin    @    42,475.5432  --  amount:    11,739.3298 , units:         0.2764
 10. [2022-02-21] sell bitcoin    @    38,514.0085  --  amount:    10,644.4465 , units:         0.2764  [portfolio:    10,644.4465]


- Number of txns     : 10
- First buy          : 2021-05-07
- Last sell          : 2022-02-21  (290 days after first buy)
- Initial investment : 10,000.00
- Portfolio value    : 10,644.45
- P/L                : 6.44 %

```
