package main

import (
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Engine struct {
	client         kubernetes.Interface
	Wl             chan v1.Namespace
	upTimeSchedule string
	logger         logrus.Logger
	loc            *time.Location
	checkInterval  int
}

func New(cs kubernetes.Interface, upTimeSchedule string, log logrus.Logger, intervalTime int) *Engine {
	loc, _ := time.LoadLocation("Europe/Paris")
	e := Engine{
		client:         cs,
		upTimeSchedule: upTimeSchedule,
		Wl:             make(chan v1.Namespace, 30),
		logger:         log,
		loc:            loc,
		checkInterval:  intervalTime,
	}
	return &e
}
