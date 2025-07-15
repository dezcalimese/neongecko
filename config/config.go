package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Theme struct {
		ForceTheme   string `json:"force_theme"`   // "day", "night", or "" for auto
		CustomColors struct {
			DayBg    string `json:"day_bg"`
			NightBg  string `json:"night_bg"`
			Accent   string `json:"accent"`
		} `json:"custom_colors"`
	} `json:"theme"`
	
	API struct {
		CacheTTL    string `json:"cache_ttl"`    // Duration string like "5m"
		Timeout     string `json:"timeout"`      // Duration string like "10s"
		RateLimit   int    `json:"rate_limit"`   // Requests per minute
	} `json:"api"`
	
	Display struct {
		Currency       string   `json:"currency"`        // "usd", "eur", etc.
		DecimalPlaces  int      `json:"decimal_places"`  // Number of decimal places for prices
		ShowHelp       bool     `json:"show_help"`       // Show help on startup
		Favorites      []string `json:"favorites"`       // List of favorite coin IDs
	} `json:"display"`
}

var DefaultConfig = Config{
	Theme: struct {
		ForceTheme   string `json:"force_theme"`
		CustomColors struct {
			DayBg    string `json:"day_bg"`
			NightBg  string `json:"night_bg"`
			Accent   string `json:"accent"`
		} `json:"custom_colors"`
	}{
		ForceTheme: "", // Auto-detect
		CustomColors: struct {
			DayBg    string `json:"day_bg"`
			NightBg  string `json:"night_bg"`
			Accent   string `json:"accent"`
		}{
			DayBg:   "#F5DEB3", // Sandy beige
			NightBg: "#191970", // Midnight blue
			Accent:  "#98FB98", // Mint green
		},
	},
	API: struct {
		CacheTTL    string `json:"cache_ttl"`
		Timeout     string `json:"timeout"`
		RateLimit   int    `json:"rate_limit"`
	}{
		CacheTTL:  "5m",
		Timeout:   "10s",
		RateLimit: 30, // 30 requests per minute
	},
	Display: struct {
		Currency       string   `json:"currency"`
		DecimalPlaces  int      `json:"decimal_places"`
		ShowHelp       bool     `json:"show_help"`
		Favorites      []string `json:"favorites"`
	}{
		Currency:      "usd",
		DecimalPlaces: 2,
		ShowHelp:      false,
		Favorites:     []string{"bitcoin", "ethereum"},
	},
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".config", "neongecko")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	
	return filepath.Join(configDir, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	
	// If config doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := SaveConfig(&DefaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return &DefaultConfig, nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func (c *Config) GetCacheTTL() time.Duration {
	duration, err := time.ParseDuration(c.API.CacheTTL)
	if err != nil {
		return 5 * time.Minute // Default fallback
	}
	return duration
}

func (c *Config) GetTimeout() time.Duration {
	duration, err := time.ParseDuration(c.API.Timeout)
	if err != nil {
		return 10 * time.Second // Default fallback
	}
	return duration
}

func (c *Config) IsFavorite(coinID string) bool {
	for _, fav := range c.Display.Favorites {
		if fav == coinID {
			return true
		}
	}
	return false
}

func (c *Config) AddFavorite(coinID string) {
	if !c.IsFavorite(coinID) {
		c.Display.Favorites = append(c.Display.Favorites, coinID)
	}
}

func (c *Config) RemoveFavorite(coinID string) {
	for i, fav := range c.Display.Favorites {
		if fav == coinID {
			c.Display.Favorites = append(c.Display.Favorites[:i], c.Display.Favorites[i+1:]...)
			break
		}
	}
}