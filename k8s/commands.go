package k8s

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func RunDocker(args ...string) string {
	cmd := strings.Join(args, " ")
	return RunCommandStreaming("docker " + cmd)
}

func RunKubectl(args ...string) (string, error) {
	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func RunKubectlInBackground(args ...string) {
	err := exec.Command("kubectl", args...).Start()
	if err != nil {
		log.Fatal(err)
	}
}

func RunCommandStreaming(build string) string {
	cmd := exec.Command("sh", "-c", build)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var allStdout string
	var allStderr string
	go func() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(stderr)
		allStderr = buf.String()
		allStderr = strings.TrimSpace(allStderr)
		wg.Done()
	}()
	go func() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(stdout)
		allStdout = buf.String()
		allStdout = strings.TrimSpace(allStdout)
		wg.Done()
	}()
	wg.Wait()
	if len(allStderr) > 0 {
		log.Fatal(allStderr)
	}
	return allStdout
}
func RunBuild(build string) {
	cmd := exec.Command("sh", "-c", build)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Start()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(os.Stdout, stdout)
		wg.Done()
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
		wg.Done()
	}()
	wg.Wait()
}
