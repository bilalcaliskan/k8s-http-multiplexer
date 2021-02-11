package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
	"time"
)

var (
	// targetPods []*TargetPod
	clientSet *kubernetes.Clientset
	// client *http.Client
	kubeConfigPath *string
	labels []string
	logger *zap.Logger
	err error
)

func init() {
	kubeConfigPath = flag.String("kubeConfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		"absolute path of the kubeconfig file, required when non inCluster environment")
	flag.Parse()

	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func main() {
	filename, _ := filepath.Abs("config/sample.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	logger.Info("successfully parsed config file", zap.Int("request_count", len(config.Requests)))

	/*log.Println("Initializing http client...")
	client = &http.Client{}*/

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

	for _, v := range config.Requests {
		labels = append(labels, v.Label)
	}
	logger.Info("initialized labels slice", zap.Any("labels", labels))
	runPodInformer(clientSet, labels, logger)

	router := mux.NewRouter()
	server := initServer(router, config, fmt.Sprintf(":%d", config.Port), time.Duration(int32(config.WriteTimeoutSeconds)),
		time.Duration(int32(config.ReadTimeoutSeconds)))
	logger.Info("server is up and running", zap.Int("port", config.Port))
	panic(server.ListenAndServe())
}