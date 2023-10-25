package tui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/felipeospina21/gql-tui-client/config"
	"github.com/felipeospina21/gql-tui-client/utils"
)

type (
	isListReady bool
	listItems   []list.Item
)

type item struct {
	title, desc string
}

type queriesModel struct {
	list     list.Model
	selected string
	ready    isListReady
	apiUrl   string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m *mainModel) newQueriesModel(apiUrl string) {
	items := []list.Item{}

	m.queries.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.queries.list.Title = "Queries"
	m.queries.apiUrl = apiUrl
}

func getQueriesFolderPath() string {
	path := config.Read().Folder.Location

	// NOTE: Match substring that starts with $ and ends with /
	envVar := regexp.MustCompile(`\$(\w+)(?:/|$)`)
	idxs := envVar.FindAllStringSubmatchIndex(path, -1)

	var parsedPath string

	for _, idx := range idxs {
		partialMatch := path[idx[2]:idx[3]]
		v := os.Getenv(partialMatch)

		if v == "" {
			log.Fatalf("Didn't find %s env variable", path[idx[0]:idx[1]])
		}

		s := strings.ReplaceAll(path, "$"+partialMatch, v)
		parsedPath = s

	}

	// NOTE: Remove duplicated / symbols
	re := regexp.MustCompile(`/+`)
	result := re.ReplaceAllString(parsedPath, "/")

	return result
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

		utils.CheckError(err)
		return listItems(queriesNames)
	}
}
