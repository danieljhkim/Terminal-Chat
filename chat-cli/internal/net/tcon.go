package net

import (
	"fmt"
	"net"
)

// creates a TCP connection to the server
func Connect(serverAddress string) (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", serverAddress, err)
	}
	return conn, nil
}