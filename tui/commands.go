package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/felipeospina21/gql-tui-client/utils"
)

// All views commands
func (m *mainModel) getGlobalCommands(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch {
	case key.Matches(msg, keys.EnvVars):
		m.currView = envVarsView
		m.envVars.textarea.Focus()

	case key.Matches(msg, keys.Quit):
		cmds = append(cmds, tea.Quit)

	case key.Matches(msg, keys.Tab):
		switch m.currView {
		case listView:
			m.currView = splitView

		case splitView:
			m.currView = responseView

		default:
			m.currView = listView
		}

		// case key.Matches(msg, keys.Back):
		// 	f, err := tea.LogToFile("debug.log", "debug")
		// 	if err != nil {
		// 		fmt.Println("error", err)
		// 	}
		// 	f.WriteString(msg.String())
		// 	defer f.Close()
		// 	m.currView = listView
	}

	return cmds
}

// Each view commands
func (m *mainModel) getPerViewCommands(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd
	switch m.currView {
	case listView:
		switch {
		case key.Matches(msg, keys.Enter):
			i, ok := m.queries.list.SelectedItem().(item)
			if ok {
				m.queries.selected = i.Title()
			}
			m.currView = fetchingView

			if utils.IsStringEmpty(m.queries.apiUrl) {

				// TODO: Add this to the error interface and the error case in tui.go
				isRespReady := func() tea.Msg {
					return responseMsg("api url not provided, please add it to .env file as URL")
				}

				cmds = append(cmds, isRespReady)
			} else {
				cmds = append(cmds, m.gqlReq(m.queries.apiUrl, "./queries/"+m.queries.selected))
			}
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

	case envVarsView:
		switch {
		case key.Matches(msg, keys.SaveEnvVars):
			content := m.envVars.textarea.Value()
			utils.OverwriteEnvVars(content)
			m.currView = listView

		case key.Matches(msg, keys.Back):
			m.currView = listView
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
		m.queries.list, cmd = m.queries.list.Update(msg)
		cmds = append(cmds, cmd)

	case responseView:
		m.response.model, cmd = m.response.model.Update(msg)
		cmds = append(cmds, cmd)

	case envVarsView:
		m.envVars.textarea, cmd = m.envVars.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}
