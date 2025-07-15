package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"neongecko/api"
	"neongecko/config"
	"neongecko/models"
)

type CoinModel struct {
	client       *api.Client
	textInput    textinput.Model
	viewport     viewport.Model
	coin         *models.Coin
	searchResults []models.Coin
	loading      bool
	err          error
	mode         string // "search" or "display"
	width        int
	height       int
}

func NewCoinModel(cfg *config.Config) CoinModel {
	ti := textinput.New()
	ti.Placeholder = "Enter coin name or symbol..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	vp := viewport.New(80, 20)
	vp.Style = BaseStyle

	return CoinModel{
		client:    api.NewClient(cfg),
		textInput: ti,
		viewport:  vp,
		mode:      "search",
	}
}

func (m CoinModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CoinModel) Reset() CoinModel {
	m.coin = nil
	m.err = nil
	m.loading = false
	m.mode = "search"
	m.textInput.SetValue("")
	m.textInput.Focus()
	m.textInput.SetCursor(0)
	return m
}

func (m CoinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update viewport size
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4 // Leave space for help text
		
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.mode == "display" {
				m.mode = "search"
				m.coin = nil
				m.err = nil
				m.textInput.SetValue("")
				m.textInput.Focus()
				m.textInput.SetCursor(0)
				return m, textinput.Blink
			}
		case "enter":
			if m.mode == "search" && m.textInput.Value() != "" {
				m.loading = true
				query := m.textInput.Value()
				return m, m.fetchCoinData(query)
			}
		}

		if m.mode == "search" {
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		} else if m.mode == "display" {
			// Handle keys in display mode
			switch msg.String() {
			case "/", "s":
				// Switch to search mode for another coin
				m.mode = "search"
				m.coin = nil
				m.err = nil
				m.textInput.SetValue("")
				m.textInput.Focus()
				m.textInput.SetCursor(0)
				return m, textinput.Blink
			}
		}

	case coinDataMsg:
		m.loading = false
		m.coin = (*models.Coin)(msg)
		m.mode = "display"
		
		// Clear the search input for next search
		m.textInput.SetValue("")
		m.textInput.SetCursor(0)
		
		return m, nil

	case errMsg:
		m.loading = false
		m.err = error(msg)
		return m, nil
	}

	return m, nil
}

func (m CoinModel) View() string {
	if m.loading {
		return BaseStyle.Render("Loading coin data...")
	}

	if m.err != nil {
		errorContent := ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\n" +
			HelpStyle.Render("Press ESC to go back to search")
		return BaseStyle.Width(m.width).Height(m.height).Render(errorContent)
	}

	switch m.mode {
	case "search":
		return m.renderSearch()
	case "display":
		return m.renderCoinDisplay()
	}

	return ""
}

func (m CoinModel) renderSearch() string {
	title := TitleStyle.Render("🔍 Search Cryptocurrency")
	searchBox := SearchStyle.Render(m.textInput.View())
	help := HelpStyle.Render("Enter coin name or symbol, then press Enter to search\nPress q to quit")

	// Center align all content
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		searchBox,
		"",
		help,
	)

	return BaseStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center).
		Render(content)
}

func (m CoinModel) renderCoinDisplay() string {
	if m.coin == nil {
		return ErrorStyle.Render("No coin data available")
	}

	// Create grid layout instead of scrollable content
	gridContent := m.renderCoinGrid()
	
	// Help text
	help := HelpStyle.Render("/,s: search • ESC: home • q: quit")
	
	// Center align the content
	content := lipgloss.JoinVertical(lipgloss.Center,
		gridContent,
		"",
		help,
	)

	return BaseStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center).
		Render(content)
}

func (m CoinModel) renderCoinGrid() string {
	if m.coin == nil {
		return ""
	}

	// Calculate responsive grid width
	maxWidth := m.width - 10 // Leave some margin
	cardWidth := 30
	if maxWidth < 80 {
		cardWidth = 25
	}

	// Header
	header := m.renderCoinHeader()

	// Create cards for different data sections
	priceCard := m.renderPriceCard(cardWidth)
	marketCard := m.renderMarketCard(cardWidth)
	supplyCard := m.renderSupplyCard(cardWidth)
	performanceCard := m.renderPerformanceCard(cardWidth)

	// Arrange in responsive grid
	var rows []string
	
	// Header takes full width
	rows = append(rows, header)
	rows = append(rows, "")

	// Arrange cards in rows based on screen width
	if m.width >= 100 {
		// Wide screen: 2x2 grid
		topRow := lipgloss.JoinHorizontal(lipgloss.Top, priceCard, "  ", marketCard)
		bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, supplyCard, "  ", performanceCard)
		rows = append(rows, topRow, "", bottomRow)
	} else {
		// Narrow screen: single column
		rows = append(rows, priceCard, "", marketCard, "", supplyCard, "", performanceCard)
	}

	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

