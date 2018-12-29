package gotcp

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
)

// SendMsg : 发送 protobuf 消息
func (sess *Session) SendMsg(cmd uint64, msg proto.Message) bool {
	data, flag, err := Encode(cmd, msg)
	if err != nil {
		return false
	}
	fmt.Println("send data =====================")
	fmt.Println("cmd:", cmd)
	fmt.Println("data len:", len(data))
	fmt.Println("flag:", flag)
	return sess.SendEx((int)(cmd), data, flag)
}
