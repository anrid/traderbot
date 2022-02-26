package coingecko

import (
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Fiat string

const (
	USD Fiat = "usd"
	EUR Fiat = "eur"
)

type Coin struct {
	ID           string     `json:"id"`
	Currency     Fiat       `json:"currency"`
	Prices       TimeSeries `json:"prices"`
	MarketCaps   TimeSeries `json:"market_caps"`
	TotalVolumes TimeSeries `json:"total_volumes"`
}

func NewCoin(id string, currency Fiat) *Coin {
	return &Coin{ID: strings.ToLower(id), Currency: currency}
}

func (c *Coin) PriceAtDate(date string) (price TimeSeriesValue, found bool) {
	for _, p := range c.Prices {
		if p.Date() == date {
			price = p
			found = true
			break
		}
	}
	return
}

type TimeSeries []TimeSeriesValue

type TimeSeriesValue struct {
	TS int64   `json:"ts"`
	V  float64 `json:"v"`
}

func (v TimeSeriesValue) Time() time.Time {
	return time.UnixMilli(int64(v.TS))
}

func (v TimeSeriesValue) Date() string {
	return v.Time().Format(dateFormat)
}

func (ts TimeSeries) Print() {
	pr := message.NewPrinter(language.English)

	for _, v := range ts {
		pr.Printf("[%s]  --  %.04f\n", v.Date(), v.V)
	}
}

func NewTimeSeries(tuples [][]interface{}) TimeSeries {
	var ts TimeSeries
	for _, tuple := range tuples {
		t := TimeSeriesValue{
			TS: int64(tuple[0].(float64)),
			V:  tuple[1].(float64),
		}
		ts = append(ts, t)
	}
	return ts
}
