package metrics

import (
	"fmt"
	config2 "k8s-http-multiplexer/internal/configuration"
	"k8s-http-multiplexer/internal/logging"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	config config2.Config
)

func init() {
	logger = logging.GetLogger()
	config = config2.GetConfig()
}

// RunMetricsServer exports metrics
func RunMetricsServer(router *mux.Router) error {
	metricServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", config.MetricsPort),
		WriteTimeout: time.Duration(int32(config.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(config.ReadTimeoutSeconds)) * time.Second,
	}
	router.Handle(config.MetricsUri, promhttp.Handler())
	logger.Info("metric server is up and running", zap.Int("port", config.MetricsPort))
	return metricServer.ListenAndServe()
}
