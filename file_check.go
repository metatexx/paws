package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

func fileCheck() int {
	var err error
	input := os.Stdin
	if fileGiven {
		var fh *os.File
		fh, err = os.Open(file)
		if err != nil {
			fmt.Printf("*** file error: %v\n", err)
			return 1
		}
		defer func() {
			_ = fh.Close()
		}()
		input = fh
	}

	// the rc
	rcChan := make(chan int, 1)
	// Run ReadFileToString function in it's own goroutine and pass back it's
	// response into dataStream channel.
	go func() {
		rcChan <- readCheck(input)
		close(rcChan)
	}()

	// Listen on dataStream channel AND a timeout channel - which ever happens first.
	var rc int
	select {
	case rc = <-rcChan:
	case <-time.After(timeout):
		rc = 5
		fmt.Println("Program execution out of time ")
	}

	if input == os.Stdin {
		err = killPipe()
		if err != nil {
			panic(err)
		}
	}

	if rc == 10 {
		rc = 0
	}
	//fmt.Println("end rc", rc)
	return rc
}

func readCheck(input io.Reader) int {
	// for checking patterns we keep the remainder if a line is not fully output so we do not miss matches
	var err error

	foundPatters := make(map[*regexp.Regexp]bool, len(okOnPattern))
	out := os.Stdout
	if passThrough {
		out = os.Stderr
	}

	if idleTime == 0 {
		idleTime = time.Hour * 24 * 365 * 100
	}

	reader := bufio.NewReader(input)

	chLine := make(chan string)
	chErr := make(chan error)

	quit := make(chan bool, 1)
	defer func() {
		// so that no more reads can happen
		quit <- true
	}()

	go func() {
		var partial string
		for {
			select {
			case <-quit:
				// this will prevent that we try to read while closing the app
				return

			default:
				// Do other stuff
				var line string
				line, err = reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						_, err2 := fmt.Fprintf(out, "*** read error: %v\n", err)
						if err2 != nil {
							panic(err2)
						}
						chErr <- err
						return
					} else {
						partial += line
						// if we read a file, we are in "tail" mode and wait on more input
						if input != os.Stdin {
							time.Sleep(250 * time.Millisecond) // Sleep a bit
						} else {
							// we can't "tail" on stdin because we can't kill the pipe early with a meaningful return code
							chErr <- io.EOF
							return
						}
					}
				} else {
					if passThrough {
						// we output as it comes in
						fmt.Print(line)
					}
					chLine <- partial + line
					partial = ""
				}
			}
		}
	}()

	for {
		select {
		case line := <-chLine:
			// for our regex testing we want full lines
			check := []byte(strings.TrimSuffix(line, "\n"))
			// check out patterns
			for _, re := range failOnPattern {
				if re.Match(check) {
					if printMatches {
						if passThrough {
							_, err = fmt.Fprintln(os.Stderr, highlight(string(check), re))
						} else {
							fmt.Println(highlight(string(check), re))
						}
					}
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
				for _, re := range okOnPattern {
					if b, ok := foundPatters[re]; ok {
						if b {
							continue
						}
					}
					if re.Match(check) {
						if printMatches {
							if passThrough {
								_, err = fmt.Fprintln(os.Stderr, highlight(string(check), re))
							} else {
								fmt.Println(highlight(string(check), re))
							}
						}
						foundPatters[re] = true
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

		case err := <-chErr:
			// if we land here with EOF (on stdin) it is always a fail
			fmt.Println("Error reading:", err)
			return 1
		case <-time.After(idleTime):
			fmt.Println("success: idle state reached")
			return 10
		}
	}
}