func (m CoinModel) renderCoinContent() string {
	if m.coin == nil {
		return ""
	}

	var content strings.Builder

	// Coin header with ASCII art style
	header := m.renderCoinHeader()
	content.WriteString(header + "\n")

	// Price and market data
	priceData := m.renderPriceData()
	content.WriteString(priceData + "\n")

	// Supply data
	supplyData := m.renderSupplyData()
	content.WriteString(supplyData + "\n")

	// Historical data
	historicalData := m.renderHistoricalData()
	content.WriteString(historicalData + "\n")

	// Performance metrics
	performance := m.renderPerformance()
	content.WriteString(performance + "\n")

	return content.String()
}

func (m CoinModel) renderCoinHeader() string {
	if m.coin == nil {
		return ""
	}

	// Create neofetch-style header
	header := fmt.Sprintf("╭─ %s (%s) ─╮", 
		strings.ToUpper(m.coin.Name), 
		strings.ToUpper(m.coin.Symbol))
	
	return TitleStyle.Render(header)
}

func (m CoinModel) renderPriceData() string {
	if m.coin == nil {
		return ""
	}

	var lines []string

	// Current price
	lines = append(lines,
		LabelStyle.Render("Current Price: ") +
		ValueStyle.Render(FormatCurrency(m.coin.CurrentPrice)))

	// Market cap
	lines = append(lines,
		LabelStyle.Render("Market Cap: ") +
		ValueStyle.Render(FormatCurrency(m.coin.MarketCap)))

	// 24h volume
	lines = append(lines,
		LabelStyle.Render("24h Volume: ") +
		ValueStyle.Render(FormatCurrency(m.coin.TotalVolume)))

	content := strings.Join(lines, "\n")
	return BoxStyle.Render(content)
}

func (m CoinModel) renderSupplyData() string {
	if m.coin == nil {
		return ""
	}

	var lines []string

	// Circulating supply
	lines = append(lines,
		LabelStyle.Render("Circulating Supply: ") +
		ValueStyle.Render(fmt.Sprintf("%.0f %s", m.coin.CirculatingSupply, strings.ToUpper(m.coin.Symbol))))

	// Total supply
	if m.coin.TotalSupply != nil {
		lines = append(lines,
			LabelStyle.Render("Total Supply: ") +
			ValueStyle.Render(fmt.Sprintf("%.0f %s", *m.coin.TotalSupply, strings.ToUpper(m.coin.Symbol))))
	} else {
		lines = append(lines,
			LabelStyle.Render("Total Supply: ") +
			ValueStyle.Render("∞"))
	}

	content := strings.Join(lines, "\n")
	return BoxStyle.Render(content)
}

func (m CoinModel) renderHistoricalData() string {
	if m.coin == nil {
		return ""
	}

	var lines []string

	// All-time high
	athDate := m.coin.AllTimeHighDate.Format("Jan 2, 2006")
	lines = append(lines,
		LabelStyle.Render("All-Time High: ") +
		ValueStyle.Render(fmt.Sprintf("%s (%s)", FormatCurrency(m.coin.AllTimeHigh), athDate)))

	// All-time low
	atlDate := m.coin.AllTimeLowDate.Format("Jan 2, 2006")
	lines = append(lines,
		LabelStyle.Render("All-Time Low: ") +
		ValueStyle.Render(fmt.Sprintf("%s (%s)", FormatCurrency(m.coin.AllTimeLow), atlDate)))

	content := strings.Join(lines, "\n")
	return BoxStyle.Render(content)
}

