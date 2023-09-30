package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/machinebox/graphql"
)

var (
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	spinnerStyle  = lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).Foreground(lipgloss.Color("63"))
	textStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	url           = "https://rickandmortyapi.com/graphql"
)

const (
	spinnerView currView = iota
	listView
)

type currView uint

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	currView currView
	list     list.Model
	spinner  spinner.Model
	loading  bool
	choice   string
	response responseMsg
	err      error
}

type (
	responseMsg string
	listItems   []list.Item
	errMsg      struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }

func newModel() model {
	// items := getQueriesList("./queries/")
	items := []list.Item{}

	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.Pulse

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Queries"

	return model{spinner: s, list: l, loading: true}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getQueriesList("./queries/"))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
			}
			return m, gqlReq(url, "./queries/"+m.choice)
		}
		switch m.currView {
		case spinnerView:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)

		case listView:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}

	case listItems:
		cmd := m.list.SetItems(msg)
		m.currView = listView
		cmds = append(cmds, cmd)

	case responseMsg:
		m.response = msg
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	if m.currView == spinnerView {
		s += lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), "loading")
	} else {
		s += docStyle.Render(m.list.View())
	}
	return s
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getQueriesList(rootDir string) tea.Cmd {
	return func() tea.Msg {
		// TODO: replace Walk with WalkDir func

		queriesNames := []list.Item{}
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !info.IsDir() {
				queriesNames = append(queriesNames, item{title: info.Name(), desc: "query"})
			}
			return nil
		})

		checkError(err)
		// TODO: remove time sleep
		time.Sleep(2 * time.Second)
		return listItems(queriesNames)
	}
}

func gqlReq(url string, file string) tea.Cmd {
	return func() tea.Msg {
		b, err := os.ReadFile(file)
		checkError(err)

		var obj map[string]interface{}
		client := graphql.NewClient(url)
		req := graphql.NewRequest(string(b))
		err = client.Run(context.Background(), req, &obj)
		checkError(err)

		f := colorjson.NewFormatter()
		f.Indent = 2
		f.KeyColor = color.New(color.FgMagenta)

		s, err := f.Marshal(obj)
		checkError(err)

		return responseMsg(s)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
