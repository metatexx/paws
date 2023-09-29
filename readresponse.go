package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"
)

func readResponseString(conn net.Conn) (string, error) {
	var response string

	// set SetReadDeadline
	err := conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		panic(fmt.Errorf("SetReadDeadline failed: %v", err))
	}
	response, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			if opErr.Timeout() {
				return "", errors.New("no data (read timeout)")
			}
			return "", fmt.Errorf("read: %s", opErr.Err.Error())
		}
		return "", fmt.Errorf("read: %s", err.Error())
	}
	return response, nil
}

func readResponseBytes(conn net.Conn, len int) ([]byte, error) {
	// set SetReadDeadline
	err := conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		panic(fmt.Errorf("SetReadDeadline failed: %v", err))
	}
	buf := make([]byte, len)
	var l int
	l, err = conn.Read(buf)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			if opErr.Timeout() {
				return nil, errors.New("no data (read timeout)")
			}
			return nil, fmt.Errorf("read: %x", opErr.Err.Error())
		}
		return nil, fmt.Errorf("read: %s", err.Error())
	}
	return buf[:l], nil
}
