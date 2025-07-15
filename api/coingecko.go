package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"neongecko/config"
	"neongecko/models"
)

const (
	BaseURL = "https://api.coingecko.com/api/v3"
)

type Client struct {
	httpClient *http.Client
	cache      *Cache
	config     *config.Config
}

func NewClient(cfg *config.Config) *Client {
	if cfg == nil {
		// Fallback to default config if none provided
		cfg = &config.DefaultConfig
	}
	
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.GetTimeout(),
		},
		cache:  NewCache(cfg.GetCacheTTL()),
		config: cfg,
	}
}

func (c *Client) GetGlobalData() (*models.GlobalData, error) {
	cacheKey := "global_data"
	
	// Check cache first
	if cached, found := c.cache.Get(cacheKey); found {
		return cached.(*models.GlobalData), nil
	}
	
	url := fmt.Sprintf("%s/global", BaseURL)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch global data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Data struct {
			TotalMarketCap         map[string]float64 `json:"total_market_cap"`
			TotalVolume           map[string]float64 `json:"total_volume"`
			MarketCapChangePercentage24h float64 `json:"market_cap_change_percentage_24h"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	globalData := &models.GlobalData{
		TotalMarketCap:         response.Data.TotalMarketCap["usd"],
		TotalVolume:           response.Data.TotalVolume["usd"],
		MarketCapChangePercentage24h: response.Data.MarketCapChangePercentage24h,
	}
	
	// Cache the result
	c.cache.Set(cacheKey, globalData)
	
	return globalData, nil
}

func (c *Client) GetCoinData(coinID string) (*models.Coin, error) {
	cacheKey := fmt.Sprintf("coin_data_%s", coinID)
	
	// Check cache first
	if cached, found := c.cache.Get(cacheKey); found {
		return cached.(*models.Coin), nil
	}
	
	url := fmt.Sprintf("%s/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false", BaseURL, coinID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coin data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
		MarketData struct {
			CurrentPrice             map[string]float64 `json:"current_price"`
			MarketCap               map[string]float64 `json:"market_cap"`
			TotalVolume             map[string]float64 `json:"total_volume"`
			CirculatingSupply       float64            `json:"circulating_supply"`
			TotalSupply             *float64           `json:"total_supply"`
			AllTimeHigh             map[string]float64 `json:"ath"`
			AllTimeHighDate         map[string]string  `json:"ath_date"`
			AllTimeLow              map[string]float64 `json:"atl"`
			AllTimeLowDate          map[string]string  `json:"atl_date"`
			PriceChangePercentage24h float64           `json:"price_change_percentage_24h"`
			PriceChangePercentage7d  float64           `json:"price_change_percentage_7d"`
			PriceChangePercentage30d float64           `json:"price_change_percentage_30d"`
			PriceChangePercentage90d float64           `json:"price_change_percentage_90d"`
		} `json:"market_data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	athDate, _ := time.Parse("2006-01-02T15:04:05.000Z", response.MarketData.AllTimeHighDate["usd"])
	atlDate, _ := time.Parse("2006-01-02T15:04:05.000Z", response.MarketData.AllTimeLowDate["usd"])

	coinData := &models.Coin{
		ID:                        response.ID,
		Symbol:                   response.Symbol,
		Name:                     response.Name,
		CurrentPrice:             response.MarketData.CurrentPrice["usd"],
		MarketCap:                response.MarketData.MarketCap["usd"],
		TotalVolume:              response.MarketData.TotalVolume["usd"],
		CirculatingSupply:        response.MarketData.CirculatingSupply,
		TotalSupply:              response.MarketData.TotalSupply,
		AllTimeHigh:              response.MarketData.AllTimeHigh["usd"],
		AllTimeHighDate:          athDate,
		AllTimeLow:               response.MarketData.AllTimeLow["usd"],
		AllTimeLowDate:           atlDate,
		PriceChangePercentage24h: response.MarketData.PriceChangePercentage24h,
		PriceChangePercentage7d:  response.MarketData.PriceChangePercentage7d,
		PriceChangePercentage30d: response.MarketData.PriceChangePercentage30d,
		PriceChangePercentage90d: response.MarketData.PriceChangePercentage90d,
	}
	
	// Cache the result
	c.cache.Set(cacheKey, coinData)
	
	return coinData, nil
}

func (c *Client) SearchCoins(query string) ([]models.Coin, error) {
	cacheKey := fmt.Sprintf("search_%s", query)
	
	// Check cache first
	if cached, found := c.cache.Get(cacheKey); found {
		return cached.([]models.Coin), nil
	}
	
	url := fmt.Sprintf("%s/search?query=%s", BaseURL, query)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search coins: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Coins []struct {
			ID     string `json:"id"`
			Symbol string `json:"symbol"`
			Name   string `json:"name"`
		} `json:"coins"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	var coins []models.Coin
	for _, coin := range response.Coins {
		coins = append(coins, models.Coin{
			ID:     coin.ID,
			Symbol: coin.Symbol,
			Name:   coin.Name,
		})
	}

	// Cache the result
	c.cache.Set(cacheKey, coins)

	return coins, nil
}