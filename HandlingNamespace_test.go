package main

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestHandlingNamespaceFunc(t *testing.T) {
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

	err = CreatingStatefulSetResources(ctx, clientset)
	if err != nil {
		t.Errorf("error creating resources, %s", err)
	}

	err = CreatingDeploymentResources(ctx, clientset)
	if err != nil {
		t.Errorf("error creating resources, %s", err)
	}

	HandlingNamespaceTestCases := []struct {
		scalingUp                bool
		wantReplicasStatefulSet  int32
		wantReplicasDeployment   int32
		wantSuspendStatusCronjob bool
		annotationNamespace      string
	}{
		{
			scalingUp:           false,
			annotationNamespace: "Suspended",
		},
		{
			scalingUp:           true,
			annotationNamespace: "Running",
		},
	}
	var log = logrus.New()
	eng := New(
		clientset,
		"1-7 08:00-20:00",
		*log)
	namespace, _ := clientset.CoreV1().Namespaces().Get(ctx, "namespace", v1.GetOptions{})

	for _, tt := range HandlingNamespaceTestCases {

		err := handlingNamespace(eng, *namespace, ctx, tt.scalingUp)
		if err != nil {
			t.Errorf("func handlingNamespace return error %s", err)
		}
		namespace, _ := clientset.CoreV1().Namespaces().Get(ctx, "namespace", v1.GetOptions{})

		v, found := namespace.Annotations["ns-scheduler/state"]
		if !found {
			t.Errorf("annotation ns-scheduler/state not found,%s", err)
		}

		if v != tt.annotationNamespace {
			t.Errorf("annotation 'ns-scheduler/state' is %s, want %s", v, tt.annotationNamespace)
		}

	}
}
