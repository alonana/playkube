package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"strings"
)

var Output = bufio.NewWriter(os.Stdout)
var Screen = new(bytes.Buffer)

const RESET = "\033[0m"

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

func Clear() {
	Output.WriteString("\033[2J")
}

func MoveCursor(x int, y int) {
	fmt.Fprintf(Screen, "\033[%d;%dH", y, x)
}

func getColor(code int) string {
	return fmt.Sprintf("\033[3%dm", code)
}

func Bold(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}

func Color(str string, color int) string {
	return fmt.Sprintf("%s%s%s", getColor(color), str, RESET)
}
func Printf(format string, a ...interface{}) {
	fmt.Fprintf(Screen, format, a...)
}

func Width() int {
	ws, err := getWinsize()

	if err != nil {
		return -1
	}

	return int(ws.Col)
}

func Height() int {
	ws, err := getWinsize()
	if err != nil {
		return -1
	}
	return int(ws.Row)
}

func Flush() {
	for idx, str := range strings.SplitAfter(Screen.String(), "\n") {
		if idx > Height() {
			return
		}

		Output.WriteString(str)
	}

	Output.Flush()
	Screen.Reset()
}

func getWinsize() (*unix.Winsize, error) {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return nil, os.NewSyscallError("GetWinsize", err)
	}

	return ws, nil
}
