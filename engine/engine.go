package engine

import (
	"k8s.io/client-go/kubernetes"
)

type Engine struct {
	client *kubernetes.Clientset
}

func New(cs *kubernetes.Clientset) *Engine {
	e := Engine{
		client: cs,
	}
	return &e
}
