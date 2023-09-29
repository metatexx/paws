package main

import (
	_ "embed"
	"fmt"
	"github.com/choria-io/fisk"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const appName = "paws"

//go:embed .version
var fullVersion string

var quiet bool
var dots bool
var timeout time.Duration
var wait time.Duration
var delay time.Duration

func main() {
	log.SetOutput(os.Stderr)
	ax := fisk.New(appName, "Port Availability Waiting System: Before you leap, let PAWS take a peep."+
		"\n\nUse --help-long for more info and --cheats for examples"+"\n\n> Anonymous kitten: 'Curiosity checked the port!'").
		Version(fullVersion).
		Author("Copyright 2023 - METATEXX GmbH authors <kontakt@metatexx.de>")

	var host string
	checkCMD := ax.Command("check", "checks the given port list").Default()
	checkCMD.Flag("host", "Default host if none is given in the ports").Short('h').Default("localhost").StringVar(&host)
	checkCMD.Flag("quiet", "Don't output anything (only return code)").Short('q').UnNegatableBoolVar(&quiet)
	checkCMD.Flag("progress", "Output dots while scanning (one dot for each port check)").Short('p').UnNegatableBoolVar(&dots)
	checkCMD.Flag("wait", "Waiting time to wait for the ports to appear").Short('w').DurationVar(&wait)
	checkCMD.Flag("timeout", "Timeout for all connect and read timeouts").Short('t').Default("250ms").DurationVar(&timeout)
	checkCMD.Flag("delay", "Minimal time between checks when waiting").Short('d').Default("250ms").DurationVar(&delay)

	portsRaw := []string{}
	checkCMD.Arg("ports", "the posts to watch for").PlaceHolder("(service:)(host:)port(-udp|-tcp)").Help("Ports to scan (if you need host but without service name you can use use ':host:port'").Required().StringsVar(&portsRaw)

	ax.Cheat("examples", "paws ssh:22")

	ax.MustParseWithUsage(os.Args[1:])

	ports := make(map[string]string)
	//ctx := ax.Context(context.Background(), nil, nil)
	var err error

	for _, portString := range portsRaw {
		parts := strings.SplitN(portString, ":", 3)
		var port string
		var service string
		if len(parts) == 3 {
			_, err = strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(parts[2], "-udp"), "-tcp"))
			ax.FatalIfError(err, "could not parse port information")
			port = parts[1] + ":" + parts[2]
			service = parts[0]
		} else if len(parts) == 2 {
			_, err = strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(parts[1], "-udp"), "-tcp"))
			ax.FatalIfError(err, "could not parse port information")
			port = parts[1]
			service = parts[0]
		} else {
			_, err = strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(parts[0], "-udp"), "-tcp"))
			ax.FatalIfError(err, "could not parse port as int")
			port = parts[0]
			service = parts[0]
		}
		if _, ok := ports[port]; ok {
			ax.Fatalf("duplicate port")
		}
		if !strings.Contains(port, ":") {
			port = net.JoinHostPort(host, port)
		}
		ports[port] = service
	}

	startTime := time.Now()

	for {
		results := tcpChecks(ports)
		allSuccess := true
		startCheck := time.Now()
		for port, resp := range results {
			if resp != "success" {
				allSuccess = false
			} else {
				delete(ports, port)
			}
			if !quiet {
				if dots {
					fmt.Print(".")
				} else {
					fmt.Printf("%s: %q\n", port, resp)
				}
			}
		}
		if allSuccess {
			if !quiet && dots {
				fmt.Println()
			}
			break
		}
		if time.Since(startTime) > wait {
			if !quiet && dots {
				fmt.Println()
			}
			os.Exit(1)
		}
		since := time.Since(startCheck)
		if since < delay {
			time.Sleep(delay - since)
		}
	}
	os.Exit(0)
}
