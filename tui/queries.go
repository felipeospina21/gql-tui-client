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
	token    string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m *mainModel) newQueriesModel(apiUrl string, token string) {
	items := []list.Item{}

	m.queries.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.queries.list.Title = "Queries"
	m.queries.apiUrl = apiUrl
	m.queries.token = token
}

func getQueriesFolderPath() string {
	path := config.Read().Queries.Location

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
		queriesNames := []list.Item{}
		err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}

			if !d.IsDir() {
				folder := getFolderName(path)
				queriesNames = append(queriesNames, item{title: d.Name(), desc: folder})
			}

			return nil
		})

		utils.CheckError(err)
		return listItems(queriesNames)
	}
}

func getFolderName(path string) string {
	index := strings.Index(path, "queries/")
	if index == -1 {
		fmt.Println("The input does not contain 'queries/'.")
	}

	// Extract everything after "queries/"
	result := path[index+len("queries/"):]
	spl := strings.Split(result, "/")

	var folder string
	b := strings.Contains(spl[len(spl)-1], ".gql")

	if len(spl) >= 2 && b {
		folder = spl[len(spl)-2]
	} else if len(spl) < 2 && b {
		folder = ""
	}

	return folder
}
