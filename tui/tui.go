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
	// URL                                 = "https://rickandmortyapi.com/graphql"
	RESPONSE_RIGHT_MARGIN               = 2
	ENV_VARS_CHAR_LIMIT                 = 5000
	useHighPerformanceRenderer          = false
	spinnerView                currView = iota
	listView
	responseView
	fetchingView
	splitView
	envVarsView
)

type (
	currView uint
)

type mainModel struct {
	currView currView
	spinner  spinnerModel
	queries  queriesModel
	envVars  envVarModel
	help     help.Model
	response response
}

func newModel() mainModel {
	m := mainModel{currView: spinnerView}
	m.response.ready = false
	m.help = help.New()
	m.help.ShowAll = true

	ev := readEnvVars()
	s := stringifyEnvVars(ev)
	apiUrl := ev["URL"]

	m.newSpinnerModel()
	m.newQueriesModel(apiUrl)
	m.newEnvVarModel(ENV_VARS_CHAR_LIMIT, s)
	return m
}

func Start() {
	// p := tea.NewProgram(newModel())
	p := tea.NewProgram(newModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m mainModel) Init() tea.Cmd {
	path := getQueriesFolderPath()
	return tea.Batch(m.spinner.model.Tick, getQueriesList(path))
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case isEditingEnvVars:
		m.currView = envVarsView

	case isListReady:
		m.currView = listView

	case isResponseReady:
		m.currView = splitView

	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		headerHeight := lipgloss.Height(m.headerView(m.queries.selected))
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		m.queries.list.SetSize(msg.Width-h, msg.Height-v)
		cmd := m.setViewportViewSize(msg, headerHeight, verticalMarginHeight)

		m.envVars.textarea.SetHeight(msg.Height - v)
		m.envVars.textarea.SetWidth(msg.Width - h)

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
		m.setResponseContent()

		isRespReady := func() tea.Msg {
			return isResponseReady(true)
		}
		cmds = append(cmds, isRespReady)

	case listItems:
		cmd := m.queries.list.SetItems(msg)
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

	renderView := map[string]string{
		"spinner":  spinnerStyle.Render(m.spinner.model.View()),
		"list":     listStyle.Render(m.queries.list.View()),
		"response": fmt.Sprintf("%s\n%s\n%s", m.headerView(m.queries.selected), m.response.model.View(), m.footerView()),
		"envVars":  varsStyle.Render(m.envVars.textarea.View()),
		// "response": lipgloss.NewStyle().Width(m.response.model.Width).Render(m.headerView(m.queriesList.selected), string(m.response.content), m.footerView()),
	}

	switch m.currView {
	case spinnerView:
		s += renderView["spinner"]

	case listView:
		s += renderView["list"]

	case fetchingView:
		s += lipgloss.JoinHorizontal(lipgloss.Top, renderView["list"], renderView["spinner"])

	case splitView:
		s += lipgloss.JoinHorizontal(lipgloss.Top, renderView["list"], renderView["response"])

	case responseView:
		s += renderView["response"]
		// s += helpStyle.Render(m.help.View(keys))

	case envVarsView:
		s += renderView["envVars"]

	}
	// s += helpStyle.Render(m.help.View(keys))
	return s
}
