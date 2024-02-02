package config

import (
	"fmt"
	"os"

	"github.com/felipeospina21/gql-tui-client/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Queries struct {
		Location string `yaml:"location"`
	}
}

// Reads configuration file from default location
func Read() Config {
	// TODO: Make configPath as cmd arg to override default location
	home := os.Getenv("HOME")
	configPath := fmt.Sprintf("%s/.config/goql/goql.yaml", home)

	data, err := os.ReadFile(configPath)
	utils.CheckError(err)

	var config Config

	err = yaml.Unmarshal(data, &config)
	utils.CheckError(err)

	return config
}
