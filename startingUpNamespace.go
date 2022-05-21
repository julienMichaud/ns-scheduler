package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func startingUpNamespace(eng *Engine, namespace v1.Namespace, ctx context.Context) error {
	res, err := eng.client.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not retrieve namespace info: %v", err)

	}
	deployments, err := eng.client.AppsV1().Deployments(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("error retrieving deployments in namespace %s", err)
	}
	for _, dep := range deployments.Items {
		eng.logger.Infof("will increase replicas of deployment %s", dep.ObjectMeta.Name)

		if err := patchDeploymentReplicas(ctx, eng.client, namespace.ObjectMeta.Name, dep.ObjectMeta.Name, true); err != nil {
			return fmt.Errorf("could not increase replicas for deployment %s, error is %s", dep.ObjectMeta.Name, err)
		}
		eng.logger.Infof("replicas of deployment %s incread.", dep.ObjectMeta.Name)

	}

	cronjobs, err := eng.client.BatchV1().CronJobs(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error retrieving cronjobs in namespace %s", err)
	}

	for _, cron := range cronjobs.Items {
		eng.logger.Infof("will enable cronjob %s", cron.ObjectMeta.Name)
		if err := patchCronjob(ctx, eng.client, namespace.ObjectMeta.Name, cron.ObjectMeta.Name, true); err != nil {
			return fmt.Errorf("could not enable cronjob %s, error is %s", cron.ObjectMeta.Name, err)
		}
		eng.logger.Infof("cronjob %s enabled", cron.ObjectMeta.Name)

	}

	res.Annotations["ns-scheduler/state"] = "Running"
	_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("cannot update namespace annotation to 'Running' for %s", namespace.ObjectMeta.Name)
	}
	return nil
}
