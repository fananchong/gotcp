package main

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strconv"
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

	this.Close()
	go func() {
		time.Sleep(time.Duration(100+rand.Int31n(50)) * time.Millisecond)
		connect()
	}()
}

func (this *Echo) OnClose() {
	fmt.Println("Echo.OnClose")
}

var g_num string

func main() {
	go http.ListenAndServe(":8001", nil)
	for i := 0; i < 100; i++ {
		connect()
	}
	for {
		time.Sleep(100 * time.Second)
	}
}

func connect() {
	echo := &Echo{}
	echo.Connect("localhost:3000", echo)
	g_num = strconv.Itoa(int(rand.Int31n(1000)))
	echo.Send([]byte(g_num), 0)
}
