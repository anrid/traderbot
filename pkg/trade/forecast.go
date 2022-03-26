package trade

import (
	"encoding/json"
	"log"
	"math"
	"sort"
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

func (fc *Forecast) ToJSON() (string, error) {
	out := [][]interface{}{
		{"Name", "Date", "Price A", "Price B", "Units A", "Units B", "Value", "HODL", "Only A", "Only B", "APR"},
	}

	dates := make(map[string]bool)

	for _, f := range fc.Farms {
		for d := range f.ChangeHistory {
			dates[d] = true
		}
	}

	var uniqueDates []string
	for d := range dates {
		uniqueDates = append(uniqueDates, d)
	}

	sort.SliceStable(uniqueDates, func(i, j int) bool {
		return uniqueDates[i] < uniqueDates[j]
	})

	for _, d := range uniqueDates {
		for j := 0; j < len(fc.Farms); j++ {
			farm := fc.Farms[j]
			if h, found := farm.ChangeHistory[d]; found {
				out = append(out, []interface{}{
					farm.Name, h.Date, h.PriceA, h.PriceB, h.UnitsA, h.UnitsB,
					h.TotalValue, h.TotalValueHODL, h.TotalValueHODLOnlyA, h.TotalValueHODLOnlyB,
					h.APR,
				})
			}
		}
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (fc *Forecast) AddLPFarm(a, b *coingecko.Market, apr, finalAPR, additionalInvestmentMonthly float64) error {
	farm, err := NewLPFarm(a, b, fc.Currency, fc.InitialInvestment, fc.StartDate, apr)
	if err != nil {
		return err
	}

	var harvestDates []string

	from := timeseries.ToTime(fc.StartDate)
	to := from.Add(time.Duration(fc.Days) * 24 * time.Hour)

	for i := 1; i <= fc.Days; /* number of days to harvest and compound yields */ i++ {
		cur := from.Add(time.Duration(i) * 24 * time.Hour)
		if cur.After(to) {
			break
		}
		harvestDates = append(harvestDates, timeseries.ToDate(cur))
	}

	// Calculate APR change based on the number of times we harvest.
	// The goal is to have the APR reach a certain target value by the final
	// harvest date.
	{
		if finalAPR > 0.0 {
			change := (apr - finalAPR) / float64(len(harvestDates))
			farm.SetAPRChangeRateAtHarvest(change)
		}
	}

	currentMonth := harvestDates[0][5:7]
	for _, d := range harvestDates {
		yield, err := farm.Harvest(d)
		if err != nil {
			log.Fatal(err)
		}

		farm.AddLP(d, yield, false) // Compound yield!

		month := d[5:7]
		if currentMonth != month {
			currentMonth = month
			if additionalInvestmentMonthly > 0 {
				// Make an additional investment into the farm with
				// outside funds, i.e. dollar-cost-average into more
				// LP.
				farm.AddLP(d, additionalInvestmentMonthly, true)
			}
		}
	}

	fc.Farms = append(fc.Farms, farm)

	return nil
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

	start := timeseries.ToTime(fc.StartDate)
	price := startingPrice

	priceChangeIndex := 0
	priceChangeDelta := 0.0
	priceChangeDaysRemaining := 0

	for day := 0; day <= fc.Days; day++ {
		date := start.Add(time.Duration(day) * 24 * time.Hour)

		if len(changes) > 0 && day > 0 {
			// Start applying price changes on the following day.
			if priceChangeDaysRemaining == 0 {
				c := changes[priceChangeIndex]

				if c.IncPct > 0 && c.IncDays > 0 {
					// Calculate price increment.
					targetPrice := price * (1 + (c.IncPct / 100))
					priceChangeDaysRemaining = c.IncDays
					priceChangeDelta = (targetPrice - price) / float64(c.IncDays)

					// Round to 2 decimals
					priceChangeDelta = math.Floor(priceChangeDelta*100) / 100
				} else if c.DecPct > 0 && c.DecDays > 0 {
					// Calculate price decrement.
					targetPrice := price * (1 - (c.DecPct / 100))
					priceChangeDaysRemaining = c.DecDays
					priceChangeDelta = (targetPrice - price) / float64(c.DecDays)

					// Round to 2 decimals
					priceChangeDelta = math.Floor(priceChangeDelta*100) / 100
				}

				priceChangeIndex++
				if priceChangeIndex >= len(changes) {
					// Rotate back to the first change.
					priceChangeIndex = 0
				}
			}

			if priceChangeDaysRemaining > 0 {
				// Apply price change.
				priceChangeDaysRemaining--
				price += priceChangeDelta

				// Round to 2 decimals
				price = math.Floor(price*100) / 100
			}
		}

		m.Prices = append(m.Prices, timeseries.ValueAt{
			TS: date.UnixMilli(),
			V:  price,
		})
	}

	return m
}
