package models

import "time"

type GlobalData struct {
	TotalMarketCap         float64 `json:"total_market_cap"`
	TotalVolume           float64 `json:"total_volume"`
	MarketCapChangePercentage24h float64 `json:"market_cap_change_percentage_24h"`
}

type Coin struct {
	ID                        string    `json:"id"`
	Symbol                   string    `json:"symbol"`
	Name                     string    `json:"name"`
	CurrentPrice             float64   `json:"current_price"`
	MarketCap                float64   `json:"market_cap"`
	TotalVolume              float64   `json:"total_volume"`
	CirculatingSupply        float64   `json:"circulating_supply"`
	TotalSupply              *float64  `json:"total_supply"`
	AllTimeHigh              float64   `json:"ath"`
	AllTimeHighDate          time.Time `json:"ath_date"`
	AllTimeLow               float64   `json:"atl"`
	AllTimeLowDate           time.Time `json:"atl_date"`
	PriceChangePercentage24h  float64   `json:"price_change_percentage_24h"`
	PriceChangePercentage7d   float64   `json:"price_change_percentage_7d"`
	PriceChangePercentage30d  float64   `json:"price_change_percentage_30d"`
	PriceChangePercentage90d  float64   `json:"price_change_percentage_90d"`
}

type APIResponse struct {
	Global *GlobalData `json:"data,omitempty"`
	Coins  []Coin      `json:",omitempty"`
}