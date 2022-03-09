package trade

import (
	"math"
	"sort"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type LPFarm struct {
	A                      *coingecko.Market
	B                      *coingecko.Market
	Currency               coingecko.Fiat
	InitialInvestment      float64
	StartDate              string
	APR                    float64
	APRChangeRateAtHarvest float64
	InitialAPR             float64
	LastHarvestDate        string
	UnitsA                 float64
	UnitsB                 float64
	InitialUnitsA          float64
	InitialUnitsB          float64
	TotalValue             float64                       // Total value of farm in given fiat currency.
	TotalValueHODL         float64                       // Total value if we simply HODL:d both assets instead, in given fiat currency.
	ChangeHistory          map[string]*LPFarmHistoryItem // All changes, e.g. when a harvest was performed or more LP was added.
}

type LPFarmHistoryItem struct {
	Date                string
	PriceA              float64
	PriceB              float64
	UnitsA              float64
	UnitsB              float64
	TotalValue          float64
	TotalValueHODL      float64
	TotalValueHODLOnlyA float64
	TotalValueHODLOnlyB float64
	APR                 float64
}

func NewLPFarm(a, b *coingecko.Market, c coingecko.Fiat, initialInvestment float64, startDate string, apr float64) (*LPFarm, error) {
	f := &LPFarm{
		A:                 a,
		B:                 b,
		Currency:          c,
		InitialInvestment: initialInvestment,
		StartDate:         startDate,
		LastHarvestDate:   startDate,
		APR:               apr,
		InitialAPR:        apr,
		ChangeHistory:     make(map[string]*LPFarmHistoryItem),
	}

	pa, pb, err := f.GetPrices(startDate)
	if err != nil {
		return nil, err
	}

	f.UnitsA = initialInvestment / 2 / pa.V
	f.UnitsB = initialInvestment / 2 / pb.V
	f.InitialUnitsA = f.UnitsA
	f.InitialUnitsB = f.UnitsB

	f.RebalanceLP(pa, pb)

	// Include initial state in change history.
	f.LogChange(startDate, pa, pb)

	f.PrintChange(startDate)

	return f, nil
}

func (f *LPFarm) LogChange(date string, priceA, priceB timeseries.ValueAt) {
	// Update change history.
	f.ChangeHistory[date] = &LPFarmHistoryItem{
		Date:                date,
		PriceA:              priceA.V,
		PriceB:              priceB.V,
		UnitsA:              f.UnitsA,
		UnitsB:              f.UnitsB,
		TotalValue:          f.TotalValue,
		TotalValueHODL:      f.TotalValueHODL,
		TotalValueHODLOnlyA: f.InitialUnitsA * 2 * priceA.V,
		TotalValueHODLOnlyB: f.InitialUnitsB * 2 * priceB.V,
		APR:                 f.APR,
	}
}

func (f *LPFarm) GetChangeHistoryAsc() (items []*LPFarmHistoryItem) {
	for _, v := range f.ChangeHistory {
		items = append(items, v)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Date < items[j].Date
	})

	return
}

func (f *LPFarm) SetAPRChangeRateAtHarvest(dailyChange float64) {
	f.APRChangeRateAtHarvest = dailyChange
}

func (f *LPFarm) GetPrices(date string) (priceA, priceB timeseries.ValueAt, err error) {
	var found bool

	priceA, found = f.A.Prices.AtDate(date)
	if !found {
		err = errors.Errorf("could not find a price for %s on date %s", f.A.Symbol, date)
		return
	}

	priceB, found = f.B.Prices.AtDate(date)
	if !found {
		err = errors.Errorf("could not find a price for %s on date %s", f.B.Symbol, date)
		return
	}

	return
}

func (f *LPFarm) RebalanceLP(priceA, priceB timeseries.ValueAt) {
	k := f.UnitsA * f.UnitsB  // k value for impermanent loss calculations.
	rt := priceA.V / priceB.V // a price in b where a and b are the two assets in the pool

	f.UnitsA = math.Sqrt(k / rt)
	f.UnitsB = math.Sqrt(k * rt)

	f.TotalValue = (f.UnitsA * priceA.V) + (f.UnitsB * priceB.V)
	f.TotalValueHODL = (f.InitialUnitsA * priceA.V) + (f.InitialUnitsB * priceB.V)
}

func (f *LPFarm) PrintChange(date string) {
	if i, found := f.ChangeHistory[date]; found {
		il := (1 - (i.TotalValue / i.TotalValueHODL)) * 100.0

		pr := message.NewPrinter(language.English)

		pr.Printf("[%s] position  : %10.02f  (IL: %6.02f , hodl: %10.02f , APR: %6.02f %% , a: %10.02f , b: %10.02f , units: %.2f / %.2f)\n",
			i.Date, i.TotalValue, il, i.TotalValueHODL, i.APR, i.PriceA, i.PriceB, i.UnitsA, i.UnitsB,
		)
	}
}

func (f *LPFarm) AddLP(date string, amount float64) error {
	pa, pb, err := f.GetPrices(date)
	if err != nil {
		return err
	}

	// Split yield 50/50 between our asset pair.
	addUnitsA := amount / 2 / pa.V
	addUnitsB := amount / 2 / pb.V

	f.UnitsA += addUnitsA
	f.UnitsB += addUnitsB

	f.RebalanceLP(pa, pb)

	f.LogChange(date, pa, pb)

	return nil
}

func (f *LPFarm) Harvest(date string) (yield float64, err error) {
	if f.StartDate >= date {
		err = errors.Errorf("harvest date %s is before start date %s", date, f.StartDate)
		return
	}

	if f.LastHarvestDate >= date {
		err = errors.Errorf("harvest date %s is before last harvest date %s", date, f.LastHarvestDate)
		return
	}

	days := timeseries.DiffDays(date, f.LastHarvestDate)
	if days > 0 {
		// pr := message.NewPrinter(language.English)
		// pr.Printf("[%s] harvest ==================\n\n", date)
		// f.Print()

		if f.APRChangeRateAtHarvest > 0.0 {
			// Apply APR change.
			f.APR -= f.APRChangeRateAtHarvest
			if f.APR < 0.0 {
				f.APR = 0.0
			}
		}

		pa, pb, err2 := f.GetPrices(date)
		if err2 != nil {
			err = err2
			return
		}

		// Rebalance to get the latest fiat value of the farm.
		f.RebalanceLP(pa, pb)

		dailyPercentageRate := (f.APR / 100 / 365) * float64(days)
		yield = f.TotalValue * dailyPercentageRate

		f.LastHarvestDate = date

		f.RebalanceLP(pa, pb)

		f.LogChange(date, pa, pb)
		f.PrintChange(date)
	}

	return
}
