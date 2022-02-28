package trade

import (
	"strings"
	"time"

	"github.com/anrid/traderbot/pkg/coingecko"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Side int

const (
	Buy Side = iota + 1
	Sell
)

type Trade struct {
	Currency coingecko.Fiat
	Side     Side
	Date     string
	Market   *coingecko.Market
	Size     float64 // a percentage expressed as a float64 in range (0.0 - 100.0]
	Price    float64
}

func NewBuyAtDate(date string, size float64, m *coingecko.Market) (*Trade, error) {
	return NewTrade(Buy, date, size, m)
}

func NewSellAtDate(date string, size float64, m *coingecko.Market) (*Trade, error) {
	return NewTrade(Sell, date, size, m)
}

func NewTrade(side Side, date string, size float64, m *coingecko.Market) (*Trade, error) {
	if size <= 0.0 || size > 100.0 {
		return nil, errors.Errorf("invalid size %f, must be a percentage expressed as a float64 in range (0.0 - 100.0]", size)
	}

	p, found := m.Prices.AtDate(date)
	if !found {
		return nil, errors.Errorf("could not find a price for `%s` at date %s", m.ID, date)
	}

	return &Trade{
		Currency: m.Currency,
		Side:     side,
		Date:     date,
		Market:   m,
		Size:     size,
		Price:    p.V,
	}, nil
}

func ExecuteTradesAndPrint(title string, initialInvestment float64, ts []*Trade) {
	if len(ts) == 0 {
		// Nothing to execute.
		return
	}

	pr := message.NewPrinter(language.English)
	m := ts[0].Market

	pr.Printf("\n\nTrading '%s' (%s) : %s\n", m.Name, strings.ToUpper(m.Symbol), title)
	pr.Printf("=============================================================\n\n")

	var totalFiat = initialInvestment
	var totalUnits float64

	var foundFirstBuy bool
	var buys int
	var sells int
	var firstBuyDate string
	var lastSellDate string

	for _, t := range ts {
		if t.Side == Buy {
			foundFirstBuy = true

			pct := t.Size / 100.0
			amount := totalFiat * pct

			units := amount / t.Price

			totalFiat -= amount
			totalUnits += units

			buys++
			if firstBuyDate == "" {
				firstBuyDate = t.Date
			}

			pr.Printf("%3d. [%s] %-4s %-10s @ %14.04f  --  amount: %14.04f , units: %14.04f\n",
				buys+sells, t.Date, "buy", t.Market.ID, t.Price, amount, units,
			)
		} else if t.Side == Sell {
			if !foundFirstBuy {
				// Skip all sell trades until we find first buy.
				continue
			}

			pct := t.Size / 100.0
			units := totalUnits * pct

			amount := units * t.Price

			totalFiat += amount
			totalUnits -= units

			sells++
			lastSellDate = t.Date

			pr.Printf("%3d. [%s] %-4s %-10s @ %14.04f  --  amount: %14.04f , units: %14.04f  [portfolio: %14.04f]\n",
				buys+sells, t.Date, "sell", t.Market.ID, t.Price, amount, units, totalFiat,
			)
		}
	}

	numTxns := buys + sells
	rangeStart, _ := time.Parse("2006-01-02", firstBuyDate)
	rangeEnd, _ := time.Parse("2006-01-02", lastSellDate)
	daysDiff := rangeEnd.Sub(rangeStart).Hours() / 24

	var totalFiatOfExistingPosition float64
	if buys > sells && totalUnits > 0 {
		// We ended on a buy; calculate the value of the position today.
		prices := ts[0].Market.Prices
		latestPrice := prices[len(prices)-1].V
		totalFiatOfExistingPosition = totalUnits * latestPrice
	}

	pr.Printf("\n\n")
	pr.Printf("- Number of txns     : %d\n", numTxns)
	pr.Printf("- First buy          : %s\n", firstBuyDate)
	pr.Printf("- Last sell          : %s  (%.f days after first buy)\n", lastSellDate, daysDiff)

	pr.Printf("- Initial investment : %.02f\n", initialInvestment)
	pr.Printf("- Portfolio value    : %.02f\n", totalFiat+totalFiatOfExistingPosition)

	pl := (totalFiat + totalFiatOfExistingPosition) / initialInvestment
	pr.Printf("- P/L                : %.02f %%\n\n\n", (pl-1)*100.0)
}
