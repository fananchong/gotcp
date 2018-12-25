package gotcp

import (
	"fmt"
	"net"
)

// GetVaildPort : 获取 1 个有效的端口
func GetVaildPort(showmsg bool) int {
	port := 10000
	for {
		port = port + 1
		address := fmt.Sprintf(":%d", port)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
		if err != nil {
			if showmsg {
				xlog.Errorln(err)
			}
			continue
		}
		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			if showmsg {
				xlog.Errorln(err)
			}
			if listener != nil {
				listener.Close()
			}
			continue
		}
		listener.Close()
		return port
	}
}
