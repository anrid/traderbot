package trade

import (
	"time"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/timeseries"
)

type Forecast struct {
	Farms             []*LPFarm
	Currency          coingecko.Fiat
	InitialInvestment float64
	Balance           float64
	StartDate         string
	Days              int
}

func NewForecast(currency coingecko.Fiat, initialInvestment float64, days int) *Forecast {
	return &Forecast{
		Currency:          currency,
		InitialInvestment: initialInvestment,
		Balance:           initialInvestment,
		StartDate:         timeseries.ToDate(time.Now()),
		Days:              days,
	}
}

type PriceChange struct {
	IncPct  float64 // Price increse percentage
	IncDays int     // Number of days that the price increases
	DecPct  float64 // Price decrese percentage
	DecDays int     // Number of days that the price decreases
}

func (fc *Forecast) CreateMarket(name, symbol string, startingPrice float64, changes []PriceChange) *coingecko.Market {
	m := &coingecko.Market{
		Currency: fc.Currency,
		ID:       symbol,
		Symbol:   symbol,
		Name:     name,
		Prices:   make(timeseries.Series, 0),
	}

	date := timeseries.ToTime(fc.StartDate)
	price := startingPrice

	changesIndex := -1
	priceChangeDelta := 0.0
	priceChangeDaysRemaining := 0

	for day := 0; day < fc.Days; day++ {
		date := date.Add(time.Duration(day) * 24 * time.Hour)

		if len(changes) > 0 {
			if priceChangeDaysRemaining == 0 && changesIndex < len(changes) {
				// Advance to next price change and calculate change delta.
				changesIndex++

				c := changes[changesIndex]

				if c.IncPct > 0 && c.IncDays > 0 {
					// Calculate price increment.
					targetPrice := price * (1 + (c.IncPct / 100))
					priceChangeDaysRemaining = c.IncDays
					priceChangeDelta = (targetPrice - price) / float64(c.IncDays)
				} else if c.DecPct > 0 && c.DecDays > 0 {
					// Calculate price decrement.
					targetPrice := price * (1 - (c.DecPct / 100))
					priceChangeDaysRemaining = c.DecDays
					priceChangeDelta = (targetPrice - price) / float64(c.DecDays)
				}
			}

			if priceChangeDaysRemaining > 0 {
				// Apply price change.
				priceChangeDaysRemaining--
				price += priceChangeDelta
			}
		}

		m.Prices = append(m.Prices, timeseries.ValueAt{
			TS: date.UnixMilli(),
			V:  price,
		})
	}

	return m
}

// func (fc *Forecast) FarmAssets(assetAName, assetBName, asset b *coingecko.Market) {

// 	NewLPFarm()

// 	fc.Farms = append(fc.Farms, f)
// }
