package main

import (
	"fmt"
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
	atomic.AddInt32(&g_counter, 1)
	this.Send(data, flag)
}

func (this *Echo) OnClose() {

}

var g_counter int32 = 0

func main() {
	s := &gotcp.Server{}
	s.RegisterSessType(Echo{})
	s.Start(":3000")

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
