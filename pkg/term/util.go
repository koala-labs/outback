package term

import (
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Clear clears your terminal screen depending on OS
func Clear() {
	value, ok := clear[runtime.GOOS]

	if !ok {
		panic("Platform is unsupported.")
	}

	value()
}
