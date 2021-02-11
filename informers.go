package main

import (
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

func runPodInformer(clientSet *kubernetes.Clientset, labels []string, logger *zap.Logger) {
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second * 30)
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// 1- Pod labelled
			// 2- Pod is not labelled

			pod := obj.(*v1.Pod)
			labelMap := pod.GetLabels()
			for _, label := range labels {
				if labelExists(labelMap, label) {
					// 1- Add pod ip to the targetPods slice

					logger.Info("label found in the labelMap", zap.String("label", label),
						zap.Any("labelMap", labelMap))
				} else {
					// 1- Do nothing

					// logger.Info("label not found in the labelMap, skipping", zap.String("label", label),
					//	zap.Any("labelMap", labelMap))
				}
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// 1- Old pod labelled, new pod labelled
			// 2- Old pod not labelled, new pod labelled
			// 3- Old pod labelled, new pod not labelled
			// 4- Old pod not labelled, new pod not labelled

			oldPod := oldObj.(*v1.Pod)
			oldLabelMap := oldPod.GetLabels()
			newPod := newObj.(*v1.Pod)
			newLabelMap := newPod.GetLabels()

			if oldPod.ResourceVersion != newPod.ResourceVersion {
				for _, label := range labels {
					if labelExists(oldLabelMap, label) && labelExists(newLabelMap, label) {
						// 1- Check ip addresses of oldPod and newPod, update targetPods slice if neccessary

						logger.Info("old pod is labelled and new pod labelled", zap.String("oldPodIp", oldPod.Status.PodIP),
							zap.String("newPodIp", newPod.Status.PodIP))
					} else if !labelExists(oldLabelMap, label) && labelExists(newLabelMap, label) {
						// 1- Add newPod ip to targetPods slice

						logger.Info("old pod is not labelled and new pod labelled", zap.String("oldPodIp",
							oldPod.Status.PodIP), zap.String("newPodIp", newPod.Status.PodIP))
					} else if labelExists(oldLabelMap, label) && !labelExists(newLabelMap, label) {
						// 1- Remove oldPod ip from targetPods slice, do nothing for newPod

						logger.Info("old pod is labelled and new pod is not labelled", zap.String("oldPodIp",
							oldPod.Status.PodIP), zap.String("newPodIp", newPod.Status.PodIP))
					} else if !labelExists(oldLabelMap, label) && !labelExists(newLabelMap, label) {
						// 1- Do nothing

						// logger.Info("old pod is not labelled and new pod is not labelled, there is nothing to do",
						// 	zap.String("oldPodIp", oldPod.Status.PodIP), zap.String("newPodIp", newPod.Status.PodIP))
					}
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			// 1- Pod is labelled
			// 2- Pod is not labelled

			pod := obj.(*v1.Pod)
			labelMap := pod.GetLabels()
			for _, label := range labels {
				if labelExists(labelMap, label) {
					// 1- Remove pod ip from targetPods slice

				} else if !labelExists(labelMap, label) {
					// 1- Do nothing

				}
			}
		},
	})

	logger.Info("starting informer factory")
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}