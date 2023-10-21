package config

import (
	"fmt"
	"os"

	"github.com/felipeospina21/gql-tui-client/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Folder struct {
		Location string `yaml:"location"`
	}
}

func Read() Config {
	// TODO: Make this dynamic, path as cmd arg or default to this.
	home := os.Getenv("HOME")
	configPath := fmt.Sprintf("%s/.config/goql/goql.yaml", home)

	data, err := os.ReadFile(configPath)
	utils.CheckError(err)

	var config Config

	err = yaml.Unmarshal(data, &config)
	utils.CheckError(err)

	return config
}
