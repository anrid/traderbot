package trade

import (
	"math"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type LPFarm struct {
	A               *coingecko.Market
	B               *coingecko.Market
	Currency        coingecko.Fiat
	StartDate       string
	APR             float64
	LastHarvestDate string
	UnitsA          float64
	UnitsB          float64
	LastPriceA      float64
	LastPriceB      float64
	InitialUnitsA   float64
	InitialUnitsB   float64
}

func NewLPFarm(a, b *coingecko.Market, c coingecko.Fiat, initialInvestment float64, startDate string, apr float64) (*LPFarm, error) {
	pa, found := a.Prices.AtDate(startDate)
	if !found {
		return nil, errors.Errorf("could not find price for asset %s on %s", a.Symbol, startDate)
	}
	pb, found := b.Prices.AtDate(startDate)
	if !found {
		return nil, errors.Errorf("could not find price for asset %s on %s", a.Symbol, startDate)
	}

	return &LPFarm{
		A:               a,
		B:               b,
		Currency:        c,
		StartDate:       startDate,
		LastHarvestDate: startDate,
		APR:             apr,
		UnitsA:          initialInvestment / 2 / pa.V,
		UnitsB:          initialInvestment / 2 / pb.V,
		LastPriceA:      pa.V,
		LastPriceB:      pb.V,
		InitialUnitsA:   initialInvestment / 2 / pa.V,
		InitialUnitsB:   initialInvestment / 2 / pb.V,
	}, nil
}

func (f *LPFarm) RebalanceLP(date string) error {
	a, found := f.A.Prices.AtDate(date)
	if !found {
		return errors.Errorf("could not find a price for %s on date %s", f.A.Symbol, date)
	}
	b, found := f.B.Prices.AtDate(date)
	if !found {
		return errors.Errorf("could not find a price for %s on date %s", f.B.Symbol, date)
	}

	k := f.UnitsA * f.UnitsB // k value for impermanent loss calculations.
	rt := a.V / b.V          // a price in b where a and b are the two assets in the pool

	f.UnitsA = math.Sqrt(k / rt)
	f.UnitsB = math.Sqrt(k * rt)

	f.LastPriceA = a.V
	f.LastPriceB = b.V

	return nil
}

func (f *LPFarm) Print() {
	pr := message.NewPrinter(language.English)

	pos := (f.UnitsA * f.LastPriceA) + (f.UnitsB * f.LastPriceB)
	ini := (f.InitialUnitsA * f.LastPriceA) + (f.InitialUnitsB * f.LastPriceB)
	il := (1 - (pos / ini)) * 100.0

	// pr.Printf("[%s] units a   : %f  (price: %f)\n", f.LastHarvestDate, f.UnitsA, f.LastPriceA)
	// pr.Printf("[%s] units b   : %f  (price: %f)\n", f.LastHarvestDate, f.UnitsB, f.LastPriceB)
	pr.Printf("[%s] position  : %10.02f  (IL: %6.02f , hodl: %10.02f , a: %10.02f , b: %10.02f)\n", f.LastHarvestDate, pos, il, ini, f.LastPriceA, f.LastPriceB)
	// pr.Printf("[%s] no LP     : %f\n", f.LastHarvestDate, ini)
	// pr.Printf("[%s] IL        : %f\n", f.LastHarvestDate, il)
}

func (f *LPFarm) Harvest(date string) error {
	if f.StartDate >= date {
		return errors.Errorf("harvest date %s is before start date %s", date, f.StartDate)
	}

	if f.LastHarvestDate >= date {
		return errors.Errorf("harvest date %s is before last harvest date %s", date, f.LastHarvestDate)
	}

	days := timeseries.DiffDays(date, f.LastHarvestDate)
	if days > 0 {
		// pr := message.NewPrinter(language.English)
		// pr.Printf("[%s] harvest ==================\n\n", date)
		// f.Print()

		yieldPct := (f.APR / 100 / 365) * float64(days)

		// pr.Println("")
		// pr.Printf("[%s] yield pct : %f\n", date, yieldPct)
		// pr.Printf("[%s] days diff : %d\n\n", date, days)

		f.UnitsA += f.UnitsA * yieldPct
		f.UnitsB += f.UnitsB * yieldPct

		f.LastHarvestDate = date

		err := f.RebalanceLP(date)
		if err != nil {
			return err
		}

		f.Print()
		// pr.Println("")
	}

	return nil
}
