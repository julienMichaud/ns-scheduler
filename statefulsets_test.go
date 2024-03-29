package main

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func CreatingStatefulSet(ctx context.Context, clientset *testclient.Clientset) (*appsv1.StatefulSet, error) {
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-statefulset",
			Namespace: "namespace",
			Annotations: map[string]string{
				"test": "test",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(2),
		},
	}

	result, err := clientset.AppsV1().StatefulSets("namespace").Create(ctx, statefulset, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetStatefulSet(ctx context.Context, clientset *testclient.Clientset) (*appsv1.StatefulSet, error) {
	result, err := clientset.AppsV1().StatefulSets("namespace").Get(ctx, "demo-statefulset", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreatingStatefulSetResources(ctx context.Context, clientset *testclient.Clientset) error {

	_, err := CreatingStatefulSet(ctx, clientset)
	if err != nil {
		return err
	}
	return nil
}

func TestStatefulsetsFunc(t *testing.T) {
	ctx := context.Background()
	clientset := testclient.NewSimpleClientset()
	err := CreatingNamespace(ctx, clientset)
	if err != nil {
		t.Errorf("error creating resources, %s", err)
	}

	err = CreatingStatefulSetResources(ctx, clientset)
	if err != nil {
		t.Errorf("error creating resources, %s", err)
	}

	statefulsetTestCases := []struct {
		scalingUp            bool
		wantReplicas         int32
		wantSchedulerState   string
		wantOriginalReplicas string
	}{
		{
			scalingUp:            false,
			wantReplicas:         0,
			wantSchedulerState:   "scaledDown",
			wantOriginalReplicas: "2",
		},
		{
			scalingUp:            true,
			wantReplicas:         2,
			wantSchedulerState:   "scaledUp",
			wantOriginalReplicas: "2",
		},
	}
	for _, tt := range statefulsetTestCases {

		err = patchStatefulSetsReplicas(ctx, clientset, "namespace", "demo-statefulset", tt.scalingUp)
		if err != nil {
			t.Errorf("func statefulset return error %s", err)
		}

		got, err := GetStatefulSet(ctx, clientset)
		if err != nil {
			t.Errorf("cannot get statefulset, %s", err)
		}

		if *got.Spec.Replicas != tt.wantReplicas {
			t.Errorf("got %v replicas want %v", *got.Spec.Replicas, tt.wantReplicas)
		}

		v, found := got.Annotations["ns-scheduler/state"]
		if !found {
			t.Errorf("annotation ns-scheduler/State not found,%s", err)
		}

		if v != tt.wantSchedulerState {
			t.Errorf("annotation ns-scheduler/State value is %s, want %s", v, tt.wantSchedulerState)
		}

		v, found = got.Annotations["ns-scheduler/originalReplicas"]
		if !found {
			t.Errorf("annotation ns-scheduler/originalReplicas not found,%s", err)
		}

		if v != tt.wantOriginalReplicas {
			t.Errorf("annotation ns-scheduler/originalReplicas value is %s, want %s", v, tt.wantOriginalReplicas)
		}

	}

}
