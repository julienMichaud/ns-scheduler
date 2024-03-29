package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (eng *Engine) Suspender(ctx context.Context, wg *sync.WaitGroup) error {

	contextLogger := eng.logger.WithFields(log.Fields{
		"go-routine": "Suspender",
	})
	contextLogger.Info("starting Suspender goroutine")
	now := time.Now().In(eng.loc)

	for {
		select {
		case n := <-eng.Wl:
			contextLogger.Infof("received namespace: %s", n.ObjectMeta.Name)
			dState := n.Annotations["ns-scheduler/state"]
			upTimeOnNamespace := n.Annotations["ns-scheduler/uptime"]

			var upTime string
			if upTimeOnNamespace != "" {
				upTime = upTimeOnNamespace
				contextLogger.Infof("found 'ns-scheduler/uptime' on namespace %s, will use it instead of default value", n.ObjectMeta.Name)
			} else {
				upTime = eng.upTimeSchedule
			}
			contextLogger := eng.logger.WithFields(log.Fields{
				"uptime-schedule": upTime,
			})

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
				}
			//if namespace annotation is Running, check if the namespace resources should be up based on the upTime variable, if not we scale down resources
			case "Running":

				contextLogger.Infof("namespace %s in running state, checking if it should be scale downed", n.ObjectMeta.Name)
				scaledown, err := shouldScaleDown(upTime, &eng.logger, now, eng.loc)
				if err != nil {
					contextLogger.Errorf("could not determine if we should scale down: %v", err)
				}
				if scaledown {
					contextLogger.Infof("namespace %s not in the upTime range specified, which mean that it should be scaled down.\n The state of the namespace will then be Suspended.", n.ObjectMeta.Name)
					if err := handlingNamespace(eng, n, ctx, false); err != nil {
						contextLogger.Errorf("error while scaling down namespace %s ressources, %s", n.ObjectMeta.Name, err)

					}
				} else {
					contextLogger.Infof("namespace %s is in the upTime range specified, not doing anything.", n.ObjectMeta.Name)
				}

			case "Suspended":

				contextLogger.Infof("namespace %s already suspended, checking if it should be revived", n.ObjectMeta.Name)
				scaledown, err := shouldScaleDown(upTime, &eng.logger, now, eng.loc)
				if err != nil {
					contextLogger.Errorf("could not determine if we should scale down: %v", err)
				}
				if scaledown {
					contextLogger.Infof("namespace %s still not in range, not doing anything.", n.ObjectMeta.Name)
				} else {
					contextLogger.Infof("namespace %s is in range, should be revived !.", n.ObjectMeta.Name)
					if err := handlingNamespace(eng, n, ctx, true); err != nil {
						contextLogger.Errorf("could not start namespace, error is %s", err)
					}
				}
			}

		case <-ctx.Done():
			// The context is over, stop processing results
			contextLogger.Infof("goroutine Suspender canceled by context")
			log.Printf("toto")
			wg.Done()
			return nil
		}
	}

}
