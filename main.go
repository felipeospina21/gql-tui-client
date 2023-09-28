package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	url           = "https://rickandmortyapi.com/graphql"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
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

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Queries"

	return model{spinner: s, list: l, loading: true}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getQueriesList("./queries/"))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case listItems:
		cmd := m.list.SetItems(msg)
		m.loading = false
		return m, cmd

	case responseMsg:
		m.response = msg
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// if !m.loading {
	// 	return docStyle.Render(m.list.View())
	// }
	if m.response != "" {
		return quitTextStyle.Render(string(m.response))
	}
	return docStyle.Render(m.list.View())
	// return docStyle.Render(m.spinner.View())
}

func main() {
	p := tea.NewProgram(newModel())

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
