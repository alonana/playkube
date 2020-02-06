package k8s

import (
	"github.com/alonana/playkube/terminal"
	"strings"
	"time"
)

type Pods struct {
	pods []PodInfo
}

func (p *Pods) PodsList() {
	out := RunKubectl("get", "pods")
	lines := strings.Split(out, "\n")
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			pod := PodParse(line)
			p.pods = append(p.pods, *pod)
		}
	}
}

func (p *Pods) PodsPrint() {
	t := time.Now()
	longest := p.longestName()
	terminal.Clear()
	terminal.MoveCursor(1, 1)
	terminal.Printf("Time: %v", t.Format("2006-01-02 15:04:05"))
	for i := 0; i < len(p.pods); i++ {
		terminal.MoveCursor(1, i+3)
		pod := p.pods[i]
		pod.Print(longest)
	}
	terminal.Flush()
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
