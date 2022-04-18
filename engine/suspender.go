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
		case "":
			fmt.Printf("didnt see annotation ns-scheduler-state on the namespace %s, first time im seeing it then\n", n.ObjectMeta.Name)
			if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				fmt.Printf("step 1, get namespac\n")
				res, err := eng.client.CoreV1().Namespaces().Get(ctx, n.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				// we set the annotation to running
				fmt.Printf("step 2, setting namespace annotation to running\n")
				res.Annotations["ns-scheduler/state"] = "Running"

				fmt.Printf("step 3,updating namespace\n")
				_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
				return err
			}); err != nil {
				fmt.Printf("cannot update namespace object\n")

				continue
			}
		//if namespace annotation is Running, we scale down every deployment of it and set ns-scheduler/state to Suspended
		case "Running":
			fmt.Printf("namespace annotation ns-scheduler/state is Running, will scale down every deployment of it and set ns-scheduler/state to Suspended\n")
			res, err := eng.client.CoreV1().Namespaces().Get(ctx, n.Name, metav1.GetOptions{})
			deployments, err := eng.client.AppsV1().Deployments(n.ObjectMeta.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				fmt.Printf("error retrieving deployments in namespace %s \n", err)
			}
			for _, dep := range deployments.Items {
				fmt.Printf("will set replicas of deployment %s to 0 \n", dep.ObjectMeta.Name)

				if err := patchDeploymentReplicas(ctx, eng.client, n.ObjectMeta.Name, dep.ObjectMeta.Name, "toto", 0); err != nil {
					fmt.Printf("could not set replicas to 0 for deployment %s, error is %s \n", dep.ObjectMeta.Name, err)
				}

			}
			res.Annotations["ns-scheduler/state"] = "Suspended"
			_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
			if err != nil {
				fmt.Printf("cannot update namespace annotation to suspensed for %s\n", n.ObjectMeta.Name)
			}

		case "Suspended":
			fmt.Printf("namespace %s already suspended, not doing anything\n", n.ObjectMeta.Name)
		}

	}
}
