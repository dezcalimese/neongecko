# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build the application
go build -o neongecko

# Run the application  
./neongecko

# Install dependencies
go mod tidy

# Install globally
go install
```

## Architecture Overview

This is a terminal-based cryptocurrency CLI application built with Go and the Charm framework (Bubble Tea, Lip Gloss, Bubbles). The architecture follows a clean separation of concerns:

### Core Components

**Main Application (`main.go`)**
- Entry point that coordinates between two main views: `homeView` and `coinView`
- Handles global keyboard navigation and view switching
- Loads configuration on startup with graceful fallback to defaults
- Uses Bubble Tea's Model-View-Update pattern for state management

**API Layer (`api/`)**
- `coingecko.go`: HTTP client for CoinGecko API with three main endpoints (global data, coin details, search)
- `cache.go`: Thread-safe in-memory caching system with TTL and automatic cleanup
- All API calls are cached to respect rate limits (default 5-minute TTL)
- Configuration-driven timeouts and cache settings

**Configuration (`config/`)**
- JSON-based config system stored at `~/.config/crypto-cli/config.json`
- Auto-creates default config on first run
- Supports theme customization, API settings, display preferences, and favorites
- Configuration is dependency-injected into API clients and UI components

**UI Layer (`ui/`)**
- `home.go`: Market overview screen showing global crypto stats
- `coin.go`: Detailed coin view with search functionality and neofetch-style display
- `styles.go`: Centralized styling with time-based theming (day/night modes based on local time)
- Uses Lip Gloss for consistent styling and color management

**Data Models (`models/`)**
- `coin.go`: Data structures for API responses (GlobalData, Coin, APIResponse)
- Handles JSON marshaling/unmarshaling for CoinGecko API

### Key Design Patterns

**Time-based Theming**: The UI automatically switches between sandy beige (day) and midnight blue (night) backgrounds based on local time (6 AM - 6 PM for day mode).

**Caching Strategy**: All API responses are cached with configurable TTL to prevent rate limiting. Cache keys are specific to the request type and parameters.

**Configuration Dependency Injection**: The config is loaded once at startup and passed to all components that need it, allowing runtime customization of behavior.

**View State Management**: The main model manages view switching between home and coin views, with each view maintaining its own state and handling its own updates.

## Configuration File

The application creates `~/.config/crypto-cli/config.json` with settings for:
- Theme colors and force theme mode
- API cache TTL and HTTP timeout
- Display preferences (currency, decimal places)
- Favorites list for quick coin access

## API Integration

Uses CoinGecko's free API (no key required) for:
- Global market statistics (`/global` endpoint)
- Individual coin data (`/coins/{id}` endpoint) 
- Coin search (`/search` endpoint)

All API calls are automatically cached and the client is configured via the config system for timeouts and cache duration.