package configuration

import (
	"io/ioutil"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	// Cfg is the representation of parsed Config
	Cfg    Config
	logger *zap.Logger
	err    error
)

func init() {
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

// ParseConfig gets the file path of config file in yaml format and parses it to Config
func ParseConfig(configFilePath string) {
	filename, _ := filepath.Abs(configFilePath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &Cfg)
	if err != nil {
		panic(err)
	}

	logger.Info("successfully parsed configuration file", zap.Int("request_count", len(Cfg.Requests)))
}
