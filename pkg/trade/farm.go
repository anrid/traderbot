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
	A         *coingecko.Market
	B         *coingecko.Market
	Currency  coingecko.Fiat
	StartDate string

	APR           float64
	APRDailyDecay float64
	InitialAPR    float64

	LastHarvestDate string
	UnitsA          float64
	UnitsB          float64
	InitialUnitsA   float64
	InitialUnitsB   float64
	LastPriceA      float64
	LastPriceB      float64
	TotalValueFiat  float64
	HODLValueFiat   float64
}

func NewLPFarm(a, b *coingecko.Market, c coingecko.Fiat, initialInvestment float64, startDate string, apr float64) (*LPFarm, error) {
	f := &LPFarm{
		A:               a,
		B:               b,
		Currency:        c,
		StartDate:       startDate,
		LastHarvestDate: startDate,
		APR:             apr,
		InitialAPR:      apr,
	}

	pa, pb, err := f.GetPrices(startDate)
	if err != nil {
		return nil, err
	}

	f.UnitsA = initialInvestment / 2 / pa.V
	f.UnitsB = initialInvestment / 2 / pb.V
	f.InitialUnitsA = f.UnitsA
	f.InitialUnitsB = f.UnitsB

	err = f.RebalanceLP(pa, pb)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *LPFarm) SetAPRDailyDecay(dailyDecayPct float64) {
	f.APRDailyDecay = dailyDecayPct
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

func (f *LPFarm) RebalanceLP(priceA, priceB timeseries.ValueAt) error {
	k := f.UnitsA * f.UnitsB  // k value for impermanent loss calculations.
	rt := priceA.V / priceB.V // a price in b where a and b are the two assets in the pool

	f.UnitsA = math.Sqrt(k / rt)
	f.UnitsB = math.Sqrt(k * rt)

	f.LastPriceA = priceA.V
	f.LastPriceB = priceB.V

	f.TotalValueFiat = (f.UnitsA * priceA.V) + (f.UnitsB * priceB.V)
	f.HODLValueFiat = (f.InitialUnitsA * priceA.V) + (f.InitialUnitsB * priceB.V)

	return nil
}

func (f *LPFarm) Print() {
	pr := message.NewPrinter(language.English)

	pos := (f.UnitsA * f.LastPriceA) + (f.UnitsB * f.LastPriceB)
	ini := (f.InitialUnitsA * f.LastPriceA) + (f.InitialUnitsB * f.LastPriceB)
	il := (1 - (pos / ini)) * 100.0

	// pr.Printf("[%s] units a   : %f  (price: %f)\n", f.LastHarvestDate, f.UnitsA, f.LastPriceA)
	// pr.Printf("[%s] units b   : %f  (price: %f)\n", f.LastHarvestDate, f.UnitsB, f.LastPriceB)
	pr.Printf("[%s] position  : %10.02f  (IL: %6.02f , hodl: %10.02f , APR: %6.02f %% , a: %10.02f , b: %10.02f)\n",
		f.LastHarvestDate, pos, il, ini, f.APR, f.LastPriceA, f.LastPriceB,
	)
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

		// Perform APR decay.
		// We assume that APR will decrease over time.
		if f.APRDailyDecay > 0.0 {
			lowerAPR := f.APR - (f.APRDailyDecay * float64(days))
			if lowerAPR > f.InitialAPR/2 {
				// Limit decay to 50% of original APR.
				f.APR = lowerAPR
			}
		}

		pa, pb, err := f.GetPrices(date)
		if err != nil {
			return err
		}

		// Rebalance to get the latest fiat value of the farm.
		err = f.RebalanceLP(pa, pb)
		if err != nil {
			return err
		}

		dailyPercentageRate := (f.APR / 100 / 365) * float64(days)
		yield := f.TotalValueFiat * dailyPercentageRate
		// Split yield 50/50 between our asset pair.
		addUnitsA := yield / 2 / pa.V
		addUnitsB := yield / 2 / pb.V

		f.UnitsA += addUnitsA
		f.UnitsB += addUnitsB

		f.LastHarvestDate = date

		err = f.RebalanceLP(pa, pb)
		if err != nil {
			return err
		}

		f.Print()
		// pr.Println("")
	}

	return nil
}
