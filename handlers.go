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

	found, configRequest := config.getRequest(r.Method, r.RequestURI)
	if !found {
		return
	}

	pods := getPods(targetPods, configRequest.Label)
	for _, pod := range pods {
		url := fmt.Sprintf("%s%s", pod.addr, configRequest.URI)
		logger.Info("making request on each pod", zap.String("url", url), zap.String("label", pod.label))
		remoteRequest, err := http.NewRequest("GET", fmt.Sprintf("%s%s", pod.addr, configRequest.URI), nil)
		if err != nil {
			logger.Error("an error occured while creating the request", zap.Error(err))
			return
		}

		for _, header := range configRequest.Headers {
			logger.Info("setting header on the remote request", zap.Any("header", header))
			remoteRequest.Header.Set(header.Key, header.Value)
		}

		if configRequest.BasicAuth {
			logger.Info("setting basic auth on remote request", zap.String("username", configRequest.Username),
				zap.String("password", configRequest.Password))
			remoteRequest.SetBasicAuth(configRequest.Username, configRequest.Password)
		}

		res, err := client.Do(remoteRequest)
		if err != nil {
			logger.Error("an error occured while making the request", zap.Error(err))
			return
		}

		if res != nil && res.StatusCode == configRequest.ExpectedResponseCode {
			successCount++
		}

		var bodyBytes []byte
		if bodyBytes, err = ioutil.ReadAll(res.Body); err != nil {
			logger.Error("an error occured while reading response body", zap.Error(err))
			return
		}

		if err = res.Body.Close(); err != nil {
			logger.Error("an error occured while closing response body", zap.Error(err))
			return
		}

		bodyString := string(bodyBytes)
		responseBody += bodyString
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if !configRequest.ReturnResponseBody {
		response := Response{
			TargetCount:  len(pods),
			SuccessCount: successCount,
		}

		var responseBytes []byte
		if responseBytes, err = json.Marshal(response); err != nil {
			logger.Error("an error occured while marshaling response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(responseBytes); err != nil {
			logger.Error("an error occured while writing response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err = w.Write([]byte(responseBody)); err != nil {
			logger.Error("an error occured while writing response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
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