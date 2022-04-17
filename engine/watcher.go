package engine

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (eng *Engine) Watcher(ctx context.Context) {

	for {
		fmt.Println("starting namespace inventory")
		ns, err := eng.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{}) // TODO: think about adding a label to filter here
		if err != nil {
			fmt.Println("cannot list namespace,%s", err)
		}

		for _, n := range ns.Items {
			fmt.Println(n.ObjectMeta.Name)

		}
	}
}
