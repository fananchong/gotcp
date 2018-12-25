package gotcp

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"time"
)

// Server : TCP 服务器类
type Server struct {
	listener  *net.TCPListener
	sessType  reflect.Type
	ctx       context.Context
	ctxCancel context.CancelFunc
	address   string
}

// RegisterSessType : 注册网络会话类型
func (server *Server) RegisterSessType(v interface{}) {
	server.sessType = reflect.ValueOf(v).Type()
}

// Start : 服务器启动
func (server *Server) Start(address string) bool {
	return server.startDetail(address, true)
}

func (server *Server) startDetail(address string, printError bool) bool {
	server.address = address
	if server.listener != nil {
		return true
	}
	err := server.bind(address)
	if err != nil {
		if printError {
			xlog.Errorln(err)
		}
		return false
	}
	xlog.Infoln("start listen", address)
	server.ctx, server.ctxCancel = context.WithCancel(context.Background())
	go server.loop()
	return true
}

// StartByUnfixedPort : 服务器启动，且会寻找有效端口去监听
func (server *Server) StartByUnfixedPort(ip string, port *uint16) bool {
	for {
		address := fmt.Sprintf("%s:%d", ip, *port)
		if ok := server.startDetail(address, false); ok {
			break
		}
		*port = *port + 1
	}
	return true
}

// Close : 关闭服务器
func (server *Server) Close() {
	if server.ctxCancel != nil {
		server.ctxCancel()
	}
	server.listener.Close()
	server.listener = nil
}

func (server *Server) loop() {
	for {
		select {
		case <-server.ctx.Done():
			xlog.Infoln("server close. address =", server.address)
			return
		default:
			conn, err := server.accept()
			if err == nil {
				sess := reflect.New(server.sessType)
				f := sess.MethodByName("Init")
				f.Call([]reflect.Value{reflect.ValueOf(server.ctx), reflect.ValueOf(conn), sess})
				f = sess.MethodByName("Start")
				f.Call([]reflect.Value{})
				f = sess.MethodByName("RemoteAddr")
				addr := f.Call([]reflect.Value{})
				xlog.Infoln("connect come in. client address =", addr)
			}
		}
	}
}

func (server *Server) bind(address string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	server.listener = listener
	return nil
}

func (server *Server) accept() (*net.TCPConn, error) {
	conn, err := server.listener.AcceptTCP()
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && !opErr.Timeout() {
			xlog.Errorln(err)
		}
		return nil, err
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(1 * time.Minute)
	conn.SetNoDelay(true)
	conn.SetWriteBuffer(128 * 1024)
	conn.SetReadBuffer(128 * 1024)
	return conn, nil
}
