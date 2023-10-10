package main

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func fileCheck() int {
	ports := make(map[string]*url.URL)
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
			if resp != "found" && resp != "verified" {
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
			return 1
		}
		since := time.Since(startCheck)
		if since < delay {
			time.Sleep(delay - since)
		}
	}
	return 0
}
