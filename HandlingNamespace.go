package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handlingNamespace(eng *Engine, namespace v1.Namespace, ctx context.Context, scalingUp bool) error {

	// deployment part
	res, err := eng.client.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("cannot get namespace infos %s", err)
	}

	deployments, err := eng.client.AppsV1().Deployments(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error retrieving deployments in namespace %s", err)
	}

	var deploymentString string

	deploymentStringError := "error while changing replicas of deployment"
	deploymentStringStatus := "replicas of deployment changed"

	if scalingUp {
		deploymentString = "will increase replicas of deployment"

	} else {
		deploymentString = "will decrease replicas of deployment"
	}

	for _, dep := range deployments.Items {
		eng.logger.Infof(deploymentString + fmt.Sprintf(" %s", dep.ObjectMeta.Name))

		if err := patchDeploymentReplicas(ctx, eng.client, namespace.ObjectMeta.Name, dep.ObjectMeta.Name, scalingUp); err != nil {
			return fmt.Errorf(deploymentStringError+" %s, error is %s", dep.ObjectMeta.Name, err)
		}
		eng.logger.Infof(deploymentStringStatus)

	}

	// statefulsets part
	statefulsets, err := eng.client.AppsV1().StatefulSets(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error retrieving statefulsets in namespace %s", err)
	}

	var statefulSetsString string

	statefulSetsStringError := "error while changing replicas of statefulSets"
	statefulSetsStringStatus := "replicas of statefulSets changed"
	if scalingUp {
		deploymentString = "will increase replicas of statefulSets"

	} else {
		statefulSetsString = "will decrease replicas of statefulSets"
	}

	for _, dep := range statefulsets.Items {
		eng.logger.Infof(statefulSetsString + fmt.Sprintf("%s", dep.ObjectMeta.Name))

		if err := patchStatefulSetsReplicas(ctx, eng.client, namespace.ObjectMeta.Name, dep.ObjectMeta.Name, scalingUp); err != nil {
			return fmt.Errorf(statefulSetsStringError+" %s, error is %s", dep.ObjectMeta.Name, err)
		}
		eng.logger.Infof(statefulSetsStringStatus)

	}

	// cronjob part
	cronjobs, err := eng.client.BatchV1().CronJobs(namespace.ObjectMeta.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error retrieving cronjobs in namespace %s", err)
	}

	var cronjobString string

	cronjobStringError := "error while changing cronjob status"
	cronjobStringStatus := "cronjob status changed"
	if scalingUp {
		deploymentString = "will enable cronjob"

	} else {
		cronjobString = "will disable cronjob"
	}

	for _, cron := range cronjobs.Items {
		eng.logger.Infof(cronjobString + fmt.Sprintf(" %s", cron.ObjectMeta.Name))
		if err := patchCronjob(ctx, eng.client, namespace.ObjectMeta.Name, cron.ObjectMeta.Name, scalingUp); err != nil {
			return fmt.Errorf(cronjobStringError+" %s, error is %s", cron.ObjectMeta.Name, err)
		}
		eng.logger.Infof(cronjobStringStatus)

	}

	var state string
	if scalingUp {
		state = "Running"
	} else {
		state = "Suspended"
	}
	res.Annotations["ns-scheduler/state"] = state
	_, err = eng.client.CoreV1().Namespaces().Update(ctx, res, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("cannot update namespace annotation to %s for %s", state, namespace.ObjectMeta.Name)
	}
	return nil
}
