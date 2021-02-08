package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func initServer(router *mux.Router, addr string, writeTimeout time.Duration, readTimeout time.Duration) *http.Server {
	registerHandlers(router)
	return &http.Server{
		Handler: router,
		Addr: addr,
		WriteTimeout: writeTimeout,
		ReadTimeout: readTimeout,
	}
}

func registerHandlers(router *mux.Router) {

}