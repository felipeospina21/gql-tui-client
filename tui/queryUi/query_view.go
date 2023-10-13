package queryui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/felipeospina21/gql-tui-client/query"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	url      = "https://rickandmortyapi.com/graphql"
)

type Model struct {
	list     list.Model
	selected string
	ready    IsReady
}

type IsReady bool

func NewModel() Model {
	items := []list.Item{}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Queries"

	return Model{list: l}
}

func (m Model) Init() tea.Cmd {
	return nil
	// return query.GetQueriesList("./queries/")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(query.Query)
			if ok {
				m.selected = i.Title()
			}
			return m, query.GqlReq(url, "./queries/"+m.selected)
		}
	case query.ListItems:
		cmd := m.list.SetItems(msg)
		ready := func() tea.Msg {
			return IsReady(true)
		}
		// m.ready = IsReady(true)
		cmds = append(cmds, cmd, ready)
		// return m, func() tea.Msg {
		// 	return IsReady(true)
		// }

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}
