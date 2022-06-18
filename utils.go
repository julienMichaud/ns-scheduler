package main

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func CreatingNamespace(ctx context.Context, clientset *testclient.Clientset) error {
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name:        "namespace",
		Annotations: map[string]string{"ns-scheduler": "true"}}}
	_, err := clientset.CoreV1().Namespaces().Create(ctx, nsSpec, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
