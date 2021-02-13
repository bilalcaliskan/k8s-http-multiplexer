package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	var successCount int
	var responseBody string

	found, req := config.getRequest(r.Method, r.RequestURI)
	if !found {
		return
	}

	logger.Info("request received", zap.String("method", req.Method), zap.String("uri", req.URI),
		zap.String("label", req.Label))

	pods := getPods(targetPods, req.Label)
	for _, pod := range pods {
		url := fmt.Sprintf("%s%s", pod.addr, req.URI)
		logger.Info("making request on each pod", zap.String("url", url))
		request, err := http.NewRequest("GET", fmt.Sprintf("%s%s", pod.addr, req.URI), nil)
		if err != nil {
			panic(err)
		}

		logger.Info("setting headers on the remote request", zap.Any("headers", req.Headers))
		for _, header := range req.Headers {
			request.Header.Set(header.Key, header.Value)
		}

		res, err := client.Do(request)
		if err != nil {
			logger.Error("an error occured while making the request", zap.Error(err))
		}

		if res != nil && res.StatusCode == req.ExpectedResponseCode {
			successCount++
		}

		defer func() {
			err = res.Body.Close()
		}()
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Fatal("an error occured while reading response body", zap.Error(err))
		}
		bodyString := string(bodyBytes)
		responseBody += bodyString
	}

	logger.Info("", zap.Bool("returnResponseBody", req.ReturnResponseBody))
	if !req.ReturnResponseBody {
		response := Response{
			TargetCount:  len(pods),
			SuccessCount: successCount,
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			logger.Fatal("an error occured while marshaling response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseBytes)
		if err != nil {
			logger.Fatal("an error occured while writing response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		_, err = w.Write([]byte(responseBody))
		if err != nil {
			logger.Fatal("an error occured while writing response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	found, req := config.getRequest(r.Method, r.RequestURI)
	if found {
		logger.Info("request received", zap.String("method", req.Method), zap.String("uri", req.URI),
			zap.String("label", req.Label))
	}

	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "hello from postHandler\n")
	if err != nil {
		panic(err)
	}
}

func registerHandlers(router *mux.Router, config Config) {
	var count int
	for _, v := range config.Requests {
		if v.Method == "GET" {
			router.HandleFunc(v.URI, getHandler).Methods("GET").Schemes("http").Name("get")
			count++
		} else if v.Method == "POST" {
			router.HandleFunc(v.URI, postHandler).Methods("POST").Schemes("http").Name("post")
			count++
		}
	}
	logger.Info("handlers registered", zap.Int("count", count))
}