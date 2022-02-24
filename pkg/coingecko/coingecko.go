package coingecko

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	apiBaseURI = "https://api.coingecko.com/api/v3"
	dateFormat = "2006-01-02"
)

func New() *CoinGecko {
	return new(CoinGecko)
}

type CoinGecko struct {
	hasSuccessfulPing bool
}

func (cg *CoinGecko) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp := struct {
		GeckoSays string `json:"gecko_says"`
	}{}

	err := cg.getJSON(ctx, apiBaseURI+"/ping", nil, &resp)
	if err != nil {
		return false
	}

	fmt.Printf("CoinGecko says: %s\n", resp.GeckoSays)
	if resp.GeckoSays != "(V3) To the Moon!" {
		log.Fatalf("CoinGecko returned an unexpected ping message: '%s'", resp.GeckoSays)
	}
	cg.hasSuccessfulPing = true

	return true
}

type Coin struct {
	ID           string      `json:"id"`
	Prices       TimeSeries  `json:"prices"`
	MarketCaps   TimeSeries  `json:"market_caps"`
	TotalVolumes TimeSeries  `json:"total_volumes"`
	Indicators   []Indicator `json:"indicators"`
}

func NewCoin(id string) *Coin {
	return &Coin{ID: strings.ToLower(id)}
}

func (cg *CoinGecko) MarketChartWithCache(coinID string, days uint, i jsoncache.InvalidateCachePeriod) (*Coin, error) {
	key := fmt.Sprintf("%s-%03d-days", coinID, days)

	c := new(Coin)
	err := jsoncache.Get(key, c, i)
	if err != nil {
		if err != jsoncache.ErrNotFound {
			return nil, err
		}

		// Perform call.
		c, err = cg.MarketChart(coinID, days)
		if err != nil {
			return nil, err
		}

		// Cache result.
		err = jsoncache.Set(key, c, i)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (cg *CoinGecko) MarketChart(coinID string, days uint) (*Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c := NewCoin(coinID)

	u := url.URL{}
	q := u.Query()
	q.Add("vs_currency", "usd")
	q.Add("days", fmt.Sprintf("%d", days))
	q.Add("interval", "daily")
	u.RawQuery = q.Encode()

	url := "/coins/" + c.ID + "/market_chart?" + u.RawQuery
	// 'https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=1&interval=daily'

	resp := struct {
		Prices       [][]interface{} `json:"prices"`
		MarketCaps   [][]interface{} `json:"market_caps"`
		TotalVolumes [][]interface{} `json:"total_volumes"`
	}{}

	err := cg.pingAndGetJSON(ctx, apiBaseURI+url, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch market chart for coin `%s`", c.ID)
	}

	// Convert CoinGecko result to TimeSeries data.
	c.Prices = NewTimeSeries(resp.Prices)
	c.MarketCaps = NewTimeSeries(resp.MarketCaps)
	c.TotalVolumes = NewTimeSeries(resp.TotalVolumes)

	return c, nil
}

func (cg *CoinGecko) pingAndGetJSON(ctx context.Context, url string, payload, response interface{}) error {
	// Ping once if we don't already have a successful ping.
	if !cg.hasSuccessfulPing {
		if !cg.Ping() {
			return errors.New("could not ping CoinGecko API")
		}
	}

	return cg.getJSON(ctx, url, payload, response)
}

func (cg *CoinGecko) getJSON(ctx context.Context, url string, payload, response interface{}) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return errors.Wrap(err, "could not marshal payload")
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, body)
	if err != nil {
		return errors.Wrap(err, "could not create new HTTP request")
	}

	req.Header.Add("accept", "application/json")
	// dump(req.Header.Get("accept"))

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not execute HTTP request")
	}
	if resp.StatusCode >= 400 {
		return errors.Errorf("got HTTP error code: %d", resp.StatusCode)
	}

	if response != nil {
		typ := resp.Header.Get("content-type")
		if !strings.Contains(typ, "application/json") {
			return errors.Errorf("excepted response to have header `content-type: application/json` but got `%s`", typ)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "could not read HTTP response body")
		}
		err = json.Unmarshal(data, response)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal HTTP response body")
		}
	}

	return nil
}

func Dump(o interface{}) {
	b, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println(string(b))
}

type TimeSeries []TimeSeriesValue

type TimeSeriesValue struct {
	TS int64   `json:"ts"`
	V  float64 `json:"v"`
}

func (v TimeSeriesValue) Time() time.Time {
	return time.UnixMilli(int64(v.TS))
}

func (v TimeSeriesValue) TimeString() string {
	return v.Time().Format(dateFormat)
}

func (ts TimeSeries) Print() {
	pr := message.NewPrinter(language.English)

	for _, v := range ts {
		pr.Printf("[%s]  --  %.04f\n", v.TimeString(), v.V)
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
