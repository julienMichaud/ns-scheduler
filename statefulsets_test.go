package main

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestStatefulsetsFunc(t *testing.T) {
	ctx := context.Background()
	clientset := testclient.NewSimpleClientset()

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-deployment",
			Namespace: "namespace",
			Annotations: map[string]string{
				"test": "test",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(2),
		},
	}

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "namespace"}}
	_, err := clientset.CoreV1().Namespaces().Create(ctx, nsSpec, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("%s", err)
	}

	_, _ = clientset.AppsV1().StatefulSets("namespace").Create(ctx, statefulset, metav1.CreateOptions{})

	err = patchStatefulSetsReplicas(ctx, clientset, "namespace", "demo-deployment", true)
	if err != nil {
		t.Errorf("func statefulset test error %s", err)
	}

}
