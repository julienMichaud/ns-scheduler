package main

import (
	"context"
	"fmt"
	"os"

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
	if err != nil {
		fmt.Println(err)
	}
	var log = logrus.New()
	upTimeSchedule := "1-7 08:00-20:00"

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "DEBUG" {
		log.Level = logrus.DebugLevel
	}

	upTimeScheduleFromEnv := os.Getenv("UPTIME_SCHEDULE")
	if upTimeScheduleFromEnv != "" {
		log.Debugf("got env variable upTimeScheduleFromEnv with value %s, using it instead of default", upTimeScheduleFromEnv)
		upTimeSchedule = upTimeScheduleFromEnv

	}

	eng := New(
		clientset,
		upTimeSchedule,
		*log)

	ctx := context.Background()

	go eng.Watcher(ctx)
	go eng.Suspender(ctx)
	select {}
}
