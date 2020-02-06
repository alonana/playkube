package k8s

import (
	"fmt"
	"github.com/alonana/playkube/terminal"
	"regexp"
	"strconv"
)

type Pod struct {
	Name              string
	Status            string
	ContainersRunning int
	ContainersTotal   int
}

func PodParse(kubectlLine string) *Pod {
	re := regexp.MustCompile("(\\S+)\\s+(\\d+)/(\\d+)\\s+(\\S+)")
	match := re.FindStringSubmatch(kubectlLine)
	running, _ := strconv.Atoi(match[2])
	total, _ := strconv.Atoi(match[3])
	return &Pod{
		Name:              match[1],
		ContainersRunning: running,
		ContainersTotal:   total,
		Status:            match[4],
	}
}

func (pod *Pod) Print(nameWidth int) {
	format := fmt.Sprintf("%%-%vv", nameWidth+9)
	terminal.PrintEx(true, terminal.BLACK, format, pod.Name)

	pod.printContainers()
	pod.printStatus()
	terminal.Print("\n")
}

func (pod *Pod) printContainers() {
	var color = terminal.GREEN
	if pod.ContainersRunning != pod.ContainersTotal {
		color = terminal.YELLOW
	}
	containers := fmt.Sprintf("%v/%v", pod.ContainersRunning, pod.ContainersTotal)
	terminal.PrintEx(false, color, "%v", containers)
}

func (pod *Pod) printStatus() {
	var color = terminal.GREEN
	if pod.Status != "Running" {
		color = terminal.YELLOW
	}
	terminal.PrintEx(false, color, " %v", pod.Status)
}

func (pod *Pod) Delete() {
	RunKubectlInBackground("delete", "pod", pod.Name)
}
