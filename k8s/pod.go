package k8s

import (
	"fmt"
	"github.com/alonana/playkube/terminal"
	"regexp"
	"strconv"
)

type PodInfo struct {
	Name              string
	ContainersRunning int
	ContainersTotal   int
}

func PodParse(kubectlLine string) *PodInfo {
	re := regexp.MustCompile("(\\S+)\\s+(\\d+)/(\\d+)")
	match := re.FindStringSubmatch(kubectlLine)
	running, _ := strconv.Atoi(match[2])
	total, _ := strconv.Atoi(match[3])
	return &PodInfo{
		Name:              match[1],
		ContainersRunning: running,
		ContainersTotal:   total,
	}
}

func (pod *PodInfo) Print(nameWidth int) {

	format := fmt.Sprintf("%%-%vv", nameWidth+9)
	terminal.Printf(format, terminal.Bold(pod.Name))

	var color = terminal.GREEN
	if pod.ContainersRunning != pod.ContainersTotal {
		color = terminal.YELLOW
	}
	containers := fmt.Sprintf("%v/%v", pod.ContainersRunning, pod.ContainersTotal)
	terminal.Printf("%v", terminal.Color(containers, color))
}
