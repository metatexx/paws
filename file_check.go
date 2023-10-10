package main

import (
	"bufio"
	"fmt"
	"github.com/google/shlex"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

func fileCheck() int {
	if command != "" {
		return runCommandChecked()
	}

	var err error
	input := os.Stdin
	if fileGiven {
		input, err = os.Open(file)
		if err != nil {
			fmt.Printf("*** file error: %v\n", err)
			return 1
		}
		defer func() {
			_ = input.Close()
		}()
	}

	return readCheck(input, false)
}

func readCheck(input io.ReadCloser, noTail bool) int {
	// for checking patterns we keep the remainder if a line is not fully output so we do not miss matches
	var remainder string
	startTime := time.Now()
	idlePoint := time.Now()

	foundPatters := make(map[*regexp.Regexp]bool, len(okOnPattern))
	out := os.Stdout
	if passThrough {
		out = os.Stderr
	}

	reader := bufio.NewReader(input)
	for {
		fmt.Println("read 1")
		line, err := reader.ReadString('\n')
		fmt.Println("read 2")
		if err != nil {
			if err != io.EOF {
				_, err2 := fmt.Fprintf(out, "*** read error: %v\n", err)
				if err2 != nil {
					panic(err2)
				}
				return -1
			} /*
				if noTail {
					// nothing found?
					return 1
				}*/
		}
		if time.Since(startTime) > timeout {
			_, err = fmt.Fprintln(out, "*** aborting: timeout reached")
			if err != nil {
				panic(err)
			}
			return 1
		}
		// If a new line is found, print it.
		if len(line) > 0 {
			idlePoint = time.Now()
			if passThrough {
				// we output as it comes in
				fmt.Print(line)
			}
			if err == io.EOF {
				remainder += line // this is a bit inefficient (for now)
			} else {
				if remainder != "" {
					if !quiet {
						line = remainder + line
					}
					remainder = ""
				}
			}
			check := []byte(strings.TrimSuffix(line, "\n"))
			// check out patterns
			for _, re := range failOnPattern {
				if re.Match(check) {
					if !quiet {
						_, err = fmt.Fprintln(out, "*** aborting: fail pattern match")
						if err != nil {
							panic(err)
						}
					}
					return 1
				}
			}
			if len(okOnPattern) > 0 {
				success := true
				for _, rx := range okOnPattern {
					if b, ok := foundPatters[rx]; ok {
						if b {
							continue
						}
					}
					if rx.Match(check) {
						if printMatches {
							if passThrough {
								_, err = fmt.Fprint(os.Stderr, highlight(string(check), rx))
							} else {
								fmt.Println(highlight(string(check), rx))
							}
						}
						foundPatters[rx] = true
						continue
					}
					success = false
				}
				if success {
					if !quiet {
						_, err = fmt.Fprintln(out, "*** success: found all patterns")
						if err != nil {
							panic(err)
						}
					}
					return 0
				}
			}
		}
		if idleTime != 0 && time.Since(idlePoint) > idleTime {
			if !quiet {
				_, err = fmt.Fprintln(out, "*** success: idle reached")
				if err != nil {
					panic(err)
				}
			}
			return 0
		}
		// If at end of file, wait for a bit before seeking more lines.
		if err == io.EOF {
			fmt.Println("read more")
			time.Sleep(250 * time.Millisecond) // Sleep a bit
		}
	}
}

func highlight(input string, re *regexp.Regexp) string {
	indices := re.FindStringSubmatchIndex(input)
	if indices == nil {
		return input
	}

	highlighted := input[:indices[0]] +
		"\033[1;31m" + // Start red highlight
		input[indices[0]:indices[1]] +
		"\033[0m" + // End highlight
		input[indices[1]:]

	return highlighted
}

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
		rcOut = readCheck(stdout, true)
		fmt.Println("outCheck", rcOut)
		_ = runner.Process.Signal(syscall.SIGINT)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rcErr = readCheck(stderr, true)
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
