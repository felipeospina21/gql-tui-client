package tui

import (
	"fmt"
	"os"

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
)

type (
	currView uint
)

type mainModel struct {
	currView    currView
	spinner     spinnerModel
	queriesList queriesList
	response    response
}

func newModel() mainModel {
	m := mainModel{currView: spinnerView}
	m.newSpinnerModel()
	m.newQueriesListModel()
	return m
}

func Start() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.model.Tick, getQueriesList("./queries/"))
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		headerHeight := lipgloss.Height(m.headerView())
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
		m.currView = responseView

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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
			// case "tab":
			// 	m.switchView()
		}

		switch m.currView {
		case listView:
			if msg.String() == "enter" {
				i, ok := m.queriesList.list.SelectedItem().(item)
				if ok {
					m.queriesList.selected = i.Title()
				}
				return m, gqlReq(url, "./queries/"+m.queriesList.selected)

			}

		case spinnerView:
			if msg.String() == "n" {
				m.nextSpinner()
				m.resetSpinner()
				cmds = append(cmds, m.spinner.model.Tick)
			}
		}
	}

	switch m.currView {
	case spinnerView:
		m.spinner.model, cmd = m.spinner.model.Update(msg)
		cmds = append(cmds, cmd)

	case listView:
		m.queriesList.list, cmd = m.queriesList.list.Update(msg)
		cmds = append(cmds, cmd)

	case responseView:
		m.response.model, cmd = m.response.model.Update(msg)
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	switch m.currView {
	case spinnerView:
		s += docStyle.Render(m.spinner.model.View())
	case listView:
		s += docStyle.Render(m.queriesList.list.View())
	case responseView:
		s += fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.response.model.View(), m.footerView())

	}

	return s
}
