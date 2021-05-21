package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"k8s-http-multiplexer/pkg/cfg"
	"net/http"
	"time"
)

var (
	client *http.Client
	logger *zap.Logger
	err    error
)

func init() {
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	client = &http.Client{}
}

// RunWebServer spins up webserver to handle incoming HTTP requests
func RunWebServer(router *mux.Router) {
	registerHandlers(router)
	webServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", cfg.Cfg.Port),
		WriteTimeout: time.Duration(int32(cfg.Cfg.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(cfg.Cfg.ReadTimeoutSeconds)) * time.Second,
	}
	logger.Info("web server is up and running", zap.Int("port", cfg.Cfg.Port))
	panic(webServer.ListenAndServe())
}
