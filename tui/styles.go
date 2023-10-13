package tui

import "github.com/charmbracelet/lipgloss"

var (
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	textStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)
