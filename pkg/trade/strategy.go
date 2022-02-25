package trade

import (
	"fmt"

	"github.com/anrid/traderbot/pkg/coingecko"
)

type EMACrossOverStrategy struct {
	ShortEMA *Indicator
	LongEMA  *Indicator
	Prices   coingecko.TimeSeries
}

func NewEMACrossOverStrategy(shortEMA, longEMA *Indicator, prices coingecko.TimeSeries) *EMACrossOverStrategy {
	strat := &EMACrossOverStrategy{
		ShortEMA: shortEMA,
		LongEMA:  longEMA,
		Prices:   prices,
	}

	var lastShort float64
	var lastLong float64

	for _, p := range prices {
		short := shortEMA.ForTimestamp(p.TS)
		long := longEMA.ForTimestamp(p.TS)
		if short == 0.0 || long == 0.0 {
			// Could not find EMA values for both short and long observation periods.
			continue
		}

		fmt.Printf("[%s] short %s = %.04f  --  long %s = %.04f\n", p.TimeString(), shortEMA.Name, short, longEMA.Name, long)

		if lastShort > 0.0 && lastLong > 0.0 {
			shortCrossedOverLongFromBelow := lastShort < lastLong && short > long
			shortCrossedOverLongFromAbove := lastShort > lastLong && short < long

			if shortCrossedOverLongFromAbove {
				// Sell signal.
				fmt.Printf("[%s] sell @ %.04f\n", p.TimeString(), p.V)
			} else if shortCrossedOverLongFromBelow {
				// Buy signal.
				fmt.Printf("[%s] buy @ %.04f\n", p.TimeString(), p.V)
			}
		}

		lastShort = short
		lastLong = long
	}

	return strat
}
