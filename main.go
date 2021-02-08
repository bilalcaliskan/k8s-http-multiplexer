package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func main() {
	filename, _ := filepath.Abs("config/sample.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	log.Printf("total %d requests found in the config!\n", len(config.Requests))

	router := mux.NewRouter()
	server := initServer(router, fmt.Sprintf(":%d", config.Port), time.Duration(int32(config.WriteTimeoutSeconds)),
		time.Duration(int32(config.ReadTimeoutSeconds)))
	log.Printf("Server is listening on port %d!", config.Port)
	log.Fatal(server.ListenAndServe())
}