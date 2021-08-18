package options

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var khmo = &K8sHttpMultiplexerOptions{}

func init() {
	khmo.addFlags(pflag.CommandLine)
	pflag.Parse()
}

// GetK8sHttpMultiplexerOptions returns the pointer of K8sHttpMultiplexerOptions
func GetK8sHttpMultiplexerOptions() *K8sHttpMultiplexerOptions {
	return khmo
}

// K8sHttpMultiplexerOptions contains frequent command line and application options.
type K8sHttpMultiplexerOptions struct {
	// KubeConfigPath is the path of the kubeconfig file to access the cluster
	KubeConfigPath string
	// ConfigFilePath is the path of the application to properly run
	ConfigFilePath string
}

func (khmo *K8sHttpMultiplexerOptions) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&khmo.KubeConfigPath, "kubeConfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		"absolute path of the kubeconfig file, required when non inCluster environment")
	fs.StringVar(&khmo.ConfigFilePath, "configFilePath", "config/sample.yaml",
		"path of the configuration file")
}
