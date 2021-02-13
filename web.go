package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func runWebServer(router *mux.Router) {
	registerHandlers(router, config)
	webServer := &http.Server{
		Handler: router,
		Addr: fmt.Sprintf(":%d", config.Port),
		WriteTimeout: time.Duration(int32(config.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(config.ReadTimeoutSeconds)) * time.Second,
	}
	logger.Info("web server is up and running", zap.Int("port", config.Port))
	panic(webServer.ListenAndServe())
}