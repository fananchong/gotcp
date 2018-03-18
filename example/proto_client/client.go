package main

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"time"

	"github.com/fananchong/gotcp"
)

type Echo struct {
	gotcp.Session
}

func (this *Echo) OnRecv(data []byte, flag byte) {

	if this.IsVerified() == false {
		this.Verify()
	}

	msg := &MsgEcho{}
	gotcp.DecodeCmd(data, flag, msg)
	if g_num != msg.GetNum() {
		fmt.Println("g_num = ", g_num)
		fmt.Println("data.num = ", msg.GetNum())
		panic("data error!")
	}

	g_num = int32(rand.Int31n(1000))
	msg = &MsgEcho{}
	msg.Num = g_num
	this.SendMsg(0, msg)

	atomic.AddInt32(&g_counter, 1)
}

func (this *Echo) OnClose() {
	fmt.Println("Echo.OnClose")
}

var g_num int32
var g_counter int32

func main() {
	go http.ListenAndServe(":8001", nil)

	echo := &Echo{}
	echo.Connect("localhost:3000", echo)
	g_num = int32(rand.Int31n(1000))

	msg := &MsgEcho{}
	msg.Num = g_num
	echo.SendMsg(0, msg)
	tick := time.NewTicker(5 * time.Second)
	pre := time.Now()
	for {
		select {
		case now := <-tick.C:
			count := atomic.SwapInt32(&g_counter, 0)
			detal := (now.UnixNano() - pre.UnixNano()) / int64(time.Second)
			fmt.Println("count = ", count/int32(detal))
			pre = now
		}
	}
}
