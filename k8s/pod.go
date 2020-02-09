package k8s

import (
	"fmt"
	"github.com/alonana/playkube/terminal"
	"regexp"
	"strconv"
	"strings"
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
	logs := RunKubectl("logs", "--tail", strconv.Itoa(rows*2), p.Name)

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
	terminal.PrintEx(false, foregroundColor, terminal.BLACK, "%-9v", p.Status)
}

func (p *Pod) Delete() {
	RunKubectlInBackground("delete", "pod", p.Name)
}
