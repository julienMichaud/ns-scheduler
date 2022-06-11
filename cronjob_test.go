package main

import (
	"context"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	testclient "k8s.io/client-go/kubernetes/fake"
)

func CreatingCronjob(ctx context.Context, clientset *testclient.Clientset) (*batchv1.CronJob, error) {
	suspend := true
	cronjob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-cronjob",
			Namespace: "namespace",
			Annotations: map[string]string{
				"test": "test",
			},
		},
		Spec: batchv1.CronJobSpec{
			Suspend: &suspend,
		},
	}

	result, err := clientset.BatchV1().CronJobs("namespace").Create(ctx, cronjob, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetCronjob(ctx context.Context, clientset *testclient.Clientset) (*batchv1.CronJob, error) {
	result, err := clientset.BatchV1().CronJobs("namespace").Get(ctx, "demo-cronjob", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreatingCronjobResources(ctx context.Context, clientset *testclient.Clientset) error {

	_, err := CreatingCronjob(ctx, clientset)
	if err != nil {
		return err
	}
	return nil
}

func TestCronjobFunc(t *testing.T) {
	ctx := context.Background()
	clientset := testclient.NewSimpleClientset()

	err := CreatingNamespace(ctx, clientset)
	if err != nil {
		t.Errorf("error creating resources, %s", err)
	}

	err = CreatingCronjobResources(ctx, clientset)
	if err != nil {
		t.Errorf("error creating resources, %s", err)
	}

	cronjobTestCases := []struct {
		scalingUp bool
		suspend   bool
	}{
		{
			scalingUp: false,
			suspend:   true,
		},
		{
			scalingUp: true,
			suspend:   false,
		},
	}
	for _, tt := range cronjobTestCases {

		err = patchCronjob(ctx, clientset, "namespace", "demo-cronjob", tt.scalingUp)
		if err != nil {
			t.Errorf("func cronjob return error %s", err)
		}

		got, err := GetCronjob(ctx, clientset)
		if err != nil {
			t.Errorf("cannot get cronjob, %s", err)
		}

		if *got.Spec.Suspend != tt.suspend {
			t.Errorf("got %v suspend status, want %v", *got.Spec.Suspend, tt.suspend)
		}

	}

}
