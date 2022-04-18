package engine

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (eng *Engine) Watcher(ctx context.Context) {
	fmt.Println("starting Watcher goroutine")

	for range time.Tick(time.Second * 30) {
		ns, err := eng.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{}) // TODO: think about adding a label to filter here
		if err != nil {
			fmt.Println("cannot list namespace,%s", err)
		}

		for _, n := range ns.Items {
			if value, ok := n.Annotations["ns-scheduler"]; ok {
				if value == "true" {
					fmt.Printf("namespace %s got the annotation ns-scheduler:true, will send it to suspender \n", n.ObjectMeta.Name)
					eng.Wl <- n
				} else {
					fmt.Printf("namespace %s got annotation ns-scheduler set but the jey is not true, not doing anything", n.ObjectMeta.Name)

				}

			}
		}
	}
}