func (m CoinModel) renderPerformance() string {
	if m.coin == nil {
		return ""
	}

	var lines []string
	lines = append(lines, LabelStyle.Render("Price Performance:"))

	// 24h change
	change24h, style24h := FormatChange(m.coin.PriceChangePercentage24h)
	lines = append(lines,
		LabelStyle.Render("  24h: ") + style24h.Render(change24h))

	// 7d change
	change7d, style7d := FormatChange(m.coin.PriceChangePercentage7d)
	lines = append(lines,
		LabelStyle.Render("  7d: ") + style7d.Render(change7d))

	// 30d change
	change30d, style30d := FormatChange(m.coin.PriceChangePercentage30d)
	lines = append(lines,
		LabelStyle.Render("  30d: ") + style30d.Render(change30d))

	// 90d change
	change90d, style90d := FormatChange(m.coin.PriceChangePercentage90d)
	lines = append(lines,
		LabelStyle.Render("  90d: ") + style90d.Render(change90d))

	content := strings.Join(lines, "\n")
	return BoxStyle.Render(content)
}

// New card-based rendering methods
func (m CoinModel) renderPriceCard(width int) string {
	if m.coin == nil {
		return ""
	}

	var lines []string
	lines = append(lines, HeaderStyle.Render("💰 Current Price"))
	lines = append(lines, "")
	lines = append(lines, ValueStyle.Render(FormatCurrency(m.coin.CurrentPrice)))

	content := strings.Join(lines, "\n")
	return BoxStyle.Width(width).Render(content)
}

func (m CoinModel) renderMarketCard(width int) string {
	if m.coin == nil {
		return ""
	}

	var lines []string
	lines = append(lines, HeaderStyle.Render("📊 Market Data"))
	lines = append(lines, "")
	lines = append(lines, 
		LabelStyle.Render("Market Cap: ") + ValueStyle.Render(FormatCurrency(m.coin.MarketCap)))
	lines = append(lines, 
		LabelStyle.Render("24h Volume: ") + ValueStyle.Render(FormatCurrency(m.coin.TotalVolume)))

	content := strings.Join(lines, "\n")
	return BoxStyle.Width(width).Render(content)
}

func (m CoinModel) renderSupplyCard(width int) string {
	if m.coin == nil {
		return ""
	}

	var lines []string
	lines = append(lines, HeaderStyle.Render("🪙 Supply Info"))
	lines = append(lines, "")
	lines = append(lines, 
		LabelStyle.Render("Circulating: ") + 
		ValueStyle.Render(fmt.Sprintf("%.0fM %s", m.coin.CirculatingSupply/1e6, strings.ToUpper(m.coin.Symbol))))

	if m.coin.TotalSupply != nil {
		lines = append(lines, 
			LabelStyle.Render("Total: ") + 
			ValueStyle.Render(fmt.Sprintf("%.0fM %s", *m.coin.TotalSupply/1e6, strings.ToUpper(m.coin.Symbol))))
	} else {
		lines = append(lines, 
			LabelStyle.Render("Total: ") + ValueStyle.Render("∞"))
	}

	content := strings.Join(lines, "\n")
	return BoxStyle.Width(width).Render(content)
}

func (m CoinModel) renderPerformanceCard(width int) string {
	if m.coin == nil {
		return ""
	}

	var lines []string
	lines = append(lines, HeaderStyle.Render("📈 Performance"))
	lines = append(lines, "")

	// 24h change
	change24h, style24h := FormatChange(m.coin.PriceChangePercentage24h)
	lines = append(lines, LabelStyle.Render("24h: ") + style24h.Render(change24h))

	// 7d change  
	change7d, style7d := FormatChange(m.coin.PriceChangePercentage7d)
	lines = append(lines, LabelStyle.Render("7d: ") + style7d.Render(change7d))

	// 30d change
	change30d, style30d := FormatChange(m.coin.PriceChangePercentage30d)
	lines = append(lines, LabelStyle.Render("30d: ") + style30d.Render(change30d))

	content := strings.Join(lines, "\n")
	return BoxStyle.Width(width).Render(content)
}

// Messages
type coinDataMsg *models.Coin

func (m CoinModel) fetchCoinData(query string) tea.Cmd {
	return func() tea.Msg {
		// First try to search for the coin
		searchResults, err := m.client.SearchCoins(query)
		if err != nil {
			return errMsg(err)
		}

		if len(searchResults) == 0 {
			return errMsg(fmt.Errorf("no coins found for '%s'", query))
		}

		// Get detailed data for the first result
		coinData, err := m.client.GetCoinData(searchResults[0].ID)
		if err != nil {
			return errMsg(err)
		}

		return coinDataMsg(coinData)
	}
}