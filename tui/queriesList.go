package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/felipeospina21/gql-tui-client/utils"
)

type (
	isListReady bool
	listItems   []list.Item
)

type item struct {
	title, desc string
}

type queriesList struct {
	list     list.Model
	selected string
	ready    isListReady
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m *mainModel) newQueriesListModel() {
	items := []list.Item{}

	m.queriesList.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.queriesList.list.Title = "Queries"
}

func getQueriesList(rootDir string) tea.Cmd {
	return func() tea.Msg {
		// // TODO: replace Walk with WalkDir func
		//
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

		utils.CheckError(err)
		// TODO: remove time sleep
		time.Sleep(2 * time.Second)
		return listItems(queriesNames)
	}
}
