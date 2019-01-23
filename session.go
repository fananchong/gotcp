package gotcp

import (
	"context"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

// ISession : 网络会话类接口
type ISession interface {
	OnRecv(data []byte, flag byte)
	OnClose()
}

const (
	cmdMaxSize    = 128 * 1024 // 消息最大长度
	cmdHeaderSize = 4          // 3字节指令长度 1字节是否压缩
	cmdVerifyTime = 10         // 连接验证超时时间
)

// Session : 网络会话类
type Session struct {
	Conn                  net.Conn
	ctx                   context.Context
	ctxCancel             context.CancelFunc
	sendBuff              *ByteBuffer
	sendChan              chan int
	sendMutex             sync.Mutex
	sendBuffSizeLimit     int // -1 没有限制
	closed                int32
	verified              bool
	verifiedChan          chan int
	Derived               ISession
	closeAfterSendingChan chan int // 这数据发送完毕后，关闭
}

// SetSendBuffSizeLimt : 设置发送缓冲区限制
func (sess *Session) SetSendBuffSizeLimt(limit int) {
	sess.sendBuffSizeLimit = limit
}

// Init : 初始化
func (sess *Session) Init(root context.Context, conn net.Conn, derived ISession) {
	sess.Derived = derived
	sess.Conn = conn
	if root == nil {
		sess.ctx, sess.ctxCancel = context.WithCancel(context.Background())
	} else {
		sess.ctx, sess.ctxCancel = context.WithCancel(root)
	}
	sess.sendBuff = NewByteBuffer()
	sess.sendChan = make(chan int, 1)
	sess.sendBuffSizeLimit = -1
	atomic.StoreInt32(&sess.closed, 0)
	sess.verified = false
	sess.verifiedChan = make(chan int, 1)
	sess.closeAfterSendingChan = make(chan int, 1)
}

// Start : 启动网络会话
func (sess *Session) Start() {
	if atomic.CompareAndSwapInt32(&sess.closed, 0, 1) {
		job := &sync.WaitGroup{}
		job.Add(2)
		go sess.sendloop(job)
		go sess.recvloop(job)
		job.Wait()
	}
}

// Close : 关闭网络会话
func (sess *Session) Close() {
	if atomic.CompareAndSwapInt32(&sess.closed, 1, 2) {
		xlog.Infoln("disconnect. remote address =", sess.RemoteAddr())
		if sess.ctxCancel != nil {
			sess.ctxCancel()
		}
		sess.Conn.Close()
		close(sess.sendChan)
		sess.Derived.OnClose()
		// sess.Derived = nil
	}
}

// CloseAfterSending : 数据发送完毕后，关闭连接
func (sess *Session) CloseAfterSending() {
	sess.closeAfterSendingChan <- 1
}

// IsClosed : 是否已关闭
func (sess *Session) IsClosed() bool {
	return atomic.LoadInt32(&sess.closed) != 1
}

// Verify : 设置已验证标记
func (sess *Session) Verify() {
	if !sess.verified {
		sess.verified = true
		sess.verifiedChan <- 1
	}
}

// IsVerified : 是否已验证
func (sess *Session) IsVerified() bool {
	return sess.verified
}

// SendEx : 发送数据 (buffer 中未包括cmd)
func (sess *Session) SendEx(cmd int, buffer []byte, flag byte) bool {
	if sess.IsClosed() {
		return false
	}
	bsize := len(buffer) + 2
	sess.sendMutex.Lock()
	if sess.sendBuffSizeLimit > 0 && sess.sendBuff.RdSize()+bsize > sess.sendBuffSizeLimit {
		sess.sendMutex.Unlock()
		xlog.Errorln("send buff size limit.")
		sess.Close()
		return false
	}
	header := [6]byte{byte(bsize), byte(bsize >> 8), byte(bsize >> 16), flag, byte(cmd), byte(cmd >> 8)}
	sess.sendBuff.Append(header[:])
	sess.sendBuff.Append(buffer)
	sess.sendMutex.Unlock()
	select {
	case sess.sendChan <- 1:
	default:
	}
	return true
}

// Send : 发送数据 (buffer 中已包括cmd)
func (sess *Session) Send(buffer []byte, flag byte) bool {
	if sess.IsClosed() {
		return false
	}
	bsize := len(buffer)
	sess.sendMutex.Lock()
	if sess.sendBuffSizeLimit > 0 && sess.sendBuff.RdSize()+bsize > sess.sendBuffSizeLimit {
		sess.sendMutex.Unlock()
		xlog.Errorln("send buff size limit.")
		sess.Close()
		return false
	}
	header := [4]byte{byte(bsize), byte(bsize >> 8), byte(bsize >> 16), flag}
	sess.sendBuff.Append(header[:])
	sess.sendBuff.Append(buffer)
	sess.sendMutex.Unlock()
	select {
	case sess.sendChan <- 1:
	default:
	}
	return true
}

// SendRaw : 发送原始数据
func (sess *Session) SendRaw(buffer []byte) bool {
	if sess.IsClosed() {
		return false
	}
	bsize := len(buffer)
	sess.sendMutex.Lock()
	if sess.sendBuffSizeLimit > 0 && sess.sendBuff.RdSize()+bsize > sess.sendBuffSizeLimit {
		sess.sendMutex.Unlock()
		xlog.Errorln("send buff size limit.")
		sess.Close()
		return false
	}
	sess.sendBuff.Append(buffer)
	sess.sendMutex.Unlock()
	select {
	case sess.sendChan <- 1:
	default:
	}
	return true
}

func (sess *Session) recvloop(job *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			xlog.Errorln("[except] ", err, "\n", string(debug.Stack()))
		}
	}()
	defer sess.Close()
	recvBuff := NewByteBuffer()
	job.Done()
	for {
		select {
		case <-sess.ctx.Done():
			return
		default:
			totalsize := recvBuff.RdSize()
			if totalsize < cmdHeaderSize {
				neednum := cmdHeaderSize - totalsize
				recvBuff.WrGrow(neednum)
				readnum, err := io.ReadAtLeast(sess.Conn, recvBuff.WrBuf(), neednum)
				if err != nil {
					xlog.Infoln("recv data fail. error =", err)
					return
				}
				recvBuff.WrFlip(readnum)
				totalsize = recvBuff.RdSize()
			}
			msgbuff := recvBuff.RdBuf()
			datasize := int(msgbuff[0]) | int(msgbuff[1])<<8 | int(msgbuff[2])<<16
			if datasize > cmdMaxSize-cmdHeaderSize {
				xlog.Errorln("data exceed the maximum. datasize =", datasize)
				return
			}
			if datasize <= 0 {
				xlog.Errorln("data length is 0 or negative. datasize =", datasize)
				return
			}
			if totalsize < cmdHeaderSize+datasize {
				neednum := cmdHeaderSize + datasize - totalsize
				recvBuff.WrGrow(neednum)
				readnum, err := io.ReadAtLeast(sess.Conn, recvBuff.WrBuf(), neednum)
				if err != nil {
					xlog.Infoln("recv data fail. error =", err)
					return
				}
				recvBuff.WrFlip(readnum)
				msgbuff = recvBuff.RdBuf()
			}

			sess.Derived.OnRecv(msgbuff[cmdHeaderSize:cmdHeaderSize+datasize], msgbuff[3])
			recvBuff.RdFlip(cmdHeaderSize + datasize)
		}
	}
}

