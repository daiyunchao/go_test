package my_net

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

const (
	MsgTypeHandshake = 1
)

const (
	MsgRouteLogin = 1
)

func TestClient(t *testing.T) {
	createConn()
}

func createConn() {
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		return
	}
	SendMsg(conn, MsgTypeHandshake, MsgRouteLogin, "这是客户端发的消息")
	for {
		buffer := make([]byte, 8)
		_, err = conn.Read(buffer)
		if err != nil {
			return
		}
		fmt.Printf("client rev: %s", buffer)
	}
}

func SendMsg(conn net.Conn, msgType int, msgRoute int64, msg string) {
	//消息类型2字节
	var data = make([]byte, 0)
	msgTypeBuff := make([]byte, 2)
	binary.BigEndian.PutUint16(msgTypeBuff, uint16(msgType))
	data = append(data, msgTypeBuff...)

	//路由8字节
	msgRouteBuff := make([]byte, 8)
	binary.BigEndian.PutUint64(msgRouteBuff, uint64(msgRoute))
	data = append(data, msgRouteBuff...)

	//消息长度4字节
	msgLenBuff := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuff, uint32(len(msg)))
	data = append(data, msgLenBuff...)

	//消息
	data = append(data, []byte(msg)...)
	_, err := conn.Write(data)
	if err != nil {
		conn.Close()
		return
	}
}
