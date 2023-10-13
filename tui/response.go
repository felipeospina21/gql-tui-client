package tui

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/felipeospina21/gql-tui-client/utils"
	"github.com/machinebox/graphql"
)

type (
	responseMsg     string
	isResponseReady bool
	errMsg          struct{ err error }
)

type response struct {
	model   viewport.Model
	ready   isResponseReady
	content responseMsg
	err     error
}

func (e errMsg) Error() string { return e.err.Error() }

func gqlReq(url string, file string) tea.Cmd {
	return func() tea.Msg {
		b, err := os.ReadFile(file)
		utils.CheckError(err)

		var obj map[string]interface{}
		client := graphql.NewClient(url)
		req := graphql.NewRequest(string(b))
		err = client.Run(context.Background(), req, &obj)
		if err != nil {
			return error(err)
		}

		f := colorjson.NewFormatter()
		f.Indent = 2
		f.KeyColor = color.New(color.FgMagenta)

		s, err := f.Marshal(obj)
		utils.CheckError(err)

		return responseMsg(s)
	}
}

func (m *mainModel) headerView(queryName string) string {
	title := titleStyle.Render(queryName)
	line := strings.Repeat("─", max(0, m.response.model.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *mainModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.response.model.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.response.model.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *mainModel) setViewportViewSize(msg tea.WindowSizeMsg, headerHeight int, verticalMarginHeight int) tea.Cmd {
	if !m.response.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		m.response.model = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.response.model.YPosition = headerHeight
		m.response.model.HighPerformanceRendering = useHighPerformanceRenderer
		m.response.model.SetContent(string(m.response.content))
		m.response.ready = true

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		m.response.model.YPosition = headerHeight + 1
	} else {
		m.response.model.Width = msg.Width
		m.response.model.Height = msg.Height - verticalMarginHeight
	}
	if useHighPerformanceRenderer {
		// Render (or re-render) the whole viewport. Necessary both to
		// initialize the viewport and when the window is resized.
		//
		// This is needed for high-performance rendering only.
		// cmds = append(cmds, viewport.Sync(m.viewport.mod))
		return viewport.Sync(m.response.model)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
