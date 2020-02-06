package k8s

import (
	"log"
	"os/exec"
)

func RunKubectl(args ...string) string {
	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func RunKubectlInBackground(args ...string) {
	err := exec.Command("kubectl", args...).Start()
	if err != nil {
		log.Fatal(err)
	}
}
