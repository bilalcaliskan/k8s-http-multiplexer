package main

import (
	"fmt"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

// TODO: Run each logic as seperate goroutine, use channels
func runPodInformer(clientSet *kubernetes.Clientset, config Config, logger *zap.Logger) {
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second * 30)
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labelMap := pod.GetLabels()
			for _, request := range config.Requests {
				if labelExists(labelMap, request.Label) && pod.Status.PodIP != "" {
					logger.Info("label found in the labelMap", zap.String("label", request.Label),
						zap.Any("labelMap", labelMap))

					containerPort := pod.Spec.Containers[0].Ports[0].ContainerPort
					if request.TargetPort != 0 {
						containerPort = request.TargetPort
					}

					// TODO: Uncomment for out of cluster
					addr := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, containerPort)
					// addr := fmt.Sprintf("http://%s:%d", "192.168.99.114", containerPort)

					targetPod := TargetPod{
						addr:  addr,
						label: request.Label,
					}

					logger.Info("adding pod to the targetPods", zap.String("targetPod.addr", targetPod.addr),
						zap.String("targetPod.label", targetPod.label))
					addTargetPod(&targetPods, &targetPod)
				} else if labelExists(labelMap, request.Label) && pod.Status.PodIP == "" {
					logger.Info("label found in the labelMap, but pod still does not have an ip address, skipping",
						zap.String("label", request.Label), zap.Any("labelMap", labelMap))
				}
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*v1.Pod)
			oldLabelMap := oldPod.GetLabels()
			newPod := newObj.(*v1.Pod)
			newLabelMap := newPod.GetLabels()

			if oldPod.ResourceVersion == newPod.ResourceVersion {
				return
			}

			for _, request := range config.Requests {
				if labelExists(oldLabelMap, request.Label) && labelExists(newLabelMap, request.Label) {
					if oldPod.Status.PodIP == "" && newPod.Status.PodIP != "" {
						logger.Info("assigned an ip address to the pod", zap.String("addr", newPod.Status.PodIP))

						containerPort := newPod.Spec.Containers[0].Ports[0].ContainerPort
						if request.TargetPort != 0 {
							containerPort = request.TargetPort
						}

						addr := fmt.Sprintf("http://%s:%d", newPod.Status.PodIP, containerPort)
						targetPod := TargetPod{
							addr:  addr,
							label: request.Label,
						}

						logger.Info("adding pod to the targetPods", zap.String("targetPod.addr", targetPod.addr),
							zap.String("targetPod.label", targetPod.label))
						addTargetPod(&targetPods, &targetPod)
					}

				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labelMap := pod.GetLabels()
			for _, request := range config.Requests {
				if labelExists(labelMap, request.Label) {
					containerPort := pod.Spec.Containers[0].Ports[0].ContainerPort
					if request.TargetPort != 0 {
						containerPort = request.TargetPort
					}

					addr := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, containerPort)
					targetPod := TargetPod{
						addr:  addr,
						label: request.Label,
					}

					if index, found := findTargetPod(targetPods, targetPod); found {
						logger.Info("pod found in the targetPods, removing", zap.String("addr", targetPod.addr),
							zap.String("label", targetPod.label))
						removeTargetPod(&targetPods, index)
						logger.Info("pod successfully removed from targetPods", zap.Any("targetPodsLength",
							len(targetPods)))
					}

				}
			}
		},
	})

	logger.Info("starting informer factory")
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}