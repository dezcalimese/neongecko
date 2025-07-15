package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"neongecko/config"
	"neongecko/ui"
)

type view int

const (
	homeView view = iota
	coinView
)

type mainModel struct {
	currentView view
	homeModel   ui.HomeModel
	coinModel   ui.CoinModel
	config      *config.Config
	width       int
	height      int
}

func initialModel() mainModel {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Failed to load config, using defaults: %v", err)
		cfg = &config.DefaultConfig
	}

	return mainModel{
		currentView: homeView,
		homeModel:   ui.NewHomeModel(cfg),
		coinModel:   ui.NewCoinModel(cfg),
		config:      cfg,
	}
}

func (m mainModel) Init() tea.Cmd {
	return m.homeModel.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Forward to current view
		switch m.currentView {
		case homeView:
			var cmd tea.Cmd
			model, cmd := m.homeModel.Update(msg)
			m.homeModel = model.(ui.HomeModel)
			return m, cmd
		case coinView:
			var cmd tea.Cmd
			model, cmd := m.coinModel.Update(msg)
			m.coinModel = model.(ui.CoinModel)
			return m, cmd
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "/", "s":
			if m.currentView == homeView {
				m.currentView = coinView
				m.coinModel = m.coinModel.Reset()
				return m, m.coinModel.Init()
			}
		case "tab":
			// Switch between views
			if m.currentView == homeView {
				m.currentView = coinView
				m.coinModel = m.coinModel.Reset()
				return m, m.coinModel.Init()
			} else {
				m.currentView = homeView
				return m, m.homeModel.Init()
			}
		case "esc":
			if m.currentView == coinView {
				m.currentView = homeView
				return m, nil
			}
		}
	}

	// Forward to current view
	switch m.currentView {
	case homeView:
		var cmd tea.Cmd
		model, cmd := m.homeModel.Update(msg)
		m.homeModel = model.(ui.HomeModel)
		return m, cmd
	case coinView:
		var cmd tea.Cmd
		model, cmd := m.coinModel.Update(msg)
		m.coinModel = model.(ui.CoinModel)
		return m, cmd
	}

	return m, nil
}

func (m mainModel) View() string {
	switch m.currentView {
	case homeView:
		return m.homeModel.View()
	case coinView:
		return m.coinModel.View()
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}