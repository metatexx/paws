package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

func tcpChecks(ports map[string]*url.URL) map[string]string {
	results := make(map[string]string, len(ports))
	for _, uri := range ports {
		results[uri.String()] = checkConn(uri)
	}
	return results
}

func checkConn(uri *url.URL) string {
	protocol := "tcp"
	address := uri.Host
	service := uri.Scheme
	if strings.HasSuffix(service, "-udp") {
		service = service[:len(service)-4]
		protocol = "udp"
	} else if strings.HasSuffix(address, "-tcp") {
		service = service[:len(service)-4]
		protocol = "tcp"
	}

	conn, err := net.DialTimeout(protocol, uri.Host, timeout)
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

	ret := "verified"

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

	case service == "mssql":
		// this is very crude and incorrect for now (but kinda works)
		return IsMSSQLServer(uri)

	case service == "http":

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

	case service == "smtp":
		// reading the info data from the server
		resp, respErr := readResponseString(conn)
		if respErr != nil {
			return respErr.Error()
		}
		if !strings.HasPrefix(resp, "220 ") {
			return fmt.Sprintf("not a SMTP (%q)", resp[:min(len(resp), 32)])
		}

	case uri.RawQuery != "":
		// reading the info data from the server
		resp, respErr := readResponseString(conn)
		if respErr != nil {
			return respErr.Error()
		}
		if !strings.HasPrefix(resp, uri.RawQuery) {
			return fmt.Sprintf("expected %q but got %q", uri.RawQuery, resp[:min(len(resp), 32)])
		}

	default:
		ret = "found"
	}
	return ret
}
