package messari

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/anrid/traderbot/pkg/jsoncache"
	"github.com/pkg/errors"
)

const (
	apiBaseURI = "https://data.messari.io/api"
)

func New(token string) *Messari {
	return &Messari{token}
}

type Messari struct {
	token string
}

func (cg *Messari) AssetsWithCache(i jsoncache.InvalidateCachePeriod) (as []*Asset, err error) {
	key := "messari-assets"

	err = jsoncache.Get(key, &as, i)
	if err != nil {
		if err != jsoncache.ErrNotFound {
			return nil, err
		}

		// Perform call.
		as, err = cg.Assets()
		if err != nil {
			return nil, err
		}

		// Cache result.
		err = jsoncache.Set(key, as, i)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Downloaded data   : %s\n", key)

	} else {
		fmt.Printf("Using cached data : %s\n", key)
	}

	return
}

func (cg *Messari) Assets() (as []*Asset, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	page := 1

	for {
		u := url.URL{}
		q := u.Query()
		q.Add("fields", "id,symbol,name,slug,metrics/market_data/price_usd,metrics/supply,metrics/all_time_high")
		q.Add("page", strconv.Itoa(page))
		q.Add("limit", "500")
		u.RawQuery = q.Encode()

		url := "/v2/assets?" + u.RawQuery
		// https://data.messari.io/api/v2/assets?fields=id,slug,symbol,metrics/market_data/price_usd

		var resp AssetsResponse
		var errResp ErrorResponse

		errResp, err = cg.getJSON(ctx, apiBaseURI+url, nil, &resp)
		if err != nil {
			if strings.Contains(errResp.Status.ErrorMessage, "Rate limit") {
				fmt.Printf("Rate limited, retrying in 10 seconds ...")
				time.Sleep(10 * time.Second)
				continue
			}
			if errResp.Status.ErrorCode == 404 {
				fmt.Printf("Got 404 error, assuming there are no more pages to fetch!")
				err = nil
				break
			}

			return nil, errors.Wrapf(err, "could not fetch assets")
		}

		if len(resp.Data) == 0 {
			return
		}

		as = append(as, resp.Data...)
		page += 1
	}

	return
}

func (cg *Messari) getJSON(ctx context.Context, url string, payload, response interface{}) (errResp ErrorResponse, err error) {
	var body io.Reader
	if payload != nil {
		var b []byte
		b, err = json.Marshal(payload)
		if err != nil {
			err = errors.Wrap(err, "could not marshal payload")
			return
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, body)
	if err != nil {
		err = errors.Wrap(err, "could not create new HTTP request")
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-messari-api-key", cg.token)

	fmt.Printf("Get JSON: %s\n", req.URL)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		err = errors.Wrap(err, "could not execute HTTP request")
		return
	}

	if response != nil {
		typ := resp.Header.Get("content-type")
		if !strings.Contains(typ, "application/json") {
			err = errors.Errorf("excepted response to have header `content-type: application/json` but got `%s`", typ)
			return
		}

		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			err = errors.Wrap(err, "could not read HTTP response body")
			return
		}

		if false {
			du := make(map[string]interface{})
			err = json.Unmarshal(data, &du)
			if err != nil {
				err = errors.Wrap(err, "could not unmarshal HTTP response body")
				return
			}
			Dump(du)
		}

		if resp.StatusCode >= 400 {
			err = json.Unmarshal(data, &errResp)
			if err != nil {
				err = errors.Wrap(err, "could not unmarshal HTTP error response body")
				return
			}
			err = errors.Errorf("got HTTP error code: %d", resp.StatusCode)
			return
		}

		err = json.Unmarshal(data, response)
		if err != nil {
			err = errors.Wrap(err, "could not unmarshal HTTP response body")
			return
		}
	}

	return
}

func Dump(o interface{}) {
	b, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println(string(b))
}

type AssetsResponse struct {
	Status struct {
		Timestamp string `json:"timestamp"` // "2018-06-02T22:51:28.209Z"
		Elapsed   int64  `json:"elapsed"`   // 10
	} `json:"status"`
	Data []*Asset `json:"data"`
}

type Asset struct {
	ID      string `json:"id"`     // "1e31218a-e44e-4285-820c-8282ee222035",
	Symbol  string `json:"symbol"` // "BTC",
	Name    string `json:"name"`   // "Bitcoin",
	Slug    string `json:"slug"`   // "bitcoin",
	Metrics struct {
		MarketData struct {
			PriceUSD float64 `json:"price_usd"`
		} `json:"market_data"`
		Supply struct {
			Y2050              float64 `json:"y_2050"`                   // 20983495.3984375,
			Y2050IssuedPct     float64 `json:"y_2050_issued_percent"`    // 20,
			YPlus10            float64 `json:"y_plus10"`                 // 8932344.3984375,
			YPlus10IssuedPct   float64 `json:"y_plus10_issued_percent"`  // 40,
			Liquid             float64 `json:"liquid"`                   // 1982345,
			Circulating        float64 `json:"circulating"`              // 17394725,
			StockToFlow        float64 `json:"stock_to_flow"`            // 0
			AnnualInflationPct float64 `json:"annual_inflation_percent"` //  1.7633
		} `json:"supply"`
		DeveloperActivity struct {
			Stars              int `json:"stars"`                 // 34996,
			Watchers           int `json:"watchers"`              // 3513,
			CommitsLast3Months int `json:"commits_last_3_months"` // 342,
			CommitsLast1Year   int `json:"commits_last_1_year"`   // 1775,
		} `json:"developer_activity"`
		AllTimeHigh struct {
			Price       float64 `json:"price"`        // 20089,
			At          string  `json:"at"`           // "2018-06-02T22:51:28.209Z",
			DaysSince   int     `json:"days_since"`   // 344,
			PercentDown float64 `json:"percent_down"` // 81.47285775644839
		} `json:"all_time_high"`
	} `json:"metrics"`
}

type ErrorResponse struct {
	Status struct {
		Timestamp    string `json:"timestamp"`     // Current ISO 8601 timestamp on the server.
		ErrorCode    int    `json:"error_code"`    // Internal error code generated or 400 if default.
		ErrorMessage string `json:"error_message"` // Corresponding error message for the code.
		Elapsed      int    `json:"elapsed"`       // Number of milliseconds taken to generate this response.
	} `json:"status"`
}
