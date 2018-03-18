package gotcp

import (
	"fmt"
	"net"
)

func GetVaildPort() int {
	port := 10000
	for {
		port = port + 1
		address := fmt.Sprintf(":%d", port)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
		if err != nil {
			continue
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			if conn != nil {
				conn.Close()
			}
			continue
		}
		conn.Close()
		return port
	}
	return 0
}
