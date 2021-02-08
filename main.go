package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	targetPods []*TargetPod
	clientSet *kubernetes.Clientset
	client *http.Client
	kubeConfigPath *string
	labels []string
)

func init() {
	kubeConfigPath = flag.String("kubeConfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		"absolute path of the kubeconfig file, required when non inCluster environment")
	flag.Parse()
}

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


	log.Println("Initializing http client...")
	client = &http.Client{}

	log.Println("Initializing kube client...")
	restConfig, err := getConfig(config.MasterUrl, *kubeConfigPath, config.InCluster)
	if err != nil {
		panic(err)
	}
	clientSet, err = getClientSet(restConfig)
	if err != nil {
		panic(err)
	}

	for _, v := range config.Requests {
		labels = append(labels, v.Label)
	}
	log.Printf("final labels slice before running pod informer = %v\n", labels)
	runPodInformer(clientSet, labels)

	router := mux.NewRouter()
	server := initServer(router, config, fmt.Sprintf(":%d", config.Port), time.Duration(int32(config.WriteTimeoutSeconds)),
		time.Duration(int32(config.ReadTimeoutSeconds)))
	log.Printf("Server is listening on port %d!", config.Port)
	log.Fatal(server.ListenAndServe())
}