package main

import "github.com/felipeospina21/gql-tui-client/tui"

func main() {
	tui.Start()
}

// func getQueriesList(rootDir string) tea.Cmd {
// 	return func() tea.Msg {
// 		// TODO: replace Walk with WalkDir func
//
// 		queriesNames := []list.Item{}
// 		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
// 			if err != nil {
// 				fmt.Println(err)
// 				return err
// 			}
// 			if !info.IsDir() {
// 				queriesNames = append(queriesNames, item{title: info.Name(), desc: "query"})
// 			}
// 			return nil
// 		})
//
// 		utils.CheckError(err)
// 		// TODO: remove time sleep
// 		time.Sleep(2 * time.Second)
// 		return listItems(queriesNames)
// 	}
// }

// func gqlReq(url string, file string) tea.Cmd {
// 	return func() tea.Msg {
// 		b, err := os.ReadFile(file)
// 		utils.CheckError(err)
//
// 		var obj map[string]interface{}
// 		client := graphql.NewClient(url)
// 		req := graphql.NewRequest(string(b))
// 		err = client.Run(context.Background(), req, &obj)
// 		utils.CheckError(err)
//
// 		f := colorjson.NewFormatter()
// 		f.Indent = 2
// 		f.KeyColor = color.New(color.FgMagenta)
//
// 		s, err := f.Marshal(obj)
// 		utils.CheckError(err)
//
// 		return responseMsg(s)
// 	}
// }

// func utils.CheckError(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }
