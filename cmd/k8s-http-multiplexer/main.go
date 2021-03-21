package main

import (
	"flag"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"k8s-http-multiplexer/pkg/cfg"
	"k8s-http-multiplexer/pkg/k8s"
	"k8s-http-multiplexer/pkg/metrics"
	"k8s-http-multiplexer/pkg/web"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
)

var (
	clientSet *kubernetes.Clientset
	kubeConfigPath, configFilePath string
	logger *zap.Logger
	err error
	router *mux.Router
)

func init() {
	flag.StringVar(&kubeConfigPath, "kubeConfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		"absolute path of the kubeconfig file, required when non inCluster environment")
	flag.StringVar(&configFilePath, "configFilePath", "config/sample.yaml", "path of the configuration file")
	flag.Parse()

	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	router = mux.NewRouter()
}

func main() {
	cfg.ParseConfig(configFilePath)

	logger.Info("initializing kube client")
	restConfig, err := k8s.GetConfig(cfg.Cfg.MasterUrl, kubeConfigPath, cfg.Cfg.InCluster)
	if err != nil {
		panic(err)
	}
	clientSet, err = k8s.GetClientSet(restConfig)
	if err != nil {
		panic(err)
	}
	logger.Info("successfully initialized kube client")

	go k8s.RunPodInformer(clientSet, logger)
	go metrics.RunMetricsServer(router, logger)
	web.RunWebServer(router)
}