package k8s

import (
	"fmt"
	"github.com/alonana/playkube/terminal"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func (p *Pod) Print(nameWidth int) {
	terminal.PrintEx(true, terminal.WHITE, terminal.BLACK, "%v", p.Name)
	terminal.Print("%v", strings.Repeat(" ", nameWidth-len(p.Name)+1))

	p.printContainers()
	p.printStatus()
}

func (p *Pod) GetLogs(rows int, cols int) []string {
	logs, err := RunKubectl("logs", "--tail", strconv.Itoa(rows*2), p.Name)
	if err != nil {
		return []string{err.Error()}
	}

	lines := strings.Split(logs, "\n")
	var colAdjustedLines []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		for start := 0; start < len(line); start += cols {
			end := start + cols
			if end > len(line) {
				end = len(line)
			}
			breakLine := line[start:end]
			colAdjustedLines = append(colAdjustedLines, breakLine)
		}
	}
	if len(colAdjustedLines) > rows {
		return colAdjustedLines[len(colAdjustedLines)-rows:]
	}
	return colAdjustedLines
}

func (p *Pod) printContainers() {
	var color = terminal.GREEN
	if p.ContainersRunning != p.ContainersTotal {
		color = terminal.YELLOW
	}
	containers := fmt.Sprintf("%v/%v", p.ContainersRunning, p.ContainersTotal)
	terminal.PrintEx(false, color, terminal.BLACK, "%v", containers)
}

func (p *Pod) printStatus() {
	var foregroundColor = terminal.GREEN
	if p.Status != "Running" {
		foregroundColor = terminal.YELLOW
	}
	terminal.Print(" ")
	terminal.PrintEx(false, foregroundColor, terminal.BLACK, "%v", p.Status)
	terminal.Print("%v", strings.Repeat(" ", 17-len(p.Status)+1))
}

func (p *Pod) Delete() {
	RunKubectlInBackground("delete", "pod", p.Name)
}

func (p *Pod) Build() {
	image, err := RunKubectl("get", "pod", p.Name, "-o", "jsonpath=\"{.spec.containers[0].image}\"")
	if err != nil {
		terminal.PrintEx(true, terminal.RED, terminal.BLACK, "%v", err)
		time.Sleep(10 * time.Second)
	} else {
		script := RunDocker("inspect", "--format", "'{{ index .Config.Labels \"BUILD_SCRIPT\"}}'", image)
		terminal.Clear()
		terminal.Print("Running build for ")
		terminal.PrintEx(true, terminal.WHITE, terminal.BLACK, "%v", script)
		RunBuild(script)
		time.Sleep(3 * time.Second)
		p.Delete()
	}
}
