package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	url                                 = "https://rickandmortyapi.com/graphql"
	useHighPerformanceRenderer          = false
	spinnerView                currView = iota
	listView
	responseView
	fetchingView
	splitView
)

type (
	currView uint
)

type mainModel struct {
	currView    currView
	spinner     spinnerModel
	queriesList queriesList
	help        help.Model
	response    response
}

func newModel() mainModel {
	m := mainModel{currView: spinnerView}
	m.response.ready = false
	m.help = help.New()
	m.help.ShowAll = true

	m.newSpinnerModel()
	m.newQueriesListModel()
	return m
}

func Start() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m mainModel) Init() tea.Cmd {
	f := getQueriesFolderPath()
	return tea.Batch(m.spinner.model.Tick, getQueriesList(f))
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		headerHeight := lipgloss.Height(m.headerView(m.queriesList.selected))
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		m.queriesList.list.SetSize(msg.Width-h, msg.Height-v)
		cmd := m.setViewportViewSize(msg, headerHeight, verticalMarginHeight)

		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case error:
		m.response.err = msg
		m.response.model.SetContent(string(m.response.err.Error()))
		isRespReady := func() tea.Msg {
			return isResponseReady(true)
		}
		cmds = append(cmds, isRespReady)

	case responseMsg:
		m.response.content = msg
		m.response.model.SetContent(string(m.response.content))

		isRespReady := func() tea.Msg {
			return isResponseReady(true)
		}
		cmds = append(cmds, isRespReady)

	case isListReady:
		m.currView = listView

	case isResponseReady:
		m.currView = splitView

	case listItems:
		cmd := m.queriesList.list.SetItems(msg)
		ready := func() tea.Msg {
			return isListReady(true)
		}
		cmds = append(cmds, cmd, ready)

	case spinner.TickMsg:
		m.spinner.model, cmd = m.spinner.model.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		cmds = append(cmds, m.getGlobalCommands(msg)...)
		cmds = append(cmds, m.getPerViewCommands(msg)...)
	}

	cmds = append(cmds, m.getUpdateViewsCommands(msg)...)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	var res string
	switch m.currView {
	case spinnerView:
		s += docStyle.Render(m.spinner.model.View())

	case listView:
		s += docStyle.Render(m.queriesList.list.View())

	case fetchingView:
		s += lipgloss.JoinHorizontal(lipgloss.Top, docStyle.Render(fmt.Sprintf("%4s", m.queriesList.list.View())), m.spinner.model.View())

	case splitView:
		res = fmt.Sprintf("%s\n%s\n%s", m.headerView(m.queriesList.selected), m.response.model.View(), m.footerView())
		s += lipgloss.JoinHorizontal(lipgloss.Top, docStyle.Render(fmt.Sprintf("%4s", m.queriesList.list.View())), res)

	case responseView:
		s += fmt.Sprintf("%s\n%s\n%s", m.headerView(m.queriesList.selected), m.response.model.View(), m.footerView())
		// s += helpStyle.Render(m.help.View(keys))

	}
	return s
}
