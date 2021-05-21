package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

// GetConfig gets required parameters and returns the rest.Config
func GetConfig(masterUrl, kubeConfigPath string, inCluster bool) (*rest.Config, error) {
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

// GetClientSet gets the rest.Config and returns the kubernetes.Clientset
func GetClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
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

// GetPods takes targetPods pointer slice and label as argument and returns a new slice with correct TargetPod instances
func GetPods(targetPods []*TargetPod, label string) []TargetPod {
	var pods []TargetPod
	for _, v := range targetPods {
		if v.Label == label {
			pods = append(pods, *v)
		}
	}
	return pods
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
