package term

import (
	"bufio"
	"fmt"
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

// PrintStdout will create a stdout pipe for the passed in command
func PrintStdout(command *exec.Cmd) error {
	stdout, err := command.StdoutPipe()

	if err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Error creating a stdout pipe")
	}

	if err := command.Start(); err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Error starting the command")
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			out := scanner.Text()
			fmt.Println(out)
		}
	}()

	if err := command.Wait(); err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Error waiting on releases of executed command")
	}

	return nil
}
