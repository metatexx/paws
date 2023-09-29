package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
)

func tcpChecks(ports map[string]string) map[string]string {
	results := make(map[string]string, len(ports))
	for in, service := range ports {
		results[in] = checkConn(in, service)
	}
	return results
}

func checkConn(in, service string) string {
	protocol := "tcp"
	address := in
	if strings.HasSuffix(in, "-udp") {
		address = in[:len(in)-4]
		protocol = "udp"
	} else if strings.HasSuffix(address, "-tcp") {
		address = in[:len(in)-4]
		protocol = "tcp"
	}

	//fmt.Println(protocol, address, timeout)
	conn, err := net.DialTimeout(protocol, address, timeout)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			return opErr.Err.Error()
		}
		return fmt.Sprintf("%#v", err)
	}
	if conn == nil {
		return "nil"
	}
	defer func() {
		_ = conn.Close()
	}()

	// depending on the service we may want to do some health / ping checking
	switch {
	case service == "mysql":
	case service == "mariadb":
		// this is very crude and incorrect for now (but kinda works)
		resp, respErr := readResponseBytes(conn, 4)
		if respErr != nil {
			return respErr.Error()
		}
		if bytes.Compare(resp[1:4], []byte("\x00\x00\x00")) != 0 {
			return fmt.Sprintf("not MySQL/MariaDB (%q)", resp)
		}

	case service == "http":

	case strings.HasPrefix(service, "TEXT/"):
		// reading the info data from the server
		resp, respErr := readResponseString(conn)
		if respErr != nil {
			return respErr.Error()
		}
		if !strings.HasPrefix(resp, service[5:]) {
			return fmt.Sprintf("expected %q but got %q", service[5:], resp[:min(len(resp), 32)])
		}

	case service == "ssh":
		// reading the info data from the server
		resp, respErr := readResponseString(conn)
		if respErr != nil {
			return respErr.Error()
		}
		if !strings.HasPrefix(resp, "SSH-") {
			return fmt.Sprintf("not a SSH server response (%q)", resp[:min(len(resp), 32)])
		}

	case service == "dns":
		// reading the info data from the server
		if protocol == "udp" {
			if !isRealDNSServerUDP(conn) {
				return "not dns"
			}
		} else if !isRealDNSServerTCP(conn) {
			return "not dns"
		}

	case service == "nats":
		// reading the info data from the server
		resp, respErr := readResponseString(conn)
		if respErr != nil {
			return respErr.Error()
		}
		if !strings.HasPrefix(resp, "INFO {") {
			return fmt.Sprintf("not a NATS server response (%q)", resp[:min(len(resp), 32)])
		}
	}
	return "success"
}
