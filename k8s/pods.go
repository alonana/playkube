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

var podStaticFieldsLength = 14

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
	rows, cols := terminal.GetWindowSize()

	// reserve one line for command
	rows--

	longestName := p.longestName()
	var logs []string
	if len(p.pods) > 0 {
		logs = p.pods[0].GetLogs(rows-1, cols-longestName-podStaticFieldsLength-1)
	}
	terminal.Clear()
	t := time.Now()
	terminal.Print("Time: %v\n\n", t.Format("2006-01-02 15:04:05"))
	for row := 0; row < rows; row++ {
		p.printRowLeftSide(longestName, rows, row)
		terminal.PrintEx(true, terminal.BLACK, terminal.WHITE, "│")
		p.printRowRightSide(row, logs)
		terminal.Print("\n")
	}
}

func (p *Pods) printRowLeftSide(longestName int, rows int, row int) {
	leftSideWidth := longestName + podStaticFieldsLength
	if row < len(p.pods) {
		pod := p.pods[row]
		pod.Print(longestName)
	} else if row == rows-2 {
		terminal.PrintEx(true, terminal.WHITE, terminal.BLACK, "<Enter>")
		terminal.Print(" refresh  ")
		terminal.PrintEx(true, terminal.WHITE, terminal.BLACK, "'del'")
		terminal.Print(" kill pods  ")
		terminal.PrintEx(true, terminal.WHITE, terminal.BLACK, "'clear'")
		terminal.Print(" clear filter                  ")
	} else if row == rows-1 {
		if len(p.filter) > 0 {
			terminal.Print("filtered by: ")
			terminal.PrintEx(true, terminal.MAGENTA, terminal.BLACK, "%v", p.filter)
			terminal.Print("%v", strings.Repeat(" ", leftSideWidth-len(p.filter)-13))
		} else {
			terminal.Print("Not Filtered")
			terminal.Print("%v", strings.Repeat(" ", leftSideWidth-12))
		}
	} else {
		terminal.Print("%v", strings.Repeat(" ", leftSideWidth))
	}
}

func (p *Pods) printRowRightSide(row int, logs []string) {
	if row == 0 {
		if len(p.pods) > 0 {
			terminal.PrintEx(true, terminal.BLUE, terminal.WHITE, "Logs of %v", p.pods[0].Name)
		}
	} else if row < len(logs) {
		terminal.Print("%v", logs[row-1])
	}
}

func (p *Pods) longestName() int {
	longest := 58 // use minimum considering legend width
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
	} else if command == "clear" {
		p.filter = ""
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
