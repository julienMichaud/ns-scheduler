package main

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Engine struct {
	client         *kubernetes.Clientset
	Wl             chan v1.Namespace
	upTimeSchedule string
	logger         logrus.Logger
}

func New(cs *kubernetes.Clientset, upTimeSchedule string, log logrus.Logger) *Engine {
	e := Engine{
		client:         cs,
		upTimeSchedule: upTimeSchedule,
		Wl:             make(chan v1.Namespace, 30),
		logger:         log,
	}
	return &e
}
