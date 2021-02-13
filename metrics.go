package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// TODO: Generate custom metrics, check below:
// https://prometheus.io/docs/guides/go-application/
// https://www.robustperception.io/prometheus-middleware-for-gorilla-mux

func runMetricsServer(router *mux.Router) {
	metricServer := &http.Server{
		Handler: router,
		Addr: fmt.Sprintf(":%d", config.MetricsPort),
		WriteTimeout: time.Duration(int32(config.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(config.ReadTimeoutSeconds)) * time.Second,
	}
	router.Handle(config.MetricsUri, promhttp.Handler())
	logger.Info("metric server is up and running", zap.Int("port", config.MetricsPort))
	panic(metricServer.ListenAndServe())
}