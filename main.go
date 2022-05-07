package main

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	eng := New(
		clientset,
		"1-7 09:00-12:20")
	ctx := context.Background()
	go eng.Watcher(ctx)
	go eng.Suspender(ctx)
	select {}
}
