package k8s

import (
	"github.com/alonana/playkube/terminal"
	"strings"
	"time"
)

type Pods struct {
	filter string
	pods   []Pod
}

func (p *Pods) KeyClick() {

}
func (p *Pods) PodsList() {
	p.pods = nil
	out := RunKubectl("get", "pods")
	lines := strings.Split(out, "\n")
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			pod := PodParse(line)
			if strings.Contains(pod.Name, p.filter) {
				p.pods = append(p.pods, *pod)
			}
		}
	}
}

func (p *Pods) PodsPrint() {
	t := time.Now()
	longest := p.longestName()
	terminal.Clear()
	terminal.Print("Time: %v\n\n", t.Format("2006-01-02 15:04:05"))
	for i := 0; i < len(p.pods); i++ {
		pod := p.pods[i]
		pod.Print(longest)
	}
	if len(p.filter) > 0 {
		terminal.PrintEx(true, terminal.MAGENTA, "filtered by: %v\n", p.filter)
	} else {
		terminal.Print("Not Filtered\n")
	}
	terminal.Print("enter to refresh, text to filter, del to delete the pods")
	terminal.Print("\n")
}

func (p *Pods) longestName() int {
	longest := 0
	for i := 0; i < len(p.pods); i++ {
		pod := p.pods[i]
		if len(pod.Name) > longest {
			longest = len(pod.Name)
		}
	}
	return longest
}

func (p *Pods) Execute(command string) {
	if len(command) == 0 {
		//do nothing, just refresh
	} else if command == "del" {
		p.Delete()
	} else {
		p.filter = command
	}
}

func (p *Pods) Delete() {
	for i := 0; i < len(p.pods); i++ {
		pod := p.pods[i]
		pod.Delete()
	}
}
