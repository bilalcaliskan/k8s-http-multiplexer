package main

import (
	"flag"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	targetPods []*TargetPod
	clientSet *kubernetes.Clientset
	client *http.Client
	kubeConfigPath, configFilePath *string
	logger *zap.Logger
	err error
	config Config
	router *mux.Router
)

func init() {
	kubeConfigPath = flag.String("kubeConfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		"absolute path of the kubeconfig file, required when non inCluster environment")
	configFilePath = flag.String("configFilePath", "config/sample.yaml", "path of the configuration file")
	flag.Parse()

	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	router = mux.NewRouter()
}

func main() {
	filename, _ := filepath.Abs(*configFilePath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	logger.Info("successfully parsed config file", zap.Int("request_count", len(config.Requests)))

	log.Println("Initializing http client...")
	client = &http.Client{}

	logger.Info("initializing kube client")
	restConfig, err := getConfig(config.MasterUrl, *kubeConfigPath, config.InCluster)
	if err != nil {
		panic(err)
	}
	clientSet, err = getClientSet(restConfig)
	if err != nil {
		panic(err)
	}
	logger.Info("successfully initialized kube client")

	go runPodInformer(clientSet, config, logger)
	go runMetricsServer(router)
	runWebServer(router)
}