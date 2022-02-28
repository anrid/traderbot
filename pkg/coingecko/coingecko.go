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
	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/pkg/errors"
)

const (
	apiBaseURI = "https://api.coingecko.com/api/v3"
)

func New(c Fiat) *CoinGecko {
	return &CoinGecko{Currency: c}
}

type CoinGecko struct {
	Currency          Fiat
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

func (cg *CoinGecko) MarketChartWithCache(coinID string, days uint, i jsoncache.InvalidateCachePeriod) (*Market, error) {
	key := fmt.Sprintf("%s-%03d-days-%s", coinID, days, cg.Currency)

	c := new(Market)
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

func (cg *CoinGecko) Markets(ids ...string) ([]*Market, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	u := url.URL{}
	q := u.Query()
	q.Add("vs_currency", string(cg.Currency))
	if len(ids) > 0 {
		q.Add("ids", strings.Join(ids, ","))
	}
	q.Add("order", "market_cap_desc")
	q.Add("per_page", "100")
	q.Add("page", "1")
	q.Add("sparkline", "false")
	u.RawQuery = q.Encode()

	url := "/coins/markets?" + u.RawQuery
	// https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=terra-luna&order=market_cap_desc&per_page=100&page=1&sparkline=false

	resp := make([]*Market, 0)

	err := cg.pingAndGetJSON(ctx, apiBaseURI+url, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch markets for ids `%s`", strings.Join(ids, ","))
	}

	for _, c := range resp {
		c.Currency = cg.Currency
	}
	return resp, nil
}

func (cg *CoinGecko) MarketChart(coinID string, days uint) (*Market, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cs, err := cg.Markets(coinID)
	if err != nil {
		return nil, err
	}
	if len(cs) == 0 {
		return nil, errors.Errorf("could not find market for coin with id '%s'", coinID)
	}
	c := cs[0]

	u := url.URL{}
	q := u.Query()
	q.Add("vs_currency", string(c.Currency))
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

	err = cg.pingAndGetJSON(ctx, apiBaseURI+url, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch market chart for coin `%s`", c.ID)
	}

	// Convert CoinGecko result to TimeSeries data.
	c.Prices = timeseries.FromTuples(resp.Prices)
	c.MarketCaps = timeseries.FromTuples(resp.MarketCaps)
	c.TotalVolumes = timeseries.FromTuples(resp.TotalVolumes)

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
