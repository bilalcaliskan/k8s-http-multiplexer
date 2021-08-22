package web

import (
	"fmt"
	"k8s-http-multiplexer/internal/configuration"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
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
		Addr:         fmt.Sprintf(":%d", configuration.Cfg.Port),
		WriteTimeout: time.Duration(int32(configuration.Cfg.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(configuration.Cfg.ReadTimeoutSeconds)) * time.Second,
	}
	logger.Info("web server is up and running", zap.Int("port", configuration.Cfg.Port))
	panic(webServer.ListenAndServe())
}
