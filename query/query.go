package query

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/felipeospina21/gql-tui-client/utils"
	"github.com/machinebox/graphql"
)

type Query struct {
	title, desc string
}

type (
	ListItems   []list.Item
	ResponseMsg string
)

func (q Query) Title() string       { return q.title }
func (q Query) Description() string { return q.desc }
func (q Query) FilterValue() string { return q.title }

func GetQueriesList(rootDir string) tea.Cmd {
	return func() tea.Msg {
		// TODO: replace Walk with WalkDir func

		queriesNames := []list.Item{}
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !info.IsDir() {
				queriesNames = append(queriesNames, Query{title: info.Name(), desc: "query"})
			}
			return nil
		})

		utils.CheckError(err)
		// TODO: remove time sleep
		time.Sleep(2 * time.Second)
		return ListItems(queriesNames)
	}
}

func GqlReq(url string, file string) tea.Cmd {
	return func() tea.Msg {
		b, err := os.ReadFile(file)
		utils.CheckError(err)

		var obj map[string]interface{}
		client := graphql.NewClient(url)
		req := graphql.NewRequest(string(b))
		err = client.Run(context.Background(), req, &obj)
		utils.CheckError(err)

		f := colorjson.NewFormatter()
		f.Indent = 2
		f.KeyColor = color.New(color.FgMagenta)

		s, err := f.Marshal(obj)
		utils.CheckError(err)

		return ResponseMsg(s)
	}
}

func X(rootDir string) []list.Item {
	queriesNames := []list.Item{}
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			queriesNames = append(queriesNames, Query{title: info.Name(), desc: "query"})
		}
		return nil
	})

	utils.CheckError(err)
	// TODO: remove time sleep
	time.Sleep(2 * time.Second)
	return ListItems(queriesNames)
}
