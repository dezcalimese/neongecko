# NeonGecko 🦎

A beautiful, responsive cryptocurrency CLI tool built with Go and Charm libraries. Display crypto market data with vibrant colors, smart grid layouts, and smooth terminal UI.

## Features

- **📊 Market Overview**: View total crypto market cap, 24h volume, and percentage changes
- **🔍 Coin Search**: Look up any cryptocurrency by name or symbol  
- **📈 Responsive Grid Layout**: Clean card-based display that adapts to terminal size
- **💰 Detailed Stats**: Price, market cap, supply data, and performance metrics in organized cards
- **🎨 Beautiful Design**: Time-based color themes with bright pastel accents
- **⚡ Smart Caching**: Intelligent API caching to respect rate limits
- **⌨️ Intuitive Navigation**: Quick search from any view, seamless switching
- **🌅 Dynamic Theming**: Sandy beige for day (6 AM - 6 PM), midnight blue for night

## Installation

### Prerequisites
- Go 1.21 or later

### Build from Source
```bash
git clone https://github.com/dezcalimese/neongecko.git
cd neongecko
go mod tidy
go build -o neongecko
./neongecko
```

### Install Globally
```bash
go install
```

## Usage

Run the application:
```bash
./neongecko
```

### Keyboard Controls

#### Navigation
- `/` or `s` - Search for a cryptocurrency (works from any view)
- `Tab` - Switch between home and search views
- `ESC` - Return to home screen from coin view, or search mode from coin display
- `q` or `Ctrl+C` - Quit application

#### Data Management
- `r` - Refresh market data
- Auto-refresh with smart caching (5-minute TTL by default)

### Color Themes

The app automatically switches themes based on your local time:
- **🌅 Day Mode (6 AM - 6 PM)**: Sandy beige background
- **🌙 Night Mode (6 PM - 6 AM)**: Midnight blue background
- **🎨 Accent Colors**: Mint green, lavender, peach, pink, powder blue

## API

This tool uses the [CoinGecko API](https://www.coingecko.com/en/api) for cryptocurrency data:
- Global market statistics
- Individual coin data
- Search functionality
- No API key required (uses free tier)

## Screenshots

### Home Screen
Display of total market cap, 24h volume, and market change percentage with time-based color themes.

### Coin View
Responsive grid layout displaying cryptocurrency information in organized cards:
- **💰 Current Price**: Live price data
- **📊 Market Data**: Market cap and 24h trading volume  
- **🪙 Supply Info**: Circulating and total supply
- **📈 Performance**: 24h, 7d, and 30d percentage changes

Layout automatically adapts:
- **Wide terminals (≥100 chars)**: 2×2 grid layout
- **Narrow terminals (<100 chars)**: Single column layout

## Project Structure

```
neongecko/
├── main.go              # Application entry point & view management
├── api/
│   ├── coingecko.go    # CoinGecko API client
│   └── cache.go        # Thread-safe caching system
├── ui/
│   ├── home.go         # Home screen UI
│   ├── coin.go         # Responsive coin detail UI with grid layout
│   └── styles.go       # Time-based color themes and styling
├── models/
│   └── coin.go         # Data models for API responses
├── config/
│   └── config.go       # Configuration management
├── CLAUDE.md           # Development guidance
└── README.md
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library  
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Features Already Implemented

- [x] **Responsive Grid Layout**: Cards automatically arrange based on terminal size
- [x] **Smart API Caching**: Thread-safe caching with configurable TTL (default 5 minutes)
- [x] **Configuration System**: JSON-based config with auto-creation and defaults
- [x] **Direct Search**: Search for new coins without returning to home view
- [x] **Time-based Theming**: Automatic day/night mode switching
- [x] **Clean State Management**: Proper state clearing when switching views

## Future Features

- [ ] Historical price charts
- [ ] Portfolio tracking  
- [ ] Price alerts
- [ ] Multiple currency support
- [ ] Favorites list with quick access
- [ ] Custom themes and color schemes

## Support

For issues or feature requests, please create an issue on GitHub.