package main

import (
	"io/ioutil"
	"k8s-http-multiplexer/internal/configuration"
	"k8s-http-multiplexer/internal/k8s"
	"k8s-http-multiplexer/internal/logging"
	"k8s-http-multiplexer/internal/metrics"
	"k8s-http-multiplexer/internal/options"
	"k8s-http-multiplexer/internal/web"
	"os"
	"strings"

	"github.com/dimiro1/banner"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

var (
	clientSet *kubernetes.Clientset
	khmo      *options.K8sHttpMultiplexerOptions
	logger    *zap.Logger
	router    *mux.Router
	config    configuration.Config
)

func init() {
	khmo = options.GetK8sHttpMultiplexerOptions()
	logger = logging.GetLogger()
	router = mux.NewRouter()
	config = configuration.GetConfig()
	bannerBytes, _ := ioutil.ReadFile("banner.txt")
	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
}

func main() {
	configuration.ParseConfig(khmo.ConfigFilePath)

	logger.Info("initializing kube client")
	restConfig, err := k8s.GetConfig(config.MasterUrl, khmo.KubeConfigPath, khmo.InCluster)
	if err != nil {
		panic(err)
	}
	clientSet, err = k8s.GetClientSet(restConfig)
	if err != nil {
		panic(err)
	}
	logger.Info("successfully initialized kube client")

	go k8s.RunPodInformer(clientSet)
	go func() {
		err := metrics.RunMetricsServer(router)
		if err != nil {
			panic(err)
		}
	}()

	if err := web.RunWebServer(router); err != nil {
		panic(err)
	}
}
