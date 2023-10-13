package spinnerui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var spinnerStyle = lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).Foreground(lipgloss.Color("63"))

type Model struct {
	Spin spinner.Model
}

func NewModel() Model {
	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.Pulse

	return Model{Spin: s}
}

func (m Model) Init() tea.Cmd {
	// return tea.Batch(m.Spin.Tick, query.GetQueriesList("./queries/"))
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spin, cmd = m.Spin.Update(msg)
		cmds = append(cmds, cmd)

		// case tea.WindowSizeMsg:
		// 	h, v := spinnerStyle.GetFrameSize()
		// 	m.spinner.
		// m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Center, m.Spin.View(), "loading")
}
