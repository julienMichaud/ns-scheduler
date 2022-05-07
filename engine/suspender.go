package engine

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (eng *Engine) Suspender(ctx context.Context) {

	fmt.Println("starting Suspender goroutine")

	for {
		n := <-eng.Wl
		fmt.Printf("received namespace: %s \n", n.ObjectMeta.Name)

		dState := n.Annotations["ns-scheduler/state"]

		switch dState {
		// if no annotation, we just set it
		case "": //put everything under into a func
			fmt.Printf("didnt see annotation ns-scheduler-state on the namespace %s, first time im seeing it then\n", n.ObjectMeta.Name)
			if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				res, err := eng.client.CoreV1().Namespaces().Get(ctx, n.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				res.Annotations["ns-scheduler/state"] = "Running"

				_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
				return err
			}); err != nil {
				fmt.Printf("cannot update namespace object\n")

				continue
			}
		//if namespace annotation is Running, check if the namespace resources should be up based on the upTime variable, if not we scale down resources
		case "Running":
			if shouldScaleDown(eng.upTimeSchedule) {
				fmt.Printf("namespace %s is in state 'Running' and is not in the upTime range specified, which mean that it should be scaled down.\n The state of the namespace will then be Suspended.\n", n.ObjectMeta.Name)
				shuttingDownNamespace(eng, n, ctx)
			} else {
				fmt.Printf("namespace %s is in state 'Running' and is in the upTime range specified, not doing anything.\n", n.ObjectMeta.Name)
			}

		case "Suspended":
			fmt.Printf("namespace %s already suspended, checking if it should be revived\n", n.ObjectMeta.Name)
			if shouldScaleDown(eng.upTimeSchedule) {
				fmt.Printf("namespace %s still not in range, not doing anything.\n", n.ObjectMeta.Name)
			} else {
				fmt.Printf("namespace %s is in range, should be revived !.\n", n.ObjectMeta.Name)
				startingUpNamespace(eng, n, ctx)
			}
		}

	}
}
