package tui

import tea "github.com/charmbracelet/bubbletea"

// All views commands
func (m *mainModel) getGlobalCommands(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "ctrl+q":
		cmds = append(cmds, tea.Quit)

	case "tab":
		switch m.currView {
		case listView:
			m.currView = splitView

		case splitView:
			m.currView = responseView

		default:
			m.currView = listView
		}
	}
	return cmds
}

// Each view commands
func (m *mainModel) getPerViewCommands(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd
	switch m.currView {
	case listView:
		switch msg.String() {
		case "enter":
			i, ok := m.queriesList.list.SelectedItem().(item)
			if ok {
				m.queriesList.selected = i.Title()
			}
			m.currView = fetchingView
			cmds = append(cmds, gqlReq(url, "./queries/"+m.queriesList.selected))
		}

	case responseView:
		if msg.String() == "backspace" {
			m.currView = listView
		}

	case spinnerView:
		if msg.String() == "n" {
			m.nextSpinner()
			m.resetSpinner()
			cmds = append(cmds, m.spinner.model.Tick)
		}

	case splitView:
		if msg.String() == "enter" {
			m.currView = responseView
		}
	}

	return cmds
}

func (m *mainModel) getUpdateViewsCommands(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

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

	return cmds
}
