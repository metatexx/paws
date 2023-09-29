package main

import (
	_ "embed"
	"fmt"
	"github.com/choria-io/fisk"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const appName = "paws"

//go:embed .version
var fullVersion string

//go:embed examples.md
var cheatExamples string

//go:embed examples/docker-mysql.sh
var cheatDocker string

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

	servicesRaw := []string{}
	checkCMD.Arg("services", "the services (ports) to check").
		PlaceHolder("(service(-udp|-tcp):)((host):)port or service(-udp|-tcp)://(user(:pass))@(host):port").
		Help("Services to scan (see 'paws cheat examples' for more information))").Required().StringsVar(&servicesRaw)

	ax.Cheat("examples", cheatExamples)
	ax.Cheat("docker", cheatDocker)

	ax.MustParseWithUsage(os.Args[1:])

	ports := make(map[string]*url.URL)
	//ctx := ax.Context(context.Background(), nil, nil)
	var err error

	for _, serviceString := range servicesRaw {
		uri, _ := url.Parse("noname://" + host + ":1")
		ax.FatalIfError(err, "internal error: could not parse definition")
		if strings.Contains(serviceString, "://") {
			// This is the URI variant for the port definition
			uri, err = url.Parse(serviceString)
			ax.FatalIfError(err, "could not parse port definition")
		} else {
			var frag string
			var parts []string
			parts = strings.SplitN(serviceString, ":", 3)
			if len(parts) == 3 {
				_, err = strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(parts[2], "-udp"), "-tcp"))
				ax.FatalIfError(err, "could not parse port information")
				uri.Scheme = parts[0]
				frag = parts[1] + ":" + parts[2]
			} else if len(parts) == 2 {
				_, err = strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(parts[1], "-udp"), "-tcp"))
				ax.FatalIfError(err, "could not parse port information")
				uri.Scheme = parts[0]
				frag = parts[1]
			} else {
				_, err = strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(parts[0], "-udp"), "-tcp"))
				ax.FatalIfError(err, "could not parse port as int")
				uri.Scheme = parts[0]
				frag = parts[0]
			}
			if !strings.Contains(frag, ":") {
				frag = net.JoinHostPort(host, frag)
			}
			uri.Host = frag
		}
		if _, ok := ports[uri.String()]; ok {
			ax.Fatalf("duplicate check %q", uri.String())
		}
		ports[uri.String()] = uri
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
					fmt.Printf("%s: %s\n", port, resp)
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
