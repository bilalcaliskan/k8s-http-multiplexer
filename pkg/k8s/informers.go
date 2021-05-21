package k8s

import (
	"fmt"
	"go.uber.org/zap"
	"k8s-http-multiplexer/pkg/cfg"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

// TargetPods keeps the pointer of TargetPod items. It is the representation of target ip addresses, ports information
var TargetPods []*TargetPod

// RunPodInformer runs the shared informer to watch Add/Update/Delete pod events
func RunPodInformer(clientSet *kubernetes.Clientset, logger *zap.Logger) {
	// TODO: Run each logic as separate goroutine, use channels
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labelMap := pod.GetLabels()
			for _, request := range cfg.Cfg.Requests {
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
						Addr:  addr,
						Label: request.Label,
					}

					logger.Info("adding pod to the targetPods", zap.String("targetPod.addr", targetPod.Addr),
						zap.String("targetPod.label", targetPod.Label))
					addTargetPod(&TargetPods, &targetPod)
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

			for _, request := range cfg.Cfg.Requests {
				if labelExists(oldLabelMap, request.Label) && labelExists(newLabelMap, request.Label) {
					if oldPod.Status.PodIP == "" && newPod.Status.PodIP != "" {
						logger.Info("assigned an ip address to the pod", zap.String("addr", newPod.Status.PodIP))

						containerPort := newPod.Spec.Containers[0].Ports[0].ContainerPort
						if request.TargetPort != 0 {
							containerPort = request.TargetPort
						}

						addr := fmt.Sprintf("http://%s:%d", newPod.Status.PodIP, containerPort)
						targetPod := TargetPod{
							Addr:  addr,
							Label: request.Label,
						}

						logger.Info("adding pod to the targetPods", zap.String("targetPod.addr", targetPod.Addr),
							zap.String("targetPod.label", targetPod.Label))
						addTargetPod(&TargetPods, &targetPod)
					}

				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			labelMap := pod.GetLabels()
			for _, request := range cfg.Cfg.Requests {
				if labelExists(labelMap, request.Label) {
					containerPort := pod.Spec.Containers[0].Ports[0].ContainerPort
					if request.TargetPort != 0 {
						containerPort = request.TargetPort
					}

					addr := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, containerPort)
					targetPod := TargetPod{
						Addr:  addr,
						Label: request.Label,
					}

					if index, found := findTargetPod(TargetPods, targetPod); found {
						logger.Info("pod found in the targetPods, removing", zap.String("addr", targetPod.Addr),
							zap.String("label", targetPod.Label))
						removeTargetPod(&TargetPods, index)
						logger.Info("pod successfully removed from targetPods", zap.Any("targetPodsLength",
							len(TargetPods)))
					}

				}
			}
		},
	})

	logger.Info("starting informer factory")
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}
