package main

import (
	"context"
	"fmt"
	"ns-scheduler/engine"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	eng := engine.New(clientset)
	ctx := context.Background()
	go eng.Watcher(ctx)
	go eng.Suspender(ctx)
	select {}
}
