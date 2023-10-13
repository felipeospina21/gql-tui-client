package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/felipeospina21/gql-tui-client/query"
	queryui "github.com/felipeospina21/gql-tui-client/tui/queryUi"
	spinnerui "github.com/felipeospina21/gql-tui-client/tui/spinnerUi"
)

const (
	spinnerView currView = iota
	listView
)

type (
	currView uint
	errMsg   struct{ err error }
)

type MainModel struct {
	currView currView
	list     tea.Model
	spinner  tea.Model
	// choice   string
	response query.ResponseMsg
	err      error
}

func (e errMsg) Error() string { return e.err.Error() }

func New() MainModel {
	return MainModel{currView: spinnerView, list: queryui.NewModel(), spinner: spinnerui.NewModel()}
}

func Start() {
	p := tea.NewProgram(New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m MainModel) Init() tea.Cmd {
	// return nil
	return query.GetQueriesList("./queries/")
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case queryui.IsReady:
		m.currView = listView
		// return m, tea.Quit

	default:
		m.currView = spinnerView
	}

	switch m.currView {
	case spinnerView:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case listView:
		m.list, cmd = m.list.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.currView {
	case spinnerView:
		return m.spinner.View()
	}
	return m.list.View()
}
