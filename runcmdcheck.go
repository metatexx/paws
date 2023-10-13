package main

import (
	"fmt"
	"github.com/google/shlex"
	"os/exec"
	"sync"
	"syscall"
)

func runCommandChecked() int {
	parts, err := shlex.Split(command)
	if err != nil {
		return 1
	}
	cmd := parts[0]
	args := parts[1:]

	runner := exec.Command(cmd, args...)
	//runner.Stdin = os.Stdin
	stdout, _ := runner.StdoutPipe()
	stderr, _ := runner.StderrPipe()

	var rcOut int
	var rcErr int
	var rcCmd int
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		rcOut = readCheck(stdout)
		fmt.Println("outCheck", rcOut)
		_ = runner.Process.Signal(syscall.SIGINT)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rcErr = readCheck(stderr)
		fmt.Println("errCheck", rcErr)
		_ = runner.Process.Signal(syscall.SIGINT)
	}()

	err = runner.Start()
	if err != nil {
		fmt.Println("after start")
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = runner.Wait()
		if err != nil {
			var sigint bool
			if exitError, ok := err.(*exec.ExitError); ok {
				// Get the process status.
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					if status.Signaled() && status.Signal() == syscall.SIGINT {
						sigint = true
					}
				}
			}
			if sigint {
				fmt.Println("The process was interrupted by SIGINT.")
			} else {
				rcCmd = runner.ProcessState.ExitCode()
			}
		}
	}()

	wg.Wait()

	_ = runner.Process.Signal(syscall.SIGINT)
	fmt.Println(rcOut, rcErr, rcCmd)
	return 0
}
