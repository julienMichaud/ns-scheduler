package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (eng *Engine) Suspender(ctx context.Context) {

	contextLogger := eng.logger.WithFields(log.Fields{
		"go-routine":      "Suspender",
		"uptime-schedule": eng.upTimeSchedule,
	})
	contextLogger.Info("starting Suspender goroutine")

	for {
		n := <-eng.Wl
		contextLogger.Infof("received namespace: %s", n.ObjectMeta.Name)
		dState := n.Annotations["ns-scheduler/state"]

		switch dState {
		// if no annotation, we just set it
		case "": //put everything under into a func
			contextLogger.Infof("didnt see annotation ns-scheduler-state on the namespace %s, first time im seeing it then", n.ObjectMeta.Name)
			if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				res, err := eng.client.CoreV1().Namespaces().Get(ctx, n.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				res.Annotations["ns-scheduler/state"] = "Running"

				_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
				return err
			}); err != nil {
				contextLogger.Error("cannot update namespace object")
				continue
			}
		//if namespace annotation is Running, check if the namespace resources should be up based on the upTime variable, if not we scale down resources
		case "Running":

			contextLogger.Infof("namespace %s in running state, checking if it should be scale downed", n.ObjectMeta.Name)
			scaledown, err := shouldScaleDown(eng.upTimeSchedule, &eng.logger)
			if err != nil {
				contextLogger.Errorf("could not determine if we should scale down: %v", err)
			}
			if scaledown {
				contextLogger.Infof("namespace %s not in the upTime range specified, which mean that it should be scaled down.\n The state of the namespace will then be Suspended.", n.ObjectMeta.Name)
				if err := shuttingDownNamespace(eng, n, ctx); err != nil {
					contextLogger.Errorf("could not shut namespace, error is %s", err)

				}
			} else {
				contextLogger.Infof("namespace %s is in the upTime range specified, not doing anything.", n.ObjectMeta.Name)
			}

		case "Suspended":

			contextLogger.Infof("namespace %s already suspended, checking if it should be revived", n.ObjectMeta.Name)
			scaledown, err := shouldScaleDown(eng.upTimeSchedule, &eng.logger)
			if err != nil {
				contextLogger.Errorf("could not determine if we should scale down: %v", err)
			}
			if scaledown {
				contextLogger.Infof("namespace %s still not in range, not doing anything.", n.ObjectMeta.Name)
			} else {
				contextLogger.Infof("namespace %s is in range, should be revived !.", n.ObjectMeta.Name)
				if err := startingUpNamespace(eng, n, ctx); err != nil {
					contextLogger.Errorf("could not start namespace, error is %s", err)
				}
			}
		}

	}
}
