package configuration

import (
	"io/ioutil"
	"k8s-http-multiplexer/internal/logging"
	"k8s-http-multiplexer/internal/options"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	// Cfg is the representation of parsed Config
	cfg    Config
	logger *zap.Logger
	khmo   *options.K8sHttpMultiplexerOptions
)

func init() {
	logger = logging.GetLogger()
	khmo = options.GetK8sHttpMultiplexerOptions()
	ParseConfig(khmo.ConfigFilePath)
}

// ParseConfig gets the file path of config file in yaml format and parses it to Config
func ParseConfig(configFilePath string) {
	filename, _ := filepath.Abs(configFilePath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	logger.Info("successfully parsed configuration file", zap.Int("request_count", len(cfg.Requests)))
}

// GetConfig returns the parsed yaml config
func GetConfig() Config {
	return cfg
}
