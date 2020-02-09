package terminal

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
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

func GetWindowSize() (int, int) {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		log.Fatalf("get window size failed: %v", err)
	}

	return int(ws.Row), int(ws.Col)
}

func Print(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func getForegroundColor(code int) string {
	return fmt.Sprintf("\033[3%dm", code)
}

func getBackgroundColor(code int) string {
	return fmt.Sprintf("\033[4%dm", code)
}

func colorize(str string, foregroundColor int, backgroundColor int) string {
	return fmt.Sprintf("%s%s%s%s", getBackgroundColor(backgroundColor), getForegroundColor(foregroundColor), str, reset)
}

func bolder(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}

func PrintEx(bold bool, foregroundColor int, backgroundColor int, format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	if bold {
		text = bolder(text)
	}

	Print(colorize(text, foregroundColor, backgroundColor))
}
