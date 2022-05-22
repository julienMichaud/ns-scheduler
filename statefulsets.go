package main

import (
	"context"
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// patchDeploymentReplicas updates the number of replicas of a given deployment
func patchStatefulSetsReplicas(ctx context.Context, cs *kubernetes.Clientset, ns, d string, scalingUp bool) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err := cs.AppsV1().StatefulSets(ns).Get(ctx, d, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// if scalingUp is false, it mean that we want to shutdown the deployment,
		// so before adjusting the replicas count, we want to save it for later
		if !scalingUp {
			result.Annotations["ns-scheduler/originalReplicas"] = strconv.Itoa(int(*result.Spec.Replicas))
			result.Annotations["ns-scheduler/state"] = "scaledDown"
			i := int32(0)
			result.Spec.Replicas = int32Ptr(i)
		} else {
			result.Annotations["ns-scheduler/state"] = "scaledUp"
			// if scaling is up, we first check if annotation ns-scheduler/originalReplicas exist on the deployment
			// if it does, we take the value and use it for the number of replicas we want
			if value, ok := result.Annotations["ns-scheduler/originalReplicas"]; ok {
				i, _ := strconv.Atoi(value)
				i32 := int32(i)
				result.Spec.Replicas = int32Ptr(i32)
			} else {
				// if the annotation doest exist, we create 3 new replicas by default (could be changed by struct ?)
				fmt.Printf("annotation 'ns-scheduler/originalReplicas' not found for statefulsets %s of namespace %s, will set replicas to 3 by default", result.ObjectMeta.Name, ns)
				i := int32(3)
				result.Spec.Replicas = int32Ptr(i)
			}
		}

		_, err = cs.AppsV1().StatefulSets(ns).Update(ctx, result, metav1.UpdateOptions{})
		return err
	})
}
