package main

import (
	"github.com/alonana/playkube/k8s"
	"time"
)

func main() {
	for {
		pods := k8s.Pods{}
		pods.PodsList()
		pods.PodsPrint()
		time.Sleep(time.Second)
	}
}
