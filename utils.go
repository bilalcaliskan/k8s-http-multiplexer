package main

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"strings"
	"time"
)

func initServer(router *mux.Router, config Config, addr string, writeTimeout time.Duration, readTimeout time.Duration) *http.Server {
	registerHandlers(router, config)
	return &http.Server{
		Handler: router,
		Addr: addr,
		WriteTimeout: writeTimeout,
		ReadTimeout: readTimeout,
	}
}

func getConfig(masterUrl, kubeConfigPath string, inCluster bool) (*rest.Config, error) {
	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags(masterUrl, kubeConfigPath)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func getClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func labelExists(labelMap map[string]string, label string) bool {
	splittedLabel := strings.Split(label, "=")
	labelKey := splittedLabel[0]
	labelValue := splittedLabel[1]
	for key, value := range labelMap {
		if labelKey == key && labelValue == value {
			return true
		}
	}
	return false
}

func findTargetPod(targetPods []*TargetPod, pod TargetPod) (int, bool) {
	for i, item := range targetPods {
		if pod.Equals(item) {
			return i, true
		}
	}
	return -1, false
}

func addTargetPod(targetPods *[]*TargetPod, pod *TargetPod) {
	_, found := findTargetPod(*targetPods, *pod)
	if !found {
		*targetPods = append(*targetPods, pod)
	}
}

func removeTargetPod(targetPods *[]*TargetPod, index int) {
	*targetPods = append((*targetPods)[:index], (*targetPods)[index+1:]...)
}