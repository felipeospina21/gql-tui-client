package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/felipeospina21/gql-tui-client/utils"
)

type (
	queryMsg     string
	isQueryReady bool
)

type queryModel struct {
	model   viewport.Model
	content queryMsg
	isReady isQueryReady
}

func readQuery(fileName string) tea.Cmd {
	return func() tea.Msg {
		b, err := os.ReadFile("./queries/" + fileName)
		// b, err := os.ReadFile(fmt.Sprintf("./queries/%s", fileName))
		utils.CheckError(err)
		return responseMsg(b)
	}
}
