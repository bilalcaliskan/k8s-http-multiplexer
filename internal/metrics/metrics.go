package metrics

import (
	"fmt"
	config2 "k8s-http-multiplexer/internal/configuration"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// TODO: Generate custom metrics, check below:
// https://prometheus.io/docs/guides/go-application/
// https://www.robustperception.io/prometheus-middleware-for-gorilla-mux

// RunMetricsServer exports metrics
func RunMetricsServer(router *mux.Router, logger *zap.Logger) {
	metricServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", config2.Cfg.MetricsPort),
		WriteTimeout: time.Duration(int32(config2.Cfg.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(config2.Cfg.ReadTimeoutSeconds)) * time.Second,
	}
	router.Handle(config2.Cfg.MetricsUri, promhttp.Handler())
	logger.Info("metric server is up and running", zap.Int("port", config2.Cfg.MetricsPort))
	panic(metricServer.ListenAndServe())
}
