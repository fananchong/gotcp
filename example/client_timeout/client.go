package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
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
}

func (this *Echo) OnClose() {
	fmt.Println("Echo.OnClose")
}

func main() {
	go http.ListenAndServe(":8001", nil)

	echo := &Echo{}
	echo.Connect("localhost:3000", echo)

	tick := time.NewTicker(5 * time.Second)
	for {
		<-tick.C
	}
}
