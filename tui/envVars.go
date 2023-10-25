package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/felipeospina21/gql-tui-client/utils"
	"github.com/joho/godotenv"
)

type (
	isEditingEnvVars = bool
	envVars          = map[string]string
)

type envVarModel struct {
	textarea  textarea.Model
	err       error
	isEditing isEditingEnvVars
}

func (m *mainModel) newEnvVarModel(charLimit int, content string) {
	m.envVars.textarea = textarea.New()
	m.envVars.textarea.CharLimit = charLimit
	m.envVars.textarea.SetValue(content)
}

func parseEnvVars(s string) envVars {
	ev := make(envVars)
	var start, end int

	for idx, char := range s {
		if string(char) == "\n" {
			end = idx
			line := s[start:end]
			start = idx + 1
			if !utils.IsStringEmpty(line) {
				ev = setKeyVal(line, ev)
			}
		}

		if idx == len(s)-1 {
			lastLine := s[start:]
			if !utils.IsStringEmpty(lastLine) {
				ev = setKeyVal(lastLine, ev)
			}
		}
	}
	return ev
}

func stringifyEnvVars(ev envVars) string {
	var s string
	for key, val := range ev {
		s += fmt.Sprintf("%s=%s\n", key, val)
	}
	return s
}

func readEnvVars() envVars {
	var m envVars
	m, err := godotenv.Read()
	utils.CheckError(err)
	return m
}

func setKeyVal(s string, ev envVars) envVars {
	v := strings.Split(s, "=")
	key := v[0]
	value := v[1]
	ev[key] = value
	return ev
}
