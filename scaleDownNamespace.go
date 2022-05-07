package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func shuttingDownNamespace(eng *Engine, namespace v1.Namespace, ctx context.Context) {
	res, err := eng.client.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
	deployments, err := eng.client.AppsV1().Deployments(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})

	if err != nil {
		fmt.Printf("error retrieving deployments in namespace %s \n", err)
	}
	for _, dep := range deployments.Items {
		fmt.Printf("will set replicas of deployment %s to 0 \n", dep.ObjectMeta.Name)

		if err := patchDeploymentReplicas(ctx, eng.client, namespace.ObjectMeta.Name, dep.ObjectMeta.Name, false); err != nil {
			fmt.Printf("could not set replicas to 0 for deployment %s, error is %s \n", dep.ObjectMeta.Name, err)
		}

	}
	res.Annotations["ns-scheduler/state"] = "Suspended"
	_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("cannot update namespace annotation to 'Suspensed' for %s\n", namespace.ObjectMeta.Name)
	}
}
