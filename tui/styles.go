package tui

import "github.com/charmbracelet/lipgloss"

var (
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	spinnerStyle  = lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).Foreground(lipgloss.Color("63"))
	textStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	titleStyle    = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)
