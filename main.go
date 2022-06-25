package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

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
	checkInterval := 30 // 30 seconds by default

	logLevel := os.Getenv("NS_SCHEDULER_LOG_LEVEL")
	if logLevel == "DEBUG" {
		log.Level = logrus.DebugLevel
	}

	upTimeScheduleFromEnv := os.Getenv("NS_SCHEDULER_UPTIME_SCHEDULE")
	if upTimeScheduleFromEnv != "" {
		log.Debugf("got env variable upTimeScheduleFromEnv with value %s, using it instead of default", upTimeScheduleFromEnv)
		upTimeSchedule = upTimeScheduleFromEnv
	}

	checkIntervalFromEnv := os.Getenv("NS_SCHEDULER_CHECK_INTERVAL")
	if checkIntervalFromEnv != "" {
		log.Debugf("got env variable CheckIntervalFromEnv with value %s, using it instead of default", checkIntervalFromEnv)
		i, _ := strconv.Atoi(checkIntervalFromEnv)
		checkInterval = i
	}

	eng := New(
		clientset,
		upTimeSchedule,
		*log,
		checkInterval)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		err := eng.Watcher(ctx, wg)
		if err != nil {
			cancel()
		}
	}()

	go func() {
		err := eng.Suspender(ctx, wg)
		if err != nil {
			cancel()
		}
	}()

	<-done // wait for SIGINT / SIGTERM
	log.Print("receive shutdown")
	cancel()
	wg.Wait()

	log.Print("controller exited properly")

}
