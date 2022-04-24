package engine

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Engine struct {
	client         *kubernetes.Clientset
	Wl             chan v1.Namespace
	upTimeSchedule string
}

func New(cs *kubernetes.Clientset, upTimeSchedule string) *Engine {
	e := Engine{
		client:         cs,
		upTimeSchedule: upTimeSchedule,
		Wl:             make(chan v1.Namespace, 30),
	}
	return &e
}
