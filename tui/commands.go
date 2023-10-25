package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/felipeospina21/gql-tui-client/utils"
)

// All views commands
func (m *mainModel) getGlobalCommands(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "ctrl+q":
		cmds = append(cmds, tea.Quit)

	case "ctrl+e":
		// m.envVars.isEditing = isEditingEnvVars(true)
		// m.envVars.textarea.Focus()
		m.currView = envVarsView
		m.envVars.textarea.Focus()

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
				cmds = append(cmds, gqlReq(m.queries.apiUrl, "./queries/"+m.queries.selected))
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
		if msg.String() == "ctrl+s" {
			// content := m.envVars.textarea.Value()
			// m := GetEnvVars(content)
			// fmt.Println(m)
			m := readEnvVars()
			s := stringifyEnvVars(m)
			fmt.Println(s)
			isRespReady := func() tea.Msg {
				// return isResponseReady(true)
				return tea.Quit()
			}
			cmds = append(cmds, isRespReady)
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
