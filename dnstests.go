package main

import (
	"encoding/binary"
	"github.com/miekg/dns"
	"net"
)

func isRealDNSServerTCP(conn net.Conn) bool {
	// Construct a DNS query
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	msg.RecursionDesired = true
	queryBytes, err := msg.Pack()
	if err != nil {
		return false
	}

	// The DNS query over TCP needs to be prefixed with a 2-byte length
	buf := make([]byte, 2+len(queryBytes))
	binary.BigEndian.PutUint16(buf, uint16(len(queryBytes)))
	copy(buf[2:], queryBytes)

	// Send the DNS query
	_, err = conn.Write(buf)
	if err != nil {
		return false
	}

	// Read and unpack the DNS response
	lengthBytes := make([]byte, 2)
	_, err = conn.Read(lengthBytes)
	if err != nil {
		return false
	}

	length := binary.BigEndian.Uint16(lengthBytes)
	responseBytes := make([]byte, length)
	_, err = conn.Read(responseBytes)
	if err != nil {
		return false
	}

	response := new(dns.Msg)
	err = response.Unpack(responseBytes)
	if err != nil {
		return false
	}

	// Simple check if it's a valid DNS response
	return response.Response
}

func isRealDNSServerUDP(conn net.Conn) bool {
	// Construct a DNS query
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	msg.RecursionDesired = true
	queryBytes, err := msg.Pack()
	if err != nil {
		return false
	}

	// Send the DNS query over UDP
	_, err = conn.Write(queryBytes)
	if err != nil {
		return false
	}

	// Read the DNS response
	responseBytes := make([]byte, 512) // 512 is a common buffer size for DNS UDP
	_, err = conn.Read(responseBytes)
	if err != nil {
		return false
	}

	response := new(dns.Msg)
	err = response.Unpack(responseBytes)
	if err != nil {
		return false
	}

	// Simple check if it's a valid DNS response
	return response.Response
}
