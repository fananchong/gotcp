package main

import (
	"fmt"
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
	atomic.AddInt32(&g_counter, 1)
	this.Send(data, flag)
}

func (this *Echo) OnClose() {
	fmt.Println("Echo.OnClose")
}

var g_counter int32 = 0

func main() {

	go http.ListenAndServe(":8000", nil)

	s := &gotcp.Server{}
	s.RegisterSessType(Echo{})
	var port uint16 = 3000
	s.StartByUnfixedPort("127.0.0.1", &port)

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
