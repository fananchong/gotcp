package main

import (
	"fmt"
	//"net/http"
	//_ "net/http/pprof"
	"log"
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
	atomic.AddInt64(&counter, 1)
	this.Send(data, flag)
}

func (this *Echo) OnClose() {
	fmt.Println("Echo.OnClose")
}

var counter int64

func main() {

	//go http.ListenAndServe(":8000", nil)

	s := &gotcp.Server{}
	s.RegisterSessType(Echo{})
	s.SetAddress("127.0.0.1", 30000)
	//s.SetUnfixedPort(true)
	s.Start()

	var t = time.Now().UnixNano() / 1e6
	for {
		select {
		case <-time.After(time.Second * 5):
			now := time.Now().UnixNano() / 1e6
			v := atomic.SwapInt64(&counter, 0)
			log.Print("count: ", float64(v)/float64((now-t)/1000), "/s")
			t = now
		}
	}
}
