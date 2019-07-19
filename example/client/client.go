package main

import (
	"fmt"
	"math/rand"
	//"net/http"
	//_ "net/http/pprof"
	"strconv"
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

	if g_num != string(data) {
		fmt.Println("g_num = ", g_num)
		fmt.Println("data = ", string(data))
		panic("data error!")
	}

	g_num = strconv.Itoa(int(rand.Int31n(1000)))
	this.Send([]byte(g_num), 0)

	atomic.AddInt32(&g_counter, 1)
}

func (this *Echo) OnClose() {
	fmt.Println("Echo.OnClose")
}

var g_num string
var g_counter int32

func main() {
	//go http.ListenAndServe(":8001", nil)

	echo := &Echo{}
	for !echo.Connect("localhost:30000", echo) {
	}
	g_num = strconv.Itoa(int(rand.Int31n(1000)))
	echo.Send([]byte(g_num), 0)
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
