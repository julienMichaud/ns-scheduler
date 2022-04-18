package engine

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Engine struct {
	client *kubernetes.Clientset
	Wl     chan v1.Namespace
}

func New(cs *kubernetes.Clientset) *Engine {
	e := Engine{
		client: cs,
		Wl:     make(chan v1.Namespace, 30),
	}
	return &e
}
