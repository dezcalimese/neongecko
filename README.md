# Crypto CLI 🚀

A beautiful, neofetch-style cryptocurrency CLI tool built with Go and Charm libraries. Display crypto market data with vibrant colors and smooth terminal UI.

## Features

- **📊 Market Overview**: View total crypto market cap, 24h volume, and percentage changes
- **🔍 Coin Search**: Look up any cryptocurrency by name or symbol  
- **📈 Detailed Stats**: Price, market cap, supply data, ATH/ATL, and performance metrics
- **🎨 Beautiful Design**: Time-based color themes with bright pastel accents
- **⌨️ Intuitive Navigation**: Vim-style keyboard shortcuts
- **🌅 Dynamic Theming**: Sandy beige for day (6 AM - 6 PM), midnight blue for night

## Installation

### Prerequisites
- Go 1.21 or later

### Build from Source
```bash
git clone <repository-url>
cd crypto-cli
go mod tidy
go build -o crypto-cli
./crypto-cli
```

### Install Globally
```bash
go install
```

## Usage

Run the application:
```bash
./crypto-cli
```

### Keyboard Controls

#### Navigation
- `/` or `s` - Search for a cryptocurrency
- `Tab` - Switch between home and search views
- `ESC` - Return to home screen
- `q` or `Ctrl+C` - Quit application

#### Data Management
- `r` - Refresh market data
- `h` - Show help (coming soon)

### Color Themes

The app automatically switches themes based on your local time:
- **🌅 Day Mode (6 AM - 6 PM)**: Sandy beige background
- **🌙 Night Mode (6 PM - 6 AM)**: Midnight blue background
- **🎨 Accent Colors**: Mint green, lavender, peach, sky blue, pink

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
Neofetch-style cryptocurrency information including:
- Current price and market cap
- Trading volume and supply data
- All-time high/low with dates
- Performance metrics (24h, 7d, 30d, 90d)

## Project Structure

```
crypto-cli/
├── main.go              # Application entry point
├── api/
│   └── coingecko.go    # CoinGecko API client
├── ui/
│   ├── home.go         # Home screen UI
│   ├── coin.go         # Coin detail UI
│   └── styles.go       # Color themes and styling
├── models/
│   └── coin.go         # Data models
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

## Future Features

- [ ] Historical price charts
- [ ] Portfolio tracking
- [ ] Price alerts
- [ ] Configuration file support
- [ ] API response caching
- [ ] Multiple currency support
- [ ] Favorites list

## Support

For issues or feature requests, please create an issue on GitHub.