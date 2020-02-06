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
