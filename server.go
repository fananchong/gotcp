package gotcp

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// Server : TCP 服务器类
type Server struct {
	listener    *net.TCPListener
	ctx         context.Context
	ctxCancel   context.CancelFunc
	sessType    reflect.Type
	address     string
	unfixedPort bool
	realPort    int32
	userdata    interface{}
}

// RegisterSessType : 注册网络会话类型
func (server *Server) RegisterSessType(v interface{}) {
	server.sessType = reflect.ValueOf(v).Type()
}

// SetAddress : 设置地址
func (server *Server) SetAddress(address string, port int32) {
	server.address = address
	server.realPort = port
}

// GetAddress 获取地址
func (server *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", server.address, server.realPort)
}

// SetUnfixedPort : 值为 True ，则寻找有效端口去监听
func (server *Server) SetUnfixedPort(v bool) {
	server.unfixedPort = v
}

// GetRealPort : 获取最终监听的端口
func (server *Server) GetRealPort() int32 {
	return server.realPort
}

// SetUserData : 设置自定义数据
func (server *Server) SetUserData(v interface{}) {
	server.userdata = v
}

// Start : 服务器启动
func (server *Server) Start() bool {
	if server.unfixedPort == false {
		address := fmt.Sprintf("%s:%d", server.address, server.realPort)
		return server.startDetail(address, true)
	}
	return server.startByUnfixedPort(server.address, &server.realPort)
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
	go server.loop(nil)
	return true
}

func (server *Server) startByUnfixedPort(ip string, port *int32) bool {
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

// GetSessionType : 获取 Session 类型
func (server *Server) GetSessionType() reflect.Type {
	return server.sessType
}

func (server *Server) loop(fn func(s interface{})) {
	for {
		select {
		case <-server.ctx.Done():
			xlog.Infoln("server close. address =", server.address)
			return
		default:
			conn, err := server.accept()
			if err == nil && server.sessType != nil {
				func() {
					defer func() {
						if err := recover(); err != nil {
							xlog.Errorln("[except] ", err, "\n", string(debug.Stack()))
						}
					}()
					sess := reflect.New(server.sessType)
					f := sess.MethodByName("Init")
					if server.userdata == nil {
						f.Call([]reflect.Value{reflect.ValueOf(server.ctx), reflect.ValueOf(conn), sess})
					} else {
						f.Call([]reflect.Value{reflect.ValueOf(server.ctx), reflect.ValueOf(conn), sess, reflect.ValueOf(server.userdata)})
					}
					f = sess.MethodByName("Start")
					f.Call([]reflect.Value{})
					f = sess.MethodByName("RemoteAddr")
					addr := f.Call([]reflect.Value{})
					xlog.Infoln("connect come in. client address =", addr)
					if fn != nil {
						fn(sess.Interface())
					}
				}()
			} else {
				if conn != nil {
					xlog.Errorln("you need call RegisterSessType, to register session type.")
					conn.Close()
				}
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
	conn.SetWriteBuffer(DefaultSendBuffSize)
	conn.SetReadBuffer(DefaultRecvBuffSize)
	return conn, nil
}

// Listen listen
func (server *Server) Listen(addr string) (err error) {
	addrs := strings.Split(addr, ":")
	if addrs[0] == "" {
		addrs[0] = "0.0.0.0"
	}
	server.address = addrs[0]
	var port int
	if port, err = strconv.Atoi(addrs[1]); err != nil {
		return
	}
	server.realPort = int32(port)
	err = server.bind(addr)
	if err != nil {
		return
	}
	server.ctx, server.ctxCancel = context.WithCancel(context.Background())
	xlog.Infoln("start listen", addr)
	return
}

// Accept accept
func (server *Server) Accept(fn func(s interface{})) error {
	server.loop(fn)
	return nil
}