func (sess *Session) sendloop(job *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			xlog.Errorln("[except] ", err, "\n", string(debug.Stack()))
		}
		sess.Close()
	}()
	tmpByte := NewByteBuffer()
	timeout := time.NewTimer(time.Second * cmdVerifyTime)
	defer timeout.Stop()
	job.Done()
	var needClose bool
	for {
		select {
		case <-sess.sendChan:
			for {
				sess.sendMutex.Lock()
				if sess.sendBuff.RdReady() {
					tmpByte.Append(sess.sendBuff.RdBuf()[:sess.sendBuff.RdSize()])
					sess.sendBuff.Reset()
				}
				sess.sendMutex.Unlock()

				if !tmpByte.RdReady() {
					if needClose {
						sess.Close()
						return
					}
					break
				}

				writenum, err := sess.Conn.Write(tmpByte.RdBuf()[:tmpByte.RdSize()])
				if err != nil {
					xlog.Infoln("send data fail. err =", err)
					return
				}
				tmpByte.RdFlip(writenum)
			}
		case <-sess.closeAfterSendingChan:
			sess.sendMutex.Lock()
			if !sess.sendBuff.RdReady() {
				sess.sendMutex.Unlock()
				sess.Close()
				return
			}
			needClose = true
			sess.sendMutex.Unlock()
		case <-sess.ctx.Done():
			return
		case <-sess.verifiedChan:
			timeout.Stop()
		case <-timeout.C:
			xlog.Infoln("verify timeout, remote address =", sess.RemoteAddr())
			return
		}
	}
}

// RemoteAddr : 远端 IP 地址
func (sess *Session) RemoteAddr() string {
	if sess.Conn == nil {
		return ""
	}
	return sess.Conn.RemoteAddr().String()
}
