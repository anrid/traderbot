package coingecko

import "github.com/anrid/traderbot/pkg/timeseries"

type Fiat string

const (
	USD Fiat = "usd"
	EUR Fiat = "eur"
)

type Market struct {
	Currency     Fiat              `json:"currency"`
	Prices       timeseries.Series `json:"prices"`
	MarketCaps   timeseries.Series `json:"market_caps"`
	TotalVolumes timeseries.Series `json:"total_volumes"`

	ID                    string  `json:"id"`                               // "usd-coin"
	Symbol                string  `json:"symbol"`                           // "usdc"
	Name                  string  `json:"name"`                             // "USD Coin"
	Image                 string  `json:"image"`                            // "https://assets.coingecko.com/coins/images/6319/large/USD_Coin_icon.png?1547042389"
	CurrentPrice          float64 `json:"current_price"`                    // 1.001
	MarketCap             float64 `json:"market_cap"`                       // 53308409071
	MarketCapRank         int     `json:"market_cap_rank"`                  // 5
	FullyDilutedValuation float64 `json:"fully_diluted_valuation"`          // null
	TotalVolume           float64 `json:"total_volume"`                     // 4557737684
	High24h               float64 `json:"high_24h"`                         // 1.005
	Low24h                float64 `json:"low_24h"`                          // 0.994159
	PriceChange24h        float64 `json:"price_change_24h"`                 // -0.00079508518
	PriceChangePct24h     float64 `json:"price_change_percentage_24h"`      // -0.0794
	MarketCapChange24h    float64 `json:"market_cap_change_24h"`            // 1414826
	MarketCapPctChange24h float64 `json:"market_cap_change_percentage_24h"` // 0.00265
	CirculatingSupply     float64 `json:"circulating_supply"`               // 53277451267.3321
	TotalSupply           float64 `json:"total_supply"`                     // 53277706933.084
	MaxSupply             float64 `json:"max_supply"`                       // null
	ATH                   float64 `json:"ath"`                              // 1.17
	ATHChangePct          float64 `json:"ath_change_percentage"`            // -14.67759
	ATHDate               string  `json:"ath_date"`                         // "2019-05-08T00:40:28.300Z"
	ATL                   float64 `json:"atl"`                              // 0.891848
	ATLChangePct          float64 `json:"atl_change_percentage"`            // 12.19189
	ATLDate               string  `json:"atl_date"`                         // "2021-05-19T13:14:05.611Z"
	LastUpdated           string  `json:"last_updated"`                     // "2022-02-26T05:01:38.509Z"
}
