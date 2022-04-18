package engine

import (
	"context"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func int32Ptr(i int32) *int32 { return &i }

// patchDeploymentReplicas updates the number of replicas of a given deployment
func patchDeploymentReplicas(ctx context.Context, cs *kubernetes.Clientset, ns, d, prefix string, repl int32) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err := cs.AppsV1().Deployments(ns).Get(ctx, d, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// if we want 0 replicas, it means that we are suspending the namespace,
		// so before adjusting the replicas count, we want to save it for later
		if repl == 0 {
			result.Annotations["ns-scheduler/originalReplicas"] = strconv.Itoa(int(*result.Spec.Replicas))
			result.Annotations["ns-scheduler/state"] = "scaledDown"
		}
		result.Spec.Replicas = int32Ptr(repl)
		_, err = cs.AppsV1().Deployments(ns).Update(ctx, result, metav1.UpdateOptions{})
		return err
	})
}
