package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

var (
	currentCmd   *exec.Cmd
	oldState     *term.State
	sourceFile   string
	tempExe      string
	waitingRerun bool
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "run" {
		fmt.Println("Usage: cref run <file.c>")
		os.Exit(1)
	}

	sourceFile = os.Args[2]
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Printf("%sError: File '%s' not found%s\n", colorRed, sourceFile, colorReset)
		os.Exit(1)
	}

	tempExe = filepath.Join(os.TempDir(), "cref_temp_"+filepath.Base(sourceFile)+".out")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cleanup()
		fmt.Printf("\r\n%sExiting...%s\n", colorYellow, colorReset)
		os.Exit(0)
	}()

	fmt.Printf("%s=== C Refresher ===%s\n", colorCyan, colorReset)
	fmt.Printf("Press %sCTRL+R%s to recompile\n\n", colorGreen, colorReset)

	compileAndRun()
	keyboardListener()
}

func keyboardListener() {
	var err error
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}

	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		if buf[0] == 18 {
			if currentCmd != nil && currentCmd.Process != nil {
				currentCmd.Process.Kill()
			}
			restoreTerminal()
			fmt.Printf("\r\n")
			compileAndRun()
			setRawMode()
		} else if waitingRerun {
			if buf[0] == 'y' || buf[0] == 'Y' {
				restoreTerminal()
				fmt.Printf("y\r\n")
				waitingRerun = false
				runProgram()
			} else if buf[0] == 'n' || buf[0] == 'N' {
				restoreTerminal()
				fmt.Printf("n\r\n")
				cleanup()
				os.Exit(0)
			}
		}
	}
}

func compileAndRun() {
	os.Remove(tempExe)

	cmd := exec.Command("clang", "-o", tempExe, sourceFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("%sCompilation failed:%s\n%s\n", colorRed, colorReset, string(output))
		return
	}

	if len(output) > 0 {
		fmt.Printf("%s%s%s", colorYellow, string(output), colorReset)
	}

	runProgram()
}

func runProgram() {
	restoreTerminal()

	currentCmd = exec.Command(tempExe)
	currentCmd.Stdout = os.Stdout
	currentCmd.Stderr = os.Stderr
	currentCmd.Stdin = os.Stdin

	err := currentCmd.Run()
	currentCmd = nil

	setRawMode()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != -1 {
				fmt.Printf("%sProgram exited with code %d%s\r\n", colorRed, exitErr.ExitCode(), colorReset)
			}
		}
	}

	fmt.Printf("%sRerun? (y/n): %s", colorCyan, colorReset)
	waitingRerun = true
}

func setRawMode() {
	if oldState != nil {
		term.MakeRaw(int(os.Stdin.Fd()))
	}
}

func restoreTerminal() {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

func cleanup() {
	restoreTerminal()
	if currentCmd != nil && currentCmd.Process != nil {
		currentCmd.Process.Kill()
	}
	os.Remove(tempExe)
}
