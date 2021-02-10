package main

import (
	"go.uber.org/zap"
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

		},
		UpdateFunc: func(oldObj, newObj interface{}) {

		},
		DeleteFunc: func(obj interface{}) {

		},
	})

	logger.Info("starting informer factory")
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}