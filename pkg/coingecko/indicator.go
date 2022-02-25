package coingecko

import (
	"fmt"
	"log"
)

type Indicator struct {
	Name string `json:"name"`
	V    map[int64]float64
}

func (i *Indicator) ForTimestamp(milli int64) (ema float64) {
	ema = i.V[milli]
	return
}

func NewEMAIndicator(days int, prices TimeSeries) *Indicator {
	in := &Indicator{
		Name: fmt.Sprintf("%d-Day EMA", days),
		V:    make(map[int64]float64),
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
	// - Assume prices of day1 = 3.0, day2 = 4.0, day3 = 5.0 (day3 being today)
	//
	// multiplier = 2 / (3 + 1) = 0.5
	// SMA = (3.0 + 4.0 + 5.0) / 3 = 4.0
	//
	// Day 3 EMA = (5.0 * multiplier) + (SMA * (1 - multiplier)) = 2.5 + 2.0 = 4.5
	// Day 2 EMA = (4.0 * multiplier) + (4.5 * (1 - multiplier)) = 2.0 + 2.25 = 4.25
	// Day 1 EMA = (3.0 * multiplier) + (4.25 * (1 - multiplier)) = 1.5 + 2.125 = 3.625
	//
	var prevDayEMA float64
	multiplier := 2.0 / (float64(days) + 1.0)
	// Dump(prices)

	for i := len(prices) - 1; i >= 0; i-- {
		priceStartIndex := i
		priceEndIndex := i - (days - 1)

		if priceEndIndex < 0 {
			// Not enough days of price data left.
			break
		}

		start := prices[priceStartIndex]
		end := prices[priceEndIndex]

		if prevDayEMA == 0 {
			// Calculate Simple Moving Average and use it as the
			// first EMA for the previous day.
			var total float64
			for j := i; j >= priceEndIndex; j-- {
				total += prices[j].V
			}
			sma := total / float64(days)
			prevDayEMA = sma
		}

		ema := (start.V * multiplier) + (prevDayEMA * (1 - multiplier))

		fmt.Printf("[%s - %s] ema = %.04f  closing = %.04f\n", start.TimeString(), end.TimeString(), ema, start.V)

		// Set EMA in values map using current start timestamp as key.
		in.V[start.TS] = ema

		prevDayEMA = ema
	}

	return in
}
