package trade

import (
	"fmt"
	"log"

	"github.com/anrid/traderbot/pkg/timeseries"
)

type Indicator struct {
	Name         string `json:"name"`
	ByTimestamp  map[int64]float64
	ByDateString map[string]float64
}

func (i *Indicator) ForTimestamp(milli int64) (ema float64) {
	ema = i.ByTimestamp[milli]
	return
}

func (i *Indicator) ForDate(date string) (ema float64) {
	ema = i.ByDateString[date]
	return
}

func NewEMAIndicator(days int, prices timeseries.Series) *Indicator {
	// Dump(prices)

	in := &Indicator{
		Name:         fmt.Sprintf("%d-Day EMA", days),
		ByTimestamp:  make(map[int64]float64),
		ByDateString: make(map[string]float64),
	}

	if days == 0 || len(prices) < days {
		// Not enough days of price data to calculate desired
		// observation period.
		log.Fatalf("not enough price data (%d entries) for observation period (%d days)", len(prices), days)
	}

	// EMA calculation:
	//
	// multipiler = [2 ÷ (number of observations + 1)]. E.g. For a 20-day moving average, the multiplier would be [2/(20+1)]= 0.0952.
	// EMA = (<closing price> x multiplier) + (<EMA previous day> x (1 - multiplier))
	//
	// Example calculation:
	// - Assume 3 days of observations
	// - Assume prices of [1.0, 2.0, 3.0, 4.0, 5.0] (daily closing prices)
	//
	// multiplier = 2 / (3 + 1) = 0.5
	// SMA = (1.0 + 2.0 + 3.0) / 3 = 2.0
	//
	// Day 4 EMA = (4.0 * multiplier) + (SMA * (1 - multiplier)) = 2.0 + 1.0 = 3.0
	// Day 5 EMA = (5.0 * multiplier) + (3.0 * (1 - multiplier)) = 2.5 + 1.5 = 4.0
	//
	multiplier := 2.0 / (float64(days) + 1.0)

	// Use prices for the first X days to calculate an SMA.
	var i int
	var total float64
	for ; i < days; i++ {
		total += prices[i].V
	}
	sma := total / float64(days)
	prevDayEMA := sma

	for ; i < len(prices); i++ {
		cur := prices[i]

		ema := (cur.V * multiplier) + (prevDayEMA * (1 - multiplier))

		// from := prices[i-days]
		// fmt.Printf("[%s - %s] ema = %.04f  closing = %.04f\n", from.TimeString(), cur.TimeString(), ema, cur.V)

		in.ByTimestamp[cur.TS] = ema
		in.ByDateString[cur.Date()] = ema

		prevDayEMA = ema
	}

	return in
}
