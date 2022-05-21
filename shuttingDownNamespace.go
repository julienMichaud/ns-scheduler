package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func shuttingDownNamespace(eng *Engine, namespace v1.Namespace, ctx context.Context) error {
	res, err := eng.client.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("cannot get namespace infos %s", err)
	}

	deployments, err := eng.client.AppsV1().Deployments(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error retrieving deployments in namespace %s", err)
	}
	for _, dep := range deployments.Items {
		eng.logger.Infof("will set replicas of deployment %s to 0 \n", dep.ObjectMeta.Name)

		if err := patchDeploymentReplicas(ctx, eng.client, namespace.ObjectMeta.Name, dep.ObjectMeta.Name, false); err != nil {
			return fmt.Errorf("could not set replicas to 0 for deployment %s, error is %s", dep.ObjectMeta.Name, err)
		}
		eng.logger.Infof("replicas of deployment %s set to 0", dep.ObjectMeta.Name)

	}

	cronjobs, err := eng.client.BatchV1().CronJobs(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error retrieving cronjobs in namespace %s", err)
	}

	for _, cron := range cronjobs.Items {
		eng.logger.Infof("will disable cronjob %s", cron.ObjectMeta.Name)
		if err := patchCronjob(ctx, eng.client, namespace.ObjectMeta.Name, cron.ObjectMeta.Name, false); err != nil {
			return fmt.Errorf("could not disable cronjob %s, error is %s", cron.ObjectMeta.Name, err)
		}
		eng.logger.Infof("cronjob %s disabled", cron.ObjectMeta.Name)

	}

	res.Annotations["ns-scheduler/state"] = "Suspended"
	_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("cannot update namespace annotation to 'Suspensed' for %s", namespace.ObjectMeta.Name)
	}
	return nil
}