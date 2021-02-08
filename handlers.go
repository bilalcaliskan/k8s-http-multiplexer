package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func getHandler(w http.ResponseWriter, r *http.Request, request Request) {

}

func postHandler(w http.ResponseWriter, r *http.Request, request Request) {

}

func registerHandlers(router *mux.Router, config Config) {
	for _, v := range config.Requests {
		if v.Method == "GET" {
			router.HandleFunc(v.URI, func(w http.ResponseWriter, r *http.Request) {
				getHandler(w, r, v)
			}).Methods("GET").Schemes("http").Name("get")
		} else if v.Method == "POST" {
			router.HandleFunc(v.URI, func(w http.ResponseWriter, r *http.Request) {
				postHandler(w, r, v)
			}).Methods("POST").Schemes("http").Name("post")
		}
	}
}