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
	config configuration.Config
)

func init() {
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	client = &http.Client{}
	config = configuration.GetConfig()
}

// RunWebServer spins up webserver to handle incoming HTTP requests
func RunWebServer() error {
	router := mux.NewRouter()
	registerHandlers(router)
	webServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", config.Port),
		WriteTimeout: time.Duration(int32(config.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(config.ReadTimeoutSeconds)) * time.Second,
	}
	logger.Info("web server is up and running", zap.Int("port", config.Port))
	return webServer.ListenAndServe()
}
