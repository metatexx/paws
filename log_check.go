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

func logCheck() int {
	var err error
	input := os.Stdin
	if fileGiven {
		input, err = os.Open(file)
		if err != nil {
			fmt.Printf("*** file error: %v", err)
			return 1
		}
		defer func() {
			_ = input.Close()
		}()
	}

	reader := bufio.NewReader(input)

	// for checking patterns we keep the remainder if a line is not fully output so we do not miss matches
	var remainder string
	startTime := time.Now()
	idlePoint := time.Now()

	foundPatters := make(map[*regexp.Regexp]bool, len(okOnPattern))
	out := os.Stdout
	if passThrough {
		out = os.Stderr
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			_, err2 := fmt.Fprintf(out, "*** read error: %v", err)
			if err2 != nil {
				panic(err2)
			}
			return 1
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
