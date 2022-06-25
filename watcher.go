package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (eng *Engine) Watcher(ctx context.Context, wg *sync.WaitGroup) error {
	contextLogger := eng.logger.WithFields(log.Fields{
		"go-routine":      "Watcher",
		"uptime-schedule": eng.upTimeSchedule,
	})
	contextLogger.Info("starting Watcher goroutine")

	ticker := time.NewTicker(time.Second * time.Duration(eng.checkInterval))
	for {
		select {
		case <-ctx.Done():
			contextLogger.Infof("goroutine watcher canceled by context")
			log.Printf("toto")
			wg.Done()
			return nil
		case <-ticker.C:
			ns, err := eng.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{}) // TODO: think about adding a label to filter here
			if err != nil {
				contextLogger.Errorf("cannot list namespace,%s", err)
				return err
			}

			for _, n := range ns.Items {
				if value, ok := n.Annotations["ns-scheduler"]; ok {
					if value == "true" {
						contextLogger.Infof("namespace %s got the annotation ns-scheduler:true, will send it to suspender", n.ObjectMeta.Name)
						eng.Wl <- n
					} else {
						contextLogger.Infof("namespace %s got annotation ns-scheduler set but the key is not true, not doing anything", n.ObjectMeta.Name)

					}

				}

			}
		}
	}
}
