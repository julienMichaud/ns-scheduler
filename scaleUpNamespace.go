package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func startingUpNamespace(eng *Engine, namespace v1.Namespace, ctx context.Context) {
	res, err := eng.client.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
	deployments, err := eng.client.AppsV1().Deployments(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})

	if err != nil {
		fmt.Printf("error retrieving deployments in namespace %s \n", err)
	}
	for _, dep := range deployments.Items {
		fmt.Printf("will increase replicas of deployment %s \n", dep.ObjectMeta.Name)

		if err := patchDeploymentReplicas(ctx, eng.client, namespace.ObjectMeta.Name, dep.ObjectMeta.Name, true); err != nil {
			fmt.Printf("could not increase replicas for deployment %s, error is %s \n", dep.ObjectMeta.Name, err)
		}

	}
	res.Annotations["ns-scheduler/state"] = "Running"
	_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("cannot update namespace annotation to 'Running' for %s\n", namespace.ObjectMeta.Name)
	}
}
