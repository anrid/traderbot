package coingecko

import (
	"fmt"
	"log"
)

type Indicator struct {
	Name string `json:"name"`
	V    map[int64]float64
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

	var prevDayEMA float64
	mul := 2.0 / (float64(days) + 1.0)
	// Dump(prices)

	for i := len(prices) - 1; i >= 0; i-- {
		priceStartIndex := i
		priceEndIndex := i - (days - 1)

		if priceEndIndex < 0 {
			// Not enough days of price data left.
			break
		}

		from := prices[priceStartIndex]
		to := prices[priceEndIndex]

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

		ema := (from.V * mul) + (prevDayEMA * (1 - mul))

		fmt.Printf("[%s - %s] ema = %.04f  closing = %.04f\n", from.TimeString(), to.TimeString(), ema, from.V)

		// Set EMA (using current timestamp as the key).
		in.V[from.TS] = ema

		prevDayEMA = ema
	}

	return in
}
