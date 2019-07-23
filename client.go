package gotcp

import (
	"net"
	"time"
)

// Connect : 连接服务器
func (sess *Session) Connect(address string, derived ISession) bool {
	if sess.IsClosed() == false {
		xlog.Errorln("close session. server address =", sess.RemoteAddr())
		sess.Close()
	}
	conn, err := connectDetail(address)
	if err == nil {
		sess.Init(nil, conn, derived)
		sess.Start()
		xlog.Infoln("connect server success. server address =", sess.RemoteAddr())
		return true
	}
	xlog.Errorln(err)
	return false
}

func connectDetail(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(1 * time.Minute)
	conn.SetNoDelay(true)
	conn.SetWriteBuffer(DefaultSendBuffSize)
	conn.SetReadBuffer(DefaultRecvBuffSize)
	return conn, nil
}
