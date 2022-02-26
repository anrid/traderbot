package trade

import (
	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/pkg/errors"
)

type EMACrossOverStrategy struct {
	ShortEMA *Indicator
	LongEMA  *Indicator
	C        *coingecko.Coin
	Trades   []*Trade
}

func NewEMACrossOverStrategy(shortEMA, longEMA *Indicator, c *coingecko.Coin) (*EMACrossOverStrategy, error) {
	strat := &EMACrossOverStrategy{
		ShortEMA: shortEMA,
		LongEMA:  longEMA,
		C:        c,
	}

	// pr := message.NewPrinter(language.English)
	var lastShort float64
	var lastLong float64

	for _, p := range c.Prices {
		short := shortEMA.ForTimestamp(p.TS)
		long := longEMA.ForTimestamp(p.TS)
		if short == 0.0 || long == 0.0 {
			// Could not find EMA values for both short and long observation periods.
			continue
		}

		// fmt.Printf("[%s] short %s = %.04f  --  long %s = %.04f\n", p.Date(), shortEMA.Name, short, longEMA.Name, long)

		if lastShort > 0.0 && lastLong > 0.0 {
			shortCrossedOverLongFromAbove := lastShort > lastLong && short < long // Sell signal.
			shortCrossedOverLongFromBelow := lastShort < lastLong && short > long // Buy signal.

			if shortCrossedOverLongFromAbove {
				// Sell signal.
				// pr.Printf("[%s] sell @ %.04f\n", p.Date(), p.V)

				t, err := NewSellAtDate(p.Date(), 100.0, c)
				if err != nil {
					return nil, errors.Wrapf(err, "could not create sell trade for %s", c.ID)
				}
				strat.Trades = append(strat.Trades, t)
			} else if shortCrossedOverLongFromBelow {
				// Buy signal.
				// pr.Printf("[%s] buy @ %.04f\n", p.Date(), p.V)

				t, err := NewBuyAtDate(p.Date(), 100.0, c)
				if err != nil {
					return nil, errors.Wrapf(err, "could not create buy trade for %s", c.ID)
				}
				strat.Trades = append(strat.Trades, t)
			}
		}

		lastShort = short
		lastLong = long
	}

	// Feature flag: Add one last sell trade if the last trade is a buy.
	addForcedSell := false
	if addForcedSell {
		if len(strat.Trades) > 0 {
			last := strat.Trades[len(strat.Trades)-1]
			if last.Side == Buy {
				latestPrice := c.Prices[len(c.Prices)-1]

				forcedSell, err := NewSellAtDate(latestPrice.Date(), 100.0, c)
				if err != nil {
					return nil, errors.Wrapf(err, "could not create one last forced sale for %s", c.ID)
				}
				strat.Trades = append(strat.Trades, forcedSell)
			}
		}
	}

	return strat, nil
}
