package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	var log = logrus.New()
	// log.Level = logrus.DebugLevel

	eng := New(
		clientset,
		"1-7 09:00-12:18",
		*log)
	ctx := context.Background()
	go eng.Watcher(ctx)
	go eng.Suspender(ctx)
	select {}
}
