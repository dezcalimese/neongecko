package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Time-based colors
	dayBg    = lipgloss.Color("#D2B48C") // Tan color
	nightBg  = lipgloss.Color("#191970") // Midnight blue
	defaultBg = lipgloss.Color("#1a1a1a") // Dark default

	// Accent colors - bright pastels
	mintGreen = lipgloss.Color("#98FB98")
	lavender  = lipgloss.Color("#E6E6FA")
	peach     = lipgloss.Color("#FFDAB9")
	skyBlue   = lipgloss.Color("#87CEEB")
	pink      = lipgloss.Color("#FFB6C1")
	
	// Status colors
	green = lipgloss.Color("#00FF00")
	red   = lipgloss.Color("#FF0000")
	white = lipgloss.Color("#FFFFFF")
	black = lipgloss.Color("#000000")
)

func GetTimeBasedBg() lipgloss.Color {
	now := time.Now()
	hour := now.Hour()
	
	if hour >= 6 && hour < 18 {
		return dayBg
	} else if hour >= 18 && hour < 6 {
		return nightBg
	}
	return defaultBg
}

func GetTextColor() lipgloss.Color {
	bg := GetTimeBasedBg()
	if bg == dayBg {
		return black
	}
	return white
}

var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
		Background(GetTimeBasedBg()).
		Foreground(GetTextColor()).
		Padding(1, 2)

	// Title style
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(mintGreen).
		Background(GetTimeBasedBg()).
		Padding(0, 1).
		MarginBottom(1)

	// Header style
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(skyBlue).
		Background(GetTimeBasedBg()).
		Padding(0, 1)

	// Box style for containers
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lavender).
		Background(GetTimeBasedBg()).
		Padding(1, 2).
		Margin(1, 0)

	// Data label style
	LabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(peach).
		Background(GetTimeBasedBg())

	// Value style
	ValueStyle = lipgloss.NewStyle().
		Foreground(GetTextColor()).
		Background(GetTimeBasedBg())

	// Positive change style
	PositiveStyle = lipgloss.NewStyle().
		Foreground(green).
		Background(GetTimeBasedBg()).
		Bold(true)

	// Negative change style
	NegativeStyle = lipgloss.NewStyle().
		Foreground(red).
		Background(GetTimeBasedBg()).
		Bold(true)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
		Foreground(pink).
		Background(GetTimeBasedBg()).
		Italic(true).
		Padding(1, 0)

	// Search input style
	SearchStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mintGreen).
		Background(GetTimeBasedBg()).
		Foreground(GetTextColor()).
		Padding(0, 1)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
		Foreground(red).
		Background(GetTimeBasedBg()).
		Bold(true)
)

func FormatChange(value float64) (string, lipgloss.Style) {
	if value > 0 {
		return fmt.Sprintf("+%.2f%%", value), PositiveStyle
	} else if value < 0 {
		return fmt.Sprintf("%.2f%%", value), NegativeStyle
	}
	return "0.00%", ValueStyle
}

func FormatCurrency(value float64) string {
	if value >= 1e12 {
		return fmt.Sprintf("$%.2fT", value/1e12)
	} else if value >= 1e9 {
		return fmt.Sprintf("$%.2fB", value/1e9)
	} else if value >= 1e6 {
		return fmt.Sprintf("$%.2fM", value/1e6)
	} else if value >= 1e3 {
		return fmt.Sprintf("$%.2fK", value/1e3)
	}
	return fmt.Sprintf("$%.2f", value)
}