package terminal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	BLACK = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

const reset = "\033[0m"

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func Clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func Print(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func getColor(code int) string {
	return fmt.Sprintf("\033[3%dm", code)
}
func colorize(str string, color int) string {
	return fmt.Sprintf("%s%s%s", getColor(color), str, reset)
}
func bolder(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}

func PrintEx(bold bool, color int, format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	if bold {
		text = bolder(text)
	}

	Print(colorize(text, color))
}
