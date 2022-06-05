package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func boolFunc(i bool) *bool { return &i }

// patchCronjob disable a given cronjob
func patchCronjob(ctx context.Context, cs kubernetes.Interface, ns, d string, scalingUp bool) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, err := cs.BatchV1().CronJobs(ns).Get(ctx, d, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// if scalingUp is false, it mean that we want to disable the cronjob
		if !scalingUp {
			result.Spec.Suspend = boolFunc(true)
		} else {
			result.Spec.Suspend = boolFunc(false)
		}

		_, err = cs.BatchV1().CronJobs(ns).Update(ctx, result, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		return err
	})

}
