package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"neongecko/api"
	"neongecko/config"
	"neongecko/models"
)

type HomeModel struct {
	client     *api.Client
	globalData *models.GlobalData
	loading    bool
	err        error
	width      int
	height     int
}

func NewHomeModel(cfg *config.Config) HomeModel {
	return HomeModel{
		client:  api.NewClient(cfg),
		loading: true,
	}
}

func (m HomeModel) Init() tea.Cmd {
	return m.fetchGlobalData
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, m.fetchGlobalData
		case "/", "s":
			// TODO: Switch to search view
			return m, nil
		case "h":
			// TODO: Show help
			return m, nil
		}

	case globalDataMsg:
		m.loading = false
		m.globalData = (*models.GlobalData)(msg)
		return m, nil

	case errMsg:
		m.loading = false
		m.err = error(msg)
		return m, nil
	}

	return m, nil
}

func (m HomeModel) View() string {
	if m.loading {
		return BaseStyle.
			Align(lipgloss.Center).
			Render("Loading global crypto data...")
	}

	if m.err != nil {
		return BaseStyle.
			Align(lipgloss.Center).
			Render(ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	if m.globalData == nil {
		return BaseStyle.
			Align(lipgloss.Center).
			Render(ErrorStyle.Render("No data available"))
	}

	// Title
	title := TitleStyle.Render("🚀 Crypto Market Overview")
	
	// Market data box
	marketData := m.renderMarketData()
	
	// Help text
	help := m.renderHelp()

	// Center align all content
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		marketData,
		"",
		help,
	)

	return BaseStyle.
		Align(lipgloss.Center).
		Render(content)
}

func (m HomeModel) renderMarketData() string {
	if m.globalData == nil {
		return ""
	}

	var lines []string

	// Total Market Cap
	lines = append(lines, 
		LabelStyle.Render("Total Market Cap: ") + 
		ValueStyle.Render(FormatCurrency(m.globalData.TotalMarketCap)))

	// Market Cap Change
	changeText, changeStyle := FormatChange(m.globalData.MarketCapChangePercentage24h)
	lines = append(lines,
		LabelStyle.Render("24h Change: ") +
		changeStyle.Render(changeText))

	// Total Volume
	lines = append(lines,
		LabelStyle.Render("24h Volume: ") +
		ValueStyle.Render(FormatCurrency(m.globalData.TotalVolume)))

	content := strings.Join(lines, "\n")
	return BoxStyle.Render(content)
}

func (m HomeModel) renderHelp() string {
	helpText := []string{
		"Navigation:",
		"• / or s - Search for a coin",
		"• r - Refresh data", 
		"• h - Show help",
		"• q or Ctrl+C - Quit",
	}
	
	return HelpStyle.Render(strings.Join(helpText, "\n"))
}

// Messages
type globalDataMsg *models.GlobalData
type errMsg error

func (m HomeModel) fetchGlobalData() tea.Msg {
	data, err := m.client.GetGlobalData()
	if err != nil {
		return errMsg(err)
	}
	return globalDataMsg(data)
}