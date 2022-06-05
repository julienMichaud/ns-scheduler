package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewFunc(t *testing.T) {
	clientset := testclient.NewSimpleClientset()
	var log = logrus.New()
	upTimeSchedule := "1-7 08:00-20:00"

	eng := New(
		clientset,
		upTimeSchedule,
		*log)

	if eng.upTimeSchedule != "1-7 08:00-20:00" {
		t.Errorf("New func did not put correct upTimeSchedule,")
	}

}
