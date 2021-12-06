package main

import (
	"io/ioutil"
	"k8s-http-multiplexer/internal/configuration"
	"k8s-http-multiplexer/internal/k8s"
	"k8s-http-multiplexer/internal/logging"
	"k8s-http-multiplexer/internal/metrics"
	"k8s-http-multiplexer/internal/options"
	"k8s-http-multiplexer/internal/web"
	"k8s.io/client-go/rest"
	"os"
	"strings"

	"github.com/dimiro1/banner"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

var (
	restConfig *rest.Config
	clientSet  *kubernetes.Clientset
	err        error
	khmo       *options.K8sHttpMultiplexerOptions
	logger     *zap.Logger
	config     configuration.Config
)

func init() {
	khmo = options.GetK8sHttpMultiplexerOptions()
	logger = logging.GetLogger()
	config = configuration.GetConfig()
	bannerBytes, _ := ioutil.ReadFile("banner.txt")
	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
}

func main() {
	logger.Info("initializing kube client")
	if restConfig, err = k8s.GetConfig(config.MasterUrl, khmo.KubeConfigPath, khmo.InCluster); err != nil {
		panic(err)
	}

	logger.Info("initializing client set")
	if clientSet, err = k8s.GetClientSet(restConfig); err != nil {
		logger.Fatal("fatal error occured while getting clientset", zap.Error(err))
	}

	logger.Info("successfully initialized kube client")

	go k8s.RunPodInformer(clientSet)
	go func() {
		if err = metrics.RunMetricsServer(); err != nil {
			logger.Fatal("fatal error occured while spinning metric server", zap.Error(err))
		}
	}()

	if err = web.RunWebServer(); err != nil {
		logger.Fatal("fatal error occured while spinning web server", zap.Error(err))
	}
}
