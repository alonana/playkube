package main

import (
	"github.com/alonana/playkube/k8s"
	"github.com/alonana/playkube/terminal"
)

func main() {
	pods := k8s.Pods{}
	go func() {
	}()
	for {
		pods.PodsList()
		pods.PodsPrint()
		command := terminal.ReadLine()
		pods.Execute(command)
	}
}
