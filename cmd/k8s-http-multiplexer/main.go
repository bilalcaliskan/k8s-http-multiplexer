package main

import (
	"github.com/dimiro1/banner"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s-http-multiplexer/pkg/cfg"
	"k8s-http-multiplexer/pkg/k8s"
	"k8s-http-multiplexer/pkg/logging"
	"k8s-http-multiplexer/pkg/metrics"
	"k8s-http-multiplexer/pkg/options"
	"k8s-http-multiplexer/pkg/web"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
)

var (
	clientSet *kubernetes.Clientset
	khmo      *options.K8sHttpMultiplexerOptions
	logger    *zap.Logger
	router    *mux.Router
)

func init() {
	khmo = options.GetK8sHttpMultiplexerOptions()
	logger = logging.GetLogger()
	router = mux.NewRouter()
	bannerBytes, _ := ioutil.ReadFile("banner.txt")
	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
}

func main() {
	cfg.ParseConfig(khmo.ConfigFilePath)

	logger.Info("initializing kube client")
	restConfig, err := k8s.GetConfig(cfg.Cfg.MasterUrl, khmo.KubeConfigPath, cfg.Cfg.InCluster)
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
